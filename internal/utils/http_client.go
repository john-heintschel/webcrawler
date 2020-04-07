package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Client interface {
	GetDocument(string) (string, error)
}

type WebCrawlerClient struct {
	HTTPClient *http.Client
}

func (v WebCrawlerClient) GetDocument(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:68.0) Gecko/20100101 Firefox/68.0")
	resp, err := v.HTTPClient.Do(req)
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
