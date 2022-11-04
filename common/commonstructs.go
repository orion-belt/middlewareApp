package common

// Generic Structs
type SNSSAI struct {
	Sst int
	Sd  string
}

type PLMN struct {
	Mcc string
	Mnc string
}

const (
	SST_VALUE_MAX = 255
)
// Common Magma Struct
type MME struct {
	Type                          string        `json:"@type,omitempty"`
	AmfDefaultSliceDifferentiator string        `json:"amfDefaultSliceDifferentiator,omitempty"`
	AmfDefaultSliceServiceType    int           `json:"amfDefaultSliceServiceType,omitempty"`
	AmfName                       string        `json:"amfName,omitempty"`
	AmfPointer                    string        `json:"amfPointer,omitempty"`
	AmfRegionID                   string        `json:"amfRegionId,omitempty"`
	AmfSetID                      string        `json:"amfSetId,omitempty"`
	AttachedEnodebTacs            []interface{} `json:"attachedEnodebTacs,omitempty"`
	CloudSubscriberdbEnabled      bool          `json:"cloudSubscriberdbEnabled,omitempty"`
	CongestionControlEnabled      bool          `json:"congestionControlEnabled,omitempty"`
	CsfbMcc                       string        `json:"csfbMcc,omitempty"`
	CsfbMnc                       string        `json:"csfbMnc,omitempty"`
	DNSPrimary                    string        `json:"dnsPrimary,omitempty"`
	DNSSecondary                  string        `json:"dnsSecondary,omitempty"`
	Enable5GFeatures              bool          `json:"enable5gFeatures,omitempty"`
	EnableDNSCaching              bool          `json:"enableDnsCaching,omitempty"`
	HssRelayEnabled               bool          `json:"hssRelayEnabled,omitempty"`
	Ipv4PCscfAddress              string        `json:"ipv4PCscfAddress,omitempty"`
	Ipv6DNSAddress                string        `json:"ipv6DnsAddress,omitempty"`
	Ipv6PCscfAddress              string        `json:"ipv6PCscfAddress,omitempty"`
	Lac                           int           `json:"lac,omitempty"`
	LogLevel                      string        `json:"logLevel,omitempty"`
	Mcc                           string        `json:"mcc,omitempty"`
	MmeCode                       int           `json:"mmeCode,omitempty"`
	MmeGid                        int           `json:"mmeGid,omitempty"`
	MmeRelativeCapacity           int           `json:"mmeRelativeCapacity,omitempty"`
	Mnc                           string        `json:"mnc,omitempty"`
	NatEnabled                    bool          `json:"natEnabled,omitempty"`
	NonEpsServiceControl          int           `json:"nonEpsServiceControl,omitempty"`
	RelayEnabled                  bool          `json:"relayEnabled,omitempty"`
	Tac                           int           `json:"tac,omitempty"`
}
