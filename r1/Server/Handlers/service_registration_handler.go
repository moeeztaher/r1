package Handlers

import (
    "encoding/json"
    "net/http"

    "r1/r1/Apis"
)
func PublishServiceHandler(w http.ResponseWriter, r *http.Request) {
    var request Apis.PublishServiceRequest

    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }


    response := request 
    respondWithJSON(w, http.StatusCreated, response)
}

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
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}
