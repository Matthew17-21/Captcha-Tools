package twocaptcha

import (
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/proxy"
)

// captchaTypeMap maps harvester captcha types to 2captcha API task types
var captchaTypeMap = map[bool]map[harvester.CaptchaType]string{
	// Without proxy (proxyless)
	false: {
		harvester.ImageCaptcha: "ImageToTextTask",
		harvester.V2Captcha:    "RecaptchaV2TaskProxyless",
		harvester.V3Captcha:    "RecaptchaV3TaskProxyless",
		harvester.CFTurnstile:  "TurnstileTaskProxyless",
		harvester.HCaptcha:     "HCaptchaTaskProxyless",
	},
	// With proxy
	true: {
		harvester.ImageCaptcha: "ImageToTextTask",
		harvester.V2Captcha:    "RecaptchaV2Task",
		harvester.V3Captcha:    "RecaptchaV3Task", // Some providers might not support this
		harvester.CFTurnstile:  "TurnstileTask",
		harvester.HCaptcha:     "HCaptchaTask",
	},
}

// Payload for 2captcha requests
type payload map[string]any

// newPayload creates and returns a new payload
func newPayload() payload {
	return make(payload)
}

// newPayloadWithClientKey creates a payload with the API key
func newPayloadWithClientKey(clientKey string) payload {
	p := newPayload()
	p.setToPayload("clientKey", clientKey)
	return p
}

// newTaskPayload creates a complete task payload based on captcha configuration
func newTaskPayload(t Twocaptcha, options ...harvester.TokenOption) payload {
	// Create base payload with API key
	p := newPayloadWithClientKey(t.Api_key)

	// Add common task data based on captcha type
	taskData := newPayload()

	// Apply all token options to the payload
	for _, opt := range options {
		opt(&taskData)
	}

	// Check if proxy is being used
	hasProxy := taskData["proxyAddress"] != nil

	// Get the appropriate task type based on captcha type and proxy usage
	taskType, exists := captchaTypeMap[hasProxy][t.CaptchaType]
	if !exists {
		// Fallback to proxyless version if combination doesn't exist
		taskType = captchaTypeMap[false][t.CaptchaType]
	}

	// Set task type based on captcha type and proxy presence
	taskData.setToPayload("type", taskType)

	// Add specific parameters based on captcha type
	switch t.CaptchaType {
	case harvester.ImageCaptcha:
		// Image captcha specific settings already handled by type
	case harvester.V2Captcha:
		taskData.setToPayload("websiteKey", t.Sitekey)
		if t.IsInvisibleCaptcha {
			taskData.setToPayload("invisible", 1)
		}
	case harvester.V3Captcha:
		taskData.setToPayload("websiteKey", t.Sitekey)
		taskData.setToPayload("minScore", t.MinScore)
		if t.Action != "" {
			taskData.setToPayload("pageAction", t.Action)
		}
	case harvester.HCaptcha, harvester.HcaptchaTurbo:
		taskData.setToPayload("sitekey", t.Sitekey)
	case harvester.CFTurnstile:
		taskData.setToPayload("websiteKey", t.Sitekey)
	}

	// Add page URL
	taskData.setToPayload("websiteURL", t.CaptchaURL)

	// Add SoftID, if specified
	if t.SoftID != 0 {
		taskData.setToPayload("soft_id", t.SoftID)
	}

	// Merge task data into main payload
	p.setToPayload("task", taskData)
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

	// Add authentication if provided
	if pxy.User != "" && pxy.Password != "" {
		p.setToPayload("proxyLogin", pxy.User)
		p.setToPayload("proxyPassword", pxy.Password)
	}
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
