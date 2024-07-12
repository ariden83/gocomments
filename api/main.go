package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	tf "github.com/wamuir/graft/tensorflow"
	"github.com/wamuir/graft/tensorflow/op"
	"log"
)

type prompt struct {
	inputTensor *tf.Tensor
	model       *tf.SavedModel
}

func main() {
	p := prompt{}
	p.loadModel()
	defer func() {
		if err := p.close(); err != nil {
			log.Printf("fail to close model: %v", err)
		}
	}()
	p.askQuestion()
}

func (p *prompt) loadModel() {
	s := op.NewScope()

	c := op.Const(s, "Hello from TensorFlow version "+tf.Version())
	graph, err := s.Finalize()
	if err != nil {
		panic(err)
	}
	// Execute the graph in a session.
	sess, err := tf.NewSession(graph, nil)
	if err != nil {
		panic(err)
	}
	output, err := sess.Run(nil, []tf.Output{c}, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(output[0].Value())

	// Charge le modèle
	model, err := tf.LoadSavedModel("/app/models", []string{"serve"}, nil)
	if err != nil {
		log.Fatalf("Erreur lors du chargement du modèle: %v", err)
	}
	defer model.Session.Close()
	p.model = model

	// Prépare les données d'entrée
	// Remplacez cela par vos propres données et la forme d'entrée
	inputTensor, err := tf.NewTensor([]float32{1.0, 2.0, 3.0})
	if err != nil {
		log.Fatalf("Erreur lors de la création du tenseur: %v", err)
	}

	p.inputTensor = inputTensor
}

func (p *prompt) askQuestion() {
	var name string
	prompt := &survey.Input{
		Message: "Enter your name:",
	}

	err := survey.AskOne(prompt, &name)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Hello, %s!\n", name)

	// Exécute la session
	output, err := p.model.Session.Run(
		map[tf.Output]*tf.Tensor{
			p.model.Graph.Operation("input_node_name").Output(0): p.inputTensor,
		},
		[]tf.Output{
			p.model.Graph.Operation("output_node_name").Output(0),
		},
		nil,
	)
	if err != nil {
		log.Fatalf("Erreur lors de l'exécution de la session: %v", err)
	}

	// Affiche le résultat
	fmt.Println(output[0].Value())

}

func (p *prompt) close() error {
	return p.model.Session.Close()
}
