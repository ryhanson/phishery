package neatprint

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
)

type NeatPrint struct {
	yellow	func(string, ...interface{}) string
	green	func(string, ...interface{}) string
	cyan	func(string, ...interface{}) string
	blue	func(string, ...interface{}) string
	red	func(string, ...interface{}) string
}

func NewNeatPrint() NeatPrint {
	if runtime.GOOS == "windows" {
		color.NoColor = true
	}

	return NeatPrint{
		yellow: color.New(color.FgHiYellow).SprintfFunc(),
		green: color.New(color.FgHiGreen).SprintfFunc(),
		cyan: color.New(color.FgHiCyan).SprintfFunc(),
		blue: color.New(color.FgHiBlue).SprintfFunc(),
		red: color.New(color.FgHiRed).SprintfFunc(),
	}
}

func (np *NeatPrint) Data(pre string, label string, value string) {
	fmt.Printf("%s %-11s: %s\n", np.blue("[%s]", pre), label, np.cyan(value))
}

func (np *NeatPrint) Info(format string, info ...interface{}) {
	format = fmt.Sprintf(np.yellow("[*]") + " %s\n", format)
	fmt.Printf(format, info...)
}

func (np *NeatPrint) Event(format string, event ...interface{}) {
	format = fmt.Sprintf(np.green("[+]") + " %s\n", format)
	fmt.Printf(format, event...)
}

func (np *NeatPrint) Error(format string, err ...interface{}) {
	format = fmt.Sprintf(np.red("[!]") + " %s\n", format)
	fmt.Printf(format, err...)
}