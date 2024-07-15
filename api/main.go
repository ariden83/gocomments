package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	tf "github.com/wamuir/graft/tensorflow"
)

type Prompt struct {
	Model *tf.SavedModel
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

	p := Prompt{}
	if err := p.loadModel(args.ModelDir); err != nil {
		log.Fatalf("fail to load model: %v", err)
	}

	defer func() {
		if err := p.close(); err != nil {
			log.Printf("fail to close model: %v", err)
		}
	}()

	if err := p.predictFromDataset(args); err != nil {
		log.Fatalf("fail to predict from dataset: %v", err)
	}
}

func (p *Prompt) loadModel(modelDir string) error {
	var err error
	p.Model, err = tf.LoadSavedModel(modelDir, []string{"serve"}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *Prompt) predictFromDataset(args Args) error {
	file, err := os.Open("/app/dataset/functions_dataset_20240624_12.jsonl")
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("err: %v", err)
		}
	}()

	records := []Record{}
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
	code := records[index].Result

	decodedCode, err := p.runPredict(text)
	if err != nil {
		return err
	}

	fmt.Println(strings.Repeat("#", 25))
	fmt.Println("QUERY:", text)
	fmt.Println()
	fmt.Println(strings.Repeat("#", 25))
	fmt.Println("ORIGINAL:")
	fmt.Println()
	fmt.Println(code)
	fmt.Println()
	fmt.Println(strings.Repeat("#", 25))
	fmt.Println("GENERATED:")
	fmt.Println()
	fmt.Println(decodedCode)

	return nil
}

func (p *Prompt) runPredict(text string) (string, error) {
	// Prétraitement du texte
	inputIds := preprocessText(text)

	// Convertir inputIds en un tableau de int32
	inputIdsInt32 := make([]int32, len(inputIds))
	for i, id := range inputIds {
		inputIdsInt32[i] = int32(id)
	}

	// Créer le tensor d'entrée pour TensorFlow
	inputTensor, err := tf.NewTensor([][]int32{inputIdsInt32})
	if err != nil {
		return "", fmt.Errorf("error creating the tensor: %v", err)
	}

	// Créer un tensor pour les masques d'attention (assume all-ones mask)
	attentionMask := make([]int32, len(inputIdsInt32))
	for i := range attentionMask {
		attentionMask[i] = 1
	}
	maskTensor, err := tf.NewTensor([][]int32{attentionMask})
	if err != nil {
		return "", fmt.Errorf("error creating the attention mask tensor: %v", err)
	}

	// Exécuter le modèle sur l'entrée
	result, err := p.Model.Session.Run(
		map[tf.Output]*tf.Tensor{
			p.Model.Graph.Operation("serving_default_input_ids").Output(0):              inputTensor,
			p.Model.Graph.Operation("serving_default_attention_mask").Output(0):         maskTensor,
			p.Model.Graph.Operation("serving_default_decoder_input_ids").Output(0):      inputTensor,
			p.Model.Graph.Operation("serving_default_decoder_attention_mask").Output(0): maskTensor,
		},
		[]tf.Output{
			p.Model.Graph.Operation("StatefulPartitionedCall").Output(0),
		},
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("error running model: %v", err)
	}

	generatedCode, ok := result[0].Value().([][][]float32)
	if !ok {
		return "", fmt.Errorf("error converting result: unexpected type %T", result[0].Value())
	}

	decodedText := convertFloat32ArrayToString(generatedCode)
	return decodedText, nil
}

func preprocessText(text string) []int {
	// Appliquer le prétraitement nécessaire pour convertir la chaîne de caractères en une liste d'IDs
	// Exemple très simpliste, en pratique cela doit être plus sophistiqué
	// Vous devrez ajuster cette partie en fonction de votre modèle
	words := strings.Fields(text)
	inputIds := make([]int, len(words))
	for i, _ := range words {
		// Ceci est un exemple de mappage des mots en IDs fictifs. Remplacez ceci par votre propre tokenizer.
		inputIds[i] = i + 1 // Vous aurez besoin d'un vrai tokenizer ici.
	}
	return inputIds
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
	return p.Model.Session.Close()
}
