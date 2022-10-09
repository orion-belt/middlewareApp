package magmanbi

import (
	// "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"middlewareApp/magmanbi/models"
)

//GetDefaultLteNetwork gets a default LTE network instance
func GetDefaultLteNetwork(networkID string) *models.LTENetwork {
	return &models.LTENetwork{
		ID:          models.NetworkID(networkID),
		Name:        "5G Network",
		Description: "OAI 5G Core Network",
		DNS: &models.NetworkDNSConfig{
			EnableCaching: swag.Bool(false),
			LocalTTL:      swag.Uint32(60),
		},
		Cellular: &models.NetworkCellularConfigs{
			Ran: &models.NetworkRanConfigs{
				BandwidthMhz: 20,
				TddConfig: &models.NetworkRanConfigsTddConfig{
					Earfcndl:               44590,
					SubframeAssignment:     2,
					SpecialSubframePattern: 7,
				},
			},
			Epc: &models.NetworkEpcConfigs{
				Mcc: "208",
				Mnc: "95",
				Tac: 1,
				// 16 bytes of \x11
				LTEAuthOp:  []byte("\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11\x11"),
				LTEAuthAmf: []byte("\x80\x00"),

				HssRelayEnabled:          swag.Bool(false),
				GxGyRelayEnabled:         swag.Bool(false),
				CloudSubscriberdbEnabled: false,
				DefaultRuleID:            "",
			},
		},
	}
}

//GetDefaultLteGateway gets a default LTE gateway instance
func GetDefaultLteGateway(gatewayID string, hardwareID string) *models.MutableLTEGateway {
	return &models.MutableLTEGateway{
		Device: &models.GatewayDevice{
			HardwareID: hardwareID,
			Key: &models.ChallengeKey{
				KeyType: "ECHO",
			},
		},
		ID:          models.GatewayID(gatewayID),
		Name:        "OAI Gateway",
		Description: "OAI 5G Gateway",
		Cellular: &models.GatewayCellularConfigs{
			Ran: &models.GatewayRanConfigs{
				Pci:             260,
				TransmitEnabled: swag.Bool(true),
			},
			Epc: &models.GatewayEpcConfigs{
				NatEnabled: swag.Bool(true),
				IPBlock:    "192.168.128.0/24",
			},
		},
		Magmad: &models.MagmadGatewayConfigs{
			CheckinInterval:         15,
			CheckinTimeout:          10,
			AutoupgradePollInterval: 300,
			AutoupgradeEnabled:      swag.Bool(true),
		},
		ConnectedENODEBSerials: []string{},
		Tier:                   "default",
		APNResources:           models.APNResources{},
	}
}

//GetDefaultTier gets the default tier
func GetDefaultTier() *models.Tier {
	return &models.Tier{
		ID:       "default",
		Gateways: models.TierGateways{},
		Images:   models.TierImages{},
		Version:  models.TierVersion("1.0"),
	}
}
