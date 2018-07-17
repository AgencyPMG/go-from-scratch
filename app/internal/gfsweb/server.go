package gfsweb

import (
	"strconv"

	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/server"
	"github.com/gogolfing/config"
)

const (
	ConfigKeyServerPort = "server.port"
)

//CreateServer returns a new Server that listens on the port specified in c
//at key ConfigKeyServerPort.
func CreateServer(c *config.Config) *server.Server {
	portStr := c.GetString(ConfigKeyServerPort)
	port, _ := strconv.Atoi(portStr)

	return &server.Server{
		ListenerFactory: server.PortListenerFactory(port),
	}
}
