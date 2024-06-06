package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "net/http"
)

// Struct to hold the service name
type ServiceInfo struct {
    Name string `json:"name"`
}

// Struct to parse the YAML file
type YamlInfo struct {
    Info struct {
        Title string `yaml:"title"`
    } `yaml:"info"`
}

// Global variable to hold the service name
var serviceName ServiceInfo

// Function to load the service name from the YAML file
func loadServiceName() error {
    yamlFile, err := ioutil.ReadFile("TS29222_CAPIF_Publish_Service_API.yaml")
    if err != nil {
        return err
    }

    var yamlInfo YamlInfo
    err = yaml.Unmarshal(yamlFile, &yamlInfo)
    if err != nil {
        return err
    }

    serviceName = ServiceInfo{Name: yamlInfo.Info.Title}
    return nil
}

// Handler for GET request
func getServiceName(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(serviceName)
}

// Handler for POST request
func postServiceName(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(serviceName)
}

func main() {
    // Load service name from YAML file
    err := loadServiceName()
    if err != nil {
        fmt.Printf("Error loading service name: %v\n", err)
        return
    }

    http.HandleFunc("/get_service_name", getServiceName)
    http.HandleFunc("/post_service_name", postServiceName)

    fmt.Println("Server is running on port 8080...")
    http.ListenAndServe(":8080", nil)
}
