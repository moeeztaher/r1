// data_registration_handlers.go

package Handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"r1/r1/Apis"
	// Adjust to your actual import path
)

// PushDataHandler handles pushing data to the specified dataPushUri.
func PushDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pushDataRequest Apis.PushDataRequest

		// Read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Unmarshal JSON body into PushDataRequest
		if err := json.Unmarshal(body, &pushDataRequest); err != nil {
			log.Println("Error unmarshalling request body:", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate dataPushUri
		if pushDataRequest.DataPushUri == "" {
			http.Error(w, "Missing dataPushUri", http.StatusBadRequest)
			return
		}

		// Forward data to the dataPushUri
		resp, err := forwardData(pushDataRequest.DataPushUri, body)
		if err != nil {
			log.Println("Error forwarding data:", err)
			http.Error(w, "Failed to push data", http.StatusInternalServerError)
			return
		}

		if resp.StatusCode != http.StatusNoContent {
			http.Error(w, "Failed to push data", resp.StatusCode)
			return
		}

		// Success
		w.WriteHeader(http.StatusNoContent)
	}
}

// forwardData sends the data payload to the specified URI.
func forwardData(uri string, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}
