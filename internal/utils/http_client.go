package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client interface {
	GetDocument(string) (string, error)
}

type WebCrawlerClient struct {
	HTTPClient *http.Client
}

func (v *WebCrawlerClient) GetDocument(url string) (string, error) {
	resp, err := v.HTTPClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("%v error returned from url %v", err, url)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("%v error when parsing response from url %v", err, url)
		}
		return string(bodyBytes), nil
	}
	return "", fmt.Errorf("%v response received", resp.StatusCode)
}


