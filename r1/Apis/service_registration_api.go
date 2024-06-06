package Apis

import "time"

type ApiStatus struct {
    AefIds []string `json:"aefIds"`
}

type Resource struct {
    ResourceName   string   `json:"resourceName"`
    CommType       string   `json:"commType"`
    Uri            string   `json:"uri"`
    CustOpName     string   `json:"custOpName"`
    CustOperations []CustOp `json:"custOperations"`
    Operations     []string `json:"operations"`
    Description    string   `json:"description"`
}

type CustOp struct {
    CommType    string   `json:"commType"`
    CustOpName  string   `json:"custOpName"`
    Operations  []string `json:"operations"`
    Description string   `json:"description"`
}

type Version struct {
    ApiVersion     string     `json:"apiVersion"`
    Expiry         time.Time  `json:"expiry"`
    Resources      []Resource `json:"resources"`
    CustOperations []CustOp   `json:"custOperations"`
}

type InterfaceDescription struct {
    Ipv4Addr       string   `json:"ipv4Addr"`
    Ipv6Addr       string   `json:"ipv6Addr"`
    Fqdn           string   `json:"fqdn"`
    Port           int      `json:"port"`
    ApiPrefix      string   `json:"apiPrefix"`
    SecurityMethods []string `json:"securityMethods"`
}

type CivicAddr struct {
    Country   string `json:"country"`
    A1        string `json:"A1"`
}

type GeoArea struct {
    Shape string `json:"shape"`
    Point struct {
        Lon float64 `json:"lon"`
        Lat float64 `json:"lat"`
    } `json:"point"`
}

type AefLocation struct {
    CivicAddr CivicAddr `json:"civicAddr"`
    GeoArea   GeoArea   `json:"geoArea"`
    DcId      string    `json:"dcId"`
}

type ServiceKpis struct {
    MaxReqRate int    `json:"maxReqRate"`
    MaxRestime int    `json:"maxRestime"`
    Availability int  `json:"availability"`
    AvalComp  string `json:"avalComp"`
}

type IpRange struct {
    Start string `json:"start"`
    End   string `json:"end"`
}

type UeIpRange struct {
    UeIpv4AddrRanges []IpRange `json:"ueIpv4AddrRanges"`
    UeIpv6AddrRanges []IpRange `json:"ueIpv6AddrRanges"`
}

type AefProfile struct {
    AefId               string                 `json:"aefId"`
    Versions            []Version              `json:"versions"`
    Protocol            string                 `json:"protocol"`
    DataFormat          string                 `json:"dataFormat"`
    SecurityMethods     []string               `json:"securityMethods"`
    DomainName          string                 `json:"domainName"`
    InterfaceDescriptions []InterfaceDescription `json:"interfaceDescriptions"`
    AefLocation         AefLocation            `json:"aefLocation"`
    ServiceKpis         ServiceKpis            `json:"serviceKpis"`
    UeIpRange           UeIpRange              `json:"ueIpRange"`
}

type PublishServiceRequest struct {
    ApiName            string       `json:"apiName"`
    ApiId              string       `json:"apiId"`
    ApiStatus          ApiStatus    `json:"apiStatus"`
    AefProfiles        []AefProfile `json:"aefProfiles"`
    Description        string       `json:"description"`
    SupportedFeatures  string       `json:"supportedFeatures"`
    ShareableInfo      struct {
        IsShareable    bool     `json:"isShareable"`
        CapifProvDoms  []string `json:"capifProvDoms"`
    } `json:"shareableInfo"`
    ServiceAPICategory string `json:"serviceAPICategory"`
    ApiSuppFeats       string `json:"apiSuppFeats"`
    PubApiPath         struct {
        CcfIds []string `json:"ccfIds"`
    } `json:"pubApiPath"`
    CcfId              string `json:"ccfId"`
}

type ProblemDetails struct {
    Type            string `json:"type"`
    Title           string `json:"title"`
    Status          int    `json:"status"`
    Detail          string `json:"detail"`
    Instance        string `json:"instance"`
    Cause           string `json:"cause"`
    InvalidParams   []struct {
        Param  string `json:"param"`
        Reason string `json:"reason"`
    } `json:"invalidParams"`
    SupportedFeatures string `json:"supportedFeatures"`
}

type ServiceInfo struct {
    Name    string `json:"name"`
    Version string `json:"version"`
}

type YamlInfo struct {
    Info struct {
        Title   string `yaml:"title"`
        Version string `yaml:"version"`
    } `yaml:"info"`
}