package Apis

// StreamingConfigurationKafka represents the Kafka streaming configuration.
type StreamingConfigurationKafka struct {
	TopicName             string                  `json:"topicName"`
	KafkaBootstrapServers []ServerAddressWithPort `json:"kafkaBootstrapServers"`
}

// ServerAddressWithPort represents the server address with port.
type ServerAddressWithPort struct {
	Hostname    string `json:"hostname"`
	PortAddress int    `json:"portAddress"`
}

// DataAvailabilityNotification represents the data availability notification.
type DataAvailabilityNotification struct {
	DataJobId               string                   `json:"dataJobId"`
	PullDeliveryDetailsHttp *PullDeliveryDetailsHttp `json:"pullDeliveryDetailsHttp,omitempty"`
	PushDeliveryDetailsHttp *PushDeliveryDetailsHttp `json:"pushDeliveryDetailsHttp,omitempty"`
}

// PullDeliveryDetailsHttp represents the HTTP pull delivery details.
type PullDeliveryDetailsHttp struct {
	DataPullUri string `json:"dataPullUri"`
}

// PushDeliveryDetailsHttp represents the HTTP push delivery details.
type PushDeliveryDetailsHttp struct {
	DataPushUri string `json:"dataPushUri"`
}

// DataJobInfo represents the data job information.
type DataJobInfo struct {
	DataJobId                       string                       `json:"dataJobId"` // Added DataJobId field
	DataDelivery                    string                       `json:"dataDelivery"`
	DataTypeId                      DataTypeId                   `json:"dataTypeId"`
	ProductionJobDefinition         interface{}                  `json:"productionJobDefinition"`
	DataDeliveryMechanism           string                       `json:"dataDeliveryMechanism"`
	DataDeliverySchemaId            string                       `json:"dataDeliverySchemaId"`
	PullDeliveryDetailsHttp         *PullDeliveryDetailsHttp     `json:"pullDeliveryDetailsHttp,omitempty"`
	DataAvailabilityNotificationUri string                       `json:"dataAvailabilityNotificationUri,omitempty"`
	PushDeliveryDetailsHttp         *PushDeliveryDetailsHttp     `json:"pushDeliveryDetailsHttp,omitempty"`
	StreamingConfigurationKafka     *StreamingConfigurationKafka `json:"streamingConfigurationKafka,omitempty"`
}

// DataTypeId represents the data type identifier.
type DataTypeId struct {
	TypeId string `json:"typeid"`
}
