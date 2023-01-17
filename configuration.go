/*
 * DevCycle Bucketing API
 *
 * Documents the DevCycle Bucketing API which provides and API interface to User Bucketing and for generated SDKs.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package devcycle

import (
	"net/http"
	"time"
)

// contextKeys are used to identify the type of value in the context.
// Since these are string, it is possible to get a short description of the
// context key for logging and debugging using key.String().

type contextKey string

func (c contextKey) String() string {
	return "auth " + string(c)
}

var (
	// ContextOAuth2 takes a oauth2.TokenSource as authentication for the request.
	ContextOAuth2 = contextKey("token")

	// ContextBasicAuth takes BasicAuth as authentication for the request.
	ContextBasicAuth = contextKey("basic")

	// ContextAccessToken takes a string oauth2 access token as authentication for the request.
	ContextAccessToken = contextKey("accesstoken")

	// ContextAPIKey takes an APIKey as authentication for the request
	ContextAPIKey = contextKey("apikey")
)

// BasicAuth provides basic http authentication to a request passed via context using ContextBasicAuth
type BasicAuth struct {
	UserName string `json:"userName,omitempty"`
	Password string `json:"password,omitempty"`
}

// APIKey provides API key based authentication to a request passed via context using ContextAPIKey
type APIKey struct {
	Key    string
	Prefix string
}

type DVCOptions struct {
	EnableEdgeDB                 bool          `json:"enableEdgeDb,omitempty"`
	EnableCloudBucketing         bool          `json:"enableCloudBucketing,omitempty"`
	EventsFlushInterval          time.Duration `json:"eventsFlushInterval,omitempty"`
	PollingInterval              time.Duration `json:"pollingInterval,omitempty"`
	RequestTimeout               time.Duration `json:"requestTimeout,omitempty"`
	DisableAutomaticEventLogging bool          `json:"disableAutomaticEventLogging,omitempty"`
	DisableCustomEventLogging    bool          `json:"disableCustomEventLogging,omitempty"`
	MaxEventQueueSize            int           `json:"maxEventsPerFlush,omitempty"`
	FlushEventQueueSize          int           `json:"minEventsPerFlush,omitempty"`
	ConfigCDNOverride            string
	EventsAPIOverride            string
}

func (o *DVCOptions) CheckDefaults() {
	if o.EventsFlushInterval <= time.Second*1 {
		o.EventsFlushInterval = time.Second * 1
	}
	if o.PollingInterval <= 1000 {
		o.PollingInterval = time.Second * 10
	}
	if o.RequestTimeout <= time.Second*5 {
		o.RequestTimeout = time.Second * 5
	}
	if o.MaxEventQueueSize <= 0 {
		o.MaxEventQueueSize = 10000
	}
	if o.FlushEventQueueSize <= 0 {
		o.FlushEventQueueSize = 1000
	}
}

type HTTPConfiguration struct {
	BasePath          string            `json:"basePath,omitempty"`
	ConfigCDNBasePath string            `json:"configCDNBasePath,omitempty"`
	EventsAPIBasePath string            `json:"eventsAPIBasePath,omitempty"`
	Host              string            `json:"host,omitempty"`
	Scheme            string            `json:"scheme,omitempty"`
	DefaultHeader     map[string]string `json:"defaultHeader,omitempty"`
	UserAgent         string            `json:"userAgent,omitempty"`
	HTTPClient        *http.Client
}

func NewConfiguration(options *DVCOptions) *HTTPConfiguration {
	configBasePath := "https://config-cdn.devcycle.com"
	if options.ConfigCDNOverride != "" {
		configBasePath = options.ConfigCDNOverride
	}

	eventsApiBasePath := "https://events.devcycle.com"
	if options.EventsAPIOverride != "" {
		eventsApiBasePath = options.EventsAPIOverride
	}

	cfg := &HTTPConfiguration{
		BasePath:          "https://bucketing-api.devcycle.com",
		ConfigCDNBasePath: configBasePath,
		EventsAPIBasePath: eventsApiBasePath,
		DefaultHeader:     make(map[string]string),
		UserAgent:         "DevCycle-Server-SDK/" + VERSION + "/go",
		HTTPClient:        http.DefaultClient,
	}
	return cfg
}

func (c *HTTPConfiguration) AddDefaultHeader(key string, value string) {
	c.DefaultHeader[key] = value
}
