package dto

import (
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
)

//Users transforms a []*user.User to a []*UserOutput.
//It delegates to User.
func Users(v interface{}) interface{} {
	users := v.([]*user.User)

	result := make([]interface{}, len(users))
	for i, u := range users {
		result[i] = User(u)
	}

	return result
}

//User transforms a *user.User to a *UserOutput.
func User(v interface{}) interface{} {
	user := v.(*user.User)

	clientIds := user.ClientIds
	if clientIds == nil {
		clientIds = []data.Id{}
	}

	return &UserOutput{
		Id:        user.Id,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Enabled:   user.Enabled,
		ClientIds: clientIds,
	}
}

//UserOutput is a marshalable type that should be used to represent a User
//outside of the API.
type UserOutput struct {
	Id        data.Id   `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Enabled   bool      `json:"enabled"`
	ClientIds []data.Id `json:"client_ids"`
}

//CreateUser is a form type that should be used for incoming create user requests
//to the API.
type CreateUser struct {
	Email     string    `json:"email" valid:"required,email"`
	Enabled   bool      `json:"enabled"`
	ClientIds []data.Id `json"client_ids"`
}

//UpdateUser is a form type that should be used for incoming udpate user requests
//to the API.
type UpdateUser struct {
	Email     *string    `json:"email" valid:"omitempty,email"`
	Enabled   *bool      `json:"enabled"`
	ClientIds *[]data.Id `json:"client_ids"`
}
