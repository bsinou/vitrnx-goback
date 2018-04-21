package main

import (
	"fmt"

	// initialise gorm and mongodb connections
	_ "github.com/bsinou/vitrnx-goback/gorm"
	_ "github.com/bsinou/vitrnx-goback/mongodb"

	// start gin router
	_ "github.com/bsinou/vitrnx-goback/route"
)

func main() {
	fmt.Println("Vitrnx 0.2 - Go Backend starting... ")
}
