package api

import (
	"net/http"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/gogolfing/cbus"
)

type API struct {
	Bus *cbus.Bus

	Users   user.QueryRepo
	Clients client.QueryRepo
}

func (a *API) Handler() http.Handler {
	return http.HandlerFunc(http.NotFound)
}
