package captchatoolsgo

// newCaptchaAnswer returns a new captcha answer
func newCaptchaAnswer(id interface{}, token string, api_key string, ss site) *CaptchaAnswer {
	return &CaptchaAnswer{
		id:           id,
		Token:        token,
		solving_site: ss,
		api_key:      api_key,
	}
}

// Returns the id, associated with the captcha answer, from solving site
func (c CaptchaAnswer) Id() interface{} {
	return c.id
}
