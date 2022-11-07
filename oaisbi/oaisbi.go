package oaisbi

import (
	"bytes"
	"encoding/json"
	"middlewareApp/common"
	"middlewareApp/config"
	"middlewareApp/logger"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
	"net/http"
	"strconv"
	// "strings"
	// "io/ioutil"
)

//#######################################################################
//#### Config Update ####################################################
//#######################################################################
var SnssaiLocal common.SNSSAI

func UpdateSnssai(Snssai common.SNSSAI) {
	var amf_cfg OaiAamfConfig
	GetAmfConfig(&amf_cfg)
	AddSlice(Snssai, &amf_cfg)
	logger.OaiSbiLog.Infoln("Snssai updated (", Snssai.Sst, ",", Snssai.Sd, ")")
}

func DeleteSnssai(Snssai common.SNSSAI) {
	var amf_cfg OaiAamfConfig
	GetAmfConfig(&amf_cfg)
	SnssaiLocal = Snssai

	DeleteSnssaiFromList(&amf_cfg, Snssai)

	reqbody, _ := json.Marshal(amf_cfg)

	AmfConfigUpdate(reqbody)
	logger.OaiSbiLog.Infoln("Snssai deleted (", Snssai.Sst, ",", Snssai.Sd, ")")
}

func InitSnssai(Snssai common.SNSSAI) {
	var amf_cfg OaiAamfConfig
	GetAmfConfig(&amf_cfg)
	SnssaiLocal = Snssai
	AddSlice(Snssai, &amf_cfg)
	logger.OaiSbiLog.Infoln("Snssai initialised (", Snssai.Sst, ",", Snssai.Sd, ")")
}

func GetAmfConfig(amf_cfg *OaiAamfConfig) {
	base_url := config.GetOaiService("oai-amf")
	url := base_url + "/namf-oai/v1/configuration"
	client := &http.Client{}
	status, data, _ := common.SendHttpRequest("GET", url, "", client)
	if status != 200 {
		logger.MagmaGwRegLog.Errorln("HTTP request failed with code :", status)
		logger.MagmaGwRegLog.Errorln("HTTP response body :", data)
	} else {
		json.Unmarshal([]byte(data), &amf_cfg)
	}
}

func AddSlice(Snssai common.SNSSAI, amf_cfg *OaiAamfConfig) {
	logger.OaiSbiLog.Debugln("Updating AMF slice")
	Sd, _ := strconv.Atoi(Snssai.Sd)

	for _, snssai := range amf_cfg.PlmnList[0].SliceList {
		if snssai.Sst == Snssai.Sst && snssai.Sd == Sd {
			logger.OaiSbiLog.Warnln("snssai already configured at AMF")
			return
		}
	}
	DeleteSnssaiFromList(amf_cfg, SnssaiLocal)

	var Slice Slice
	Slice.Sst = Snssai.Sst
	Slice.Sd = Sd

	amf_cfg.PlmnList[0].SliceList = append(amf_cfg.PlmnList[0].SliceList, Slice)
	reqbody, _ := json.Marshal(amf_cfg)

	SnssaiLocal.Sst = Snssai.Sst
	SnssaiLocal.Sd = Snssai.Sd

	AmfConfigUpdate(reqbody)

	// Add siice to SMF
	// Add slice to UPF
}

func AmfConfigUpdate(reqbody []byte) {
	//
	base_url := config.GetOaiService("oai-amf")
	url := base_url + "/namf-oai/v1/configuration"
	client := &http.Client{}

	// status, data, _ := common.SendHttpRequest("PUT", url,string(reqbody),client,[]string{"Content-Type", "application/json"})
	// if status != 200 {
	// 	logger.OaiSbiLog.Errorln("HTTP request failed with code :", status)
	// 	logger.OaiSbiLog.Errorln("HTTP response body :", data)
	// } else {
	// 	logger.OaiSbiLog.Infoln("AMF config updated successfully: ", data)
	// }
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqbody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	logger.OaiSbiLog.Infoln("AMF config updated successfully: ", resp.Body)
}

func DeleteSnssaiFromList(amf_cfg *OaiAamfConfig, Snssai common.SNSSAI) {
	index := 0
	del_index := 255
	Sd, _ := strconv.Atoi(Snssai.Sd)
	for _, snssai := range amf_cfg.PlmnList[0].SliceList {
		if snssai.Sst == Snssai.Sst && snssai.Sd == Sd {
			del_index = index
		}
		index++
	}

	if del_index != 255 {
		amf_cfg.PlmnList[0].SliceList[del_index] = amf_cfg.PlmnList[0].SliceList[len(amf_cfg.PlmnList[0].SliceList)-1] // Copy last element to index i.
		amf_cfg.PlmnList[0].SliceList = amf_cfg.PlmnList[0].SliceList[:len(amf_cfg.PlmnList[0].SliceList)-1]           // Erase last element
	}
}

//#######################################################################
//#### Subscriber Update ################################################
//#######################################################################
func UpdateSubscriber(imsi string, sub Subscriber) {
	// imsi = strings.ReplaceAll(imsi, "IMSI", "")
	base_url := config.GetOaiService("oai-udr")
	url := base_url + "/nudr-dr/v1/subscription-data/" + imsi + "/authentication-data/authentication-subscription"
	client := &http.Client{}

	// sub := CreateSubscriberProfile(imsi)
	reqbody, _ := json.Marshal(sub)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqbody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	logger.OaiSbiLog.Infoln("Subscriber ", imsi, " updated successfully: ", resp.Body)
}

// func UpdateSubscriber(imsi string) {
// 	imsi = strings.ReplaceAll(imsi, "IMSI", "")
// 	base_url := config.GetOaiService("oai-udr")
// 	url := base_url+"/nudr-dr/v1/subscription-data/"+imsi+"/authentication-data/authentication-subscription"
// 	client := &http.Client{}

// 	sub := CreateSubscriberProfile(imsi)
// 	reqbody, _ := json.Marshal(sub)
// 	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(reqbody))
// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
//     logger.OaiSbiLog.Infoln("Subscriber updated successfully: ", resp.Body)
// }

func CreateSubscriberProfile(imsi string) Subscriber {
	logger.OaiSbiLog.Infoln("Generating subscriber prifile for imsi -", imsi)
	var Subscriber Subscriber
	Subscriber.AuthenticationMethod = "5G_AKA"
	Subscriber.EncPermanentKey = "0C0A34601D4F07677303652C0462535B"
	Subscriber.ProtectionParameterID = "0C0A34601D4F07677303652C0462535B"
	Subscriber.SequenceNumber.Sqn = "000000000020"
	Subscriber.SequenceNumber.SqnScheme = "NON_TIME_BASED"
	Subscriber.SequenceNumber.LastIndexes.Ausf = 0
	Subscriber.AuthenticationManagementField = "8000"
	Subscriber.AlgorithmID = "milenage"
	Subscriber.EncOpcKey = "63bfa50ee6523365ff14c1f45f88737d"
	Subscriber.EncTopcKey = "63bfa50ee6523365ff14c1f45f88737d"
	Subscriber.VectorGenerationInHss = false
	Subscriber.N5GcAuthMethod = ""
	Subscriber.RgAuthenticationInd = false
	Subscriber.Supi = imsi
	return Subscriber
}

func HandleSubscriberUpdate(subscriberlist []Subscriber, update *protos.DataUpdateBatch) (Subscriber, bool) {
	var profile Subscriber
	if len(update.Updates) > len(subscriberlist) {
		logger.OaiSbiLog.Infoln("Handle add subscriber")
		for _, subi := range update.Updates {
			var imsi_found bool
			imsi_found = false
			new_imsi := subi.GetKey()
			for _, subj := range subscriberlist {
				if subi.GetKey() == subj.Supi {
					imsi_found = true
					break
				}
			}
			if !imsi_found {
				profile = CreateSubscriberProfile(new_imsi)
				UpdateSubscriber(subi.GetKey(), profile)
				logger.OaiSbiLog.Infoln("Updated new subscriber ->", new_imsi)
				return profile, true
			}
		}
	} else {
		logger.OaiSbiLog.Infoln("Handle remove subscriber")
		for _, subi := range subscriberlist {
			var imsi_found bool
			imsi_found = false
			new_imsi := subi.Supi
			for _, subj := range update.Updates {
				if subi.Supi == subj.GetKey() {
					imsi_found = true
					break
				}
			}
			if !imsi_found {
				RemoveSubscriberFromUdr(new_imsi)
				logger.OaiSbiLog.Infoln("Remove existing subscriber -", new_imsi)
				profile = CreateSubscriberProfile(new_imsi)
				return profile, false
			}
		}
	}
	return profile, true
}

func RemoveSubscriberFromUdr(imsi string) {
	base_url := config.GetOaiService("oai-udr")
	url := base_url + "/nudr-dr/v1/subscription-data/" + imsi + "/authentication-data/authentication-subscription"
	client := &http.Client{}

	// sub := CreateSubscriberProfile(imsi)
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	logger.OaiSbiLog.Infoln("Subscriber ", imsi, " removed successfully: ", resp.Body)
}

func RemoveSubscriberFromList(subscriberlist *[]Subscriber, subscriber Subscriber) {
	index := 0
	del_index := 255
	var subscriberlisttmp []Subscriber
	subscriberlisttmp = *subscriberlist
	for _, sub := range subscriberlisttmp {
		if sub.Supi == subscriber.Supi {
			del_index = index
		}
		index++
	}
	if del_index != 255 {
		subscriberlisttmp[del_index] = subscriberlisttmp[len(subscriberlisttmp)-1] // Copy last element to index i.
		subscriberlisttmp = subscriberlisttmp[:len(subscriberlisttmp)-1]           // Erase last element
	}
	println("New length", len(subscriberlisttmp))
}
