package magmanbi

type OaiAgwConfig struct {
	id          string
	name        string
	description string
	device      *DeviceConfig
	tier        string
}

type MagmadGatewayConfigs struct {
	autoupgrade_enabled       bool
	autoupgrade_poll_interval int
	checkin_interval          int
	checkin_timeout           int
	dynamic_services          []string
}

type DeviceConfig struct {
	hardware_id string
	key         *ChallengeKey
}

type ChallengeKey struct {
	key      string
	key_type string
}

const (
	ECHO                  string = "ECHO"
	SOFTWARE_ECDSA_SHA256        = "SOFTWARE_ECDSA_SHA256"
)

type NetworkEpcConfig struct {
	lte_auth_amf string
	lte_auth_op  string
	mcc          string
	mnc          string
	tac          int
}

type NetworkRANConfig struct {
	//ToDo
}
