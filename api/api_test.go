package api

import (
	"context"
	"regexp"
	"testing"
)

func TestGetVersion(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.GetVersion(ctx)
	if err != nil {
		t.Error(err)
	}
	if matched, _ := regexp.Match("[0-9]+\\.[0-9]+\\.[0-9]+", []byte(data.Version)); !matched {
		t.Errorf("Version did not match regex '[0-9]+\\.[0-9]+\\.[0-9]+', got %v", data.Version)
	}
}

func TestGetContents(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.GetContents(ctx, "/", nil)
	if err != nil {
		t.Error(err)
	}
	if data.Type != "directory" {
		t.Errorf("Root content returned should be directory, got %v", data.Type)
	}
}
