package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api/dto"
	"github.com/gogolfing/cbus"
	"github.com/gorilla/mux"
)

//ErrBadRequestRouteParameter is a sentinel error indicating a route parameter
//could not be parsed or found.
var ErrBadRequestRouteParameter = errors.New("api: bad request route parameter")

//API allows for handling all incoming requests dealing with the application's API.
//Its Handler method should be used for serving API requests.
//
//All fields must be non-nil and safe for use by multiple goroutintes due to their
//use inside a http.Server and http.Handler.
type API struct {
	//Bus is a command bus that is used to execute commands to update entities
	//within the API.
	Bus *cbus.Bus

	//Users is a query repository used to retrieve Users.
	Users user.QueryRepo

	//Clients is a query repository used to retrieve Clients.
	Clients client.QueryRepo
}

//Handler returns an http.Handler that serves all requests for a.
func (a *API) Handler() http.Handler {
	router := mux.NewRouter()

	//User routes.
	router.HandleFunc("/users", a.listUsers).
		Methods(http.MethodGet)

	router.HandleFunc("/users", a.createUser).
		Methods(http.MethodPost)

	router.HandleFunc("/users/{"+routeParamUserId+"}", a.getUser).
		Methods(http.MethodGet)

	router.HandleFunc("/users/{"+routeParamUserId+"}", a.updateUser).
		Methods(http.MethodPatch)

	//Client routes.
	router.HandleFunc("/clients", a.listClients).
		Methods(http.MethodGet)

	router.HandleFunc("/clients", a.createClient).
		Methods(http.MethodPost)

	router.HandleFunc("/clients/{"+routeParamClientId+"}", a.getClient).
		Methods(http.MethodGet)

	router.HandleFunc("/clients/{"+routeParamClientId+"}", a.updateClient).
		Methods(http.MethodPatch)

	router.HandleFunc("/clients/{"+routeParamClientId+"}", a.deleteClient).
		Methods(http.MethodDelete)

	return router
}

//idFrom is a helper method to retrieve data.Id(s) from route parameters.
func (a *API) idFrom(w http.ResponseWriter, r *http.Request, varName string) (data.Id, bool) {
	vars := mux.Vars(r)

	varId, ok := vars[varName]
	if !ok {
		a.sendError(w, ErrBadRequestRouteParameter, http.StatusBadRequest)
		return data.EmptyId(), false
	}

	id, err := data.ParseId(varId)
	if err != nil {
		a.sendError(w, err, http.StatusBadRequest)
		return data.EmptyId(), false
	}

	return id, true
}

//parseForm is a helper method to parse form types that this package knows how
//to use.
func (a *API) parseForm(w http.ResponseWriter, r *http.Request, f interface{}) bool {
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(f)
	if err != nil {
		a.sendError(w, err, http.StatusBadRequest)
		return false
	}

	//Do some sort of validation here on f.

	return true
}

//sendData is a helper method to send successful responses by marshaling data
//and sending a response with status.
func (a *API) sendData(w http.ResponseWriter, data interface{}, status int) {
	transformed, err := dto.Transform(data)
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError)
		return
	}

	a.sendResponse(w, transformed, status)
}

//sendError is a helper method to send failure responses by marshaling err into
//a simple value and sending status.
func (a *API) sendError(w http.ResponseWriter, err error, status int) {
	a.sendResponse(
		w,
		&apiError{
			Error: err.Error(),
		},
		status,
	)
}

//sendResponse is a helper method to send all responses by marshaling resp and
//sending status.
func (a *API) sendResponse(w http.ResponseWriter, resp interface{}, status int) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t") //Do this for prettier output for our example application.
	enc.Encode(resp)
}

//apiError is a simple type that knows how to marshal an error.
type apiError struct {
	Error string `json:"error"`
}
