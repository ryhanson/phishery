package gp

import (
	"fmt"
	"github.com/fatih/color"
)

type NeatPrint struct {}

func (np *NeatPrint) Info(info string) {
	yellow := color.New(color.FgHiYellow).SprintFunc()
	fmt.Printf("%s %s\n", yellow("[*]"), info)
}

func (np *NeatPrint) Event(event string) {
	green := color.New(color.FgHiGreen).SprintFunc()
	fmt.Printf("%s %s\n", green("[+]"), event)
}

func (np *NeatPrint) Data(pre string, label string, value string) {
	cyan := color.New(color.FgCyan).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	fmt.Printf("%s %-11s : %s\n", blue(pre), label, cyan(value))
}

func (np *NeatPrint) Error(error string) {
	red := color.New(color.FgHiRed).SprintFunc()
	fmt.Printf("%s %s\n", red("[!]"), error)
}