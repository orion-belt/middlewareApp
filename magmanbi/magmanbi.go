package magmanbi

import (
	"bytes"
	"context"
	"encoding/json"
	"middlewareApp/config"
	"middlewareApp/common"
	"middlewareApp/apiconv"
	"middlewareApp/logger"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
	"os"
	"os/exec"
	"strings"
	"time"
)

type rawMconfigMsg struct {
	ConfigsByKey map[string]json.RawMessage
}

const (
	nbi_service          = "magma"
	ServiceName          = "streamer"
	SubscriberStreamName = "subscriberdb"
	ConfigStreamName     = "configs"
)

var nbi_base_url string
var nbi_stream_url string
var nbi_stream_authority string
var networkID string
var gatewayID string
var hardwareID string
var stream_interval time.Duration

func Init() {
	// Initialze global variables
	nbi_base_url = config.GetCloudHttpUrl(nbi_service)
	nbi_stream_url = config.GetCloudGrpcUrl(nbi_service)
	nbi_stream_authority = config.GetCloudAuthority(nbi_service, ServiceName)
	networkID = config.GetCloudNetworkId(nbi_service)
	gatewayID = config.GetCloudGatewayId(nbi_service)
	stream_interval = config.GetCloudStreamInterval(nbi_service)
	hardwareID = config.GetHardwareId(nbi_service)

	// Register OAI Gateway (5GCN) to Magma Orchestreator
	logger.MagmaGwRegLog.Infoln("Get registered networks at Orchestrator")
	if RegisterOaiNetwork() {
		time.Sleep((2 * time.Second))

		for {
			if IsNetworkReady(networkID) {
				for {
					if IsGatewayReady(gatewayID) {
						// Start concurrent stream updates for config, subscriber etc.
						logger.MagmaGwRegLog.Infoln("Generate gateway certs")
						if GenerateGatewayCerts() {
							go StreamConfigUpdates()
							go StreamSubscriberUpdates()
							select {} // Run these threads only
						}
					}
					time.Sleep((2 * time.Second))
				}
			}
			time.Sleep((2 * time.Second))
		}

	}
}

func RegisterOaiNetwork() bool {
	if !config.IsRegisterGateway(nbi_service) {
		logger.MagmaGwRegLog.Infoln("Gateway registration disabled")
		return true
	}

	url := nbi_base_url + "networks"

	status, data, _ := SendHttpRequest("GET", url, "")
	if status != 200 {
		logger.MagmaGwRegLog.Errorln("HTTP request failed with code :", status)
		logger.MagmaGwRegLog.Errorln("HTTP response body :", data)
		logger.MagmaGwRegLog.Errorln("<<<< Make sure controller address and certificates are correctly provided ")
		return false
	} else {
		logger.MagmaGwRegLog.Infoln("Registered networks at Orchestrator: ", data)
		if strings.Contains(string(data), networkID) {
			logger.MagmaGwRegLog.Infoln("OAI network is already registered")
			RegisterOaiGateway()
		} else {
			logger.MagmaGwRegLog.Infoln("OAI network is not registered yet")
			RegisterNetwork()
			RegisterTier()
			RegisterGateway()
		}
	}
	return true
}

func RegisterNetwork() {
	logger.MagmaGwRegLog.Infoln("Registering OAI network")
	// Create test network
	var jsonData []byte
	jsonData, _ = json.Marshal(GetDefaultLteNetwork(networkID))
	data, _ := PrettyString([]byte(jsonData))

	logger.MagmaGwRegLog.Infoln("Network Config:- \n", data)

	url := nbi_base_url + "lte"
	status, data, err := SendHttpRequest("POST", url, string(jsonData))
	if status != 201 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed, Details- :", status, err)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		return
	} else {
		url := nbi_base_url + "networks"
		_, data, _ := SendHttpRequest("GET", url, "")
		logger.MagmaGwRegLog.Infoln("Registered networks at Orchestrator: ", data)
		logger.MagmaGwRegLog.Infoln("OAI network registered successfully")
	}
}

func RegisterTier() {
	logger.MagmaGwRegLog.Infoln("Registering OAI tier")
	// Create test network
	var jsonData []byte
	jsonData, _ = json.Marshal(GetDefaultTier())

	url := nbi_base_url + "networks/" + networkID + "/tiers"
	status, data, err := SendHttpRequest("POST", url, string(jsonData))
	if status != 201 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed, Details- :", status, err)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		return
	} else {
		_, data, _ := SendHttpRequest("GET", url, "")
		logger.MagmaGwRegLog.Infoln("Tier details: ", data)
		logger.MagmaGwRegLog.Infoln("OAI tier created successfully")
	}
}
func RegisterOaiGateway() {
	url := nbi_base_url + "networks/" + networkID + "/gateways"
	status, data, _ := SendHttpRequest("GET", url, "")
	if status != 200 {
		logger.MagmaGwRegLog.Errorln("HTTP request failed with code :", status)
		logger.MagmaGwRegLog.Errorln("HTTP response body :", data)
	} else {
		logger.MagmaGwRegLog.Infoln("Registered gateways at Orchestrator: ", data)
		if strings.Contains(string(data), gatewayID) {
			logger.MagmaGwRegLog.Infoln("OAI gateway is already registered")
		} else {
			logger.MagmaGwRegLog.Infoln("OAI gateway is not registered yet")
			RegisterGateway()
		}
	}
}
func RegisterGateway() {
	var jsonData []byte
	jsonData, _ = json.Marshal(GetDefaultLteGateway(gatewayID, hardwareID))

	url := nbi_base_url + "lte/" + networkID + "/gateways"
	status, data, err := SendHttpRequest("POST", url, string(jsonData))
	if status != 201 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed, Details- :", status, err)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		return
	} else {
		_, data, _ := SendHttpRequest("GET", url, "")
		data, _ = PrettyString([]byte(data))
		logger.MagmaGwRegLog.Infoln("Gateway details: \n", data)
		logger.MagmaGwRegLog.Infoln("OAI gateway registered successfully")
	}
}

func StreamConfigUpdates() {
	logger.MagmaGwRegLog.Infoln("Stareaming updates from Orcheatrator")
	var mme common.MME

	conn, _ := GetCloudConnection(nbi_stream_authority, nbi_stream_url)
	streamerClient := protos.NewStreamerClient(conn)
	for {
		stream, _ := streamerClient.GetUpdates(context.Background(), &protos.StreamRequest{GatewayId: hardwareID, StreamName: ConfigStreamName})
		actualMarshaled, _ := stream.Recv()
		// println(actualMarshaled.String())

		cfg := &protos.GatewayConfigs{}
		protos.UnmarshalMconfig(actualMarshaled.Updates[0].GetValue(), cfg)

		newCfg := &rawMconfigMsg{ConfigsByKey: map[string]json.RawMessage{}}
		json.Unmarshal(actualMarshaled.Updates[0].GetValue(), newCfg)
		err := json.Unmarshal([]byte(string(newCfg.ConfigsByKey["mme"])), &mme)
		if err != nil {
			logger.MagmaGwRegLog.Errorln("Error parsing mme config -", err)
		} else {
			apiconv.CheckForUpdate(&mme)
		}

		data, _ := PrettyString([]byte(string(newCfg.ConfigsByKey["mme"])))
		logger.MagmaGwRegLog.Infoln("\n", data)
		logger.MagmaGwRegLog.Infoln("Stareaming config updates from Orcheatrator [StreamInterval : ", stream_interval*time.Second, "]")
		time.Sleep((stream_interval * time.Second))
	}
}

func StreamSubscriberUpdates() {
	logger.MagmaGwRegLog.Infoln("Stareaming subscriber updates from Orcheatrator")

	conn, _ := GetCloudConnection(nbi_stream_authority, nbi_stream_url)

	streamerClient := protos.NewStreamerClient(conn)
	for {
		stream, _ := streamerClient.GetUpdates(context.Background(), &protos.StreamRequest{GatewayId: hardwareID, StreamName: SubscriberStreamName, ExtraArgs: nil})
		actualMarshaled, _ := stream.Recv()
		num_sub := len(actualMarshaled.Updates)
		logger.MagmaGwRegLog.Infoln("Number of subscribes", num_sub)
		for i := 0; i < num_sub; i++ {
			logger.MagmaGwRegLog.Infoln("Subscriber ", actualMarshaled.Updates[i].Key, " information", (actualMarshaled.Updates[i].String()))
		}
		logger.MagmaGwRegLog.Infoln("Stareaming subscriber updates from Orcheatrator [StreamInterval : ", (stream_interval)*time.Second, "]")
		time.Sleep(((stream_interval) * time.Second))
	}
}

func PrettyString(str []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, str, "", "  "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func GenerateGatewayCerts() bool {
	base_path, _ := os.Getwd()
	_, err := exec.Command(base_path+"/magmanbi/scripts/generate_gateway_certs.sh", hardwareID).Output()
	if err != nil {
		logger.MagmaGwRegLog.Errorln("Certificate generation failed --> Verify if correct key type is used")
		return false
	}
	logger.MagmaGwRegLog.Infoln("Gateway certificates generated sucessfully")
	return true
}

func IsNetworkReady(NetworkID string) bool {
	url := nbi_base_url + "lte"
	_, data, _ := SendHttpRequest("GET", url, "")
	if strings.Contains(string(data), NetworkID) {
		logger.MagmaGwRegLog.Infoln("Network ready for stream")
		return true
	} else {
		logger.MagmaGwRegLog.Warnln("Network not ready for stream")
		return false
	}
}

func IsGatewayReady(gatewayID string) bool {
	url := nbi_base_url + "networks/" + networkID + "/gateways"
	_, data, _ := SendHttpRequest("GET", url, "")
	if strings.Contains(string(data), gatewayID) {
		logger.MagmaGwRegLog.Infoln("Gateway ready for stream")
		return true
	} else {
		logger.MagmaGwRegLog.Warnln("Gateway not ready for stream")
		return false
	}
}
