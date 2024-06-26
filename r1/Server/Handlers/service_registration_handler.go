package Handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"github.com/gorilla/mux"
	"r1/r1/Apis"
	//"io/ioutil"
	//"os"
	//"path/filepath"
	//"log"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func PublishServiceHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var api Apis.PublishServiceAPI

		vars := mux.Vars(r)
		apfId := vars["apfId"]

		rappFilter := bson.D{{"apf_id", apfId}}
		rappCursor, err := rappCollection.Find(context.TODO(), rappFilter)
		if err != nil {
			http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
			return
		}
		defer rappCursor.Close(context.TODO())

		var rapp Apis.Rapp
		rappExists := false

		for rappCursor.Next(context.TODO()) {
			if err = rappCursor.Decode(&rapp); err != nil {
				http.Error(w, "Failed to decode document", http.StatusInternalServerError)
				return
			}
			rappExists = true
		}

		if !rappExists {
			http.Error(w, "Invoker ID not found", http.StatusNotFound)
			return
		}

		// Parse the request body into the ServiceAPI struct
		err = json.NewDecoder(r.Body).Decode(&api)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if api.APIName == "" {
			http.Error(w, "apiName field is required", http.StatusBadRequest)
			return
		}

		if api.APIID == "" {
			api.APIID = uuid.New().String()
		}

		// Define the update operation using $push
		rappUpdate := bson.D{
			{"$push", bson.D{
				{"authorized_services", api.APIID},
			}},
		}
		_, err = rappCollection.UpdateOne(context.TODO(), rappFilter, rappUpdate)
		if err != nil {
			http.Error(w, "Failed to insert document into MongoDB", http.StatusInternalServerError)
			return
		}

		_, err = serviceCollection.InsertOne(context.TODO(), api)
		if err != nil {
			http.Error(w, "Failed to insert document into MongoDB", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		// TODO: Change the location to the actual location of created service
		w.Header().Set("Location", "https://api.service-apis/"+api.APIID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Service API published successfully"})
	}
}

func GetServiceAPIsHandler(serviceCollection *mongo.Collection, rappCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var serviceAPIs []Apis.GetServiceAPI
		invokerId := r.URL.Query().Get("api-invoker-id")

		rappFilter := bson.D{{"apf_id", invokerId}}
		rappCursor, err := rappCollection.Find(context.TODO(), rappFilter)
		if err != nil {
			http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
			return
		}
		defer rappCursor.Close(context.TODO())

		var rapp Apis.Rapp
		rappExists := false

		for rappCursor.Next(context.TODO()) {
			if err = rappCursor.Decode(&rapp); err != nil {
				http.Error(w, "Failed to decode document", http.StatusInternalServerError)
				return
			}
			rappExists = true
		}

		if !rappExists {
			http.Error(w, "Invoker ID not found", http.StatusNotFound)
			return
		}

		if rapp.IsAuthorized == false {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(rapp)

		serviceFilter := bson.D{{"apiid", bson.D{{"$in", rapp.AuthorizedServices}}}}

		// check if invokerId is present in the rapps collection
		serviceCursor, err := serviceCollection.Find(context.TODO(), serviceFilter)
		if err != nil {
			http.Error(w, "Failed to retrieve documents from MongoDB", http.StatusInternalServerError)
			return
		}
		defer serviceCursor.Close(context.TODO())

		//Iterate through the cursor and decode each document into a ServiceAPI struct
		for serviceCursor.Next(context.TODO()) {
			var api Apis.GetServiceAPI
			if err = serviceCursor.Decode(&api); err != nil {
				http.Error(w, "Failed to decode document", http.StatusInternalServerError)
				return
			}
			serviceAPIs = append(serviceAPIs, api)
		}

		// Check for cursor errors
		if err := serviceCursor.Err(); err != nil {
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
