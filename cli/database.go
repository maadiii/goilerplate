package cli

import (
	"goldfish/app"
	"goldfish/domain/services"

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
		app, err := app.New()
		if err != nil {
			return err
		}
		dbs := app.DBSession
		defer app.Close()

		//usecases.MigrateDB(dbs)
		services.MigrateDB(dbs)

		return nil
	},
}

var dropCli = &cobra.Command{
	Use:   DROP_USE,
	Short: DROP_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		app, err := app.New()
		if err != nil {
			return err
		}
		dbs := app.DBSession
		defer app.Close()

		//usecases.DropDB(dbs)
		services.DropDB(dbs)

		return nil
	},
}

var basedataCli = &cobra.Command{
	Use:   BASEDATA_USE,
	Short: BASEDATA_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		app, err := app.New()
		if err != nil {
			return err
		}
		dbs := app.DBSession
		defer app.Close()

		//usecases.InsertBaseData(dbs)
		services.InsertBaseData(dbs)
		return nil
	},
}

func init() {
	dbCli.AddCommand(migrateCli, basedataCli, dropCli)
	rootCli.AddCommand(dbCli)
}
