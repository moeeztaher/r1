package Handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"r1/r1/Apis"
	"log"
	"fmt"
)

func CreateSubscriptionHandler(subscriptionsCollection *mongo.Collection, subscribersCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		subscriberID := params["subscriberId"]
	
		var subscription Apis.Subscription
		err := json.NewDecoder(r.Body).Decode(&subscription)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
		subscriptionID, err := insertSubscription(&subscription, subscriptionsCollection)
		if err != nil {
			http.Error(w, "Failed to create subscription", http.StatusInternalServerError)
			return
		}
	
		err = updateSubscriberWithSubscription(subscriberID, subscriptionID, subscribersCollection)
		if err != nil {
			http.Error(w, "Failed to update subscriber with subscription", http.StatusInternalServerError)
			return
		}
	
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(subscription)
}
}

func insertSubscription(subscription *Apis.Subscription, subscriptionsCollection *mongo.Collection) (primitive.ObjectID, error) {

	result, err := subscriptionsCollection.InsertOne(context.Background(), subscription)
	if err != nil {
		log.Println("Error inserting subscription:", err)
		return primitive.NilObjectID, err
	}
	return result.InsertedID.(primitive.ObjectID), nil
}


func updateSubscriberWithSubscription(subscriberID string, subscriptionID primitive.ObjectID, subscribersCollection *mongo.Collection) error {

	subscriptionObjectID := subscriptionID
	var subscriber Apis.Subscriber
	err := subscribersCollection.FindOne(context.Background(), bson.M{"subscriberId": subscriberID}).Decode(&subscriber)
	if err != nil {
		subscriber = Apis.Subscriber{
			SubscriberID: subscriberID,
			SubscriptionIds: []primitive.ObjectID{subscriptionObjectID},
		}
		_, err := subscribersCollection.InsertOne(context.Background(), subscriber)
		if err != nil {
			log.Println("Error inserting new subscriber:", err)
			return err
		}
		return nil
	}

	update := bson.M{"$addToSet": bson.M{"subscriptionIds": subscriptionObjectID}}
	_, err = subscribersCollection.UpdateOne(context.Background(), bson.M{"subscriberId": subscriberID}, update)
	if err != nil {
		log.Println("Error updating subscriber with subscription:", err)
		return err
	}
	return nil
}


func DeleteSubscriptionHandler(subscriptionsCollection *mongo.Collection, subscribersCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	subscriberId := params["subscriberId"]
	subscriptionId := params["subscriptionId"]

	subscriptionObjectId, err := primitive.ObjectIDFromHex(subscriptionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid subscriptionId",
		})
		return
	}

	deleteResult, err := subscriptionsCollection.DeleteOne(context.TODO(), bson.M{"_id": subscriptionObjectId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to delete subscription",
		})
		return
	}

	if deleteResult.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Subscription not found",
		})
		return
	}

	updateResult, err := subscribersCollection.UpdateOne(context.TODO(),
		bson.M{"subscriberid": subscriberId},
		bson.M{"$pull": bson.M{"subscriptionIds": subscriptionObjectId}},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to update subscriber",
		})
		return
	}

	if updateResult.ModifiedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Subscriber not found or update did not modify any document",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Subscription %s deleted successfully for subscriber %s", subscriptionId, subscriberId),
	})
}
}

func UpdateSubscriptionHandler(subscriptionsCollection *mongo.Collection, subscribersCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	subscriberId := params["subscriberId"]
	subscriptionId := params["subscriptionId"]

	var updatedSubscription Apis.Subscription
	err := json.NewDecoder(r.Body).Decode(&updatedSubscription)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request payload",
		})
		return
	}

	count, err := subscribersCollection.CountDocuments(context.TODO(), bson.M{"subscriberid": subscriberId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to check subscriber",
		})
		return
	}
	if count == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Subscriber not found",
		})
		return
	}

	subscriptionObjectId, err := primitive.ObjectIDFromHex(subscriptionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid subscriptionId",
		})
		return
	}

	deleteResult, err := subscriptionsCollection.DeleteOne(context.TODO(), bson.M{"_id": subscriptionObjectId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to delete existing subscription",
		})
		return
	}

	if deleteResult.DeletedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Subscription not found",
		})
		return
	}

	_, err = subscriptionsCollection.InsertOne(context.TODO(), bson.M{
		"_id":                     subscriptionObjectId, // Reuse the same ObjectId
		"events":                  updatedSubscription.Events,
		"eventFilters":            updatedSubscription.EventFilters,
		"eventReq":                updatedSubscription.EventReq,
		"notificationDestination": updatedSubscription.NotificationDestination,
		"requestTestNotification": updatedSubscription.RequestTestNotification,
		"websockNotifConfig":      updatedSubscription.WebsockNotifConfig,
		"supportedFeatures":       updatedSubscription.SupportedFeatures,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to insert updated subscription",
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Subscription %s created/updated successfully for subscriber %s", subscriptionId, subscriberId),
	})
}
}

func PatchSubscriptionHandler(subscriptionsCollection *mongo.Collection, subscribersCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	subscriberId := params["subscriberId"]
	subscriptionId := params["subscriptionId"]

	var patchData map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&patchData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid request payload",
		})
		return
	}

	count, err := subscribersCollection.CountDocuments(context.TODO(), bson.M{"subscriberid": subscriberId})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to check subscriber",
		})
		return
	}
	if count == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Subscriber not found",
		})
		return
	}

	subscriptionObjectId, err := primitive.ObjectIDFromHex(subscriptionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Invalid subscriptionId",
		})
		return
	}
	updateFields := bson.M{}
	for key, value := range patchData {
		updateFields[key] = value
	}
	updateFields["_id"] = subscriptionObjectId // Ensure _id remains the same

	updateResult, err := subscriptionsCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": subscriptionObjectId},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Failed to update subscription",
		})
		return
	}

	if updateResult.ModifiedCount == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": "Subscription not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Subscription %s patched successfully for subscriber %s", subscriptionId, subscriberId),
	})
}
}