package captchatoolsgo

type CaptchaAnswer struct {
	Token        string      // the actual captcha answer
	UserAgent    string      // The user agent that was used to solve the captcha answer. Can be an empty string if the service doesn't return anything
	id           interface{} // id from the solving site
	api_key      string      // api key to the solving site
	solving_site site        // site used to get the captcha answer
	capType      captchaType // type of the captcha
}

// newCaptchaAnswer returns a new captcha answer
func newCaptchaAnswer(id interface{}, token string, api_key string, capType captchaType, ss site, ua string) *CaptchaAnswer {
	return &CaptchaAnswer{
		id:           id,
		Token:        token,
		solving_site: ss,
		api_key:      api_key,
		capType:      capType,
		UserAgent:    ua,
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
	case AnticaptchaSite:
		err = report_anticaptcha(was_correct, &c)
	}
	return err
}
