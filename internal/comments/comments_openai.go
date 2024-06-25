package comments

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/ast"
	"io"
	"log"
	"moul.io/http2curl"
	"net/http"
)

type OpenAIConfig struct {
	// Do we use OPENAI to generate function comments
	Active *bool   `yaml:"active"`
	APIKey *string `yaml:"api_key"`
	URL    string  `yaml:"url"`
}

type openAI struct {
	OpenAIConfig
}

func (o *openAI) isActive() bool {
	if o.Active == nil || *o.Active == false {
		return false
	}
	if o.APIKey == nil || *o.APIKey == "" {
		log.Fatal("Please set your OpenAI API key in the openai-api_key variable.")
	}
	if o.URL == "" {
		log.Fatal("Please set the OpenAI API URL in the openai-url variable.")
	}
	return true
}

func (o *openAI) commentFunc(fn *ast.FuncDecl) (string, error) {
	return o.callOpenAI(GenerateFuncCode(fn))
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (o *openAI) callOpenAI(functionCode string) (string, error) {
	prompt := fmt.Sprintf("Generate a detailed comment in English for the following Go function. The comment should be written in a way that is helpful for other developers. Include the purpose of the function, a description of its parameters and return values, potential error conditions, and any side effects or important details. Here is the function :\n%s", functionCode)

	requestBody, err := json.Marshal(map[string]interface{}{
		"max_tokens":  150,
		"temperature": 0.7,
		"model":       "gpt-3.5-turbo",
		"messages": []OpenAIMessage{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: prompt},
		},
	})

	if err != nil {
		return "", fmt.Errorf("error creating request body: %v", err)
	}

	req, err := http.NewRequest("POST", o.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+*o.APIKey)

	client := &http.Client{}

	command, _ := http2curl.GetCurlCommand(req)
	fmt.Println(fmt.Sprintf("%s", command))

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("error decoding response: %v", err)
	}

	if choices, ok := response["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if text, ok := choice["text"].(string); ok {
				return text, nil

			} else {
				return "", errors.New("error: no text found in response choice")
			}
		} else {
			return "", errors.New("error: invalid choice format")
		}
	}

	return "", errors.New("error: no choices found in response")
}

func (o *openAI) commentConst(string, bool) (string, error) {
	var constComment string
	return constComment, nil
}

func (o *openAI) commentVar(name, declType, explainVar string, exported bool) (string, error) {
	var varComment string
	return varComment, nil
}

func (o *openAI) commentType(genDecl *ast.GenDecl) (string, error) {
	var typeComment string
	return typeComment, nil
}
