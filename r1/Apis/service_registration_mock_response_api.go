package Apis

import "time"

// Define all nested structs to match the JSON structure
type ApiStatusMock struct {
    AefIds []string `json:"aefIds"`
}

type ResourceMock struct {
    ResourceName   string `json:"resourceName"`
    CommType       string `json:"commType"`
    Uri            string `json:"uri"`
    CustOpName     string `json:"custOpName"`
    CustOperations []struct {
        CommType   string   `json:"commType"`
        CustOpName string   `json:"custOpName"`
        Operations []string `json:"operations"`
        Description string  `json:"description"`
    } `json:"custOperations"`
    Operations  []string `json:"operations"`
    Description string   `json:"description"`
}

type VersionMock struct {
    ApiVersion    string     `json:"apiVersion"`
    Expiry        time.Time  `json:"expiry"`
    Resources     []Resource `json:"resources"`
    CustOperations []struct {
        CommType   string   `json:"commType"`
        CustOpName string   `json:"custOpName"`
        Operations []string `json:"operations"`
        Description string  `json:"description"`
    } `json:"custOperations"`
}

type InterfaceDescriptionMock struct {
    Ipv4Addr       string   `json:"ipv4Addr"`
    Ipv6Addr       string   `json:"ipv6Addr"`
    Fqdn           string   `json:"fqdn"`
    Port           int      `json:"port"`
    ApiPrefix      string   `json:"apiPrefix"`
    SecurityMethods []string `json:"securityMethods"`
}

type CivicAddrMock struct {
    Country     string `json:"country"`
    A1          string `json:"A1"`
    A2          string `json:"A2"`
    A3          string `json:"A3"`
    A4          string `json:"A4"`
    A5          string `json:"A5"`
    A6          string `json:"A6"`
    PRD         string `json:"PRD"`
    POD         string `json:"POD"`
    STS         string `json:"STS"`
    HNO         string `json:"HNO"`
    HNS         string `json:"HNS"`
    LMK         string `json:"LMK"`
    LOC         string `json:"LOC"`
    NAM         string `json:"NAM"`
    PC          string `json:"PC"`
    BLD         string `json:"BLD"`
    UNIT        string `json:"UNIT"`
    FLR         string `json:"FLR"`
    ROOM        string `json:"ROOM"`
    PLC         string `json:"PLC"`
    PCN         string `json:"PCN"`
    POBOX       string `json:"POBOX"`
    ADDCODE     string `json:"ADDCODE"`
    SEAT        string `json:"SEAT"`
    RD          string `json:"RD"`
    RDSEC       string `json:"RDSEC"`
    RDBR        string `json:"RDBR"`
    RDSUBBR     string `json:"RDSUBBR"`
    PRM         string `json:"PRM"`
    POM         string `json:"POM"`
    UsageRules  string `json:"usageRules"`
    Method      string `json:"method"`
    ProvidedBy  string `json:"providedBy"`
}

type PointMock struct {
    Lon float64 `json:"lon"`
    Lat float64 `json:"lat"`
}

type GeoAreaMock struct {
    Shape string `json:"shape"`
    Point PointMock  `json:"point"`
}

type AefLocationMock struct {
    CivicAddr CivicAddr `json:"civicAddr"`
    GeoArea   GeoArea   `json:"geoArea"`
    DcId      string    `json:"dcId"`
}

type ServiceKpisMock struct {
    MaxReqRate   int    `json:"maxReqRate"`
    MaxRestime   int    `json:"maxRestime"`
    Availability int    `json:"availability"`
    AvalComp     string `json:"avalComp"`
    AvalGraComp  string `json:"avalGraComp"`
    AvalMem      string `json:"avalMem"`
    AvalStor     string `json:"avalStor"`
    ConBand      int    `json:"conBand"`
}

type UeIpRangeMock struct {
    UeIpv4AddrRanges []struct {
        Start string `json:"start"`
        End   string `json:"end"`
    } `json:"ueIpv4AddrRanges"`
    UeIpv6AddrRanges []struct {
        Start string `json:"start"`
        End   string `json:"end"`
    } `json:"ueIpv6AddrRanges"`
}

type AefProfileMock struct {
    AefId                 string                `json:"aefId"`
    Versions              []Version             `json:"versions"`
    Protocol              string                `json:"protocol"`
    DataFormat            string                `json:"dataFormat"`
    SecurityMethods       []string              `json:"securityMethods"`
    DomainName            string                `json:"domainName"`
    InterfaceDescriptions []InterfaceDescription `json:"interfaceDescriptions"`
    AefLocation           AEFLocation           `json:"aefLocation"`
    ServiceKpis           ServiceKPIs           `json:"serviceKpis"`
    UeIpRange             UEIPRange             `json:"ueIpRange"`
}

type ShareableInfoMock struct {
    IsShareable    bool     `json:"isShareable"`
    CapifProvDoms  []string `json:"capifProvDoms"`
}

type PubApiPathMock struct {
    CcfIds []string `json:"ccfIds"`
}

type ServiceInfoMock struct {
    ApiName          string          `json:"apiName"`
    ApiId            string          `json:"apiId"`
    ApiStatus        APIStatus       `json:"apiStatus"`
    AefProfiles      []AEFProfile    `json:"aefProfiles"`
    Description      string          `json:"description"`
    SupportedFeatures string         `json:"supportedFeatures"`
    ShareableInfo    ShareableInfoMock   `json:"shareableInfo"`
    ServiceApiCategory string        `json:"serviceAPICategory"`
    ApiSuppFeats     string          `json:"apiSuppFeats"`
    PubApiPath       PubApiPathMock      `json:"pubApiPath"`
    CcfId            string          `json:"ccfId"`
}
