package usercmd

import (
	"context"
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
	"github.com/gogolfing/cbus"
)

//Handler is a type that knows how to handle all related User commands for the
//command bus in the application.
type Handler struct {
	//Users is the repo to use to manager Users.
	Users user.Repo
}

//CreateUser attempts to create a new User and add it to h.Users.
//Cmd must be of type *CreateUserCommand.
//The result, if not nil and without error, is a *user.User.
func (h *Handler) CreateUser(ctx context.Context, cmd cbus.Command) (interface{}, error) {
	createUser := cmd.(*CreateUserCommand)

	user, err := createUser.newUser()
	if err != nil {
		return nil, err
	}

	err = h.Users.Add(ctx, *user)

	return user, err
}

//CreateUserCommand is a Command that should be used to create a new User.
type CreateUserCommand struct {
	//Email is the new User's Email field.
	Email string

	//Enabled is the new User's Enabeld field.
	Enabled bool

	//ClientIds is the new User's ClientIds field.
	ClientIds []data.Id
}

func (c *CreateUserCommand) newUser() (*user.User, error) {
	id, err := data.NewId()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &user.User{
		Id:        id,
		Email:     c.Email,
		CreatedAt: now,
		UpdatedAt: now,
		Enabled:   c.Enabled,
		ClientIds: c.ClientIds,
	}, nil
}

//UpdateUser attempts to update a User from cmd.
//Cmd must be of type *UpdateUserCommand.
//The result, if not nil and without error, is a *user.User.
func (h *Handler) UpdateUser(ctx context.Context, cmd cbus.Command) (interface{}, error) {
	updateUser := cmd.(*UpdateUserCommand)

	user, err := h.Users.Get(ctx, updateUser.Id)
	if err != nil {
		return nil, err
	}

	updateUser.updateUser(user)

	err = h.Users.Set(ctx, *user)

	return user, err
}

//UpdateUserCommand is a Command that should be used to update single User
//entities.
type UpdateUserCommand struct {
	//Id is the Id of the User to update.
	Id data.Id

	//Email, if not nil, is the email to set on the User.
	Email *string

	//Enabled, if not nil, is the enabled value to set on the User.
	Enabled *bool

	//ClientIds, if not nil, is the slice of client ids to set on the User.
	ClientIds *[]data.Id
}

func (c *UpdateUserCommand) updateUser(user *user.User) {
	//Do stuff for all User updates.
	user.UpdatedAt = time.Now()

	//Do stuff specific to this instance of c.
	if c.Email != nil {
		user.Email = *c.Email
	}
	if c.Enabled != nil {
		user.Enabled = *c.Enabled
	}
	if c.ClientIds != nil {
		user.ClientIds = *c.ClientIds
	}
}
