package apiconv

import (
	"middlewareApp/common"
	"middlewareApp/logger"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
	"middlewareApp/oaisbi"
)

//#######################################################################
//#### Config Update ####################################################
//#######################################################################

var DefaultConfigInit bool
var DefaultSubscriberInit bool
var Snssai common.SNSSAI
var SubscribersList []oaisbi.Subscriber

func Init() {
	DefaultConfigInit = false
	DefaultSubscriberInit = false
}

func IsVAlidSnssai(snssai_tmp common.SNSSAI) bool {
	if snssai_tmp.Sst > 0 && snssai_tmp.Sst <= common.SST_VALUE_MAX {
		return true
	}
	if DefaultConfigInit {
		oaisbi.DeleteSnssai(Snssai)
		DefaultConfigInit = false
	}
	return false
}

func IsPlmnUpdated(int, string) bool {
	return true
}

func IsSnssaiUpdated(snssai_tmp common.SNSSAI) bool {
	if !DefaultConfigInit {
		Snssai.Sst = snssai_tmp.Sst
		Snssai.Sd = snssai_tmp.Sd
		oaisbi.InitSnssai(Snssai)
		DefaultConfigInit = true
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

	if (IsVAlidSnssai(snssai_tmp) || DefaultConfigInit) && IsSnssaiUpdated(snssai_tmp) {
		logger.ApiConv.Infoln("Trigger Snssai update")
		Snssai.Sst = mme.AmfDefaultSliceServiceType
		Snssai.Sd = mme.AmfDefaultSliceDifferentiator
		oaisbi.UpdateSnssai(Snssai)
	}
}

//#######################################################################
//#### Subscriber Update ################################################
//#######################################################################


func CheckForSubscriberUpdate(actualMarshaled *protos.DataUpdateBatch) {
	num_sub := len(actualMarshaled.Updates)
	if !DefaultSubscriberInit {
		for index := 0; index < num_sub; index++ {
			imsi := actualMarshaled.Updates[index].Key
			subscriber := oaisbi.CreateSubscriberProfile(imsi)
			SubscribersList = append(SubscribersList, subscriber)
			oaisbi.UpdateSubscriber(imsi, subscriber)
			DefaultSubscriberInit = true
		}
		return
	} else {
		if num_sub == len(SubscribersList) {
			logger.MagmaGwRegLog.Infoln("No update in subscriber list")
			return
		} else {
			logger.MagmaGwRegLog.Infoln("Handle subscriber update")
			subscriber, isAdd := oaisbi.HandleSubscriberUpdate(SubscribersList, actualMarshaled)
			if isAdd {
				SubscribersList = append(SubscribersList, subscriber)
			} else {
				index := 0
				del_index := 255

				for _, sub := range SubscribersList {
					if sub.Supi == subscriber.Supi{
						del_index = index
					}
					index++
				}
				if del_index != 255 {
					SubscribersList[del_index] = SubscribersList[len(SubscribersList)-1] // Copy last element to index i.
					SubscribersList = SubscribersList[:len(SubscribersList)-1]           // Erase last element
				}
			}
		}
	}
}
