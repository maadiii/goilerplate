package cli

import (
	"fmt"
	"goldfish/app"
	"goldfish/controllers"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/spf13/cobra"
)

var cryptoCli = &cobra.Command{
	Use:   CRYPTO_USE,
	Short: CRYPTO_SHORT,
}

var generateSymmetricKeyCli = &cobra.Command{
	Use:   SYMMETRIC_USE,
	Short: SYMMETRIC_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		len_, err := cli.Flags().GetInt(LENGTH_FLAG)
		if err != nil {
			return err
		}

		key := app.GenerateSymmetricKey(len_)
		fmt.Println(string(key))

		return nil
	},
}

var createJWTCli = &cobra.Command{
	Use:   JWT_USE,
	Short: JWT_SHORT,
	RunE: func(cli *cobra.Command, args []string) error {
		a, err := app.New()
		if err != nil {
			return err
		}
		defer a.Close()
		maxage, err := cli.Flags().GetInt(SECONDS_FLAG)
		if err != nil {
			return err
		}
		a.Config.JWT.MaxAge = uint(maxage)
		token, err := a.CreateJWT(uuid.New(), EMPTY, EMPTY, false, EMPTY)
		if err != nil {
			return err
		}

		se := securecookie.New(a.Config.SecretKey, a.Config.BlockSecretKey)
		encoded, err := se.Encode(controllers.ACCESS_COOKIE, token)
		if err != nil {
			return err
		}

		fmt.Println(encoded)

		return nil
	},
}

func init() {
	generateSymmetricKeyCli.PersistentFlags().
		IntP(LENGTH_FLAG, LENGTH_FLAG_SHORT, 16, LENGTH_FLAG_MESSAGE)
	createJWTCli.PersistentFlags().
		IntP(SECONDS_FLAG, SECONDS_FLAG_SHORT, 86400, SECONDS_FLAG_MESSAGE)
	cryptoCli.AddCommand(generateSymmetricKeyCli, createJWTCli)
	rootCli.AddCommand(cryptoCli)
}
