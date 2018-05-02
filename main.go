package main

import (
	"fmt"

	"github.com/bsinou/vitrnx-goback/conf"
	"github.com/bsinou/vitrnx-goback/route"

	// initialise gorm and mongodb connections
	_ "github.com/bsinou/vitrnx-goback/gorm"
	_ "github.com/bsinou/vitrnx-goback/mongodb"
)

func main() {
	ts := conf.BuildTimestamp
	if ts == "" {
		ts = "Now"
	}
	fmt.Printf("Vitrnx Gobackend %s built %s\n ==> Starting in %s mode...\n\n", conf.VitrnxVersion, ts, conf.Env)
	// start gin router
	route.StartRouter()
}
