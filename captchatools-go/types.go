package captchatoolsgo

// General type declarations
const (
	CapmonsterSite  site = iota // The int 1 will represent Capmonter
	AnticaptchaSite             // The int 2 will represent Anticaptcha
	TwoCaptchaSite              // The int 3 will represent 2captcha
)

type (
	site int

	// Configurations for the captchas you are solving.
	// For more a detailed documentation, visit
	// https://github.com/Matthew17-21/Captcha-Tools
	Config struct {
		Api_key            string  // The API Key for the captcha solving site.
		Sitekey            string  // Sitekey from the site where captcha is loaded.
		CaptchaURL         string  // URL where the captcha is located.
		CaptchaType        string  // Type of captcha you are solving. Visit https://github.com/Matthew17-21/Captcha-Tools for types
		Action             string  // Action that is associated with the V3 captcha.
		IsInvisibleCaptcha bool    // If the captcha is invisible or not.
		MinScore           float32 // Minimum score for v3 captchas.
		SoftID             int     // SoftID for 2captcha. Developers get reward 10% of spendings of their software users.
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
		GetToken() (string, error) // Function to get a captcha token
		GetBalance() (float32, error)
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
		Task      struct {
			WebsiteURL  string  `json:"websiteURL"`
			WebsiteKey  string  `json:"websiteKey"`
			Type        string  `json:"type"`
			IsInvisible bool    `json:"isInvisible,omitempty"`
			MinScore    float32 `json:"minScore,omitempty"`
			PageAction  string  `json:"pageAction,omitempty"`
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
		ErrorID  int `json:"errorId"`
		Solution struct {
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
