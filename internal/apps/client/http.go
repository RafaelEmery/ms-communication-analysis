package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	domain "github.com/RafaelEmery/performance-analysis-server/internal"
)

var ErrRequestFailed = errors.New("request failed")

func HandleHTTP(baseURL string, data InteractionData) error {
	endpoint, method := getRequestData(data.Resource, baseURL)

	for i := 0; i < data.RequestQuantity; i++ {
		start := time.Now()
		resp, err := doRequest(endpoint, method)
		if err != nil {
			log.Default().Println(err)
			if data.RequestQuantity == 1 {
				return err
			}
			continue
		}
		defer resp.Body.Close()

		elapsed := time.Since(start).String()
		log.Default().Printf("[%d] %s - %s", resp.StatusCode, endpoint, elapsed)
	}

	return nil
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
	var r *http.Response
	var err error

	if method == http.MethodPost {
		var product domain.Product

		payload, err := json.Marshal(product.Fake())
		if err != nil {
			log.Default().Println(err)
			return nil, err
		}
		body := bytes.NewBuffer(payload)
		r, err = http.Post(endpoint, "application/json", body)
		if r.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, ErrRequestFailed
			}

			bodyString := string(body)
			return nil, fmt.Errorf("request failed: %s", bodyString)
		}

		return r, err
	}
	if method == http.MethodGet {
		r, err = http.Get(endpoint)
		if r.StatusCode != http.StatusOK {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return nil, ErrRequestFailed
			}

			bodyString := string(body)
			return nil, fmt.Errorf("request failed: %s", bodyString)
		}

		return r, err
	} else {
		return nil, fmt.Errorf("the method %s is not allowed", method)
	}
}
