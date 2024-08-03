package Handlers

import (
	"context"
	"net/http"
	"r1/r1/Apis"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetAllDataTypesHandler retrieves all DmeTypeInformation objects with optional filtering.
func GetAllDataTypesHandler(dataTypeCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		namespace := query.Get("identity-namespace")
		name := query.Get("identity-name")
		category := query["data-category"]

		filter := bson.M{}
		if namespace != "" {
			filter["dataTypeId.namespace"] = namespace
		}
		if name != "" {
			filter["dataTypeId.name"] = name
		}
		if len(category) > 0 {
			filter["dataTypeId.dataCategory"] = bson.M{"$all": category}
		}

		cursor, err := dataTypeCollection.Find(context.TODO(), filter)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error querying database")
			return
		}
		defer cursor.Close(context.TODO())

		var results []Apis.DataTypeProdCapRegistration
		if err = cursor.All(context.TODO(), &results); err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error processing results")
			return
		}

		respondWithJSON(w, http.StatusOK, results)
	}
}

// GetDataTypeByIdHandler retrieves a specific DmeTypeInformation object by its dataTypeId.
func GetDataTypeByIdHandler(dataTypeProdCapsCollection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		dataTypeId := vars["dataTypeId"]

		// Create filter to find the document by nested datatypeinformation.datatypeid.typeid.
		filter := bson.M{"datatypeinformation.datatypeid.typeid": dataTypeId}
		var result Apis.DataTypeProdCapRegistration
		err := dataTypeProdCapsCollection.FindOne(context.TODO(), filter).Decode(&result)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				respondWithError(w, http.StatusNotFound, "DataType not found")
			} else {
				respondWithError(w, http.StatusInternalServerError, "Error querying database")
			}
			return
		}

		respondWithJSON(w, http.StatusOK, result)
	}
}
