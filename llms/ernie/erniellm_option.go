package ernie

import "github.com/tmc/langchaingo/callbacks"

const (
	ernieAPIKey    = "ERNIE_API_KEY"    //nolint:gosec
	ernieSecretKey = "ERNIE_SECRET_KEY" //nolint:gosec
)

type ModelName string

const (
	ModelNameERNIEBot                 = "ERNIE-Bot"
	ModelNameERNIEBotTurbo            = "ERNIE-Bot-turbo"
	ModelNameERNIEBotPro              = "ERNIE-Bot-pro"
	ModelNameBloomz7B                 = "BLOOMZ-7B"
	Modelnamellama27bchat             = "Llama-2-7b-chat"
	Modelnamellama213bchat            = "Llama-2-13b-chat"
	Modelnamellama270bchat            = "Llama-2-70b-chat"
	ModelNameERNIEBot408K             = "ERNIE-4.0-8K"
	ModelNameERNIEBot408KPreview      = "ERNIE-4.0-8K-Preview"
	ModelNameERNIEBot408KLatest       = "ERNIE-4.0-8K-Latest"
	ModelNameERNIEBot408K0613         = "ERNIE-4.0-8K-0613"
	ModelNameERNIEBot408K0329         = "ERNIE-4.0-8K-0329"
	ModelNameERNIEBot40Turbo8K        = "ERNIE-4.0-Turbo-8K"
	ModelNameERNIEBot40Turbo8KPreview = "ERNIE-4.0-Turbo-8K-Preview"
	ModelNameERNIEBot358K             = "ERNIE-3.5-8K"
	ModelNameERNIEBot358KPreview      = "ERNIE-3.5-8K-Preview"
	ModelNameERNIEBot35128K           = "ERNIE-3.5-128K"
	ModelNameERNIEBot358K0613         = "ERNIE-3.5-8K-0613"
	ModelNameERNIEBot358K0329         = "ERNIE-3.5-8K-0329"
)

type options struct {
	apiKey           string
	secretKey        string
	accessToken      string
	modelName        ModelName
	callbacksHandler callbacks.Handler
}

type Option func(*options)

// WithAKSK passes the ERNIE API Key and Secret Key to the client. If not set, the keys
// are read from the ERNIE_API_KEY and ERNIE_SECRET_KEY environment variable.
// eg:
//
//	export ERNIE_API_KEY={Api Key}
//	export ERNIE_SECRET_KEY={Serect Key}
//
// Api Key,Serect Key from https://console.bce.baidu.com/qianfan/ais/console/applicationConsole/application
// More information available: https://cloud.baidu.com/doc/WENXINWORKSHOP/s/flfmc9do2
func WithAKSK(apiKey, secretKey string) Option {
	return func(opts *options) {
		opts.apiKey = apiKey
		opts.secretKey = secretKey
	}
}

// WithAccessToken usually used for dev, Prod env recommend use WithAKSK.
func WithAccessToken(accessToken string) Option {
	return func(opts *options) {
		opts.accessToken = accessToken
	}
}

// WithModelName passes the Model Name to the client. If not set, use default ERNIE-Bot.
func WithModelName(modelName ModelName) Option {
	return func(opts *options) {
		opts.modelName = modelName
	}
}

// WithCallbackHandler passes the callback Handler to the client.
func WithCallbackHandler(callbacksHandler callbacks.Handler) Option {
	return func(opts *options) {
		opts.callbacksHandler = callbacksHandler
	}
}
