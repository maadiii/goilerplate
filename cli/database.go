package cli

import (
	"goilerplate/infrastructure/application"

	"github.com/spf13/cobra"
)

var dbCli = &cobra.Command{
	Use:   DB_USE,
	Short: DB_SHORT,
}

var migrateCli = &cobra.Command{
	Use:   MIGRATE_USE,
	Short: MIGRATE_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		app, err := application.New()
		if err != nil {
			return err
		}
		defer app.Close()

		app.MigrateDB()

		return nil
	},
}

var dropCli = &cobra.Command{
	Use:   DROP_USE,
	Short: DROP_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		app, err := application.New()
		if err != nil {
			return err
		}
		defer app.Close()

		app.DropDB()

		return nil
	},
}

var basedataCli = &cobra.Command{
	Use:   BASEDATA_USE,
	Short: BASEDATA_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		app, err := application.New()
		if err != nil {
			return err
		}
		defer app.Close()

		app.InsertBaseData()

		return nil
	},
}

func init() {
	dbCli.AddCommand(migrateCli, basedataCli, dropCli)
	rootCli.AddCommand(dbCli)
}
