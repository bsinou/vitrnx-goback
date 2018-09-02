package main

import (

	// initialise gorm connections
	_ "github.com/bsinou/vitrnx-goback/gorm"
	// idem with mongodb
	_ "github.com/bsinou/vitrnx-goback/mongodb"

	"github.com/bsinou/vitrnx-goback/cmd"
)

func main() {
	cmd.Execute()
}
