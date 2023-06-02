package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/AlexandreMarcq/gozimbra/cmd/account"
	"github.com/AlexandreMarcq/gozimbra/internal/cmd_utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	config        string
	logDir        string
	outputFile    string
	defaultOutput string = fmt.Sprintf("./out/gozimbra-%s.csv", cmd_utils.GetFormattedTime())
	platform      string
	noLog         bool
	noUi          bool
	stdout        bool
)

func Execute() error {
	return NewRootCmd().Execute()
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gozimbra",
		Short: "CLI tool for Zimbra administration",
	}

	cobra.OnInitialize(initConfig)

	cmd.PersistentFlags().StringVarP(&config, "config", "c", "./gozimbra.yaml", "configuration file")
	cmd.PersistentFlags().StringVarP(&logDir, "log", "l", "./log", "log directory")
	cmd.PersistentFlags().StringVarP(&outputFile, "output", "o", defaultOutput, "output file")
	cmd.PersistentFlags().StringVarP(&platform, "platform", "p", "", "zimbra platform")
	cmd.PersistentFlags().BoolVar(&noLog, "no-log", false, "disables logs")
	cmd.PersistentFlags().BoolVar(&noUi, "no-ui", false, "disables UI")
	cmd.PersistentFlags().BoolVar(&stdout, "stdout", false, "outputs to stdout")

	cmd.MarkFlagFilename("config")

	viper.BindPFlag("defaults.log", cmd.PersistentFlags().Lookup("log"))
	viper.BindPFlag("defaults.no-log", cmd.PersistentFlags().Lookup("no-log"))
	viper.BindPFlag("defaults.no-ui", cmd.PersistentFlags().Lookup("no-ui"))
	viper.BindPFlag("defaults.platform", cmd.PersistentFlags().Lookup("platform"))
	viper.BindPFlag("defaults.stdout", cmd.PersistentFlags().Lookup("stdout"))

	cmd.AddCommand(account.NewAccountCmd())

	return cmd
}

func initConfig() {
	if config != "" {
		viper.SetConfigFile(config)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("gozimbra.yaml")
	}

	err := viper.ReadInConfig()
	if err != nil {
		cobra.CheckErr(err)
	}

	var logFile io.Writer
	if viper.GetBool("defaults.no-log") {
		logFile = io.Discard
	} else {
		logInfo, err := os.Stat(logDir)
		if err != nil {
			cobra.CheckErr(err)
		}
		if !logInfo.IsDir() {
			cobra.CheckErr("log flag must be given a directory")
		}

		logFile, err = os.Create(fmt.Sprintf("%s/gozimbra-%s.log", viper.GetString("defaults.log"), cmd_utils.GetFormattedTime()))
		if err != nil {
			cobra.CheckErr(err)
		}
	}

	log.SetOutput(logFile)
}
