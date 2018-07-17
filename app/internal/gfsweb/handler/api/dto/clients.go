package dto

import (
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
)

func Clients(v interface{}) interface{} {
	clients := v.([]*client.Client)

	result := make([]interface{}, len(clients))
	for i, c := range clients {
		result[i] = Client(c)
	}

	return result
}

func Client(v interface{}) interface{} {
	client := v.(*client.Client)

	return &ClientOutput{
		Id:        client.Id,
		Name:      client.Name,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}

type ClientOutput struct {
	Id        data.Id   `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateClient struct {
	Name string `json:"name" valid:"required"`
}

type UpdateClient struct {
	Name *string `json:"name" valid:"omitempty,required"`
}
