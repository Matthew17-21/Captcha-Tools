package captchatoolsgo

import (
	"fmt"
	"strings"
)

type Proxy struct {
	Ip       string
	Port     string
	User     string
	Password string
}

// NewProxy returns a proxy that can be used to solve captchas
func NewProxy(proxy string) (*Proxy, error) {
	if strings.TrimSpace(proxy) == "" {
		return nil, ErrProxyEmpty
	}

	splitted := strings.Split(proxy, ":")
	pLen := len(splitted)
	p := &Proxy{}
	if pLen >= 2 {
		p.Ip = splitted[0]
		p.Port = splitted[1]
	}
	if pLen >= 4 {
		p.User = splitted[2]
		p.Password = splitted[3]
	}
	return p, nil
}

// IsUserAuth returns if a proxy is user authenticated
func (p Proxy) IsUserAuth() bool {
	return p.User != "" && p.Password != ""
}

// Returns the proxy as a string, unformatted
//
// Example: would return "ip:port" || "ip:port:user:pass"
func (p Proxy) String() string {
	var formatted string = p.Ip + ":" + p.Port
	if p.IsUserAuth() {
		formatted = formatted + ":" + p.User + ":" + p.Password
	}
	return formatted
}

// StringFormatted returns the proxy as a string in the correct format
//
// Example: returns "user:pass@ip:port" || "ip:port"
func (p Proxy) StringFormatted() string {
	var formatted string = p.Ip + ":" + p.Port
	if p.IsUserAuth() {
		formatted = fmt.Sprintf("%v:%v@%v", p.User, p.Password, formatted)
	}
	return formatted
}
