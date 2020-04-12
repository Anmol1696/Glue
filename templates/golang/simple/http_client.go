/*
Http client can be used to make backend calls to any http request
*/
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// Different kinds of bankend api are identified and initialized here
// Actual Url of these apis should come from the config
const (
	BackendAPI = "BackendAPI"
)

// Custom errors raised from the backend
var (
	ErrGettingData = errors.New("Error getting data")
)

// Struct for request for Http client calls
// Holds information for making a http reqeust
type HttpStruct struct {
	RequestId string            `json:"request_id,omitempty"`
	Api       string            `json:"api,omitempty"` // Holds kind of backend api
	Uri       string            `json:"uri,omitempty"`
	Method    string            `json:"method,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
	// Data for body of request
	Data []byte `json:"data,omitempty"`
}

// Struct holds the client, logger and config
// All method made from the client are defined as methods
type httpClient struct {
	client *http.Client
	logger *zap.Logger
	config *Config
}

// Initialization of the client
// Client is initialized on app startup hence fail fast incase of setup failure
//
// For Oauth2 client use following code snipet
// ---
//    OAuth2ClientCredsConfig := &clientcredentials.Config{
//            ClientID:     config.KeycloakClientId,
//            ClientSecret: config.KeycloakClientSecret,
//            TokenURL:     strings.Join([]string{config.KeycloakUrl, config.KeycloakTokenUri}, "/"),
//        }
//        client := OAuth2ClientCredsConfig.Client(context.Background())
// ---
func NewHttpClient(logger *zap.Logger, config *Config) (*httpClient, error) {
	client := &http.Client{}

	return &httpClient{client: client, logger: logger, config: config}, nil
}

// Create request to the backend api, log incase of error log with errorStr
func (hc *httpClient) makeRequest(r *HttpStruct, errorStr string) ([]byte, error) {
	// Calculate latency of every request
	start := time.Now()
	defer func() {
		hc.logger.Info("http client request",
			zap.Duration("latency", time.Since(start)),
			zap.String("api", r.Api),
			zap.String("method", r.Method),
			zap.String("uri", r.Uri),
			zap.String("request_id", r.RequestId))
	}()

	// Create url from reqStruct
	var baseUrl string
	switch r.Api {
	// Replace SomeApi with an actual backend api you are calling
	case BackendAPI:
		baseUrl = hc.config.BackendURL
	default:
		err := errors.New("Internal server error. Trying to make request to unknown api")
		hc.logger.Error("Unknown api in struct", zap.Error(err))
		return nil, err
	}
	url := strings.Join([]string{strings.TrimSuffix(baseUrl, "/"), strings.TrimPrefix(r.Uri, "/")}, "/")

	hc.logger.Debug("making request",
		zap.String("url", url),
		zap.Any("request", r))

	// create http request
	req, _ := http.NewRequest(strings.ToUpper(r.Method), url, bytes.NewBuffer(r.Data))
	if r.Headers != nil {
		for key, val := range r.Headers {
			req.Header.Add(key, val)
		}
	}
	if len(r.Data) != 0 {
		req.Header.Add("Content-Type", "application/json")
	}

	resp, err := hc.client.Do(req)
	if err != nil {
		hc.logger.Error("Error from http response", zap.String("while processing", errorStr), zap.Error(err))
		return nil, err
	}

	var respBody []byte
	// Read response body
	respBody, _ = ioutil.ReadAll(resp.Body)

	// Check response is a success
	if resp.StatusCode >= http.StatusMultipleChoices {
		hc.logger.Warn("Unsucessfull response",
			zap.String("url", url),
			zap.Any("request", r),
			zap.Int("response_status", resp.StatusCode),
			zap.String("response_body", string(respBody)),
		)

		return nil, fmt.Errorf("Uncessfull response from api: %s", errorStr)
	}

	hc.logger.Debug("Sucessfull response",
		zap.String("url", url),
		zap.Any("request", r),
		zap.Int("response_status", resp.StatusCode),
		zap.String("response_body", string(respBody)),
	)

	return respBody, nil
}

// Sample get request to the backend at an endpoint
// Replace below function with actual function to backend
func (hc *httpClient) getSomeData(ctx context.Context) (map[string]string, error) {
	reqStruct := &HttpStruct{
		RequestId: getRequestId(ctx),
		Api:       BackendAPI,
		Uri:       "/endpoint",
		Method:    "get",
	}
	body, err := hc.makeRequest(reqStruct, fmt.Sprintf("Getting some data from: %s", "/endpoint"))
	if err != nil {
		return nil, err
	}

	ss := map[string]string{}
	err = json.Unmarshal(body, &ss)
	if err != nil {
		hc.logger.Error("Unable to unmarshal response", zap.String("response string", string(body)))
		return nil, ErrGettingData
	}

	return ss, nil
}
