package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type Prompt struct{}

type Record struct {
	Text   string `json:"name"`
	Result string `json:"comment"`
}

func main() {
	p := Prompt{}

	if err := p.predictFromDataset(); err != nil {
		log.Fatalf("fail to predict from dataset: %v", err)
	}
}

func (p *Prompt) predictFromDataset() error {
	file, err := os.Open("/app/dataset/functions_dataset_20240624_12.jsonl")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("err: %v", err)
		}
	}()

	var records []Record
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		var record Record
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			log.Printf("error deserializing the row: %v", err)
			continue
		}
		records = append(records, record)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(records))
	text := records[index].Text
	commentOriginal := records[index].Result

	commentPredicted, err := p.runPredict(text)
	if err != nil {
		return err
	}

	fmt.Println(strings.Repeat("#", 25))
	fmt.Println("QUERY:", text)
	fmt.Println(strings.Repeat("#", 25))
	fmt.Println("ORIGINAL:")
	fmt.Println(commentOriginal)
	fmt.Println(strings.Repeat("#", 25))
	fmt.Println("GENERATED:")
	fmt.Println(commentPredicted)

	return nil
}

type TokenizeRequest struct {
	Text    string `json:"text"`
	Version int    `json:"version"`
}

type TokenizeResponse struct {
	Comment string `json:"comment"`
}

func (p *Prompt) runPredict(text string) (string, error) {
	requestBody, err := json.Marshal(TokenizeRequest{
		Text:    text,
		Version: 9,
	})
	if err != nil {
		return "", err
	}

	// Vérifier que le conteneur est en cours d'exécution
	for {
		resp, err := http.Get("http://tokenizer_container:5000/ping")
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

	resp, err := http.Post("http://tokenizer_container:5000/tokenize", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to call tokenizer API: status code %d", resp.StatusCode)
	}

	var tokenizeResponse TokenizeResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &tokenizeResponse); err != nil {
		return "", err
	}

	return tokenizeResponse.Comment, nil
}
