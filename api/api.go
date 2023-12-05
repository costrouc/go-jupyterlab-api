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
	url := fmt.Sprintf("contents/%s", path)
	if options != nil {
		url = fmt.Sprintf("contents/%s?%s", path, options.Encode())
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
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("contents/%s", path)
	data, err := c.Request(ctx, http.MethodPost, url, "application/json", body)
	if err != nil {
		return nil, err
	}

	var result CreateContentsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) PatchContents(ctx context.Context, path string, options *PatchContentsBody) (*PatchContentsResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("contents/%s", path)
	data, err := c.Request(ctx, http.MethodPatch, url, "application/json", body)
	if err != nil {
		return nil, err
	}

	var result PatchContentsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) PutContents(ctx context.Context, path string, options *PutContentsBody) (*PutContentsResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("contents/%s", path)
	data, err := c.Request(ctx, http.MethodPut, url, "application/json", body)
	if err != nil {
		return nil, err
	}

	var result PutContentsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) DeleteContents(ctx context.Context, path string) error {
	url := fmt.Sprintf("contents/%s", path)
	_, err := c.Request(ctx, http.MethodDelete, url, "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Contents Checkpoint API

func (c *ClientConfig) GetSessions(ctx context.Context) (*GetSessionsResponse, error) {
	data, err := c.Request(ctx, http.MethodGet, "sessions", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetSessionsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateSession(ctx context.Context, session *Session) (*CreateSessionResponse, error) {
	body, err := json.Marshal(session)
	if err != nil {
		return nil, err
	}

	data, err := c.Request(ctx, http.MethodPost, "sessions", "application/json", body)
	if err != nil {
		return nil, err
	}

	var result CreateSessionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetSession(ctx context.Context, session string) (*GetSessionResponse, error) {
	url := fmt.Sprintf("sessions/%s", session)
	data, err := c.Request(ctx, http.MethodGet, url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetSessionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) PatchSession(ctx context.Context, session string, options *Session) (*PatchSessionResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("sessions/%s", session)
	data, err := c.Request(ctx, http.MethodPatch, url, "application/json", body)
	if err != nil {
		return nil, err
	}

	var result PatchSessionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) DeleteSession(ctx context.Context, session string) error {
	url := fmt.Sprintf("sessions/%s", session)
	_, err := c.Request(ctx, http.MethodDelete, url, "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}
