package gp

import (
	"fmt"
	"github.com/fatih/color"
	"runtime"
)

type NeatPrint struct {}

func (np *NeatPrint) Info(info string) {
	if runtime.GOOS == "windows" {
		fmt.Printf("[*] %s\n", info)
	} else {
		yellow := color.New(color.FgHiYellow).SprintFunc()
		fmt.Printf("%s %s\n", yellow("[*]"), info)
	}

}

func (np *NeatPrint) Event(event string) {
	if runtime.GOOS == "windows" {
		fmt.Printf("[+] %s\n", event)
	} else {
		green := color.New(color.FgHiGreen).SprintFunc()
		fmt.Printf("%s %s\n", green("[+]"), event)
	}
}

func (np *NeatPrint) Data(pre string, label string, value string) {
	if runtime.GOOS == "windows" {
		fmt.Printf("%s %-11s : %s\n", pre, label, value)
	} else {
		cyan := color.New(color.FgCyan).SprintFunc()
		blue := color.New(color.FgBlue).SprintFunc()
		fmt.Printf("%s %-11s : %s\n", blue(pre), label, cyan(value))
	}
}

func (np *NeatPrint) Error(error string) {
	if runtime.GOOS == "windows" {
		fmt.Printf("[!] %s\n", error)
	} else {
		red := color.New(color.FgHiRed).SprintFunc()
		fmt.Printf("%s %s\n", red("[!]"), error)
	}
}