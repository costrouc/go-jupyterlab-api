package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CreateClient(config *ClientConfig) (*ClientConfig, error) {
	clientConfig := ClientConfig{
		ApiToken: "",
		ApiURL:   "http://localhost:8888/api",
	}

	if config.ApiToken != "" {
		clientConfig.ApiToken = config.ApiToken
	} else {
		apiToken, ok := os.LookupEnv("JUPYTERLAB_API_TOKEN")
		if !ok {
			return nil, errors.New("api token not defined can be set via JUPYTERLAB_API_TOKEN")
		}
		clientConfig.ApiToken = apiToken
	}
	return &clientConfig, nil
}

func (c *ClientConfig) Request(ctx context.Context, method string, path string, contentType string, requestBody []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.ApiURL, path)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	req = req.WithContext(ctx)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiToken))
	req.Header.Set("Content-Type", contentType)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("response returned status code of %d instead of 2XX", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *ClientConfig) GetVersion(ctx context.Context) (*GetVersionResponse, error) {
	data, err := c.Request(ctx, http.MethodGet, "", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetVersionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetContents(ctx context.Context, path string, options *GetContentsParams) (*GetContentsResponse, error) {
	url := "contents"
	if options != nil {
		url = fmt.Sprintf("%s?%s", url, options.Encode())
	}

	data, err := c.Request(ctx, http.MethodGet, url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetContentsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateContents(ctx context.Context, path string, options *CreateContentsBody) (*CreateContentsResponse, error) {
	_, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (c *ClientConfig) ListKernels(ctx context.Context) {

}

func (c *ClientConfig) CreateKernel(ctx context.Context, kernelSpec string) {

}
