package magmanbi

import (
	"context"
	"time"
    "os"
    "os/exec"
	"bytes"
	"encoding/json"
	"middlewareApp/logger"
	"strings"
	"middlewareApp/magmanbi/orc8r/lib/go/registry"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
)

type rawMconfigMsg struct {
	ConfigsByKey map[string]json.RawMessage
}
const (
	BASE_URL = "https://127.0.0.1:9443/magma/v1/"
	ServiceName = "streamer"
	hardwareID = "c29f9ded-0d34-4e64-9ee5-c66d202081d6"
)

func Init() {
	// Start GRPC server for Streamer service of Magma Orchestreator
	logger.MagmaLog.Infoln("Initaiating gRPC NBI for Magma Orchestrator")

	// Register OAI Gateway (5GCN) to Magma Orchestreator
	logger.MagmaGwRegLog.Infoln("Get registered networks at Orchestrator")
	RegisterOaiNetwork()
	GenerateGatewayCerts()
	time.Sleep((2 * time.Second))
	StreamUpdates()
}

func RegisterOaiNetwork() {

	url := BASE_URL + "networks"
	status, data, _ := SendHttpRequest("GET", url, "")
	if status != 200 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed with code :", status)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		return
	} else {
		logger.MagmaGwRegLog.Infoln("Registered networks at Orchestrator: ", data)
		if strings.Contains(string(data), "oai") {
			logger.MagmaGwRegLog.Infoln("OAI network is already registered")
			return
		} else {
			logger.MagmaGwRegLog.Infoln("OAI network is not registered yet")
			RegisterNetwork()
			RegisterTier()
			RegisterGateway()
		}
	}
}

func RegisterNetwork() {
	logger.MagmaGwRegLog.Infoln("Registering OAI network")
	// Create test network
	networkID := "oai"
	var jsonData []byte
	jsonData, _ = json.Marshal(GetDefaultLteNetwork(networkID))
	data, _ := PrettyString([]byte(jsonData))

	logger.MagmaGwRegLog.Infoln("Network Config:- \n", data)

	url := BASE_URL + "lte"
	status, data, err := SendHttpRequest("POST", url, string(jsonData))
	if status != 201 {
		logger.MagmaGwRegLog.Infoln("HTTP request failed, Details- :", status, err)
		logger.MagmaGwRegLog.Infoln("HTTP response body :", data)
		return
	} else {
		url := BASE_URL + "networks"
		_, data, _ := SendHttpRequest("GET", url, "")
		logger.MagmaGwRegLog.Infoln("Registered networks at Orchestrator: ", data)
		logger.MagmaGwRegLog.Infoln("OAI network registered successfully")
	}
}

func RegisterTier() {
	logger.MagmaGwRegLog.Infoln("Registering OAI tier")
	// Create test network
	networkID := "oai"
	var jsonData []byte
	jsonData, _ = json.Marshal(GetDefaultTier())
	// logger.MagmaGwRegLog.Infoln("Tier Config:- ", string(jsonData)	)

	url := BASE_URL + "networks/"+networkID+"/tiers"
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

func RegisterGateway () {
	gatewayID := "oaigw1"
	// hardwareID := "c29f9ded-0d34-4e64-9ee5-c66d202081d6"
	networkID := "oai"

	var jsonData []byte
	jsonData, _ = json.Marshal(GetDefaultLteGateway(gatewayID, hardwareID))

	url := BASE_URL + "lte/"+networkID+"/gateways"
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

func PrettyString(str []byte) (string, error) {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, str, "", "  "); err != nil {
		return "", err
	}
	return prettyJSON.String(), nil
}

func GenerateGatewayCerts(){
	base_path, _:= os.Getwd()
	_, err := exec.Command(base_path+"/magmanbi/scripts/generate_gateway_certs.sh", hardwareID).Output()
    if err != nil {
        logger.MagmaGwRegLog.Panicln(err)
		return
    }
    logger.MagmaGwRegLog.Infoln("Gateway certificates generated sucessfully")
}

func StreamUpdates(){
       logger.MagmaGwRegLog.Infoln("Stareaming updates from Orcheatrator")

		conn, _ := registry.Get().GetCloudConnection(ServiceName)
		streamerClient := protos.NewStreamerClient(conn)
		for {
		stream, _ := streamerClient.GetUpdates(context.Background(), &protos.StreamRequest{GatewayId: "c29f9ded-0d34-4e64-9ee5-c66d202081d6", StreamName: "configs"})
		actualMarshaled, _ := stream.Recv()
		// println(actualMarshaled.String())
	
		cfg := &protos.GatewayConfigs{}
		protos.UnmarshalMconfig(actualMarshaled.Updates[0].GetValue(), cfg)
	
		newCfg := &rawMconfigMsg{ConfigsByKey: map[string]json.RawMessage{}}
		json.Unmarshal(actualMarshaled.Updates[0].GetValue(), newCfg)
		data, _ := PrettyString([]byte(string(newCfg.ConfigsByKey["mme"])))
		logger.MagmaGwRegLog.Infoln("\n",data)
		logger.MagmaGwRegLog.Infoln("Stareaming updates from Orcheatrator [StreamInterval : 5]")
		time.Sleep((3 * time.Second))
	}
}