package gfsweb

import (
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api"
	"github.com/gogolfing/cbus"
)

func CreateAPI(bus *cbus.Bus, repos *Repos) *api.API {
	return &api.API{}
}
