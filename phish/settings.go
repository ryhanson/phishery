package phish

import (
	"os"
	"encoding/json"
)

type Settings struct {
	IP 	  	string
	Port 	  	string
	SSLKey	  	string
	SSLCert   	string
	BasicRealm	string
	ResponseFile	string
	ResponseBody	string
	ResponseStatus	int
	ResponseHeaders [][]string
}

func loadSettings(jsonFile string) Settings {
	file, err := os.Open(jsonFile)
	if err != nil {
		neat.Error("Error loading settings: %s", err)
		os.Exit(1)
	}

	settings := Settings{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&settings); err != nil {
		neat.Error("Error decoding settings: %s", err)
		os.Exit(1)
	}

	return settings
}