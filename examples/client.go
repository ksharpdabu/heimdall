package examples

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gojektech/heimdall"
	"github.com/pkg/errors"
)

const (
	baseURL = "http://localhost:9090"
)

func httpClientUsage() error {
	timeoutInMillis := 100
	httpClient := heimdall.NewHTTPClient(timeoutInMillis)
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")

	httpClient.SetRetryCount(2)
	httpClient.SetRetrier(heimdall.NewRetrier(heimdall.NewConstantBackoff(10)))

	response, err := httpClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}

func hystrixClientUsage() error {
	timeoutInMillis := 100

	hystrixConfig := heimdall.NewHystrixConfig("MyCommand", heimdall.HystrixCommandConfig{
		Timeout:                1100,
		MaxConcurrentRequests:  100,
		ErrorPercentThreshold:  25,
		SleepWindow:            10,
		RequestVolumeThreshold: 10,
	})

	hystrixClient := heimdall.NewHystrixHTTPClient(timeoutInMillis, hystrixConfig)
	headers := http.Header{}
	response, err := hystrixClient.Get(baseURL, headers)
	if err != nil {
		return errors.Wrap(err, "failed to make a request to server")
	}

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	fmt.Printf("Response: %s", string(respBody))
	return nil
}
