package config

import (
	"github.com/google/uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"middlewareApp/logger"
	"os"
	"strconv"
	"time"
)

// Map represents a map generated from a service YML file.
type Map struct {
	RawMap map[interface{}]interface{}
}

const (
	CONFIG_PATH = "/config/config.yaml"
)

var MwConfig *Config

func LoadConfig(f string) error {
	content, err := ioutil.ReadFile(f)
	if err != nil {
		logger.ConfigLog.Errorln("Failed to read", f, "file:", err)
		return err
	}
	MwConfig = &Config{}

	err = yaml.Unmarshal(content, MwConfig)
	if err != nil {
		logger.ConfigLog.Errorln("Failed to unmarshal:", err)
		return err
	}

	err = MwConfig.Validate()
	if err != nil {
		logger.ConfigLog.Errorln("Invalid Configuration:", err)
	}
	return nil
}

func GetCloudGrpcUrl(service string) string {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			BASE_URL := nbi_tmp.CloudAddress + ":" + strconv.Itoa(nbi_tmp.CloudStreamerPort)
			return BASE_URL
		}
	}
	return ""
}


func GetCloudAuthority(service string, servicename string) string {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			BASE_URL := servicename+"-"+nbi_tmp.CloudAddress
			return BASE_URL
		}
	}
	return ""
}

func GetCloudHttpUrl(service string) string {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			BASE_URL := "https://" + nbi_tmp.CloudAddress + ":" + strconv.Itoa(nbi_tmp.CloudHTTPSPort) + "/magma/v1/"
			return BASE_URL
		}
	}
	return ""
}

func IsRegisterGateway(service string) bool {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			return nbi_tmp.RegisterGateway
		}
	}
	return false
}

func GetCloudNetworkId(service string) string {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			return nbi_tmp.NetworkID
		}
	}
	return ""
}

func GetCloudGatewayId(service string) string {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			return nbi_tmp.GatewayID
		}
	}
	return ""
}

func GetCloudStreamInterval(service string) time.Duration {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			return time.Duration(nbi_tmp.StreamInterval)
		}
	}
	return 0
}
func GetHardwareId(service string) string {
	for _, nbi_tmp := range MwConfig.Configuration.Nbis {
		if nbi_tmp.Name == service {
			if !nbi_tmp.GeneratehwId {
				return nbi_tmp.HardwareId
			} else {
				return GetUUID()
			}
		}
	}
	return ""
}

func GetUUID() string {
	id := uuid.New()
	return id.String()
}

func GetGatewayCerds() (string, string) {
	base_path, _ := os.Getwd()
	clientCaFile := base_path + "/magmanbi/.certs/gateway.crt"
	clientKeyFile := base_path + "/magmanbi/.certs/gateway.key"
	return clientCaFile, clientKeyFile
}

func GetOaiService (service string) string {
	for _, sbi_tmp := range MwConfig.Configuration.Sbi.ServicesList {
		if sbi_tmp.Name == service {
			BASE_URL := "https://" + sbi_tmp.Address + ":" + strconv.Itoa(sbi_tmp.Port)
			return BASE_URL
		}
	}
	return ""
}