package twocaptcha

import (
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/proxy"
)

// Payload for 2captcha requests
type payload map[string]any

// newPayload creates and returns a new payload
func newPayload() payload {
	return make(payload)
}

func newPayloadWithClientKey(clientKey string) payload {
	p := newPayload()
	p.setToPayload("clientKey", clientKey)
	return p
}

// SetB64Img sets a base64-encoded image for image captcha solving
func (p *payload) SetB64Img(img string) {
	p.setToPayload("body", img)
}

// SetProxy sets the proxy configuration for solving captchas
// that require specific IP addresses or geo-locations
func (p *payload) SetProxy(pxy *proxy.Proxy) {
	if pxy == nil {
		return
	}
	p.setToPayload("proxyAddress", pxy.Ip)
	p.setToPayload("proxyPort", pxy.Port)
	if pxy.IsUserAuth() {
		p.setToPayload("proxyLogin", pxy.User)
		p.setToPayload("proxyPassword", pxy.Password)
	}
	// TODO: Set proxy type here
}

// SetProxyType sets the type of proxy (HTTP, HTTPS, SOCKS4, SOCKS5)
// to be used for the captcha solving request
func (p *payload) SetProxyType(s string) {
	p.setToPayload("proxyType", s)
}

// SetUserAgent sets the user agent string that will be passed
// to the service and used to solve the captcha
func (p *payload) SetUserAgent(userAgent string) {
	p.setToPayload("userAgent", userAgent)
}

// SetRqData sets custom data for hCaptcha enterprise challenges,
// typically found as rqdata in network requests for invisible captchas
func (p *payload) SetRqData(rq string) {
	p.setToPayload("data", rq)
}

func (p *payload) setToPayload(key string, value any) {
	(*p)[key] = value
}
