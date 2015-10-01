package api

import "github.com/AgencyPMG/go-from-scratch/app/server/context"

type API struct {
}

func (a *API) Handler() context.Handler {
}

func (a *API) middlewares(context.Handler) context.Handler {
}
