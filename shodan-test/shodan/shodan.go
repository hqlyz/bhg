package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const BaseURL = "https://api.shodan.io"

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) APIInfo() (*APIInfo, error) {
	resp, err := http.Get(fmt.Sprintf("%s/api-info?key=%s", BaseURL, c.apiKey))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var apiInfo APIInfo
	if err = json.NewDecoder(resp.Body).Decode(&apiInfo); err != nil {
		return nil, err
	}
	return &apiInfo, nil
}

func (c *Client) HostSearch(q string) (*HostSearch, error) {
	resp, err := http.Get(fmt.Sprintf("%s/shodan/host/search?key=%s&query=%s", BaseURL, c.apiKey, q))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var hostSearch HostSearch
	if err = json.NewDecoder(resp.Body).Decode(&hostSearch); err != nil {
		return nil, err
	}
	return &hostSearch, nil
}