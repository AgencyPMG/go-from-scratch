package gfsweb

import (
	"strconv"

	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/server"
	"github.com/gogolfing/config"
)

const (
	ConfigKeyServerPort = "server.port"
)

func CreateServer(c *config.Config) *server.Server {
	portStr := c.GetString(ConfigKeyServerPort)
	port, _ := strconv.Atoi(portStr)

	return &server.Server{
		ListenerFactory: server.PortListenerFactory(port),
	}
}
