package api

import (
	"context"
	"regexp"
	"strings"
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
	data, err := client.GetContents(ctx, "", nil)
	if err != nil {
		t.Error(err)
	}
	if data.Type != "directory" {
		t.Errorf("Root content returned should be directory, got %v", data.Type)
	}
}

func TestCreateContents(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.CreateContents(ctx, "", &CreateContentsBody{Ext: ".test.txt"})
	if err != nil {
		t.Error(err)
	}
	if !strings.HasSuffix(data.Name, ".test.txt") {
		t.Errorf("Created new file %s does not have extension '.test.txt', got %s", data.Name, data.Name)
	}
}

func TestCreateRenameContents(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	createData, err := client.CreateContents(ctx, "", &CreateContentsBody{Ext: ".test.txt"})
	if err != nil {
		t.Error(err)
	}
	if !strings.HasSuffix(createData.Name, ".test.txt") {
		t.Errorf("Created new file %s does not have extension '.test.txt', got %s", createData.Name, createData.Name)
	}

	patchData, err := client.PatchContents(ctx, createData.Name, &PatchContentsBody{Path: "file.txt"})
	if err != nil {
		t.Error(err)
	}
	if patchData.Name != "file.txt" {
		t.Errorf("Renamed file %s does not have name %s", patchData.Name, patchData.Name)
	}
}

func TestPutContents(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.PutContents(ctx, "hello.txt", &PutContentsBody{Content: "hello world", Format: "text", Type: "file"})
	if err != nil {
		t.Error(err)
	}
	if data.Name != "hello.txt" {
		t.Errorf("Created new file does not have correct name hello.txt, got %s", data.Name)
	}
}

func TestCreateDeleteContents(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	createData, err := client.CreateContents(ctx, "", &CreateContentsBody{Ext: ".test.txt"})
	if err != nil {
		t.Error(err)
	}
	if !strings.HasSuffix(createData.Name, ".test.txt") {
		t.Errorf("Created new file %s does not have extension '.test.txt', got %s", createData.Name, createData.Name)
	}

	err = client.DeleteContents(ctx, createData.Name)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateListGetDeleteSessions(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	id := "fe03e6be-47e0-407d-8c3f-46acc6b18d9d"
	createData, err := client.CreateSession(ctx, &Session{Id: id, Kernel: map[string]string{"name": "python3"}})
	if err != nil {
		t.Error(err)
	}
	if createData.Id != id {
		t.Errorf("Expected created session to have id %s", id)
	}

	listData, err := client.GetSessions(ctx)
	if err != nil {
		t.Error(err)
	}

	foundSession := false
	for _, session := range *listData {
		if session.Id == id {
			foundSession = true
		}
	}
	if !foundSession {
		t.Errorf("Session %s not found in list", id)
	}

	getData, err := client.GetSession(ctx, id)
	if err != nil {
		t.Error(err)
	}
	if getData.Id != id {
		t.Errorf("Session %s not found", id)
	}

	err = client.DeleteSession(ctx, id)
	if err != nil {
		t.Error(err)
	}
}
