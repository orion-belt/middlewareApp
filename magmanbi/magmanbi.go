package magmanbi

import (
	"bytes"
	"context"
	"encoding/json"
	"middlewareApp/logger"
	"middlewareApp/magmanbi/orc8r/lib/go/protos"
	"middlewareApp/magmanbi/orc8r/lib/go/registry"
	"os"
	//	"fmt"
	"os/exec"
	"strings"
	//"time"
)

type rawMconfigMsg struct {
	ConfigsByKey map[string]json.RawMessage
}

const (
	BASE_URL    = "https://127.0.0.1:9443/magma/v1/"
	ServiceName = "streamer"
	hardwareID  = "c29f9ded-0d34-4e64-9ee5-c66d202081d6"
)

type SliceUpdatedValues struct {
	Updatedsstsd bool
	Sst int
	Sd  string
}

var Slices SliceUpdatedValues

func DefaultInit() {
	sliceJson := `{"updatedsstsd":false, "sst":0,"sd":"0"}`
	json.Unmarshal([]byte(sliceJson), &Slices)
}

type MmeStruct struct {
	Types                         string
	LogLevel                      string
	Mcc                           string //int
	Mnc                           string //int
	Tac                           int
	MmeGid                        int
	MmeCode                       int
	Lac                           int
	MmeRelativeCapacity           int
	NatEnabled                    bool
	AmfDefaultSliceServiceType    int
	AmfDefaultSliceDifferentiator string //int
	AmfName                       string
	AmfRegionId                   string //int
	AmfSetId                      string //int
	AmfPointer                    string //int
}

func Init() {
	// Start GRPC server for Streamer service of Magma Orchestreator
	logger.MagmaLog.Infoln("Initaiating gRPC NBI for Magma Orchestrator")

	// Register OAI Gateway (5GCN) to Magma Orchestreator
	logger.MagmaGwRegLog.Infoln("Get registered networks at Orchestrator")
	RegisterOaiNetwork()
	GenerateGatewayCerts()
	DefaultInit()
	//GetPlmn ()
	//time.Sleep((2 * time.Second))
	//StreamUpdates()
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

	url := BASE_URL + "networks/" + networkID + "/tiers"
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

func RegisterGateway() {
	gatewayID := "oaigw1"
	// hardwareID := "c29f9ded-0d34-4e64-9ee5-c66d202081d6"
	networkID := "oai"

	var jsonData []byte
	jsonData, _ = json.Marshal(GetDefaultLteGateway(gatewayID, hardwareID))

	url := BASE_URL + "lte/" + networkID + "/gateways"
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

func GenerateGatewayCerts() {
	base_path, _ := os.Getwd()
	_, err := exec.Command(base_path+"/magmanbi/scripts/generate_gateway_certs.sh", hardwareID).Output()
	if err != nil {
		logger.MagmaGwRegLog.Panicln(err)
		return
	}
	logger.MagmaGwRegLog.Infoln("Gateway certificates generated sucessfully")
}

func StreamUpdates() {
	logger.MagmaGwRegLog.Infoln("Streaming updates from Orcheatrator")

	var mmes MmeStruct
	
	conn, _ := registry.Get().GetCloudConnection(ServiceName)
	streamerClient := protos.NewStreamerClient(conn)

	stream, _ := streamerClient.GetUpdates(context.Background(), &protos.StreamRequest{GatewayId: "c29f9ded-0d34-4e64-9ee5-c66d202081d6", StreamName: "configs"})
	actualMarshaled, _ := stream.Recv()

	cfg := &protos.GatewayConfigs{}
	protos.UnmarshalMconfig(actualMarshaled.Updates[0].GetValue(), cfg)

	newCfg := &rawMconfigMsg{ConfigsByKey: map[string]json.RawMessage{}}
	json.Unmarshal(actualMarshaled.Updates[0].GetValue(), newCfg)

	err := json.Unmarshal([]byte(string(newCfg.ConfigsByKey["mme"])), &mmes)
	if err != nil {
		panic(err)
	}
	
	if (!(Slices.Sst == mmes.AmfDefaultSliceServiceType)) || (!strings.EqualFold(Slices.Sd, mmes.AmfDefaultSliceDifferentiator)) {
		Slices.Sst = mmes.AmfDefaultSliceServiceType
		Slices.Sd = mmes.AmfDefaultSliceDifferentiator
		Slices.Updatedsstsd = true
	} else {
		Slices.Updatedsstsd = false
	}

	if Slices.Updatedsstsd {
		logger.MagmaGwRegLog.Infoln("PLMN values have been updated on the Orchestrator side \n")
		logger.MagmaGwRegLog.Infoln("The new Slice Service Type (sst) value =", Slices.Sst)
		logger.MagmaGwRegLog.Infoln("The new Slice Differentiator (sd) value =", Slices.Sd)
	}
}