package comments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type LocalAIConfig struct {
	Active          *bool  `yaml:"active"`
	URL             string `yaml:"url"`
	APIModelVersion int    `yaml:"api_model_version"`
}

type localAI struct {
	LocalAIConfig
}

func (o *localAI) isActive() bool {
	if o.Active == nil || *o.Active == false {
		return false
	}
	if o.URL == "" {
		log.Fatal("please set the local API URL in the localai-url variable")
	}
	if o.APIModelVersion == 0 {
		o.APIModelVersion = 1
	}
	return true
}

func (o *localAI) commentFunc(fn *ast.FuncDecl) (string, error) {
	return o.callLocalAI(GenerateFuncCode(fn))
}

type TokenizeRequest struct {
	Text    string `json:"text"`
	Version int    `json:"version"`
}

type TokenizeResponse struct {
	Comment string `json:"comment"`
}

func (o *localAI) callLocalAI(functionCode string) (string, error) {
	requestBody, err := json.Marshal(TokenizeRequest{
		Text:    functionCode,
		Version: o.APIModelVersion,
	})
	if err != nil {
		return "", err
	}

	for {
		resp, err := http.Get(o.URL + "/ping")
		if err != nil {
			log.Printf("fail to ping tokenizer API: %+v", err)
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		time.Sleep(10 * time.Second)
	}

	resp, err := http.Post(o.URL+"/tokenize", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("fail to close reponse: %+v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to call tokenizer API: status code %d", resp.StatusCode)
	}

	var tokenizeResponse TokenizeResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &tokenizeResponse); err != nil {
		return "", err
	}

	return addDoubleSlash(tokenizeResponse.Comment), nil
}

// addDoubleSlash adds "// " at the beginning of each line in the input string.
// // If the last line is empty, it will not add the prefix and will not keep it.
// // Additionally, it adds a period at the end of the last line if it does not already have one..
func addDoubleSlash(input string) string {
	input = strings.TrimSuffix(input, "\n")

	if !strings.HasSuffix(input, ".") {
		input += "."
	}

	// Split the input string by newlines to get each line separately.
	lines := strings.Split(input, "\n")

	// Loop through each line and prepend "// " to it.
	for i, line := range lines {
		lines[i] = "// " + line
	}

	// Join the modified lines back together with newlines.
	return strings.Join(lines, "\n")
}

func (o *localAI) commentConst(name string, exported bool) (string, error) {
	var exportedTxt string
	if !exported {
		exportedTxt = "private "
	}
	explainConst := convertVarToCamelCaseTo(name)
	txt := fmt.Sprintf("// %s is a %sconstant%s.", name, exportedTxt, explainConst)

	return txt, nil
}

func (o *localAI) commentVar(name, declType, explainVar string, exported bool) (string, error) {
	var exportedTxt string
	if !exported {
		exportedTxt = "private "
	}
	txt := fmt.Sprintf("// %s is a %svariable of type %s%s.", name, exportedTxt, declType, explainVar)

	return txt, nil
}

func (o *localAI) commentType(_ *ast.GenDecl) (string, error) {
	return "", nil
}
