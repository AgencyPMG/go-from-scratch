package api

import (
	"net/http"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user/usercmd"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api/dto"
)

const (
	routeParamUserId = "user_id"
)

//getUser attempts to retrieve and send a single User.
func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	user, ok := a.getUserEntity(w, r)
	if !ok {
		return
	}

	a.sendData(w, user, http.StatusOK)
}

//listUsers attempts to retrieve and send all Users.
func (a *API) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.Users.List(r.Context())
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError)
		return
	}

	a.sendData(w, users, http.StatusOK)
}

//createUser attempts to create a new User from a dto.CreateUser form and
//executing a usercmd.CreateUserCommand.
func (a *API) createUser(w http.ResponseWriter, r *http.Request) {
	f := &dto.CreateUser{}

	if ok := a.parseForm(w, r, f); !ok {
		return
	}

	command := &usercmd.CreateUserCommand{
		Email:     f.Email,
		Enabled:   f.Enabled,
		ClientIds: f.ClientIds,
	}

	user, err := a.Bus.ExecuteContext(r.Context(), command)
	a.userCommandResponse(w, user, err, http.StatusCreated)
}

//updateUser attempts to update an existing User from a dto.UpdateUser form
//and executing a usercmd.UpdateUserCommand.
func (a *API) updateUser(w http.ResponseWriter, r *http.Request) {
	f := &dto.UpdateUser{}

	id, _ := a.idFrom(w, r, routeParamUserId)

	if ok := a.parseForm(w, r, f); !ok {
		return
	}

	command := &usercmd.UpdateUserCommand{
		Id:        id,
		Email:     f.Email,
		Enabled:   f.Enabled,
		ClientIds: f.ClientIds,
	}

	user, err := a.Bus.ExecuteContext(r.Context(), command)
	a.userCommandResponse(w, user, err, http.StatusOK)
}

//userCommandResponse is a helper method to send the correct response from the
//result of a user command.
func (a *API) userCommandResponse(w http.ResponseWriter, result interface{}, err error, okStatus int) {
	if err != nil {
		//We would normally inspect the error to figure out what status code to send.
		a.sendError(w, err, http.StatusInternalServerError)
		return
	}
	a.sendData(w, result, okStatus)
}

//getUserEntity is a helper method to retrieve a User from the route parameter.
//If false is returned that indicates a User could not be retrieved because of
//an actual failure or because the user does not exist.
//The appropraite response is sent if false is returned.
func (a *API) getUserEntity(w http.ResponseWriter, r *http.Request) (*user.User, bool) {
	id, ok := a.idFrom(w, r, routeParamUserId)
	if !ok {
		return nil, ok
	}

	user, err := a.Users.Get(r.Context(), id)
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError)
		return nil, false
	}

	return user, true
}
