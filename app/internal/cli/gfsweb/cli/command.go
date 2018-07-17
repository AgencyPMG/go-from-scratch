package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb"
)

//Main is the the entrypoint client code can use to invoke gfsweb as if
//it were entered through the command line.
func Main(ctx context.Context, _ []string, _ io.Reader, _, outErr io.Writer) int {
	err := main(ctx)

	if err != nil {
		fmt.Fprintln(outErr, err)
		return 1
	}

	return 0
}

//main buils and runs a gfsweb.App with ctx.
//An error is returned either from building or running the application.
func main(ctx context.Context) error {
	c, err := newConfig()
	if err != nil {
		return err
	}

	ab := gfsweb.NewAppBuilder()

	app, err := ab.Build(c)
	if err != nil {
		return err
	}
	defer app.Close()

	err = app.Run(ctx)

	if err != nil {
		return err
	}

	return nil
}
