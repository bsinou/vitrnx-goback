package main

import (

	// Initialise gorm connections.
	_ "github.com/bsinou/vitrnx-goback/gorm"

	"github.com/bsinou/vitrnx-goback/cmd"
)

func main() {
	// Bootstrap Cobra13 framework.
	cmd.Execute()
}
