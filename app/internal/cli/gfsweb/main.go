package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AgencyPMG/go-from-scratch/app/internal/cli/gfsweb/cli"
)

const (
	//CommandName is used to send to our cli package.
	CommandName = "gfsweb"
)

func main() {
	//Create a Context for the cli package (and rest of the application) to
	//operate within.
	//Canceling it should stop termination anywhere in the program.
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)

	//We want to listen to the Interrupt and Terminate signal.
	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	//Wait until we receive a signal.
	//That is the only thing we care about to stop the application.
	go func() {
		defer cancel()
		<-signals
	}()

	//Get an exit code from our cli package and exit with it.
	//This returns after ctx is done.
	exitCode := cli.Main(
		ctx,
		os.Args[1:],
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)

	os.Exit(exitCode)
}
