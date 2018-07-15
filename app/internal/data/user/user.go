package user

import (
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
)

//User is a domain type that represents a user of the application.
type User struct {
	//Id is the User's id.
	Id data.Id

	//Email is the User's email.
	//It must be unique across all users of the application.
	Email string

	//CreatedAt is the time at which this User was first created.
	CreatedAt time.Time

	//UpdatedAt is the most recent time this User was updated.
	UpdatedAt time.Time

	//Enabled is aboolean value indicating whether or not this User is enabled
	//and is allowed to use the application.
	Enabled bool

	//ClientIds is a slice of Ids that represent all of the Clients that this User
	//has access to.
	ClientIds []data.Id
}

//QueryRepo provides methods for retrieving Users.
type QueryRepo interface {
	//Get should return the User whose Id equals id.
	//
	//An error should be returned if the User does not exist or there was an error
	//attempting to load the User.
	Get(id data.Id) (*User, error)

	//List should return all Users in the application.
	//They should be sorted by their lexicographic Email order.
	List() ([]*User, error)
}

//Repo provides method for creating and updating Users
//as well as promotes the QueryRepo interface.
type Repo interface {
	//QueryRepo is promoted here to indicate a Repo contains all query methods.
	QueryRepo

	//Add should add u to the underlying storage repository.
	Add(u User) error

	//Set should update all stored fields of u in the underlying storage repository.
	//The update should use u.Id for determining which entity to update.
	Set(u User) error
}
