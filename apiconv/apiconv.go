package apiconv

import (
	"middlewareApp/common"
	"middlewareApp/logger"
)

func CheckForUpdate (mme *common.MME){
	logger.MagmaGwRegLog.Infoln("MCC, MNC - ",mme.Mcc, mme.Mnc)
}