package oaisbi

import (
	//"time"
	//"bytes"
	"encoding/json"
	"middlewareApp/logger"
	"middlewareApp/magmanbi"
	"strconv"
	//"strings"
)

type rawMconfigMsg struct {
	ConfigsByKey map[string]json.RawMessage
}

const (
	BASE_URL = "http://192.168.70.132/namf-oai/v1/"
    START_BRACE = "{"
	PLMN_LISTS_STRING = "\"plmn_list\":" 
	END_BRACE =  "}"	
)

type Guami struct {
	AmfPointer string `json:"AmfPointer"`
	AmfSetID   string `json:"AmfSetID"`
	Mcc        string `json:"mcc"`
	Mnc        string `json:"mnc"`
	RegionID   string `json:"regionID"`
}

type Guami_List struct {
	AmfPointer string `json:"AmfPointer"`
	AmfSetID   string `json:"AmfSetID"`
	Mcc        string `json:"mcc"`
	Mnc        string `json:"mnc"`
	RegionID   string `json:"regionID"`
}

type N2 struct {
	Addr4    string `json:"addr4"`
	Addr6    string `json:"addr6"`
	If_name  string `json:"if_name"`
	Mtu      int    `json:"mtu"`
	Network4 string `json:"network4"`
	Port     int    `json:"port"`
}

type N11 struct {
	Addr4          string `json:"addr4"`
	Addr6          string `json:"addr6"`
	If_name        string `json:"if_name"`
	Mtu            int    `json:"mtu"`
	Network4       string `json:"network4"`
	Port           int    `json:"port"`
	Sbi_http2_port int    `json:"sbi_http2_port"`
}

type Slice_List struct {
	Sd  int `json:"sd"` //string
	Sst int `json:"sst"`
}

type Plmn_List struct {
	Mcc        string       `json:"mcc"`
	Mnc        string       `json:"mnc"`
	Slice_List []Slice_List `json:"slice_list"`
	Tac        uint32       `json:"tac"`
}

type Support_Features struct {
	Enable_External_Ausf   bool `json:"enable_external_ausf"`
	Enable_external_nrf    bool `json:"enable_external_nrf"`
	Enable_external_nssf   bool `json:"enable_external_nssf"`
	Enable_external_udm    bool `json:"enable_external_udm"`
	Enable_nf_registration bool `json:"enable_nf_registration"`
	Enable_nrf_selection   bool `json:"enable_nrf_selection"`
	Enable_smf_selection   bool `json:"enable_smf_selection"`
	Use_fqdn_dns           bool `json:"use_fqdn_dns"`
	Use_http2              bool `json:"use_http2"`
}

type Function struct {
	Api_version string `json:"api_version"`
	Fqdn        string `json:"fqdn"`
	Ipv4_addr   string `json:"ipv4_addr"`
	Port        uint   `json:"port"`
}

type Auth_Para struct {
	Mysql_db     string `json:"mysql_db"`
	Mysql_pass   string `json:"mysql_pass"`
	Mysql_server string `json:"mysql_server"`
	Mysql_user   string `json:"mysql_user"`
	Random       string `json:"random"`
}

type Smf_Pool struct {
	Fqdn       string `json:"fqdn"`
	Http2_port uint32 `json:"http2_port"`
	Id         int    `json:"id"`
	Ipv4       string `json:"ipv4"`
	Port       string `json:"port"` //int
	Selected   bool   `json:"selected"`
	Version    string `json:"version"` //int
}

type Amf_Configuration struct {
	Amf_name              string           `json:"amf_name"`
	Ausf                  Function         `json:"ausf"`
	Auth_Para             Auth_Para        `json:"auth_para"`
	Guami                 Guami            `json:"guami"`
	Guami_List            []Guami_List     `json:"guami_list"`
	Instance              uint             `json:"instance"`
	Is_emergency_support  string           `json:"is_emergency_support"`
	N11                   N11              `json:"n11"`
	N2                    N2               `json:"n2"`
	Nrf                   Function         `json:"nrf"`
	Pid_dir               string           `json:"pid_dir"`
	Plmn_List             []Plmn_List      `json:"plmn_list"`
	Relative_amf_capacity int              `json:"relative_amf_capacity"`
	Smf_Pool              []Smf_Pool       `json:"smf_pool"`
	Support_Features      Support_Features `json:"support_features"`
	Udm                   Function         `json:"udm"`
}


func GetPlmn() Amf_Configuration {
	var Amf Amf_Configuration
	url := BASE_URL + "configuration"
	status, data, _ := magmanbi.SendHttpRequest("GET", url, "")
	if status != 200 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed with code:", status)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		panic(status)
	} else {
		err := json.Unmarshal([]byte(data), &Amf)
		if err != nil {
			panic(err)
		}
		logger.MagmaGwRegLog.Infoln("Unmarshal AMF : ", Amf, "\n")
		return Amf
	}
}

func UpdateAmfPlmnForAllElements() {
	if magmanbi.Slices.Updatedsstsd {
		Amf := GetPlmn()
		logger.MagmaGwRegLog.Infoln("Updating All The Slice_Lists")

		for i := 0; i < len(Amf.Plmn_List); i++ {
			for j := 0; j < len(Amf.Plmn_List[i].Slice_List); j++ {
				Amf.Plmn_List[i].Slice_List[j].Sd, _ = strconv.Atoi(magmanbi.Slices.Sd)
				Amf.Plmn_List[i].Slice_List[j].Sst = magmanbi.Slices.Sst
			}
		}
		UpdateSlice(Amf)
	} else {
		logger.MagmaGwRegLog.Infoln("Nothing to Update on the AMF")
	}
}

func UpdateAmfPlmnForSpecificElement(plmnPosition int, slice_listPosition int) {
	if magmanbi.Slices.Updatedsstsd {
		Amf := GetPlmn()
		logger.MagmaGwRegLog.Infoln("Updating The Slice_list: (", plmnPosition, ",", slice_listPosition, ")\n")

		if plmnPosition >= len(Amf.Plmn_List) {
			panic("Out of range")
		} else {
			if slice_listPosition >= len(Amf.Plmn_List[plmnPosition].Slice_List) {
				panic("Out of range")
			} else {
				Amf.Plmn_List[plmnPosition].Slice_List[slice_listPosition].Sd, _ = strconv.Atoi(magmanbi.Slices.Sd)
				Amf.Plmn_List[plmnPosition].Slice_List[slice_listPosition].Sst = magmanbi.Slices.Sst
			}
		}
		UpdateSlice(Amf)
	} else {
		logger.MagmaGwRegLog.Infoln("Nothing to Update on the AMF")
	}
}

func UpdateSlice(Amf Amf_Configuration) {
	url := BASE_URL + "configuration"
	sliceData, _ := json.Marshal(Amf.Plmn_List)
	
	newplmn := START_BRACE + PLMN_LISTS_STRING + string(sliceData) + END_BRACE


	status, data, err := magmanbi.SendHttpRequest("PUT", url, newplmn)
	if status != 200 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed, Details- :", status, err)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
	} else {
		logger.MagmaGwRegLog.Infoln("AMF has been updated:", data)
		//GetPlmn()
	}
}