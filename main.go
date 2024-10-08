package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"r1/r1/Apis"
	"r1/r1/Server/Handlers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	serviceCollection := client.Database("test").Collection("services")
	rappCollection := client.Database("test").Collection("rapps")
	subscriptionsCollection := client.Database("test").Collection("subscriptions")
	subscribersCollection := client.Database("test").Collection("subscribers")
	dataTypeProdCapsCollection := client.Database("test").Collection("dataTypeProdCaps")
	dataJobsCollection := client.Database("test").Collection("dataJobs")

	//For testing purpose: insert a few documents into the rapps collection
	newRapps := []interface{}{
		Apis.Rapp{ApfId: "testrapp1", IsAuthorized: true, AuthorizedServices: []string{}},
		Apis.Rapp{ApfId: "testrapp2", IsAuthorized: false, AuthorizedServices: []string{}},
	}
	_, err = rappCollection.InsertMany(context.TODO(), newRapps)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/allServiceAPIs", Handlers.ServiceDiscoveryHandler(serviceCollection, rappCollection)).Methods("GET")
	r.HandleFunc("/{subscriberId}/subscriptions", Handlers.CreateSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("POST")
	r.HandleFunc("/{subscriberId}/subscriptions/{subscriptionId}", Handlers.DeleteSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("DELETE")
	r.HandleFunc("/{subscriberId}/subscriptions/{subscriptionId}", Handlers.UpdateSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("PUT")
	r.HandleFunc("/{subscriberId}/subscriptions/{subscriptionId}", Handlers.PatchSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("PATCH")
	r.HandleFunc("/{apfId}/service-apis", Handlers.PublishServiceHandler(serviceCollection, rappCollection)).Methods("POST")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.GetSpecificServiceAPIHandler(serviceCollection)).Methods("GET")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.UpdateServiceAPIHandler(serviceCollection, rappCollection)).Methods("PUT")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.DeleteServiceAPIHandler(serviceCollection, rappCollection)).Methods("DELETE")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.PatchServiceAPIHandler(serviceCollection, rappCollection)).Methods("PATCH")

	// Data registration API routes
	r.HandleFunc("/rapps/{rAppId}/datatypeprodcaps", Handlers.RegisterDmeTypeProdCapHandler(rappCollection, dataTypeProdCapsCollection)).Methods("POST")
	r.HandleFunc("/rapps/{rAppId}/datatypeprodcaps/{registrationId}", Handlers.DeregisterDmeTypeProdCapHandler(rappCollection, dataTypeProdCapsCollection)).Methods("DELETE")

	r.HandleFunc("/datatypes", Handlers.GetAllDataTypesHandler(dataTypeProdCapsCollection)).Methods("GET")
	r.HandleFunc("/datatypes/{dataTypeId}", Handlers.GetDataTypeByIdHandler(dataTypeProdCapsCollection)).Methods("GET")

	// Data Access API routes
	r.HandleFunc("/{consumerId}/dataJobs", Handlers.CreateDataJobHandler(dataJobsCollection)).Methods("POST")
	r.HandleFunc("/{consumerId}/dataJobs/{dataJobId}", Handlers.DeleteDataJobHandler(dataJobsCollection)).Methods("DELETE")
	r.HandleFunc("/notifyDataAvailability", Handlers.NotifyDataAvailabilityHandler(dataJobsCollection)).Methods("POST")

	// Register the Push Data handler
	r.HandleFunc("/api/v1/push-data", Handlers.PushDataHandler()).Methods("POST")
	r.HandleFunc("/api/v1/pull-data", Handlers.PullDataHandler()).Methods("GET")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")
}
