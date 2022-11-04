package apiconv

import (
	"middlewareApp/common"
	"middlewareApp/oaisbi"
	"middlewareApp/logger"
)

var DefaultSnssaiInit bool
var Snssai common.SNSSAI

func IsPlmnUpdated(int, string) bool {
	return true
}

func IsSnssaiUpdated(Sst int, Sd string) bool {
	if !DefaultSnssaiInit {
		Snssai.Sst = Sst
		Snssai.Sd = Sd
		oaisbi.InitSnssai(Snssai)
		DefaultSnssaiInit = true
		return false
	}
	if Snssai.Sst == Sst{
		if (Sd !="" && Snssai.Sd != Sd) {return true}
		return false
	}
	return true
}

func CheckForUpdate(mme *common.MME) {
	if mme.AmfDefaultSliceServiceType <= common.SST_VALUE_MAX && IsSnssaiUpdated(mme.AmfDefaultSliceServiceType, mme.AmfDefaultSliceDifferentiator) {
		logger.ApiConv.Infoln("Trigger Snssai update")
		Snssai.Sst = mme.AmfDefaultSliceServiceType
		Snssai.Sd = mme.AmfDefaultSliceDifferentiator
		oaisbi.UpdateSnssai(Snssai)
	}
}
