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

// report submits to the captcha solving service whether or not the answer
// they provided was correct or not.
//
// If the answer is reported incorrect and the solving service accepts it,
// a refund will be credited
func (c CaptchaAnswer) Report(was_correct bool) error {
	var err error
	switch c.solving_site {
	case TwoCaptchaSite:
		err = report_2captcha(was_correct, &c)
	}
	return err
}
