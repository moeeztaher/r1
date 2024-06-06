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

func GetServiceApisHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    apfId := vars["apfId"]

    filteredRequests := serviceRequests //add filtering logic

    respondWithJSON(w, http.StatusOK, filteredRequests)
}

func UpdateServiceApiHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    apfId := vars["apfId"]
    serviceApiId := vars["serviceApiId"]

    var updatedRequest Apis.PublishServiceRequest

    err := json.NewDecoder(r.Body).Decode(&updatedRequest)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    for i, request := range serviceRequests {
        if request.ID == serviceApiId {
            serviceRequests[i] = updatedRequest
            respondWithJSON(w, http.StatusOK, updatedRequest)
            return
        }
    }

    respondWithError(w, http.StatusNotFound, "Service not found")
}

func DeleteServiceApiHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    apfId := vars["apfId"]
    serviceApiId := vars["serviceApiId"]

    for i, request := range serviceRequests {
        if request.ID == serviceApiId {
            serviceRequests = append(serviceRequests[:i], serviceRequests[i+1:]...)
            w.WriteHeader(http.StatusNoContent)
            return
        }
    }

    respondWithError(w, http.StatusNotFound, "Service not found")
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
