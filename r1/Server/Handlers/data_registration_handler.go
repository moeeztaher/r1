package Handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"r1/r1/Apis"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// POST handler for registering a new DME type production capability
func RegisterDmeTypeProdCapHandler(rappCollection, dataTypeProdCapsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rAppId := vars["rAppId"]

		var registration Apis.DataTypeProdCapRegistration
		err := json.NewDecoder(r.Body).Decode(&registration)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Insert the registration into the dataTypeProdCaps collection
		insertResult, err := dataTypeProdCapsCollection.InsertOne(context.TODO(), registration)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error inserting into database")
			return
		}

		// Assert the InsertedID to ObjectID and get the hex value
		objectId, ok := insertResult.InsertedID.(primitive.ObjectID)
		if !ok {
			respondWithError(w, http.StatusInternalServerError, "Error inserting into database")
			return
		}

		// Update the Rapp collection to include the new registration
		filter := bson.M{"apf_id": rAppId}
		update := bson.M{"$push": bson.M{"dataTypeProdCaps": objectId.Hex()}}

		opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
		var updatedRapp bson.M
		err = rappCollection.FindOneAndUpdate(context.TODO(), filter, update, opts).Decode(&updatedRapp)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error updating database")
			return
		}

		locationHeader := fmt.Sprintf("/rapps/%s/datatypeprodcaps/%s", rAppId, objectId.Hex())
		w.Header().Set("Location", locationHeader)
		respondWithJSON(w, http.StatusCreated, registration)
	}
}

// DELETE handler for deregistering a DME type production capability
func DeregisterDmeTypeProdCapHandler(rappCollection, dataTypeProdCapsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rAppId := vars["rAppId"]
		registrationId := vars["registrationId"]

		// Convert registrationId to ObjectID
		objectId, err := primitive.ObjectIDFromHex(registrationId)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid registration ID")
			return
		}

		// Check if the Rapp entry contains the registrationId
		rappFilter := bson.M{"apf_id": rAppId, "dataTypeProdCaps": registrationId}
		rappResult := rappCollection.FindOne(context.TODO(), rappFilter)

		var rappDoc bson.M
		err = rappResult.Decode(&rappDoc)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				respondWithError(w, http.StatusNotFound, "Registration ID not found in Rapp")
			} else {
				respondWithError(w, http.StatusInternalServerError, "Error checking Rapp collection")
			}
			return
		}

		// Remove the registration from the dataTypeProdCaps collection
		filter := bson.M{"_id": objectId}
		_, err = dataTypeProdCapsCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error deleting from database")
			return
		}

		// Update the Rapp collection to remove the registration reference
		update := bson.M{"$pull": bson.M{"dataTypeProdCaps": registrationId}}
		updateResult, err := rappCollection.UpdateOne(context.TODO(), rappFilter, update)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error updating database")
			return
		}

		// Check if any document was modified
		if updateResult.MatchedCount == 0 {
			respondWithError(w, http.StatusNotFound, "Rapp entry not found or registration ID already removed")
			return
		}

		// If ModifiedCount is zero, it means the document was found but not modified (could be already removed)
		if updateResult.ModifiedCount == 0 {
			respondWithError(w, http.StatusNotFound, "Registration ID not found in Rapp entry")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// Helper functions to respond with JSON or error
func respondWithError(w http.ResponseWriter, code int, message string) {
	details := Apis.ProblemDetails{
		Type:   "error",
		Title:  http.StatusText(code),
		Status: code,
		Detail: message,
	}
	respondWithJSON(w, code, details)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
