package main

import (
	"github.com/urfave/cli"
	"middlewareApp/logger"
	"middlewareApp/magmanbi"
	"middlewareApp/oaisbi"
	"os"
	"time"
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

	cfg := c.String("cfg")
	if cfg == "" {
		logger.AppLog.Warnln("No configuration file provided. Using default configuration file:")
		logger.AppLog.Infoln("Application Usage:", c.App.Usage)
	}

	go magmanbi.Init()
	for {
		logger.AppLog.Infoln("\n\n")
		logger.AppLog.Infoln("=====================================================================")
		time.Sleep((5 * time.Second))
		go magmanbi.StreamUpdates()
		if UpdateAllSlice_Lists {
			go oaisbi.UpdateAmfPlmnForAllElements()
		}else{
			go oaisbi.UpdateAmfPlmnForSpecificElement(PlmnPos, SlicePos)
		}
		//go oaisbi.UpdateSlice()
		go oaisbi.GetPlmn()
		logger.AppLog.Infoln("=====================================================================")
		logger.AppLog.Infoln("\n\n")
	}
	select {} // block forever
	// return nil
}
