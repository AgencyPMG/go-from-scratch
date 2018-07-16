package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/AgencyPMG/go-from-scratch/app/internal/cli/gfsweb/cli"
)

const (
	CommandName = "gfsweb"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)

	signal.Notify(
		signals,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		defer cancel()
		<-signals
	}()

	exitCode := cli.Main(
		ctx,
		os.Args[1:],
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)

	os.Exit(exitCode)
}
