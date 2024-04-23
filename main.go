package main

import (
	"newLogger/utils"
)

func main() {
	utils.ParseArgs()
	utils.ParseConfigurations()

	rf := utils.NewFileReader()
	rf.Paint()
}
