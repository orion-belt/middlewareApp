package magmasbi

import (
	//"time"
	//"bytes"
	"encoding/json"
	"middlewareApp/logger"
	"middlewareApp/magmanbi"
	"strings"
)

type rawMconfigMsg struct {
	ConfigsByKey map[string]json.RawMessage
}

const (
	BASE_URL = "http://192.168.70.132/namf-oai/v1/"
	//ServiceName = "streamer"
)

func UpdateAmfPlmn() {
	logger.MagmaGwRegLog.Infoln("Updating PLMN values (SST and SD)")
	newdata := GetPlmn()
	var mcc, mnc, sst, sd, tac string
	if magmanbi.UPDATEDSSTSD {
		if (newdata != nil) {
			mccPosition, foundmcc := magmanbi.Find(newdata, "mcc")
			mncPosition, foundmnc := magmanbi.Find(newdata, "mnc")
			tacPosition, foundtac := magmanbi.Find(newdata, "tac")

			if foundmcc {
				quote := []string(strings.Split(([]string(strings.Split(newdata[mccPosition], ": "))[2]), "\""))
				if quote[0] != ""{
					mcc = quote[0]
				}else{
					mcc = quote[1]
				}

				logger.MagmaGwRegLog.Infoln("mcc:", mcc, "\n")
			} else {
				mcc = "208"
			}

			if foundmnc {
				mnc = []string(strings.Split(newdata[mncPosition], ": "))[1]
				quote := []string(strings.Split(([]string(strings.Split(newdata[mncPosition], ": "))[1]), "\""))
				if quote[0] != ""{
					mnc = quote[0]
				}else{
					mnc = quote[1]
				}
				logger.MagmaGwRegLog.Infoln("mnc:", mnc, "\n")
			} else {
				mnc = "95"
			}

			if foundtac {
				tac = []string(strings.Split(([]string(strings.Split(newdata[tacPosition], ": "))[1]), "\n"))[0]
				logger.MagmaGwRegLog.Infoln("tac:", tac, "\n")
			} else {
				tac = "40960"
			}
		} else {
			mcc = "208"
			mnc = "95"
			tac = "40960"
		}

		sd = magmanbi.SD
		sst = magmanbi.SST

		plmnConfigurationList := "{\"plmn_list\":[{ \"mcc\":\"" + mcc +
			"\",\"mnc\":\"" + mnc +
			"\",\"slice_list\":[{\"sst\":" + sst +
			", \"sd\":" + sd +
			"}],\"tac\":" + tac +
			"}]}"

		url := BASE_URL + "configuration"
		logger.MagmaGwRegLog.Infoln("plmnConfigurationList:", plmnConfigurationList)
		//status, data, err := magmanbi.SendHttpRequest("PUT", url, "{\"plmn_list\":[{ \"mcc\":\"208\",\"mnc\":\"95\",\"slice_list\":[{\"sst\":20, \"sd\":20}],\"tac\":40960}]}")
		status, data, err := magmanbi.SendHttpRequest("PUT", url, plmnConfigurationList)
		if status != 200 {
			logger.MagmaGwRegLog.Infoln("HTTP request failed, Details- :", status, err)
			logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		} else {
			logger.MagmaGwRegLog.Infoln("AMF has been updated:", data)
		}
	} else {
		logger.MagmaGwRegLog.Infoln("AMF has the same values", newdata)
	}

}

func GetPlmn() []string {
	url := BASE_URL + "configuration"
	status, data, _ := magmanbi.SendHttpRequest("GET", url, "")
	if status != 200 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed with code:", status)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		return nil
	} else {
		data, err := magmanbi.PrettyString([]byte(data))
		if err != nil {
			logger.MagmaGwRegLog.Panicln(err)
			return nil
		} else {
			newdata := []string(strings.Split(data, ","))

			plmnConfigPositionsStart, foundplmn := magmanbi.Find(newdata, "plmn_list")
			plmnConfigPositionEnd, foundtac := magmanbi.Find(newdata, "tac")
			if foundplmn && foundtac {
				plmnConfig := newdata[plmnConfigPositionsStart : plmnConfigPositionEnd+1]
				logger.MagmaGwRegLog.Infoln("PLMN: \n", plmnConfig)
				return plmnConfig
			} else{
				return nil
			}
		}
	}
}

// func PrettyString(str []byte) (string, error) {
// 	var prettyJSON bytes.Buffer
// 	if err := json.Indent(&prettyJSON, str, "", "  "); err != nil {
// 		return "", err
// 	}
// 	return prettyJSON.String(), nil
// }
