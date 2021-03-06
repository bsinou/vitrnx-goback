// Package cmd enables simple management of the backend via CLI
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bsinou/vitrnx-goback/conf"
)

var ()

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "vitrnx-goback",
	Short: "Simple Go Backend for the vitrnx Project",
	Long: `
`,
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

		if args[0] == "" {
			cmd.Println("no instance name provided, cannot launch the vitrnx Backend")
			os.Exit(1)
		}
		conf.VitrnxInstanceID = args[0]

		// Override environment type if provided
		if viper.GetString(conf.KeyEnvType) != "" {
			// TODO check if env value is valid
			conf.Env = viper.GetString(conf.KeyEnvType)
		}

		// Load  configuration
		viper.SetConfigName(conf.BaseName) // name of config file (without extension)
		for _, path := range conf.GetKnownConfFolderPaths() {
			fmt.Println("Adding config path: " + path)
			viper.AddConfigPath(path)
		}
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file
			panic(fmt.Errorf("fatal error config file: %s", err))
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	viper.SetEnvPrefix("vitrnx")
	viper.AutomaticEnv()

	flags := rootCmd.PersistentFlags()
	flags.StringP(conf.KeyEnvType, "e", "", "Override default environment mode defined by the build. Valid values are: dev, test, staging, prod")
	viper.BindPFlag(conf.KeyEnvType, rootCmd.PersistentFlags().Lookup(conf.KeyEnvType))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
