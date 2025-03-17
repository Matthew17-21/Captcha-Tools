package twocaptcha

import (
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/harvester"
	"github.com/Matthew17-21/Captcha-Tools/captchatools-go/proxy"
)

var proxylessCapthaMap = map[harvester.CaptchaType]string{
	harvester.ImageCaptcha: "ImageToTextTask",
	harvester.V2Captcha:    "RecaptchaV2TaskProxyless",
	harvester.V3Captcha:    "RecaptchaV3TaskProxyless",
	harvester.CFTurnstile:  "TurnstileTaskProxyless",
}

var proxyCaptchaMap = map[harvester.CaptchaType]string{
	harvester.ImageCaptcha: "ImageToTextTask",
	harvester.CFTurnstile:  "TurnstileTask",
	harvester.V3Captcha:    "RecaptchaV3TaskProxyless", // 2captha doesn't support with proxy
	harvester.V2Captcha:    "RecaptchaV2Task",
}

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

func newTaskPayload(t Twocaptcha) payload {
	// Create base payload with API key
	p := newPayloadWithClientKey(t.Api_key)

	// Add common task data based on captcha type
	taskData := newPayload()

	// Set method based on captcha type
	switch t.CaptchaType {
	case harvester.ImageCaptcha:
		taskData.setToPayload("type", "ImageToTextTask")
	case harvester.V2Captcha:
		taskData.setToPayload("type", "RecaptchaV2TaskProxyless")
		taskData.setToPayload("method", "userrecaptcha")
		taskData.setToPayload("websiteKey", t.Sitekey)
		if t.IsInvisibleCaptcha {
			taskData.setToPayload("invisible", 1)
		}
	case harvester.V3Captcha:
		taskData.setToPayload("type", "RecaptchaV3TaskProxyless")
		taskData.setToPayload("websiteKey", t.Sitekey)
		taskData.setToPayload("minScore", t.MinScore)
		// TODO: Add pageAction if passed
	case harvester.HCaptcha, harvester.HcaptchaTurbo:
		// TODO: Add captcha type
		taskData.setToPayload("method", "hcaptcha")
		taskData.setToPayload("sitekey", t.Sitekey)
	case harvester.CFTurnstile:
		taskData.setToPayload("type", "TurnstileTaskProxyless")
		taskData.setToPayload("websiteURL", t.Sitekey)
	}

	// TODO: If using a proxy, set the proxy and change the type on taskdata

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
