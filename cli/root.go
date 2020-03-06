package cli

import (
	"goilerplate/app"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/sys/unix"
)

var rootCli = &cobra.Command{
	Use:   app.NAME,
	Short: app.DESCRIPTION,
	PersistentPreRun: func(cli *cobra.Command, args []string) {
		if !terminal.IsTerminal(unix.Stdout) {
			logrus.SetFormatter(&logrus.JSONFormatter{})
		} else {
			logrus.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: time.RFC3339Nano,
			})
		}

		if verbose, _ := cli.Flags().GetBool(VERBOSE_FLAG); verbose {
			logrus.SetLevel(logrus.DebugLevel)
		}
	},
}

func Execute() {
	_ = rootCli.Execute()
}

var configFile string

func init() {
	rootCli.PersistentFlags().BoolP(
		VERBOSE_FLAG,
		VERBOSE_FLAG_SHORT,
		false,
		VERBOSE_FLAG_MESSAGE,
	)
	rootCli.PersistentFlags().StringVarP(
		&configFile,
		CONFIG_FLAG,
		CONFIG_FLAG_SHORT,
		EMPTY,
		CONFIG_FLAG_USAGE,
	)
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigFile(configFile)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warning(
			"Config file is not set, you're using internal " +
				"config that use for debuging.",
		)
	}
}
