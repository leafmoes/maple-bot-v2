package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CFKVClient struct {
	*http.Client
	BaseURL     string
	accountId   string
	namespaceId string
	apiToken    string
}

type WorkersKVPair struct {
	Key           string      `json:"key"`
	Value         string      `json:"value"`
	Expiration    int         `json:"expiration,omitempty"`
	ExpirationTTL int         `json:"expiration_ttl,omitempty"`
	Metadata      interface{} `json:"metadata,omitempty"`
	Base64        bool        `json:"base64,omitempty"`
}

type Response struct {
	Success  bool           `json:"success"`
	Errors   []ResponseInfo `json:"errors"`
	Messages []ResponseInfo `json:"messages"`
}

type ResponseInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewCFClient(accountId string, namespaceId string, apiToken string) *CFKVClient {
	return &CFKVClient{
		&http.Client{},
		"https://api.cloudflare.com/client/v4",
		accountId,
		namespaceId,
		apiToken,
	}
}

func (c *CFKVClient) Read(key string) ([]byte, error) {
	url := fmt.Sprintf("/accounts/%s/storage/kv/namespaces/%s/values/%s", c.accountId, c.namespaceId, url.PathEscape(key))
	res, err := c.makeRequestWithHeaders(http.MethodGet, url, nil, nil, 10)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *CFKVClient) Write(key string, value []byte) (bool, error) {
	url := fmt.Sprintf("/accounts/%s/storage/kv/namespaces/%s/values/%s", c.accountId, c.namespaceId, url.PathEscape(key))
	res, err := c.makeRequestWithHeaders(http.MethodGet, url, value, nil, 10)
	if err != nil {
		return false, err
	}
	result := Response{}
	if err := json.Unmarshal(res, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal json response: %w", err)
	} else if !result.Success {
		return false, nil
	}
	return true, err
}

func (c *CFKVClient) makeRequestWithHeaders(method string, url string, params interface{}, header http.Header, maxRetries int) ([]byte, error) {
	var err error
	var req *http.Request
	var reqBody io.Reader
	var resp *http.Response
	var respErr error
	var respBody []byte
	for i := 0; i <= maxRetries; i++ {
		if params != nil {
			if r, ok := params.(io.Reader); ok {
				reqBody = r
			} else if paramBytes, ok := params.([]byte); ok {
				reqBody = bytes.NewReader(paramBytes)
			} else {
				var jsonBody []byte
				jsonBody, err = json.Marshal(params)
				if err != nil {
					return nil, fmt.Errorf("error marshalling params to JSON: %w", err)
				}
				reqBody = bytes.NewReader(jsonBody)
			}
		}
		req, err = http.NewRequest(method, c.BaseURL+url, reqBody)
		if header != nil {
			req.Header = header
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+c.apiToken)
		resp, respErr := c.Do(req)
		if respErr != nil || resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode >= 500 {
			if resp != nil && resp.StatusCode == http.StatusTooManyRequests {
				respErr = errors.New("exceeded available rate limit retries")
			}

			if respErr == nil {
				respErr = fmt.Errorf("received %s response (HTTP %d), please try again later", strings.ToLower(http.StatusText(resp.StatusCode)), resp.StatusCode)
			}
			continue
		} else {
			respBody, err = io.ReadAll(resp.Body)
			defer resp.Body.Close()
			if err != nil {
				return nil, fmt.Errorf("could not read response body: %w", err)
			}

			break
		}
	}
	if respErr != nil {
		return nil, respErr
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("failed to access cf kv: %w", err)
	}
	return respBody, nil
}
