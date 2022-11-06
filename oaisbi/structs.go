package oaisbi

type OaiAamfConfig struct {
	AmfName string `json:"amf_name"`
	Ausf    struct {
		APIVersion string `json:"api_version"`
		Fqdn       string `json:"fqdn"`
		Ipv4Addr   string `json:"ipv4_addr"`
		Port       int    `json:"port"`
	} `json:"ausf"`
	AuthPara struct {
		MysqlDb     string `json:"mysql_db"`
		MysqlPass   string `json:"mysql_pass"`
		MysqlServer string `json:"mysql_server"`
		MysqlUser   string `json:"mysql_user"`
		Random      string `json:"random"`
	} `json:"auth_para"`
	Guami struct {
		AmfPointer string `json:"AmfPointer"`
		AmfSetID   string `json:"AmfSetID"`
		Mcc        string `json:"mcc"`
		Mnc        string `json:"mnc"`
		RegionID   string `json:"regionID"`
	} `json:"guami"`
	GuamiList []struct {
		AmfPointer string `json:"AmfPointer"`
		AmfSetID   string `json:"AmfSetID"`
		Mcc        string `json:"mcc"`
		Mnc        string `json:"mnc"`
		RegionID   string `json:"regionID"`
	} `json:"guami_list"`
	Instance           int    `json:"instance"`
	IsEmergencySupport string `json:"is_emergency_support"`
	N11                struct {
		Addr4        string `json:"addr4"`
		Addr6        string `json:"addr6"`
		IfName       string `json:"if_name"`
		Mtu          int    `json:"mtu"`
		Network4     string `json:"network4"`
		Port         int    `json:"port"`
		SbiHTTP2Port int    `json:"sbi_http2_port"`
	} `json:"n11"`
	N2 struct {
		Addr4    string `json:"addr4"`
		Addr6    string `json:"addr6"`
		IfName   string `json:"if_name"`
		Mtu      int    `json:"mtu"`
		Network4 string `json:"network4"`
		Port     int    `json:"port"`
	} `json:"n2"`
	Nrf struct {
		APIVersion string `json:"api_version"`
		Fqdn       string `json:"fqdn"`
		Ipv4Addr   string `json:"ipv4_addr"`
		Port       int    `json:"port"`
	} `json:"nrf"`
	PidDir   string `json:"pid_dir"`
	PlmnList []struct {
		Mcc       string `json:"mcc"`
		Mnc       string `json:"mnc"`
		SliceList []struct {
			Sst int `json:"sst"`
			Sd  int `json:"sd,omitempty"`
		} `json:"slice_list"`
		Tac int `json:"tac"`
	} `json:"plmn_list"`
	RelativeAmfCapacity int `json:"relative_amf_capacity"`
	SmfPool             []struct {
		Fqdn      string `json:"fqdn"`
		HTTP2Port int    `json:"http2_port"`
		ID        int    `json:"id"`
		Ipv4      string `json:"ipv4"`
		Port      string `json:"port"`
		Selected  bool   `json:"selected"`
		Version   string `json:"version"`
	} `json:"smf_pool"`
	SupportFeatures struct {
		EnableExternalAusf   bool `json:"enable_external_ausf"`
		EnableExternalNrf    bool `json:"enable_external_nrf"`
		EnableExternalNssf   bool `json:"enable_external_nssf"`
		EnableExternalUdm    bool `json:"enable_external_udm"`
		EnableNfRegistration bool `json:"enable_nf_registration"`
		EnableNrfSelection   bool `json:"enable_nrf_selection"`
		EnableSmfSelection   bool `json:"enable_smf_selection"`
		UseFqdnDNS           bool `json:"use_fqdn_dns"`
		UseHTTP2             bool `json:"use_http2"`
	} `json:"support_features"`
}

type Slice struct {
	Sst int `json:"sst"`
	Sd  int `json:"sd,omitempty"`
}

type SubscriberList struct {
	Subscribers []Subscriber
}

type Subscriber struct {
	AuthenticationMethod  string `json:"authenticationMethod"`
	EncPermanentKey       string `json:"encPermanentKey"`
	ProtectionParameterID string `json:"protectionParameterId"`
	SequenceNumber        struct {
		Sqn         string `json:"sqn"`
		SqnScheme   string `json:"sqnScheme"`
		LastIndexes struct {
			Ausf int `json:"ausf"`
		} `json:"lastIndexes"`
	} `json:"sequenceNumber"`
	AuthenticationManagementField string `json:"authenticationManagementField"`
	AlgorithmID                   string `json:"algorithmId"`
	EncOpcKey                     string `json:"encOpcKey"`
	EncTopcKey                    string `json:"encTopcKey"`
	VectorGenerationInHss         bool   `json:"vectorGenerationInHss"`
	N5GcAuthMethod                string `json:"n5gcAuthMethod"`
	RgAuthenticationInd           bool   `json:"rgAuthenticationInd"`
	Supi                          string `json:"supi"`
}
