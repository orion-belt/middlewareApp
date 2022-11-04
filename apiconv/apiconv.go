package apiconv

import (
	"middlewareApp/common"
	"middlewareApp/logger"
	"middlewareApp/oaisbi"
)

//#######################################################################
//#### Config Update ####################################################
//#######################################################################

var DefaultSnssaiInit bool
var Snssai common.SNSSAI

func Init() {
	DefaultSnssaiInit = false
}

func IsVAlidSnssai(snssai_tmp common.SNSSAI) bool {
	if snssai_tmp.Sst > 0 && snssai_tmp.Sst <= common.SST_VALUE_MAX {
		return true
	}
	if DefaultSnssaiInit {
		oaisbi.DeleteSnssai(Snssai)
		DefaultSnssaiInit = false
	}
	return false
}

func IsPlmnUpdated(int, string) bool {
	return true
}

func IsSnssaiUpdated(snssai_tmp common.SNSSAI) bool {
	if !DefaultSnssaiInit {
		Snssai.Sst = snssai_tmp.Sst
		Snssai.Sd = snssai_tmp.Sd
		oaisbi.InitSnssai(Snssai)
		DefaultSnssaiInit = true
		return false
	}
	if Snssai.Sst == snssai_tmp.Sst {
		if snssai_tmp.Sd != "" && Snssai.Sd != snssai_tmp.Sd {
			logger.ApiConv.Infoln("SNSSAI-SD updated", Snssai.Sst, snssai_tmp.Sd)
			return true
		}
		return false
	}
	logger.ApiConv.Infoln("SNSSAI-SST updated", Snssai.Sd, snssai_tmp.Sd)
	return true
}

func CheckForConfigUpdate(mme *common.MME) {
	var snssai_tmp common.SNSSAI
	snssai_tmp.Sst = mme.AmfDefaultSliceServiceType
	snssai_tmp.Sd = mme.AmfDefaultSliceDifferentiator

	if (IsVAlidSnssai(snssai_tmp) || DefaultSnssaiInit) && IsSnssaiUpdated(snssai_tmp) {
		logger.ApiConv.Infoln("Trigger Snssai update")
		Snssai.Sst = mme.AmfDefaultSliceServiceType
		Snssai.Sd = mme.AmfDefaultSliceDifferentiator
		oaisbi.UpdateSnssai(Snssai)
	}
}

//#######################################################################
//#### Subscriber Update ################################################
//#######################################################################

func CheckForSubscriberUpdate(imsi string){
	oaisbi.UpdateSubscriber(imsi)

}