package main

import (
	"github.com/jindrichskupa/vocative-api/app"
	"github.com/jindrichskupa/vocative-api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(config.ListenAddress())
}
