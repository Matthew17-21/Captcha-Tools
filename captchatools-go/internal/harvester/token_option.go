package harvester

import "github.com/Matthew17-21/Captcha-Tools/captchatools-go/proxy"

// ConfigurableSetter defines the interface for types that can be configured
// with token options. Implementing this interface allows a type to be used
// with the functional options pattern for token requests.
type ConfigurableSetter interface {
	// SetB64Img sets a base64-encoded image for image captcha solving
	SetB64Img(string)

	// SetProxy sets the proxy configuration for solving captchas
	// that require specific IP addresses or geo-locations
	SetProxy(*proxy.Proxy)

	// SetProxyType sets the type of proxy (HTTP, HTTPS, SOCKS4, SOCKS5)
	// to be used for the captcha solving request
	SetProxyType(string)

	// SetUserAgent sets the user agent string that will be passed
	// to the service and used to solve the captcha
	SetUserAgent(string)

	// SetRqData sets custom data for hCaptcha enterprise challenges,
	// typically found as rqdata in network requests for invisible captchas
	SetRqData(string)
}

// TokenOption is a function that configures a ConfigurableSetter object.
// It implements the functional options pattern for customizing captcha
// token requests with additional parameters.
type TokenOption func(ConfigurableSetter)

// WithB64Img creates an option to set a base64-encoded image for
// solving image-based captchas.
//
// Parameters:
//
//	b64Img - Base64-encoded string of the captcha image
func WithB64Img(b64Img string) TokenOption {
	return func(c ConfigurableSetter) {
		c.SetB64Img(b64Img)
	}
}

// WithProxy creates an option to set proxy information for the
// captcha solving request. Some captchas require specific IP ranges
// or geographical locations.
//
// Parameters:
//
//	p - Proxy configuration including IP, port, and optional credentials
func WithProxy(p *proxy.Proxy) TokenOption {
	return func(c ConfigurableSetter) {
		c.SetProxy(p)
	}
}

// WithProxyType creates an option to set the proxy protocol type.
//
// Parameters:
//
//	s - Proxy type: "HTTP", "HTTPS", "SOCKS4", or "SOCKS5"
func WithProxyType(s string) TokenOption {
	return func(c ConfigurableSetter) {
		c.SetProxyType(s)
	}
}

// WithUserAgent creates an option to set the user agent string.
// This can be important for captchas that validate browser fingerprints.
//
// Parameters:
//
//	userAgent - User agent string for the HTTP requests
func WithUserAgent(userAgent string) TokenOption {
	return func(c ConfigurableSetter) {
		c.SetUserAgent(userAgent)
	}
}

// WithRqData creates an option to set custom rqdata for hCaptcha.
// This is typically used for hCaptcha Enterprise or invisible implementations.
//
// Parameters:
//
//	rq - rqdata string extracted from the captcha implementation
func WithRqData(rq string) TokenOption {
	return func(c ConfigurableSetter) {
		c.SetRqData(rq)
	}
}
