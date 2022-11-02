package main

import (
	"middlewareApp/logger"
	"middlewareApp/magmanbi"
	"middlewareApp/config"
	"os"
	"github.com/urfave/cli"
)

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
		base_path,_ := os.Getwd()
		configPath := base_path+config.CONFIG_PATH 
		logger.AppLog.Warnln("No configuration file provided. Using default configuration file:",configPath)
		logger.AppLog.Infoln("Application Usage:", c.App.Usage)
		cfg = base_path+config.CONFIG_PATH 
	} 

	if err := config.LoadConfig(cfg); err != nil {
		logger.AppLog.Errorln("Failed to load config:", err)
		return err
	}

	// Initialise middleware services
	go magmanbi.Init()
	select {} // block forever
	// return nil
}
