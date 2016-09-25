package gp

const VERSION = "1.0"

type AuthInfo struct {
	Host	  string `json:"host"`
	Request   string `json:"request"`
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

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