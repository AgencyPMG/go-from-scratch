package clientcmd

import (
	"context"
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
	"github.com/gogolfing/cbus"
)

//Handler is a handler type that understands how to work with Client commands.
type Handler struct {
	//Clients is the Client repository to use within Command handling.
	Clients client.Repo
}

//CreateClient attempts to create a new Client and add it to h.Clients.
//Cmd must be of type *CreateClientCommand.
//The result, if not nil and without error, will be a *client.Client.
func (h *Handler) CreateClient(ctx context.Context, cmd cbus.Command) (interface{}, error) {
	createClient := cmd.(*CreateClientCommand)

	client, err := createClient.newClient()
	if err != nil {
		return nil, err
	}

	err = h.Clients.Add(ctx, *client)

	return client, err
}

//CreateClientCommand is Command that should be used to create a new Client and
//add it to a repository.
type CreateClientCommand struct {
	//Name is the name of the Client to create.
	Name string
}

func (c *CreateClientCommand) newClient() (*client.Client, error) {
	id, err := data.NewId()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &client.Client{
		Id:        id,
		Name:      c.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

//UpdateClient attempts to retrieve and set a Client using h.Clients.
//Cmd must be of type *UpdateClientCommand.
//The result, if not nil and without error, will be a *client.Client.
func (h *Handler) UpdateClient(ctx context.Context, cmd cbus.Command) (interface{}, error) {
	updateClient := cmd.(*UpdateClientCommand)

	client, err := h.Clients.Get(ctx, updateClient.Id)
	if err != nil {
		return nil, err
	}

	updateClient.updateClient(client)

	err = h.Clients.Set(ctx, *client)

	return client, err
}

//UpdateClientCommand is a Command to update a Client.
type UpdateClientCommand struct {
	//Id is the Id of the Client to update.
	Id data.Id

	//Name, if not nil, is the name to set on the Client.
	Name *string
}

func (c *UpdateClientCommand) updateClient(client *client.Client) {
	//Do stuff for all User updates.
	client.UpdatedAt = time.Now()

	//Do stuff specific to this instance of c.
	if c.Name != nil {
		client.Name = *c.Name
	}
}

//DeleteClient attempts to remove a Client from h.Clients.
//Cmd must be of type *DeleteClientCommand.
//The result will always be nil.
//A non-nil error will be returned to indicate failure.
func (h *Handler) DeleteClient(ctx context.Context, cmd cbus.Command) (interface{}, error) {
	deleteClient := cmd.(*DeleteClientCommand)

	return nil, h.Clients.Remove(ctx, deleteClient.Id)
}

//DeleteClientCommand is a Command to delete (remove) a Client.
type DeleteClientCommand struct {
	//Id is the Id of the Client to delete.
	Id data.Id
}
