package server

import (
	"log"

	"github.com/AgencyPMG/go-from-scratch/app/server/api"
	"github.com/AgencyPMG/go-from-scratch/app/server/login"
)

type Server struct {
	Logger *log.Logger
	API    *api.API
	Login  *login.Login
}
