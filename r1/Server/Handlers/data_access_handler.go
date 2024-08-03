package Handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"r1/r1/Apis"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// CreateDataJobHandler handles POST requests to create a new data job.
func CreateDataJobHandler(dataJobsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		consumerId := vars["consumerId"]

		var dataJobInfo Apis.DataJobInfo
		if err := json.NewDecoder(r.Body).Decode(&dataJobInfo); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Generate a new UUID and convert it to a string
		dataJobID := uuid.New().String()
		dataJobInfo.DataJobId = dataJobID

		fmt.Printf("Generated DataJobId: %s\n", dataJobID)
		fmt.Printf("DataJobInfo: %+v\n", dataJobInfo)

		// Insert the data job into the database
		_, err := dataJobsCollection.InsertOne(context.TODO(), dataJobInfo)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error inserting data job into database")
			return
		}

		location := fmt.Sprintf("/%s/dataJobs/%s", consumerId, dataJobID)
		w.Header().Set("Location", location)
		respondWithJSON(w, http.StatusCreated, dataJobInfo)
	}
}

// DeleteDataJobHandler handles DELETE requests to remove a data job.
func DeleteDataJobHandler(dataJobsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		consumerId := vars["consumerId"]
		dataJobId := vars["dataJobId"]

		// Create a filter to find the data job by both consumerId and dataJobId
		filter := bson.M{
			"consumerId": consumerId,
			"dataJobId":  dataJobId,
		}

		result, err := dataJobsCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error deleting data job from database")
			return
		}

		if result.DeletedCount == 0 {
			respondWithError(w, http.StatusNotFound, "Data job not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

// NotifyDataAvailabilityHandler handles POST requests to notify about data availability.
func NotifyDataAvailabilityHandler(dataJobsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dataAvailabilityNotification Apis.DataAvailabilityNotification

		err := json.NewDecoder(r.Body).Decode(&dataAvailabilityNotification)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
