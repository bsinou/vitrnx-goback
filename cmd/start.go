package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bsinou/vitrnx-goback/auth"
	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/mongodb"
	"github.com/bsinou/vitrnx-goback/route"
)

var (
// wg sync.WaitGroup
)

// StartCmd launches the vitrnx backend process.
var StartCmd = &cobra.Command{
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

		// TODO Implement a better way to initialise services and manage clean shutdown
		gorm.InitGormRepo()
		// TODO enhance: launch a sync with firebase on each startup
		// it is not too gravious this is reintrant
		auth.ListExistingUsers(nil)

		mongodb.InitMongoConnection()
		// start gin router
		route.StartRouter()

		// wg.Add(1)
		// wg.Wait()
	},
}

func init() {
	RootCmd.AddCommand(StartCmd)
}
