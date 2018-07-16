package client

import (
	"time"

	"github.com/AgencyPMG/go-from-scratch/app/internal/data"
)

//Client is a domain type that represents a client inside the application.
//
//This could be thought of as a business project or business unit.
//In out line of work, they are called clients.
//
//Different Users are allowed to access different Clients.
//Furthermore, there would normally be "lower level" domain types that belong
//to a Client.
type Client struct {
	//Id is the Client's id.
	Id data.Id

	//Name is the Client's name.
	//
	//Name must be unique across all Clients in the application.
	Name string

	//CreatedAt is the time at which this Client was created.
	CreatedAt time.Time

	//UpdatedAt is the time at which this Client was most recently updated.
	UpdatedAt time.Time
}
