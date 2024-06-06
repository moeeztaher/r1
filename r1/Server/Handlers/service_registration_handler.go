package Handlers

import (
    "encoding/json"
    "net/http"
    "fmt"
    "github.com/gorilla/mux"
    "r1/r1/Apis"
    "io/ioutil"
    "os"
    "path/filepath"
)

var serviceRequests []Apis.PublishServiceRequest

func PublishServiceHandler(w http.ResponseWriter, r *http.Request) {
    var request Apis.PublishServiceRequest

    err := json.NewDecoder(r.Body).Decode(&request)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    serviceInfo, err := loadServiceMockResponse()
    if err != nil {
        fmt.Println("Error loading service name:", err)
        respondWithError(w, http.StatusInternalServerError, "Error loading service name")
        return
    }

    serviceRequests = append(serviceRequests, request)
    respondWithJSON(w, http.StatusCreated, serviceInfo)
}

func GetServiceApisHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    apfId := vars["apfId"]
    _ = apfId
    filteredRequests := serviceRequests 

    respondWithJSON(w, http.StatusOK, filteredRequests)
}

func UpdateServiceApiHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    apfId := vars["apfId"]
    _ = apfId
    serviceApiId := vars["serviceApiId"]

    var updatedRequest Apis.PublishServiceRequest

    err := json.NewDecoder(r.Body).Decode(&updatedRequest)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    for i, request := range serviceRequests {
        if request.ApiId == serviceApiId {
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
    _ = apfId
    serviceApiId := vars["serviceApiId"]

    for i, request := range serviceRequests {
        if request.ApiId == serviceApiId {
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


func loadServiceMockResponse() (Apis.ServiceInfoMock, error) {
    dir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error:", err)
        return Apis.ServiceInfoMock{}, err
    }
    fmt.Println("Current directory:", dir)
    filename := "examples/TS29222_CAPIF_Publish_Service_API.json"
    fullPath := filepath.Join(dir, filename)
    jsonFile, err := ioutil.ReadFile(fullPath)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return Apis.ServiceInfoMock{}, err
    }

    var serviceInfo Apis.ServiceInfoMock
    err = json.Unmarshal(jsonFile, &serviceInfo)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return Apis.ServiceInfoMock{}, err
    }

    return serviceInfo, nil
}