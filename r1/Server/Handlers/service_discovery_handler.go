package Handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var serviceAPIsCollection *mongo.Collection

func InitServiceAPIsCollection(collection *mongo.Collection) {
	serviceAPIsCollection = collection
}

// AllServiceAPIsHandler handles GET requests to fetch all service APIs
// AllServiceAPIsHandler handles GET requests to fetch all service APIs with filtering options
// ServiceDiscoveryHandler handles GET requests to discover service APIs with filtering options
func ServiceDiscoveryHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		filter := bson.M{}

		// Filter by apiInvokerId
		if apiInvokerID := params.Get("api-invoker-id"); apiInvokerID != "" {
			filter["apiid"] = apiInvokerID
		}

		// Filter by apiName
		if apiName := params.Get("api-name"); apiName != "" {
			filter["apiname"] = apiName
		}

		// Filter by apiVersion (adjust field name if needed)
		if apiVersion := params.Get("api-version"); apiVersion != "" {
			filter["aefprofiles.versions.apiversion"] = apiVersion
		}

		// Filter by commType
		if commType := params.Get("comm-type"); commType != "" {
			filter["commtype"] = commType
		}

		// Filter by protocol
		if protocol := params.Get("protocol"); protocol != "" {
			filter["protocol"] = protocol
		}

		// Filter by aefId
		if aefID := params.Get("aef-id"); aefID != "" {
			// Assuming aefID is directly under apistatus.aefids array
			filter["apistatus.aefids"] = aefID
		}

		// Filter by dataFormat (adjust field name if needed)
		if dataFormat := params.Get("data-format"); dataFormat != "" {
			filter["description"] = dataFormat // Example mapping; adjust as per your schema
		}

		// Filter by apiCat (adjust field name if needed)
		if apiCat := params.Get("api-cat"); apiCat != "" {
			filter["serviceapicategory"] = apiCat
		}

		// Filter by reqApiProvName (adjust field name if needed)
		if reqApiProvName := params.Get("req-api-prov-name"); reqApiProvName != "" {
			filter["apiprovname"] = reqApiProvName
		}

		// Filter by supportedFeatures (adjust field name if needed)
		if supportedFeatures := params.Get("supported-features"); supportedFeatures != "" {
			filter["supportedfeatures"] = supportedFeatures
		}

		// Filter by apiSupportedFeatures (adjust field name if needed)
		if apiSupportedFeatures := params.Get("api-supported-features"); apiSupportedFeatures != "" {
			filter["apisuppfeats"] = apiSupportedFeatures
		}

		// Filter by ueIpAddr (adjust field name if needed)
		if ueIPAddr := params.Get("ue-ip-addr"); ueIPAddr != "" {
			// Assuming ueIPAddr is an object in your schema; adjust accordingly
			filter["ueipaddr.ipv4Addr"] = ueIPAddr // Example mapping; adjust as per your schema
		}

		// Filter by serviceKPIs (adjust field names and structure if needed)
		if maxReqRate := params.Get("service-kpis.maxReqRate"); maxReqRate != "" {
			filter["servicekpi.maxReqRate"] = maxReqRate
		}
		if maxResTime := params.Get("service-kpis.maxRestime"); maxResTime != "" {
			filter["servicekpi.maxRestime"] = maxResTime
		}
		if availability := params.Get("service-kpis.availability"); availability != "" {
			filter["servicekpi.availability"] = availability
		}
		if avalComp := params.Get("service-kpis.avalComp"); avalComp != "" {
			filter["servicekpi.avalComp"] = avalComp
		}
		if avalGraComp := params.Get("service-kpis.avalGraComp"); avalGraComp != "" {
			filter["servicekpi.avalGraComp"] = avalGraComp
		}
		if avalMem := params.Get("service-kpis.avalMem"); avalMem != "" {
			filter["servicekpi.avalMem"] = avalMem
		}
		if avalStor := params.Get("service-kpis.avalStor"); avalStor != "" {
			filter["servicekpi.avalStor"] = avalStor
		}
		if conBand := params.Get("service-kpis.conBand"); conBand != "" {
			filter["servicekpi.conBand"] = conBand
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cur, err := collection.Find(ctx, filter)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Error finding service APIs", err)
			return
		}
		defer cur.Close(ctx)

		var serviceAPIs []map[string]interface{}

		for cur.Next(ctx) {
			var result map[string]interface{}
			err := cur.Decode(&result)
			if err != nil {
				handleError(w, http.StatusInternalServerError, "Error decoding service API", err)
				return
			}
			serviceAPIs = append(serviceAPIs, result)
		}

		if err := cur.Err(); err != nil {
			handleError(w, http.StatusInternalServerError, "Error iterating through service APIs", err)
			return
		}

		jsonBytes, err := json.Marshal(serviceAPIs)
		if err != nil {
			handleError(w, http.StatusInternalServerError, "Error encoding service APIs to JSON", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)
	}
}



func handleError(w http.ResponseWriter, statusCode int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorMessage := fmt.Sprintf(`{"error": "%s", "message": "%s"}`, message, err.Error())
	w.Write([]byte(errorMessage))
}

