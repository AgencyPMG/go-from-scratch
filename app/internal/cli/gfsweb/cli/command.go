package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/AgencyPMG/go-from-scratch/app/internal/gfsweb"
)

func Main(ctx context.Context, _ []string, _ io.Reader, _, outErr io.Writer) int {
	err := main(ctx)

	if err != nil {
		fmt.Fprintln(outErr, err)
		return 1
	}

	return 0
}

func main(ctx context.Context) error {
	c, err := newConfig()
	if err != nil {
		return err
	}

	// c.Values().EachKeyValue(func(key config.Key, value interface{}) {
	// 	fmt.Println(key, value)
	// })

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
