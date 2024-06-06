package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "r1/r1/Server/Handlers"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/publish", Handlers.PublishServiceHandler).Methods("POST")

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
