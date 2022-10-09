package main

import (
	"os"
  "./logger"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "middlewareApp"
	app.Usage = "./gnbsim --cfg [gnbsim configuration file]"
	// app.Action = action
//	app.Flags = getCliFlags()

	logger.AppLog.Infoln("App Name:", app.Name)

	if err := app.Run(os.Args); err != nil {
		logger.AppLog.Errorln("Failed to run middlewareApp:", err)
		return
	}
