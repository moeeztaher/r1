package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"r1/r1/Server/Handlers"
	"r1/r1/Apis"
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

	// For testing purpose: insert a few test rapps into the rapps collection
	newRapps := []interface{}{
		Apis.Rapp{ApfId: "testrapp1", IsAuthorized: true, AuthorizedServices: []string{}},
		Apis.Rapp{ApfId: "testrapp2", IsAuthorized: false, AuthorizedServices: []string{}},
	}
	_, err = rappCollection.InsertMany(context.TODO(), newRapps)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/allServiceAPIs", Handlers.ServiceDiscoveryHandler(serviceCollection)).Methods("GET")
	r.HandleFunc("/{subscriberId}/subscriptions", Handlers.CreateSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("POST")
	r.HandleFunc("/{subscriberId}/subscriptions/{subscriptionId}", Handlers.DeleteSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("DELETE")
	r.HandleFunc("/{subscriberId}/subscriptions/{subscriptionId}", Handlers.UpdateSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("PUT")
	r.HandleFunc("/{subscriberId}/subscriptions/{subscriptionId}", Handlers.PatchSubscriptionHandler(subscriptionsCollection, subscribersCollection)).Methods("PATCH")
	r.HandleFunc("/{apfId}/service-apis", Handlers.PublishServiceHandler(serviceCollection, rappCollection)).Methods("POST")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.GetSpecificServiceAPIHandler(serviceCollection)).Methods("GET")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.UpdateServiceAPIHandler(serviceCollection)).Methods("PUT")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.DeleteServiceAPIHandler(serviceCollection)).Methods("DELETE")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.PatchServiceAPIHandler(serviceCollection)).Methods("PATCH")

	
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")
}
