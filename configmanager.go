package devcycle

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var pollingStop = make(chan bool)

type EnvironmentConfigManager struct {
	environmentKey string
	configETag     string
	localBucketing *DevCycleLocalBucketing
	firstLoad      bool
	SDKEvents      chan SDKEvent
}

func (e *EnvironmentConfigManager) Initialize(environmentKey string, options *DVCOptions) {
	e.environmentKey = environmentKey
	e.SDKEvents = make(chan SDKEvent)

	if options.PollingInterval == 0 {
		options.PollingInterval = time.Second * 30
	}
	if options.RequestTimeout == 0 {
		options.RequestTimeout = time.Second * 10
	}

	ticker := time.NewTicker(options.PollingInterval)
	e.firstLoad = true

	e.fetchConfig()
	go func() {
		for {
			select {
			case <-pollingStop:
				ticker.Stop()
				log.Fatal("Stopping config polling.")
				return
			case <-ticker.C:
				e.fetchConfig()
			}
		}
	}()
}

func (e *EnvironmentConfigManager) fetchConfig() {
	resp, err := http.Get(e.getConfigURL())
	if err != nil {
		e.SDKEvents <- SDKEvent{Success: false, Message: "Could not make HTTP Request to CDN.", Error: err}
	}
	switch resp.StatusCode {
	case http.StatusOK:
		err := e.setConfig(resp)
		if err != nil {
			e.SDKEvents <- SDKEvent{Success: false, Message: "Failed to set config.", Error: err}
			return
		}
		break
	case http.StatusNotModified:
		log.Println("Config not modified. Using cached config. %s", e.configETag)
		break
	case http.StatusForbidden:
		pollingStop <- true
		log.Println("403 Forbidden - SDK key is likely incorrect. Aborting polling.")
		return
	case http.StatusInternalServerError:
	case http.StatusBadGateway:
	case http.StatusServiceUnavailable:
		// Retryable Errors. Continue polling.
		log.Println("Retrying config fetch. Status:" + resp.Status)
		break
	default:
		log.Println("Unexpected response code: %d", resp.StatusCode)
		log.Println("Body: %s", resp.Body)
		log.Println("URL: %s", e.getConfigURL())
		log.Println("Headers: %s", resp.Header)
		log.Println("Could not download configuration. Using cached version if available %s", resp.Header.Get("ETag"))
		e.SDKEvents <- SDKEvent{Success: false,
			Message: "Unexpected response code - Aborting Polling. Code: " + strconv.Itoa(resp.StatusCode), Error: nil}
		pollingStop <- true
		break
	}
}

func (e *EnvironmentConfigManager) setConfig(response *http.Response) error {
	raw, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = e.localBucketing.StoreConfig(e.environmentKey, string(raw))
	if err != nil {
		return err
	}
	e.configETag = response.Header.Get("ETag")
	log.Println("Config set. ETag: %s", e.configETag)
	e.SDKEvents <- SDKEvent{Success: true, Message: "Config set. ETag: " + e.configETag, Error: nil}
	if e.firstLoad {
		e.firstLoad = false
		log.Println("DevCycle SDK Initialized.")
		e.SDKEvents <- SDKEvent{Success: true, Message: "DevCycle SDK Initialized.", Error: nil, FirstInitialization: true}
	}
	return nil
}

func (e *EnvironmentConfigManager) getConfigURL() string {
	return fmt.Sprintf("https://config-cdn.devcycle.com/config/v1/server/%s.json", e.environmentKey)
}
