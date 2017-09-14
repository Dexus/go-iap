package roku

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

const (
	ProductionURL string = "https://apipub.roku.com"
)

// Config is a configuration to initialize client
type Config struct {
	IsProduction bool
	DevToken     string
	TimeOut      time.Duration
}

// The IAPValidationResponse type has the response properties
type IAPValidationResponse struct {
	TransactionID        string  `json:"transactionId"`
	PurchaseDate         string  `json:"purchaseDate"`
	ChannelName          string  `json:"channelName"`
	ProductName          string  `json:"productName"`
	ProductID            string  `json:"productId"`
	Amount               float32 `json:"amount"`
	Currency             string  `json:"currency"`
	Quantity             int     `json:"quantity"`
	ExpirationDate       string  `json:"expirationDate"`
	OriginalPurchaseDate string  `json:"originalPurchaseDate"`
	Status               string  `json:"status"`
	ErrorMessage         string  `json:"errorMessage"`
}

type IAPResponseError struct {
	Status       string `json:"status"`
	Message      string `json:"errorMessage"`
	ErrorDetails string `json:"errorDetails"`
	ErrorCode    string `json:"errorCode"`
}

// IAPClient is an interface to call validation API in Amazon App Store
type IAPClient interface {
	Verify(string) (IAPValidationResponse, error)
}

// Client implements IAPClient
type Client struct {
	URL      string
	DevToken string
	TimeOut  time.Duration
}

// New creates a client object
func New(devToken string) IAPClient {
	client := Client{
		URL:      ProductionURL,
		DevToken: devToken,
		TimeOut:  time.Second * 5,
	}
	return client
}

// NewWithConfig creates a client with configuration
func NewWithConfig(config Config) Client {
	if config.TimeOut == 0 {
		config.TimeOut = time.Second * 5
	}

	client := Client{
		URL:      ProductionURL,
		DevToken: config.DevToken,
		TimeOut:  config.TimeOut,
	}
	if config.IsProduction {
		client.URL = ProductionURL
	}

	return client
}

// Verify sends receipts and gets validation result
func (c Client) Verify(transactionID string) (IAPValidationResponse, error) {
	result := IAPValidationResponse{}
	url := fmt.Sprintf("%v/listen/transaction-service.svc/validate-transaction/%v/%v", c.URL, c.DevToken, transactionID)
	res, body, errs := gorequest.New().
		Get(url).
		Timeout(c.TimeOut).
		End()

	if errs != nil {
		return result, fmt.Errorf("%v", errs)
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		responseError := IAPResponseError{}
		json.NewDecoder(strings.NewReader(body)).Decode(&responseError)
		return result, errors.New(responseError.Message)
	}

	err := json.NewDecoder(strings.NewReader(body)).Decode(&result)

	return result, err
}
