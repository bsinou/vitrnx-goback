package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/route"
)

var (
// wg sync.WaitGroup
)

// startCmd launches the vitrnx backend process.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the vitrnx backend",
	Long:  ``,

	PreRun: func(cmd *cobra.Command, args []string) {
	},

	Run: func(cmd *cobra.Command, args []string) {

		// Real start of the backend. Should be enhanced
		ts := conf.BuildTimestamp
		if ts == "" {
			ts = "just now..."
		} else {
			ts = "on " + ts
		}

		cmd.Print(fmt.Sprintf("\n\n%s - Vitrnx Go Backend v%s (built %s)\n ==> Starting in %s mode.\n\n", conf.VitrnxInstanceID, conf.VitrnxVersion, ts, conf.Env))
		cmd.Print(fmt.Sprintf("Current admin is %s\n", viper.GetString(conf.KeyAdminEmail)))

		// TODO Implement a better way to initialise services and manage clean shutdown
		gorm.InitGormRepo()

		// // TODO enhance: launch a sync with firebase on each startup
		// // it is not too gravious this is reintrant
		// auth.ListExistingUsers(nil)

		// start gin router
		route.StartRouter()

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
