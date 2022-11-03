package main

import (
	"middlewareApp/config"
	"middlewareApp/logger"
	"middlewareApp/magmanbi"
	"middlewareApp/oaisbi"
	"os"
	"time"

	"github.com/urfave/cli"
)

var UpdateAllSlice_Lists = true
var PlmnPos = 0
var SlicePos = 1

func main() {
	app := cli.NewApp()
	app.Name = "Openairinterface middlewareApp"
	app.Description = "Middleware application for OAI 5GCN Orchestration"
	app.Usage = "./middlewareApp --config [middlewareApp configuration file]"
	app.Author = "Openairinterface"
	app.Email = "contact@openairinterface.org"
	app.Action = AppInit

	logger.AppLog.Infoln("")
	logger.AppLog.Infoln("=====================================================================")
	logger.AppLog.Infoln("App Name        :", app.Name)
	logger.AppLog.Infoln("Description     :", app.Description)
	logger.AppLog.Infoln("Auther          :", app.Author)
	logger.AppLog.Infoln("Contact         :", app.Email)
	logger.AppLog.Infoln("=====================================================================")
	logger.AppLog.Infoln("")

	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorln("Failed to run middlewareApp:", err)
		return
	}
}

func AppInit(c *cli.Context) error {

	// Load configuration
	cfg := c.String("config")
	if cfg == "" {
		base_path, _ := os.Getwd()
		configPath := base_path + config.CONFIG_PATH
		logger.AppLog.Warnln("No configuration file provided. Using default configuration file:", configPath)
		logger.AppLog.Infoln("Application Usage:", c.App.Usage)
		cfg = base_path + config.CONFIG_PATH
	}

	if err := config.LoadConfig(cfg); err != nil {
		logger.AppLog.Errorln("Failed to load config:", err)
		return err
	}
	// Initialise middleware services
	go magmanbi.Init()
	if magmanbi.RegisterOaiNetwork() {
		time.Sleep((2 * time.Second))
		for {
			if magmanbi.IsNetworkReady(magmanbi.NetworkID) {
				for {
					if magmanbi.IsGatewayReady(magmanbi.GatewayID) {
						// Start concurrent stream updates for config, subscriber etc.
						logger.MagmaGwRegLog.Infoln("Generate gateway certs")
						if magmanbi.GenerateGatewayCerts() {
							logger.AppLog.Infoln("\n\n")
							time.Sleep((5 * time.Second))
							go magmanbi.StreamConfigUpdates()
							go magmanbi.StreamSubscriberUpdates()
							if UpdateAllSlice_Lists {
								go oaisbi.UpdateAmfPlmnForAllElements()
							} else {
								go oaisbi.UpdateAmfPlmnForSpecificElement(PlmnPos, SlicePos)
							}
							go oaisbi.GetPlmn()
							select {} // Run these threads only
						}
					}
					time.Sleep((2 * time.Second))
				}
			}
			time.Sleep((2 * time.Second))
		}
	}
	return nil
}
