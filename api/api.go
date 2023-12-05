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

func (c *ClientConfig) GetStatus(ctx context.Context) (*GetStatusResponse, error) {
	data, err := c.Request(ctx, http.MethodGet, "status", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetStatusResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetMe(ctx context.Context) (*GetMeResponse, error) {
	data, err := c.Request(ctx, http.MethodGet, "me", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetMeResponse
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

func (c *ClientConfig) GetKernelSpecs(ctx context.Context) (*GetKernelSpecsResponse, error) {
	data, err := c.Request(ctx, http.MethodGet, "kernelspecs", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetKernelSpecsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetKernels(ctx context.Context) (*GetKernelsResponse, error) {
	data, err := c.Request(ctx, http.MethodGet, "kernels", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetKernelsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateKernel(ctx context.Context, options CreateKernelBody) (*CreateKernelResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	data, err := c.Request(ctx, http.MethodPost, "kernels", "application/json", body)
	if err != nil {
		return nil, err
	}
	var result CreateKernelResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetKernel(ctx context.Context, kernel string) (*GetKernelResponse, error) {
	url := fmt.Sprintf("kernels/%s", kernel)
	data, err := c.Request(ctx, http.MethodGet, url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetKernelResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) DeleteKernel(ctx context.Context, kernel string) error {
	url := fmt.Sprintf("kernels/%s", kernel)
	_, err := c.Request(ctx, http.MethodDelete, url, "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) InterruptKernel(ctx context.Context, kernel string) error {
	url := fmt.Sprintf("kernels/%s/interrupt", kernel)
	_, err := c.Request(ctx, http.MethodPost, url, "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) RestartKernel(ctx context.Context, kernel string) error {
	url := fmt.Sprintf("kernels/%s/restart", kernel)
	_, err := c.Request(ctx, http.MethodPost, url, "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}

// TODO: GET and PATCH /api/config/...

func (c *ClientConfig) GetTerminals(ctx context.Context) (*GetTerminalsResponse, error) {
	data, err := c.Request(ctx, http.MethodGet, "terminals", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetTerminalsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateTerminal(ctx context.Context) (*CreateTerminalResponse, error) {
	data, err := c.Request(ctx, http.MethodPost, "terminals", "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result CreateTerminalResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetTerminal(ctx context.Context, terminal string) (*GetTerminalResponse, error) {
	url := fmt.Sprintf("terminals/%s", terminal)
	data, err := c.Request(ctx, http.MethodGet, url, "application/json", nil)
	if err != nil {
		return nil, err
	}

	var result GetTerminalResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) DeleteTerminal(ctx context.Context, terminal string) error {
	url := fmt.Sprintf("terminals/%s", terminal)
	_, err := c.Request(ctx, http.MethodDelete, url, "application/json", nil)
	if err != nil {
		return err
	}
	return nil
}
