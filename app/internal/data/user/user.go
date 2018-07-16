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
