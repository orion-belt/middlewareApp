package main

import (
	"middlewareApp/logger"
	"middlewareApp/magmanbi"
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

	cfg := c.String("cfg")
	if cfg == "" {
		logger.AppLog.Warnln("No configuration file provided. Using default configuration file:")
		logger.AppLog.Infoln("Application Usage:", c.App.Usage)
	}

	go magmanbi.Init()
	select {} // block forever
	// return nil
}
