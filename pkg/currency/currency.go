package currency

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	currencyCode = "usd"
)

type Client struct {
	currencyCode string
}

type ExchangeRate struct {
	No            string  `json:"no"`
	EffectiveDate string  `json:"effectiveDate"`
	Mid           float64 `json:"mid"`
}

type ExchangeRatesResponse struct {
	Table    string         `json:"table"`
	Currency string         `json:"currency"`
	Code     string         `json:"code"`
	Rates    []ExchangeRate `json:"rates"`
}

func NewClient() *Client {
	return &Client{
		currencyCode: currencyCode,
	}
}

func (c Client) GetLast100() (*ExchangeRatesResponse, error) {
	startTime := time.Now()
	clientURL := fmt.Sprintf("https://api.nbp.pl/api/exchangerates/rates/a/%s/last/100/?format=json", c.currencyCode)
	httpRequest, err := http.NewRequest(http.MethodGet, clientURL, nil)
	httpRequest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	httpClient := http.Client{}

	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	duration := time.Since(startTime)
	log.Printf("Request GET %s took %d ms", clientURL, duration.Milliseconds())

	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotModified) {
		return nil, fmt.Errorf("incorrect status respone code: expected 200 but fount %d", resp.StatusCode)
	}

	if contentType := resp.Header.Get("Content-Type"); !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("incorrect Content-Type: expected application/json but found %s", contentType)
	}

	var responseStruct ExchangeRatesResponse
	err = json.NewDecoder(resp.Body).Decode(&responseStruct)
	if err != nil {
		return nil, err
	}
	return &responseStruct, nil
}
