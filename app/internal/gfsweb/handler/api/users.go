package api

import (
	"net/http"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user/usercmd"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api/dto"
)

func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	user, ok := a.getUserEntity(w, r)
	if !ok {
		return
	}

	a.sendData(w, user, http.StatusOK)
}

func (a *API) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.Users.List(r.Context())
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError)
		return
	}

	a.sendData(w, users, http.StatusOK)
}

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
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError) //TODO note here for correct status code.
		return
	}

	a.sendData(w, user, http.StatusCreated)
}

func (a *API) updateUser(w http.ResponseWriter, r *http.Request) {
	f := &dto.UpdateUser{}

	id, _ := a.idFrom(w, r, "user_id")

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
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError) //TODO note here for correct status code.
		return
	}

	a.sendData(w, user, http.StatusCreated)
}

func (a *API) getUserEntity(w http.ResponseWriter, r *http.Request) (*user.User, bool) {
	id, ok := a.idFrom(w, r, "user_id")
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
