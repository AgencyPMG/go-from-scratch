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
To build the executable.
```bash
./bin/dev/build gfswebb
# OR
go build ./app/internal/cli/gfswebb
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
