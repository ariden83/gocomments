package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Prompt struct {
	currentModelVersion int
	datasetDirectory    string
	datasetPrefix       string
	datasetExtension    string
}

type Record struct {
	Text   string `json:"name"`
	Result string `json:"comment"`
}

func main() {
	p := Prompt{
		currentModelVersion: 10,
		datasetDirectory:    "/app/dataset/",
		datasetPrefix:       "functions_dataset",
		datasetExtension:    ".jsonl",
	}

	if err := p.predictFromDataset(); err != nil {
		log.Fatalf("fail to predict from dataset: %v", err)
	}
}

func (p *Prompt) getLastDatasetFile() (string, error) {
	files, err := ioutil.ReadDir(p.datasetDirectory)
	if err != nil {
		return "", err
	}

	var latestFile string
	var latestTime time.Time

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), p.datasetPrefix) && strings.HasSuffix(file.Name(), p.datasetExtension) {
			fileTime := file.ModTime()
			if fileTime.After(latestTime) {
				latestTime = fileTime
				latestFile = file.Name()
			}
		}
	}

	if latestFile == "" {
		return "", fmt.Errorf("no dataset files found in directory %s with prefix %s and extension %s", p.datasetDirectory, p.datasetPrefix, p.datasetExtension)
	}

	return filepath.Join(p.datasetDirectory, latestFile), nil
}

func (p *Prompt) predictFromDataset() error {
	latestFile, err := p.getLastDatasetFile()
	if err != nil {
		return fmt.Errorf("could not get latest dataset file: %v", err)
	}

	file, err := os.Open(latestFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("err: %v", err)
		}
	}()

	start := 0
	limit := start + 50
	i := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i++
		if i < start {
			continue
		}
		if i > limit {
			break
		}
		line := scanner.Text()
		var record Record
		if err := json.Unmarshal([]byte(line), &record); err != nil {
			log.Printf("error deserializing the row: %v", err)
			continue
		}

		fmt.Println(strings.Repeat("#", 25))
		fmt.Println("QUERY:", record.Text)
		fmt.Println(strings.Repeat("#", 25))
		fmt.Println("ORIGINAL:")
		fmt.Println(record.Result)

		j := 1
		for j <= p.currentModelVersion {
			start := time.Now()
			commentPredicted, err := p.runPredict(record.Text, j)
			if err != nil {
				return err
			}
			elapsed := time.Since(start) //
			fmt.Println(strings.Repeat("#", 25))
			fmt.Println(fmt.Sprintf("GENERATED (took %d ms) with version %d of model", elapsed.Milliseconds(), j))
			fmt.Println(commentPredicted)
			j++
		}
	}

	return scanner.Err()
}

type TokenizeRequest struct {
	Text    string `json:"text"`
	Version int    `json:"version"`
}

type TokenizeResponse struct {
	Comment string `json:"comment"`
}

func (p *Prompt) runPredict(text string, version int) (string, error) {
	requestBody, err := json.Marshal(TokenizeRequest{
		Text:    text,
		Version: version,
	})
	if err != nil {
		return "", err
	}

	// Vérifier que le conteneur est en cours d'exécution
	for {
		resp, err := http.Get("://:5000/ping")
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

	resp, err := http.Post("://:5000/tokenize", "application/json", bytes.NewBuffer(requestBody))
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &tokenizeResponse); err != nil {
		return "", err
	}

	return tokenizeResponse.Comment, nil
}
