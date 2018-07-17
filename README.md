# go-from-scratch
A public working example of how to build a Golang web application.

## Installation
We are assuming you already have Golang isntalled.

Clone this repo to `$GOPATH/src/github.com/AgencyPMG/go-from-scratch`.

Docker Compose is required for backing services.
See [here](https://docs.docker.com/compose/install/) for installation instructions.

We use [Glide](https://glide.sh) to manage our Golang dependencies.

We use [DBSchema](https://github.com/gogolfing/dbschema) to manage our database schema.

## Running the App
There is a helper script `./bin/dev/build` to build executables in the
`app/internal/cli` directory.
To build the executable:
```bash
./bin/dev/build gfsweb
# OR
go build ./app/internal/cli/gfsweb
```

Next, bring up the backing services:
```bash
./bin/dev/start
```
to get the Postgres container running through Docker Compose.

Next, source the environment needed for configuration:
```bash
source ./etc/config/.env
```

Finally, you can run the executable with
```bash
./gfsweb
```

## Using the App
Once the executable `gfsweb` is running, you can interact with it via curl or other
web development utilities.

For instance, use `curl localhost:8080/users` to get a list of all Users in the
application.

This is meant to be an exploratory application for example purposes.
Go ahead and look into the code and play around with it to figure out how to create,
edit, and delete other entities.
