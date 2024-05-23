package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"io"
	"log"
	"net/http"
)

func newAnthropic() commentsProcess {
	return &anthropic{}
}

type AnthropicConfig struct {
	// Do we use Anthropic Claude to generate function comments
	Active       *bool   `yaml:"active"`
	URL          string  `yaml:"url"`
	SSORegion    *string `yaml:"sso-region"`
	AccessKey    *string `yaml:"access-key"`
	SecretKey    *string `yaml:"secret-key"`
	SessionToken *string `yaml:"session-token"`
}

type anthropic struct {
	AnthropicConfig
}

func (a *anthropic) isActive() bool {
	if a.Active == nil || *a.Active == false {
		return false
	}
	if a.URL == "" {
		log.Fatal("Please set your AWS SSO Region in the anthropic-url variable.")
	}
	if a.AccessKey == nil || *a.AccessKey == "" {
		log.Fatal("Please set your AWS SSO Region in the anthropic-access-key variable.")
	}
	return true
}

func (a *anthropic) commentFunc(fn *ast.FuncDecl) (string, error) {
	functionCode := generateFuncCode(fn)
	var funcComment string

	payload := RequestPayload{
		Prompt:      fmt.Sprintf("Generate a detailed comment in English for the following Go function:\n%s", functionCode),
		MaxTokens:   150,
		Model:       "claude-v1",
		Temperature: 0.7,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return funcComment, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", a.URL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return funcComment, fmt.Errorf("error creating request: %v", err)
	}

	// Set the appropriate headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.AccessKey))

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return funcComment, fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("fail to close reader")
		}
	}(resp.Body)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return funcComment, fmt.Errorf("error reading response body: %v", err)
	}

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return funcComment, fmt.Errorf("request failed with status %d: %s\n", resp.StatusCode, string(body))
	}

	// Parse the response
	var responsePayload ResponsePayload
	err = json.Unmarshal(body, &responsePayload)
	if err != nil {
		return funcComment, fmt.Errorf("error unmarshalling response: %v", err)
	}

	// Print the response
	return responsePayload.Completion, nil
}

func (a *anthropic) commentConst(string, bool) (string, error) {
	var constComment string
	return constComment, nil
}

func (a *anthropic) commentVar(name, declType, explainVar string, exported bool) (string, error) {
	var varComment string
	return varComment, nil
}

func (a *anthropic) commentType(genDecl *ast.GenDecl) (string, error) {
	var typeComment string
	return typeComment, nil
}
