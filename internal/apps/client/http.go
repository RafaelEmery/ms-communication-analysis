package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
)

func HandleHTTP(baseURL string, data InteractionData) {
	endpoint, method := getRequestData(data.Resource, baseURL)

	for i := 0; i < data.RequestQuantity; i++ {
		start := time.Now()
		resp, err := doRequest(endpoint, method)
		if err != nil {
			log.Default().Println(err)
			continue
		}
		defer resp.Body.Close()

		elapsed := time.Since(start).String()
		log.Default().Printf("[%d] %s - %s", resp.StatusCode, endpoint, elapsed)
	}
}

func getRequestData(resource, baseURL string) (string, string) {
	var (
		endpoint      string
		requestMethod string
	)
	switch resource {
	case createResource:
		endpoint = baseURL + "/products"
		requestMethod = http.MethodPost
	case reportResource:
		endpoint = baseURL + "/products/report"
		requestMethod = http.MethodGet
	case getByDiscountResource:
		endpoint = baseURL + "/products/discount"
		requestMethod = http.MethodGet
	}

	return endpoint, requestMethod
}

func doRequest(endpoint, method string) (*http.Response, error) {
	if method == http.MethodPost {
		var product domain.Product

		payload, err := json.Marshal(product.Fake())
		if err != nil {
			log.Default().Println(err)
			return nil, err
		}
		body := bytes.NewBuffer(payload)

		return http.Post(endpoint, "application/json", body)
	}
	if method == http.MethodGet {
		return http.Get(endpoint)
	}

	return nil, fmt.Errorf("the method %s is not allowed", method)
}
