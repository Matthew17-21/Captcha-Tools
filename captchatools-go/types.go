package captchatoolsgo

// General type declarations
const (
	CapmonsterSite  site = iota // The int 1 will represent Capmonter
	AnticaptchaSite             // The int 2 will represent Anticaptcha
	TwoCaptchaSite              // The int 3 will represent 2captcha
)

const (
	V2Captcha    captchaType = "v2"
	V3Captcha    captchaType = "v3"
	HCaptcha     captchaType = "hcaptcha"
	ImageCaptcha captchaType = "image"
)

type (
	site        int
	captchaType string

	// Configurations for the captchas you are solving.
	// For more a detailed documentation, visit
	// https://github.com/Matthew17-21/Captcha-Tools
	Config struct {
		Api_key            string      // The API Key for the captcha solving site.
		Sitekey            string      // Sitekey from the site where captcha is loaded.
		CaptchaURL         string      // URL where the captcha is located.
		CaptchaType        captchaType // Type of captcha you are solving. Visit https://github.com/Matthew17-21/Captcha-Tools for types
		Action             string      // Action that is associated with the V3 captcha.
		IsInvisibleCaptcha bool        // If the captcha is invisible or not.
		MinScore           float32     // Minimum score for v3 captchas.
		SoftID             int         // SoftID for 2captcha. Developers get reward 10% of spendings of their software users.
	}

	/*
		- type Harvester will be used to represent a captcha harvester.

		- In order to to have the same functionality/implementation as the
		Python version, type Harvester has the `childHarvester` field, which is of type interface.
		This allows us to set the field as a pointer to a real captcha harvester
		and use the `GetToken` method that each captcha harvester struct has.
		(Polymorphism)

		- For documentation, visit https://github.com/Matthew17-21/Captcha-Tools
	*/
	// Interface that will allow us to interact with the methods from the
	// individual structs
	Harvester interface {
		GetToken(additional ...*AdditionalData) (*CaptchaAnswer, error) // Function to get a captcha token
		GetBalance() (float32, error)
	}

	AdditionalData struct {
		B64Img    string // base64 encoded image
		Proxy     *Proxy // A proxy in correct formatting - such as user:pass@ip:port
		ProxyType string // Type of your proxy: HTTP, HTTPS, SOCKS4, SOCKS5
		UserAgent string // UserAgent that will be passed to the service and used to solve the captcha
	}

	CaptchaAnswer struct {
		Token        string      // the actual captcha answer
		id           interface{} // id from the solving site
		api_key      string      // api key to the solving site
		solving_site site        // site used to get the captcha answer
		capType      captchaType // type of the captcha
	}

	Proxy struct {
		Ip       string
		Port     string
		User     string
		Password string
	}

	Anticaptcha struct {
		config *Config
	}

	Capmonster struct {
		config *Config
	}

	Twocaptcha struct {
		config *Config
	}
)

// Payload type declarations
type (
	// This struct will be the payload to get the queue ID from capmonster
	capmonsterIDPayload struct {
		ClientKey string `json:"clientKey"`
		SoftID    int    `json:"softId,omitempty"`
		Task      struct {
			WebsiteURL  string      `json:"websiteURL"`
			WebsiteKey  string      `json:"websiteKey"`
			Type        captchaType `json:"type"`
			IsInvisible bool        `json:"isInvisible,omitempty"`
			MinScore    float32     `json:"minScore,omitempty"`
			PageAction  string      `json:"pageAction,omitempty"`
			Body        string      `json:"body,omitempty"`
		} `json:"task"`
	}
	capmonsterCapAnswerPayload struct {
		ClientKey string `json:"clientKey"`
		TaskID    int    `json:"taskId"`
	}
)

// Response type declarations
type (
	capmonsterIDResponse struct {
		ErrorID   int    `json:"errorId"`
		ErrorCode string `json:"errorCode"`
		TaskID    int    `json:"taskId"`
	}

	capmonsterTokenResponse struct {
		ErrorID   int    `json:"errorId"`
		ErrorCode string `json:"errorCode"`
		Solution  struct {
			Text               string `json:"text"`
			GRecaptchaResponse string `json:"gRecaptchaResponse"`
		} `json:"solution"`
		Status string `json:"status"`
	}
	capmonsterBalanceResponse struct {
		Balance   float32 `json:"balance"`
		ErrorCode string  `json:"errorCode"`
		ErrorID   int     `json:"errorId"`
	}

	twoCapIDPayload struct {
		Key       string  `json:"key"`
		Method    string  `json:"method"`
		Googlekey string  `json:"googlekey"`
		Pageurl   string  `json:"pageurl"`
		JSON      int     `json:"json"`
		Sitekey   string  `json:"sitekey,omitempty"`
		Invisible int     `json:"invisible,omitempty"`
		Version   string  `json:"version,omitempty"`
		Action    string  `json:"action,omitempty"`
		MinScore  float32 `json:"min_score,omitempty"`
		SoftID    int     `json:"soft_id,omitempty"`
		Body      string  `json:"body,omitempty"`      // Base64-encoded captcha image
		UserAgent string  `json:"userAgent,omitempty"` // userAgent that will be used to solve the captcha
		Proxy     string  `json:"proxy,omitempty"`     // Proxy to use to solve captchas from
		ProxyType string  `json:"proxytype,omitempty"` // Type of the proxy
	}
	twocaptchaResponse struct {
		Status  int    `json:"status"`
		Request string `json:"request"`
	}

	anticaptchaBalanceResponse struct {
		ErrorID          int     `json:"errorId"`
		ErrorCode        string  `json:"errorCode"`
		ErrorDescription string  `json:"errorDescription"`
		Balance          float32 `json:"balance"`
	}
)
