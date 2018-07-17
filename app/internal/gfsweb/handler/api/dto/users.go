package dto

import (
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/user"
)

func Users(v interface{}) interface{} {
	users := v.([]*user.User)

	result := make([]interface{}, len(users))
	for i, u := range users {
		result[i] = User(u)
	}

	return result
}

func User(v interface{}) interface{} {
	user := v.(*user.User)

	return &UserOutput{
		Id:        user.Id,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Enabled:   user.Enabled,
	}
}

type UserOutput struct {
	Id        data.Id   `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Enabled   bool      `json:"enabled"`
}

type CreateUser struct {
	Email     string    `json:"email" valid:"required,email"`
	Enabled   bool      `json:"enabled"`
	ClientIds []data.Id `json"client_ids"`
}

type UpdateUser struct {
	Email     *string    `json:"email" valid:"omitempty,email"`
	Enabled   *bool      `json:"enabled"`
	ClientIds *[]data.Id `json:"client_ids"`
}
