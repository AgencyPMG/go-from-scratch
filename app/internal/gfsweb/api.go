package gfsweb

import (
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api"
	"github.com/gogolfing/cbus"
)

//CreateAPI returns a new api.API populated with bus and the query repositories
//in repos.
func CreateAPI(bus *cbus.Bus, repos *Repos) *api.API {
	return &api.API{
		Bus:     bus,
		Users:   repos.Users,
		Clients: repos.Clients,
	}
}
