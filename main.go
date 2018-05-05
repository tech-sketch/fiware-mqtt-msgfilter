/*
Package main : entry point of fiware-distinct-message

	license: Apache license 2.0
	copyright: Nobuyuki Matsui <nobuyuki.matsui@gmail.com>
*/
package main

import (
	"github.com/tech-sketch/fiware-mqtt-msgfilter/conf"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/router"
	"github.com/tech-sketch/fiware-mqtt-msgfilter/utils"
)

func main() {
	logger := utils.NewLogger("main")
	config := conf.NewConfig()
	handler, err := router.NewHandler(config)
	if err != nil {
		logger.Errorf("NewHandler raise error: %s", err)
		return
	}
	handler.Run(config.ListenPort)
}
