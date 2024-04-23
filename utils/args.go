package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Args struct {
	Grep        string
	IgnoreGrep  string
	File        string
	IsFile      bool
	LinesOffset int
	Tail        bool
	Pattern     string
}

func isPipe() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeDevice) == 0
}

func newArgs() *Args {
	return &Args{
		LinesOffset: 10,
		Tail:        true,

		Grep:       "",
		IgnoreGrep: "",

		File:   "",
		IsFile: !isPipe(),

		Pattern: "level=\"(.*?)\"",
	}
}

// INFO: Help Manual
func (Args) help() {
	fileName := os.Args[0]
	path := strings.Split(fileName, string(os.PathSeparator))
	if len(path) > 1 {
		fileName = path[len(path)-1]
	}

	fmt.Printf("Usage:\t%s [OPTIONS] [FILE]\n", fileName)

	fmt.Println("OPTIONS: ")
	fmt.Println("\t-g, --grep\t[PATTERN]\t Find Regular expression")
	fmt.Println("\t-i, --ignore\t[PATTERN]\t Ignore Regular expression")
	fmt.Println("\t-h, --help\t\t\t show this manual page")

	fmt.Println("\nExamples:")
	fmt.Printf("\t%s [OPTIONS] [FILE]\n", fileName)
	fmt.Printf("\t%s [OPTIONS] < [FILE]\n", fileName)
	fmt.Printf("\tcat [FILE] | %s [OPTIONS] \n", fileName)
	fmt.Printf("\ttail [FILE] | %s [OPTIONS] \n", fileName)
	fmt.Printf("\ttail -f [FILE] | %s [OPTIONS] \n", fileName)
	fmt.Printf("\ttail -f [FILE] | %s [OPTIONS] \n", fileName)

	fmt.Println("Thank You for Using this app :)")
	fmt.Println("Need New feature, Contact gershm@omc.co.il")

}

func (a Args) logFileExists() bool {
	_, err := os.Stat(a.File)
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

var args *Args

// INFO:  Read Arguments
func ParseArgs() {
	args = newArgs()

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-g" || arg == "--grep" {
			i++
			args.Grep = os.Args[i]
		} else if arg == "-i" || arg == "--ignore" {
			i++
			args.IgnoreGrep = os.Args[i]
		} else if arg == "-p" || arg == "--pattern" {
			i++
			args.Pattern = os.Args[i]
		} else if arg == "-h" || arg == "--help" {
			args.help()
			os.Exit(0)
		}
	}

	if args.IsFile {
		if len(os.Args) <= 1 {
			args.help()
			os.Exit(0)
		}
		args.File = os.Args[len(os.Args)-1]
		if !args.logFileExists() {
			fmt.Println("The File doesn't exists")
			os.Exit(0)
		}
	}
}
