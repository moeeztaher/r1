// data_registration_handlers.go

package Handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"r1/r1/Apis"
	// Adjust to your actual import path
)

// PullDataHandler handles pulling data from the specified dataPullUri.
func PullDataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pullDataRequest Apis.PullDataRequest

		// Read request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Unmarshal JSON body into PullDataRequest
		if err := json.Unmarshal(body, &pullDataRequest); err != nil {
			log.Println("Error unmarshalling request body:", err)
			http.Error(w, "Invalid request format", http.StatusBadRequest)
			return
		}

		// Validate dataPullUri
		if pullDataRequest.DataPullUri == "" {
			http.Error(w, "Missing dataPullUri", http.StatusBadRequest)
			return
		}

		// Fetch data from the dataPullUri
		resp, err := fetchData(pullDataRequest.DataPullUri)
		if err != nil {
			log.Println("Error fetching data:", err)
			http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Handle HTTP response status and content
		if resp.StatusCode == http.StatusOK {
			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println("Error reading response body:", err)
				http.Error(w, "Error reading data", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.WriteHeader(http.StatusOK)
			w.Write(data)
		} else if resp.StatusCode == http.StatusAccepted {
			w.WriteHeader(http.StatusAccepted)
		} else {
			http.Error(w, "Failed to fetch data", resp.StatusCode)
		}
	}
}

// fetchData sends a GET request to the specified URI and returns the response.
func fetchData(uri string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	return client.Do(req)
}
