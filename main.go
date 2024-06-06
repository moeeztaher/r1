package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "r1/r1/Server/Handlers"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/service-apis", Handlers.PublishServiceHandler).Methods("POST")
    r.HandleFunc("/service-apis", Handlers.PublishServiceHandler).Methods("GET")

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
