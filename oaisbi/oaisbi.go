package oaisbi

import (
	"middlewareApp/common"
	"middlewareApp/logger"
	"middlewareApp/config"
)

func UpdateSnssai (Snssai common.SNSSAI) {
	logger.OaiSbiLog.Infoln("Snssai updated (", Snssai.Sst,",", Snssai.Sd,")")

}

func InitSnssai (Snssai common.SNSSAI){
	base_url := config.GetOaiService("oai-amf")
	println("Here  ",string(base_url))
	logger.OaiSbiLog.Infoln("Snssai initialised (", Snssai.Sst,",", Snssai.Sd,")")
}
