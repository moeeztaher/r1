// data_registration_api.go

package Apis

// PushDataRequest represents the structure of the push data request.
type PushDataRequest struct {
	DataPushUri string `json:"dataPushUri"`
	Data        []byte `json:"data"`
}
