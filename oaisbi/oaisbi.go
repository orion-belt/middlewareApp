package oaisbi

import (
	"middlewareApp/common"
	"middlewareApp/logger"
	"middlewareApp/config"
    "strconv"
	"net/http"
	"encoding/json"
	"bytes"
	// "io/ioutil"
)

var SnssaiLocal common.SNSSAI

func UpdateSnssai (Snssai common.SNSSAI) {
	var amf_cfg OaiAamfConfig
	GetAmfConfig(&amf_cfg)
	AddSlice(Snssai, &amf_cfg)
	logger.OaiSbiLog.Infoln("Snssai updated (", Snssai.Sst,",", Snssai.Sd,")")
}

func DeleteSnssai (Snssai common.SNSSAI) {
	var amf_cfg OaiAamfConfig
	GetAmfConfig(&amf_cfg)
	SnssaiLocal = Snssai

	DeleteSnssaiFromList(&amf_cfg, Snssai)

	reqbody, _ := json.Marshal(amf_cfg)

	AmfConfigUpdate(reqbody)
	logger.OaiSbiLog.Infoln("Snssai deleted (", Snssai.Sst,",", Snssai.Sd,")")
}

func InitSnssai (Snssai common.SNSSAI){
	var amf_cfg OaiAamfConfig
	GetAmfConfig(&amf_cfg)
	SnssaiLocal = Snssai
	AddSlice(Snssai, &amf_cfg)
	logger.OaiSbiLog.Infoln("Snssai initialised (", Snssai.Sst,",", Snssai.Sd,")")
}

func GetAmfConfig (amf_cfg *OaiAamfConfig) {
	base_url := config.GetOaiService("oai-amf")
	url := base_url+"/namf-oai/v1/configuration"
	client := &http.Client{}
	status, data, _ := common.SendHttpRequest("GET", url, "",client)
	if status != 200 {
		logger.MagmaGwRegLog.Errorln("HTTP request failed with code :", status)
		logger.MagmaGwRegLog.Errorln("HTTP response body :", data)
	} else {
		json.Unmarshal([]byte(data), &amf_cfg)
	}
}

func AddSlice (Snssai common.SNSSAI, amf_cfg *OaiAamfConfig){
logger.OaiSbiLog.Debugln("Updating AMF slice")
Sd, _:= strconv.Atoi(Snssai.Sd)

for _, snssai := range amf_cfg.PlmnList[0].SliceList {
	if (snssai.Sst == Snssai.Sst && snssai.Sd == Sd) {
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
	url := base_url+"/namf-oai/v1/configuration"
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

func DeleteSnssaiFromList (amf_cfg *OaiAamfConfig, Snssai common.SNSSAI){
	index := 0
	del_index := 255
	Sd, _:= strconv.Atoi(Snssai.Sd)
	for _, snssai := range amf_cfg.PlmnList[0].SliceList {
		if (snssai.Sst == Snssai.Sst && snssai.Sd == Sd) {
			del_index = index
		}
		index++
	}

	if del_index != 255{
		amf_cfg.PlmnList[0].SliceList[del_index] = amf_cfg.PlmnList[0].SliceList[len(amf_cfg.PlmnList[0].SliceList)-1] // Copy last element to index i.
		amf_cfg.PlmnList[0].SliceList = amf_cfg.PlmnList[0].SliceList[:len(amf_cfg.PlmnList[0].SliceList)-1] // Erase last element
	}
}