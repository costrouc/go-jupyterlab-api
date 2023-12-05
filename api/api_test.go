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

func TestGetStatus(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	_, err = client.GetStatus(ctx)
	if err != nil {
		t.Error(err)
	}
}

func TestGetMe(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	_, err = client.GetMe(ctx)
	if err != nil {
		t.Error(err)
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

func TestGetKernelSpecs(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.GetKernelSpecs(ctx)
	if err != nil {
		t.Error(err)
	}
	if data.Default != "python3" {
		t.Errorf("Expected default kernelspec to be python3, got %s", data.Default)
	}
}

func TestCreateListGetDeleteKernels(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	createData, err := client.CreateKernel(ctx, CreateKernelBody{Name: "python3"})
	if err != nil {
		t.Error(err)
	}
	id := createData.Id

	listKernels, err := client.GetKernels(ctx)
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, kernel := range *listKernels {
		if kernel.Id == id {
			found = true
		}
	}
	if !found {
		t.Errorf("Expecting to find kernel with id %s didn't find", id)
	}

	getData, err := client.GetKernel(ctx, id)
	if err != nil {
		t.Error(err)
	}
	if getData.Id != id {
		t.Errorf("Expecting to find kernel with id %s didn't find, got %s", id, getData.Id)
	}

	err = client.DeleteKernel(ctx, id)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateListGetDeleteTerminals(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "faketoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	createData, err := client.CreateTerminal(ctx)
	if err != nil {
		t.Error(err)
	}
	name := createData.Name

	listTerminals, err := client.GetTerminals(ctx)
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, terminal := range *listTerminals {
		if terminal.Name == name {
			found = true
		}
	}
	if !found {
		t.Errorf("Expecting to find terminal with name %s didn't find", name)
	}

	getData, err := client.GetTerminal(ctx, name)
	if err != nil {
		t.Error(err)
	}
	if getData.Name != name {
		t.Errorf("Expecting to find kernel with id %s didn't find, got %s", name, getData.Name)
	}

	err = client.DeleteTerminal(ctx, name)
	if err != nil {
		t.Error(err)
	}
}
