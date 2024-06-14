package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"r1/r1/Server/Handlers"
)

func main() {
	// Set client options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	serviceCollection := client.Database("test").Collection("services")
	rappCollection := client.Database("test").Collection("rapps")

	// Insert a few test rapps into the rapps collection
	//newRapps := []interface{}{
	//	Rapp{ApfId: "testrapp1", IsAuthorized: true, AuthorizedServices: []string{}},
	//	Rapp{ApfId: "testrapp2", IsAuthorized: false, AuthorizedServices: []string{}},
	//}
	//_, err = rappCollection.InsertMany(context.TODO(), newRapps)
	//if err != nil {
	//	panic(err)
	//}

	r := mux.NewRouter()
	r.HandleFunc("/allServiceApis", Handlers.GetServiceAPIsHandler(serviceCollection, rappCollection)).Methods("GET")
	r.HandleFunc("/{apfId}/service-apis", Handlers.PublishServiceHandler(serviceCollection, rappCollection)).Methods("POST")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.GetSpecificServiceAPIHandler(serviceCollection)).Methods("GET")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.UpdateServiceAPIHandler(serviceCollection)).Methods("PUT")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.DeleteServiceAPIHandler(serviceCollection)).Methods("DELETE")
	r.HandleFunc("/{apfId}/service-apis/{serviceApiId}", Handlers.PatchServiceAPIHandler(serviceCollection)).Methods("PATCH")

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	// Close the connection once done
	err = client.Disconnect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection to MongoDB closed.")
}
