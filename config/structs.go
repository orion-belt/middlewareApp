package config

import (
	"fmt"
)

type Config struct {
	Info          *Info          `yaml:"info"`
	Configuration *Configuration `yaml:"configuration"`
	Logger        *Logger        `yaml:"logger"`
}

type Info struct {
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

type Configuration struct {
	Nbis map[string]*NbiService `yaml:"nbi"`
	Sbi  SbiService             `yaml:"sbi"`
}

type Logger struct {
	LogLevel string `yaml:"logLevel"`
}

type NbiService struct {
	Name              string `yaml:"name"`
	NetworkID         string `yaml:"networkId"`
	GatewayID         string `yaml:"gatewayId"`
	HardwareId        string `yaml:"hardwareId"`
	CloudAddress      string `yaml:"cloudAddress"`
	CloudStreamerPort int    `yaml:"cloudStreamerPort"`
	CloudHTTPSPort    int    `yaml:"cloudHttpsPort"`
	StreamInterval    int    `yaml:"streamInterval"`
	RegisterGateway   bool   `yaml:"registerGateway"`
	GeneratehwId      bool   `yaml:"generatehwId"`
}

type Services struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `yaml:"port"`
}
type SbiService struct {
	EnableReporting bool       `yaml:"enableReporting"`
	EnableHttp2     bool       `yaml:"enableHttp2"`
	ServicesList    []Services `yaml:"services"`
}

func (c *Config) Validate() (err error) {

	if c.Info == nil {
		return fmt.Errorf("Info field missing")
	}

	if c.Configuration == nil {
		return fmt.Errorf("Configuration field missing")
	}

	for _, nbi_tmp := range c.Configuration.Nbis {
		if nbi_tmp.Name != "magma" {
			return fmt.Errorf("NBI not supported. Supported NBI - magma")
		}
	}

	for _, sbi_tmp := range c.Configuration.Sbi.ServicesList {
		if sbi_tmp.Name == "" {
			return fmt.Errorf("service name missing")
		}
	}
	return err
}
