package user

import (
	"context"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
)

//QueryRepo provides methods for retrieving Users.
type QueryRepo interface {
	//Get should return the User whose Id equals id.
	//
	//An error should be returned if the User does not exist or there was an error
	//attempting to load the User.
	Get(ctx context.Context, id data.Id) (*User, error)

	//List should return all Users in the application.
	//They should be sorted by their lexicographic Email order.
	List(ctx context.Context) ([]*User, error)
}

//Repo provides method for creating and updating Users
//as well as promotes the QueryRepo interface.
type Repo interface {
	//QueryRepo is promoted here to indicate a Repo contains all query methods.
	QueryRepo

	//Add should add u to the underlying storage repository.
	Add(ctx context.Context, u User) error

	//Set should update all stored fields of u in the underlying storage repository.
	//The update should use u.Id for determining which entity to update.
	Set(ctx context.Context, u User) error
}
