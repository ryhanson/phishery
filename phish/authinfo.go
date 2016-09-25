package phish

type AuthInfo struct {
	Host	  string `json:"host"`
	Request   string `json:"request"`
	IPAddress string `json:"ipAddress"`
	UserAgent string `json:"userAgent"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}