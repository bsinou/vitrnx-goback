package main

import (

	// initialise gorm and mongodb connections
	_ "github.com/bsinou/vitrnx-goback/gorm"
	_ "github.com/bsinou/vitrnx-goback/mongodb"

	"github.com/bsinou/vitrnx-goback/cmd"
)

func main() {
	// ts := conf.BuildTimestamp
	// if ts == "" {
	// 	ts = "Now"
	// }
	// fmt.Printf("Vitrnx Gobackend %s built %s\n ==> Starting in %s mode...\n\n", conf.VitrnxVersion, ts, conf.Env)
	// // start gin router
	// route.StartRouter()
	cmd.Execute()
}
