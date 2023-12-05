package api

import (
	"fmt"
	"net/url"
)

type ClientConfig struct {
	ApiToken string
	ApiURL   string
}

type GetVersionResponse struct {
	Version string `json:"version"`
}

type GetContentsParams struct {
	Type    string
	Format  string
	Content int
	Hash    int
}

func (r *GetContentsParams) Encode() string {
	v := url.Values{}
	if r.Type != "" {
		v.Set("type", r.Type)
	}
	if r.Format != "" {
		v.Set("type", r.Format)
	}
	if r.Content != 1 {
		v.Set("content", fmt.Sprintf("%d", r.Content))
	}
	if r.Hash != 0 {
		v.Set("hash", fmt.Sprintf("%d", r.Hash))
	}
	return v.Encode()
}

type Content struct {
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	LastModified  string    `json:"last_modified"`
	Created       string    `json:"created"`
	Content       []Content `json:"content"`
	Format        string    `json:"format"`
	Mimetype      string    `json:"mimetype"`
	Size          int       `json:"size"`
	Type          string    `json:"type"`
	Writeable     bool      `json:"writeable"`
	Hash          string    `json:"hash"`
	HashAlgorithm string    `json:"hash_algorithm"`
}

type GetContentsResponse Content

type CreateContentsBody struct {
	CopyFrom string `json:"copy_from,omitempty"`
	Ext      string `json:"ext,omitempty"`
	Type     string `json:"type,omitempty"`
}

type CreateContentsResponse Content

type PatchContentsBody struct {
	Path string `json:"path"`
}

type PatchContentsResponse Content

type PutContentsBody struct {
	Content string `json:"content"`
	Format  string `json:"format"` // json, text, base64
	Name    string `json:"name"`
	Path    string `json:"path"`
	Type    string `json:"type"` // notebook, file, directory
}

type PutContentsResponse Content

type Session struct {
	Id     string      `json:"id"`
	Kernel interface{} `json:"kernel"`
	Name   string      `json:"name"`
	Path   string      `json:"path"`
	Type   string      `json:"type"`
}

type GetSessionsResponse []Session

type CreateSessionResponse Session

type GetSessionResponse Session

type PatchSessionResponse Session
