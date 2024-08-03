package Apis

type DmeTypeIdStruct struct {
	TypeId string `json:"typeId"`
}

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type DeliverySchema struct {
	SchemaType string `json:"schemaType"`
	SchemaURL  string `json:"schemaUrl"`
}

type DataDeliveryMechanism struct {
	MechanismType string `json:"mechanismType"`
	Details       string `json:"details"`
}

type DataTypeInformation struct {
	DataTypeId             DmeTypeIdStruct         `json:"dataTypeId"`
	Metadata               Metadata                `json:"metadata"`
	DataProductionSchema   interface{}             `json:"dataProductionSchema,omitempty"`
	DataDeliverySchemas    []DeliverySchema        `json:"dataDeliverySchemas"`
	DataDeliveryMechanisms []DataDeliveryMechanism `json:"dataDeliveryMechanisms"`
}

type InterfaceDefinition struct {
	Endpoint string `json:"endpoint"`
}

type DataTypeProdCapRegistration struct {
	DataTypeInformation      DataTypeInformation `json:"dataTypeInformation"`
	DataRequestEndpoint      InterfaceDefinition `json:"dataRequestEndpoint,omitempty"`
	DataSubscriptionEndpoint InterfaceDefinition `json:"dataSubscriptionEndpoint,omitempty"`
	DataProducerConstraints  interface{}         `json:"dataProducerConstraints,omitempty"`
}

type ProblemDetails struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail"`
}
