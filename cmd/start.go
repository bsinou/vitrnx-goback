package cmd

import (
	"fmt"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/gorm"
	"github.com/bsinou/vitrnx-goback/mongodb"
	"github.com/bsinou/vitrnx-goback/route"
	"github.com/spf13/cobra"
)

var (
// wg sync.WaitGroup
)

// StartCmd represents the start command
var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the VitrnX Backend",
	Long:  ``,

	PreRun: func(cmd *cobra.Command, args []string) {
	},

	Run: func(cmd *cobra.Command, args []string) {

		// Real start of the backend. Should be enhanced
		ts := conf.BuildTimestamp
		if ts == "" {
			ts = "Now"
		}

		cmd.Print(fmt.Sprintf("Vitrnx Go Backend v%s for [%s] built on %s\n ==> Starting in %s mode...\n\n", conf.VitrnxInstanceID, conf.VitrnxVersion, ts, conf.Env))

		// TODO Implement a better way to initialise services and manage clean shutdown
		gorm.InitGormRepo()
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