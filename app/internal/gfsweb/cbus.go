package gfsweb

import (
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user/usercmd"
	"github.com/gogolfing/cbus"
)

func CreateCBus(repos *Repos) *cbus.Bus {
	bus := &cbus.Bus{}

	RegisterUserCommands(bus, repos.Users)

	return bus
}

func RegisterUserCommands(bus *cbus.Bus, users user.Repo) {
	h := &usercmd.Handler{
		Users: users,
	}

	bus.Handle(&usercmd.CreateUserCommand{}, cbus.HandlerFunc(h.CreateUser))
	bus.Handle(&usercmd.UpdateUserCommand{}, cbus.HandlerFunc(h.UpdateUser))
}
