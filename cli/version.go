package cli

import (
	"fmt"
	"goilerplate/app"

	"github.com/spf13/cobra"
)

func init() {
	rootCli.AddCommand(&cobra.Command{
		Use:   VERSION_USE,
		Short: VERSION_SHORT,
		Run: func(cli *cobra.Command, args []string) {
			fmt.Println(app.Name, app.Version)
		},
	})
}
