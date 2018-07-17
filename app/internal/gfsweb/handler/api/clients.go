package api

import (
	"net/http"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client/clientcmd"
	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb/handler/api/dto"
)

const (
	routeParamClientId = "client_id"
)

//getClient retrieves and sends a single Client.
func (a *API) getClient(w http.ResponseWriter, r *http.Request) {
	client, ok := a.getClientEntity(w, r)
	if !ok {
		return
	}

	a.sendData(w, client, http.StatusOK)
}

//listClients retrieves and sends all Clients.
func (a *API) listClients(w http.ResponseWriter, r *http.Request) {
	clients, err := a.Clients.List(r.Context())
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError)
		return
	}

	a.sendData(w, clients, http.StatusOK)
}

//createClient attempts to create a new Client from a dto.CreateClient form
//and executing a clientcmd.CreateClientCommand.
func (a *API) createClient(w http.ResponseWriter, r *http.Request) {
	f := &dto.CreateClient{}
	if ok := a.parseForm(w, r, f); !ok {
		return
	}

	command := &clientcmd.CreateClientCommand{
		Name: f.Name,
	}

	client, err := a.Bus.ExecuteContext(r.Context(), command)
	a.clientCommandResponse(w, client, err, http.StatusCreated)
}

//updateClient attempts to update an existing client from a dto.UpdateClient form
//and executing a clientcmd.UpdateClientCommand.
func (a *API) updateClient(w http.ResponseWriter, r *http.Request) {
	clientId, ok := a.idFrom(w, r, routeParamClientId)
	if !ok {
		return
	}

	f := &dto.UpdateClient{}
	if ok := a.parseForm(w, r, f); !ok {
		return
	}

	command := &clientcmd.UpdateClientCommand{
		Id:   clientId,
		Name: f.Name,
	}

	client, err := a.Bus.ExecuteContext(r.Context(), command)
	a.clientCommandResponse(w, client, err, http.StatusOK)
}

//deleteClient attempts to delete an existing Client the route parameter id
//and executing a clientcmd.DeleteClientCommand.
func (a *API) deleteClient(w http.ResponseWriter, r *http.Request) {
	//We still want to retrive the client in order to make sure it exists.
	//If this operation fails, it would ideally result in a 404 or special entity
	//not found error messaging.
	//
	//The delete command may not necessarily be able to indicate there was no
	//entity to delete.
	client, ok := a.getClientEntity(w, r)
	if !ok {
		return
	}

	command := &clientcmd.DeleteClientCommand{
		Id: client.Id,
	}

	_, err := a.Bus.ExecuteContext(r.Context(), command)
	a.clientCommandResponse(w, client, err, http.StatusOK)
}

//getClientEntity is a helper method to retrieve a Client from a.Clients.
//If false is returned, that indicates the Client could not be retrieved -
//through an actual failure or because it does not exist.
//A response is sent if false is returned.
func (a *API) getClientEntity(w http.ResponseWriter, r *http.Request) (*client.Client, bool) {
	clientId, ok := a.idFrom(w, r, routeParamClientId)
	if !ok {
		return nil, ok
	}

	client, err := a.Clients.Get(r.Context(), clientId)
	if err != nil {
		//Again, we can examine the error the see what type of status to send.
		a.sendError(w, err, http.StatusInternalServerError)
		return nil, false
	}

	return client, true
}

//clientCommandResponse is a helper method to send the correct response from the
//result of executing a client command.
func (a *API) clientCommandResponse(w http.ResponseWriter, result interface{}, err error, okStatus int) {
	if err != nil {
		a.sendError(w, err, http.StatusInternalServerError)
		return
	}

	a.sendData(w, result, okStatus)
}
