package RegistrationServiceHandlers

import (
    "encoding/json"
    "net/http"
    "fmt"
    "r1/r1/Apis"
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"

    "github.com/gorilla/mux"

    "./Handlers/errorHandler"
)

/*func PublishServiceHandler(w http.ResponseWriter, r *http.Request) {
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
    }*/



func PublishServiceHandler1(collection *mongo.Collection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // HTTP 401
        if !ErrorHandler.checkAuthorization(w,r) {
            return
        }
        // HTTP 403
        if !ErrorHandler.checkPermissions(w,r) {
            return
        }

        var api Apis.PublishServiceAPI
        //apfId := params["apfId"]

        // Parse the request body into the ServiceAPI struct
        err := json.NewDecoder(r.Body).Decode(&api)
        if err != nil {
            // Handle JSON decoding error (HTTP 400)
			errorResponse := Apis.ErrorResponse{
				Type:             "https://example.com/errors/json-decoding",
				Title:            "JSON Decoding Error",
				Status:           http.StatusBadRequest,
				Detail:           "Failed to decode JSON body",
				Cause:            err.Error(),
				SupportedFeatures: "string",
			}
			ErrorHandler.writeErrorResponse(w, errorResponse, http.StatusBadRequest)
            return
        }

        // Insert the received API into the MongoDB collection
        _, err = collection.InsertOne(context.TODO(), api)
        if err != nil {
            // Handle MongoDB insertion error (HTTP 500)
			errorResponse := Apis.ErrorResponse{
				Type:             "https://example.com/errors/mongodb-insertion",
				Title:            "MongoDB Insertion Error",
				Status:           http.StatusInternalServerError,
				Detail:           "Failed to insert document into MongoDB",
				Cause:            err.Error(),
				SupportedFeatures: "string",
			}
			ErrorHandler.writeErrorResponse(w, errorResponse, http.StatusInternalServerError)
            return
        }

        // Respond with a success message HTTP 201
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "Service API published successfully"})
    }
}

func GetServiceAPIsHandler(collection *mongo.Collection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var serviceAPIs []Apis.GetServiceAPI

        // Find all documents in the collection
        cursor, err := collection.Find(context.TODO(), bson.D{})
        if err != nil {
            http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
            return
        }
        defer cursor.Close(context.TODO())

        // Iterate through the cursor and decode each document into a ServiceAPI struct
        for cursor.Next(context.TODO()) {
            var api Apis.GetServiceAPI
            if err = cursor.Decode(&api); err != nil {
                http.Error(w, "Failed to decode document", http.StatusInternalServerError)
                return
            }
            serviceAPIs = append(serviceAPIs, api)
        }

        // Check for cursor errors
        if err := cursor.Err(); err != nil {
            http.Error(w, "Cursor error", http.StatusInternalServerError)
            return
        }

        // Respond with the retrieved service APIs
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(serviceAPIs)
    }
}

func GetSpecificServiceAPIHandler(collection *mongo.Collection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        params := mux.Vars(r)
        //apfId := params["apfId"]
        serviceApiId := params["serviceApiId"]

        var serviceAPI Apis.GetServiceAPI

        // Find the document with the specified apfId and serviceApiId
        filter := bson.M{"apiid": serviceApiId}
        err := collection.FindOne(context.TODO(), filter).Decode(&serviceAPI)
        if err != nil {
            http.Error(w, "Failed to retrieve document from MongoDB", http.StatusInternalServerError)
            return
        }

        // Respond with the retrieved service API
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(serviceAPI)
    }
}

func UpdateServiceAPIHandler(collection *mongo.Collection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        serviceApiId := vars["serviceApiId"]

        var apiData Apis.ApiData
        err := json.NewDecoder(r.Body).Decode(&apiData)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        filter := bson.M{"apiid": serviceApiId}
        update := bson.M{
            "$set": apiData,
        }

        _, err = collection.UpdateOne(context.Background(), filter, update)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        response := map[string]string{"message": fmt.Sprintf("Service API %s updated", serviceApiId)}
        json.NewEncoder(w).Encode(response)
    }
}

func DeleteServiceAPIHandler(collection *mongo.Collection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    serviceApiId := vars["serviceApiId"]

    filter := bson.M{"apiid": serviceApiId}

    result, err := collection.DeleteOne(context.Background(), filter)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    if result.DeletedCount == 0 {
        http.Error(w, "No documents found with the given ID", http.StatusNotFound)
        return
    }

    response := map[string]string{"message": fmt.Sprintf("Service API %s deleted", serviceApiId)}
    json.NewEncoder(w).Encode(response)
    }
}

// Handler to update a specific service API
func PatchServiceAPIHandler(collection *mongo.Collection) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        serviceApiId := vars["serviceApiId"]

        var patchReq Apis.PatchRequest
        if err := json.NewDecoder(r.Body).Decode(&patchReq); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        filter := bson.M{"apiid": serviceApiId}

        update := bson.M{"$set": patchReq}

        result, err := collection.UpdateOne(context.Background(), filter, update)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        if result.MatchedCount == 0 {
            http.Error(w, "No documents found with the given ID", http.StatusNotFound)
            return
        }

        response := map[string]string{"message": fmt.Sprintf("Service API %s updated", serviceApiId)}
        json.NewEncoder(w).Encode(response)
    }
}

/*func GetServiceApisHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    apfId := vars["apfId"]
    _ = apfId
    filteredRequests := serviceRequests 

    respondWithJSON(w, http.StatusOK, filteredRequests)
}*/

/*func UpdateServiceApiHandler(w http.ResponseWriter, r *http.Request) {
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
}*/

/*func DeleteServiceApiHandler(w http.ResponseWriter, r *http.Request) {
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
}*/

/*func respondWithError(w http.ResponseWriter, code int, message string) {
    details := Apis.ProblemDetails{
        Type:   "error",
        Title:  http.StatusText(code),
        Status: code,
        Detail: message,
    }
    respondWithJSON(w, code, details)
}*/

/*func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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
}*/


/*func loadServiceMockResponse() (Apis.ServiceInfoMock, error) {
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
}*/