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
	"os/exec"
	"strings"
	"time"

	tf "github.com/wamuir/graft/tensorflow"
)

type Prompt struct {
	model         *tf.SavedModel
	tokenizerName string
}

type Record struct {
	Text   string `json:"name"`
	Result string `json:"comment"`
}

type Args struct {
	ModelDir string
}

func main() {
	args := Args{ModelDir: "/app/models"}

	p := Prompt{
		tokenizerName: "Salesforce/codet5-base",
	}

	if err := p.loadModel(args.ModelDir); err != nil {
		log.Fatalf("fail to load model: %v", err)
	}

	defer func() {
		if err := p.close(); err != nil {
			log.Printf("fail to close model: %v", err)
		}
	}()

	if err := p.predictFromDataset(); err != nil {
		log.Fatalf("fail to predict from dataset: %v", err)
	}
}

func (p *Prompt) loadModel(modelDir string) error {
	var err error
	p.model, err = tf.LoadSavedModel(modelDir, []string{"serve"}, nil)
	if err != nil {
		return err
	}
	return nil
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

// Fonction pour décoder les tokens générés
func (p *Prompt) decodeTokens(tokenizerPath string, generatedTokens []float32) (string, error) {
	// Implémenter le décodage selon votre tokenizer spécifique.
	// Cette partie dépend de votre tokenizer utilisé pour l'encodage.

	// Exemple simplifié :
	decodedText := "Decoded text from tokens"
	return decodedText, nil
}

func (p *Prompt) waitTokenizerContainer() error {
	for {
		cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", "tokenizer_container")
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			return err
		}
		if out.String() == "true\n" {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

type tokenizeResp struct {
	InputIDs      [][]int32 `json:"input_ids"`
	AttentionMask [][]int32 `json:"attention_mask"`
}

func (p *Prompt) tokenizeWithPython(text string) (tokenizeResp, error) {
	var tokenize tokenizeResp
	/*if err := p.waitTokenizerContainer(); err != nil {
		return tokenize, err
	}*/

	cmd := exec.Command("docker", "exec", "tokenizer_container", "python3", "tokenizer_script.py",
		text, p.tokenizerName)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return tokenize, err
	}

	if err := json.Unmarshal(out.Bytes(), &tokenize); err != nil {
		return tokenize, err
	}

	log.Printf("tokenize found: %+v", tokenize)

	return tokenize, nil
}

func convertFloat32ArrayToString(arr [][][]float32) string {
	var sb strings.Builder
	for _, sentence := range arr {
		for _, word := range sentence {
			for _, char := range word {
				sb.WriteString(fmt.Sprintf("%c", int(char)))
			}
		}
	}
	return sb.String()
}

func (p *Prompt) close() error {
	return p.model.Session.Close()
}
