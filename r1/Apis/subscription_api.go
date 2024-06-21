package Apis
import ("time"
"go.mongodb.org/mongo-driver/bson/primitive")

type Subscriber struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	SubscriberID   string             `json:"subscriberId,omitempty"`
	SubscriptionIds []primitive.ObjectID `json:"subscriptionIds,omitempty" bson:"subscriptionIds,omitempty"`
}

type Subscription struct {
    Events                 []string `json:"events"`
    EventFilters           []EventFilter `json:"eventFilters"`
    EventReq               EventReq `json:"eventReq"`
    NotificationDestination string `json:"notificationDestination"`
    RequestTestNotification bool   `json:"requestTestNotification"`
    WebsockNotifConfig     WebsockNotifConfig `json:"websockNotifConfig"`
    SupportedFeatures      string `json:"supportedFeatures"`
}

type EventFilter struct {
    ApiIds        []string `json:"apiIds"`
    ApiInvokerIds []string `json:"apiInvokerIds"`
    AefIds        []string `json:"aefIds"`
}

type EventReq struct {
    ImmRep           bool          `json:"immRep"`
    NotifMethod      string        `json:"notifMethod"`
    MaxReportNbr     int           `json:"maxReportNbr"`
    MonDur           time.Time     `json:"monDur"`
    RepPeriod        int           `json:"repPeriod"`
    SampRatio        int           `json:"sampRatio"`
    PartitionCriteria []string      `json:"partitionCriteria"`
    GrpRepTime       int           `json:"grpRepTime"`
    NotifFlag        string        `json:"notifFlag"`
    NotifFlagInstruct NotifFlagInstruct `json:"notifFlagInstruct"`
    MutingSetting    MutingSetting `json:"mutingSetting"`
}

type NotifFlagInstruct struct {
    BufferedNotifs string `json:"bufferedNotifs"`
    Subscription   string `json:"subscription"`
}

type MutingSetting struct {
    MaxNoOfNotif         int `json:"maxNoOfNotif"`
    DurationBufferedNotif int `json:"durationBufferedNotif"`
}

type WebsockNotifConfig struct {
    WebsocketUri        string `json:"websocketUri"`
    RequestWebsocketUri bool   `json:"requestWebsocketUri"`
}