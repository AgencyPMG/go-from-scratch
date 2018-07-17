package gfsweb

import (
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client/clientcmd"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user/usercmd"
	"github.com/gogolfing/cbus"
)

func CreateCBus(repos *Repos) *cbus.Bus {
	bus := &cbus.Bus{}

	RegisterUserCommands(bus, repos.Users)
	RegisterClientCommands(bus, repos.Clients)

	return bus
}

func RegisterUserCommands(bus *cbus.Bus, users user.Repo) {
	h := &usercmd.Handler{
		Users: users,
	}

	bus.Handle(&usercmd.CreateUserCommand{}, cbus.HandlerFunc(h.CreateUser))
	bus.Handle(&usercmd.UpdateUserCommand{}, cbus.HandlerFunc(h.UpdateUser))
}

func RegisterClientCommands(bus *cbus.Bus, clients client.Repo) {
	h := &clientcmd.Handler{
		Clients: clients,
	}

	bus.Handle(&clientcmd.CreateClientCommand{}, cbus.HandlerFunc(h.CreateClient))
	bus.Handle(&clientcmd.UpdateClientCommand{}, cbus.HandlerFunc(h.UpdateClient))
	bus.Handle(&clientcmd.DeleteClientCommand{}, cbus.HandlerFunc(h.DeleteClient))
}
