package main

import (
    "log"
    "net/http"
    "context"
    "fmt"

    "github.com/gorilla/mux"
    "r1/r1/Server/Handlers/registrationServiceHandler"
    "r1/r1/Server/Handlers/errorHandler"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"

    "github.com/go-redis/redis/v8"
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

    // Get a handle for your collection
    collection := client.Database("test2").Collection("PublishAPITest")

    r := mux.NewRouter()

    r.Use(ErrorHandler.CheckContentType)
    r.Use(ErrorHandler.CheckContentLength)
    r.Use(ErrorHandler.CheckPayloadSize)

    r.HandleFunc("/service-apis", RegistrationServiceHandlers.PublishServiceHandler1(collection)).Methods("POST")
    r.HandleFunc("/service-apis", RegistrationServiceHandlers.GetServiceAPIsHandler(collection)).Methods("GET")
    r.HandleFunc("/service-apis/{serviceApiId}", RegistrationServiceHandlers.GetSpecificServiceAPIHandler(collection)).Methods("GET")
    r.HandleFunc("/service-apis/{serviceApiId}", RegistrationServiceHandlers.UpdateServiceAPIHandler(collection)).Methods("PUT")
    r.HandleFunc("/service-apis/{serviceApiId}", RegistrationServiceHandlers.DeleteServiceAPIHandler(collection)).Methods("DELETE")
    r.HandleFunc("/service-apis/{serviceApiId}", RegistrationServiceHandlers.PatchServiceAPIHandler(collection)).Methods("PATCH")

    r.NotFoundHandler = http.HandlerFunc(ErrorHandler.NotFoundHandler)

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))

    // Close the connection once done
    err = client.Disconnect(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Connection to MongoDB closed.")
}
