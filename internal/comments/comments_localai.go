package comments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"io"
	"log"
	"net/http"
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
			time.Sleep(1 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			break
		}
		log.Printf("wait initiliazation of tokenizer_container")
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

	return tokenizeResponse.Comment, nil
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
