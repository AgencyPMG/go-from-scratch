package dto

import (
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
	"github.com/AgencyPMG/go-from-scratch/app/internal/data/client"
)

//Clients transforms a []*client.Client to a []*ClientOutput.
//It delegates to Client.
func Clients(v interface{}) interface{} {
	clients := v.([]*client.Client)

	result := make([]interface{}, len(clients))
	for i, c := range clients {
		result[i] = Client(c)
	}

	return result
}

//Create transforms a single *client.Client to a *ClientOutput.
func Client(v interface{}) interface{} {
	client := v.(*client.Client)

	return &ClientOutput{
		Id:        client.Id,
		Name:      client.Name,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}
}

//ClientOutput is a marshalable type that should be used for external
//representations of client.Client(s) outside the API.
type ClientOutput struct {
	Id        data.Id   `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//CreateClient is a form type that should be used for incoming create client
//requests to the API.
type CreateClient struct {
	Name string `json:"name" valid:"required"`
}

//UpdateClient is a form type that should be used for incoming update client
//requests to the API.
type UpdateClient struct {
	Name *string `json:"name" valid:"omitempty,required"`
}
