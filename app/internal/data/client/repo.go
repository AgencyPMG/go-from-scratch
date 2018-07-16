package client

import (
	"context"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
)

//QueryRepo provides the methods for retrieving Clients.
type QueryRepo interface {
	//Get should return the Client whose Id equals id.
	Get(ctx context.Context, id data.Id) (*Client, error)

	//List should return all Clients in the repository.
	List(ctx context.Context) ([]*Client, error)

	//ListById should return all Clients whose Id is in ids.
	ListByIds(ctx context.Context, ids []data.Id) ([]*Client, error)
}

//Repo provides methods for retrieving and manipulating Clients.
type Repo interface {
	//QueryRepo is promoted to include the retrival methods.
	QueryRepo

	//Add should add c to the underlying repo's storage.
	Add(ctx context.Context, c Client) error

	//Set should udpate all stored fields of c in the underlying storage repository.
	//The update should use c.Id to determine with entity to update.
	Set(ctx context.Context, c Client) error

	//Remove should remove the Client with id from the underlying storage.
	Remove(ctx context.Context, id data.Id) error
}
