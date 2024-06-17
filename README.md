# gocomments

![languague: Go](https://img.shields.io/badge/language-go-007d9c)

`gocomments` is a tool to update to automatically adding comments to go functions and variables to code such as

```
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const openaiAPIURL = "https://api.openai.com/v1/engines/davinci-codex/completions"

func init() {
	fmt.Println("Init func")
}

func InitDatabase() {}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	go myLongRunningProcess(ctx)

	<-ctx.Done()

	fmt.Println("Processus terminé ou interrompu.")
}

func GetEventStatus(eventID string, status string) error {
	// Implementation goes here
	return nil
}

func SetEventStatus(eventID string, status string) error {
	// Implementation goes here
	return nil
}

func myLongRunningProcess(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():

			fmt.Println("Processus annulé.")
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

type toto int

type tata toto

var (
	Var1 int

	Var2 string

	Var3 toto
)

var Var4 string

var var5 string

const var6 = "tata"

const Var7 = "toto"

type Test struct {
	Key  string
	KeyB int
	KeyC toto
	KeyD *toto
	KeyE *int
	keyg bool
}

func (t Test) Tota() {}

func Totbczzsd() {}

func (t Test) Totc(ctx context.Context, value string) (string, error) {
	return "test", nil
}

func Totd(ctx context.Context, value string) (string, error) {
	return "test", nil
}

type Config struct {
}

type Adapter struct {
}

func New(config Config, logger log.Logger) (Adapter, error) {
	// Implementation goes here.
	var5 = "toto"
	Var4 = "tata"
	Var1 = 1
	Var2 = "test"
	Var3 = 1
	u := tata(5)
	log.Printf("%s %s %d", var6, Var7, u)
	return Adapter{}, nil
}

func IsValid(token string) bool {
	return true
}

func HasPermission(user string) bool {
	return true
}

```

to give this result

```
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// openaiAPIURL is a private constant that indicates the endpoint URL for accessing to openai api.
const openaiAPIURL = "https://api.openai.com/v1/engines/davinci-codex/completions"

func init() {
	fmt.Println("Init func")
}

// InitDatabase is a method to initializes the database.
// It does not take any arguments.
func InitDatabase()	{}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		cancel()
	}()

	go myLongRunningProcess(ctx)

	<-ctx.Done()

	fmt.Println("Processus terminé ou interrompu.")
}

// GetEventStatus is a method that retrieve the event status that take an eventID of type string, a status of type string
// It's return an error if fails, otherwise nil.
//
// Example:
//   err := GetEventStatus(my-string, my-string)
//   if err != nil {
//       log.Fatalf("Error: %v", err)
//   }
func GetEventStatus(eventID string, status string) error {

	return nil
}

// SetEventStatus is a method which update the event status that take an eventID of type string, a status of type string
// It's return an error if fails, otherwise nil.
//
// Example:
//   err := SetEventStatus(my-string, my-string)
//   if err != nil {
//       log.Fatalf("Error: %v", err)
//   }
func SetEventStatus(eventID string, status string) error {

	return nil
}

// myLongRunningProcess is a private method which execute my long running process that take a ctx of type context.Context
// It's return an error if fails, otherwise nil.
//
// Example:
//   ctx := context.Background()
//   err := myLongRunningProcess(ctx)
//   if err != nil {
//       log.Fatalf("Error: %v", err)
//   }
func myLongRunningProcess(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():

			fmt.Println("Processus annulé.")
			return nil
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

// toto is a type alias for the int type.
// It allows you to create a new type with the same
// underlying type as int, but with a different name.
// This can be useful for improving code readability
// and providing more semantic meaning to your types.
type toto int

// tata is a type alias for the toto type.
// It allows you to create a new type with the same
// underlying type as int, but with a different name.
// This can be useful for improving code readability
// and providing more semantic meaning to your types.
type tata toto

var (
	// Var1 is a variable of type int.
	Var1	int

	// Var2 is a variable of type string.
	Var2	string

	// Var3 is a variable of type toto.
	Var3	toto
)

// Var4 is a variable of type string.
var Var4 string

// var5 is a private variable of type string.
var var5 string

// var6 is a private constant.
const var6 = "tata"

// Var7 is a constant.
const Var7 = "toto"

// Test represents a structure for 
// It contains information about a Key, a KeyB, a KeyC, a private keyg and
// optionally a KeyD, a KeyE.
type Test struct {
	Key	string
	KeyB	int
	KeyC	toto
	KeyD	*toto
	KeyE	*int
	keyg	bool
}

// Tota is a method that belongs to the Test struct.
// It does not take any arguments.
func (t Test) Tota()	{}

// Totbczzsd is a method.
// It does not take any arguments.
func Totbczzsd()	{}

// Totc is a method that belongs to the Test struct that take a ctx of type context.Context, a value of type string
// and returns a string
// It's return an error if fails, otherwise nil.
//
// Example:
//   ctx := context.Background()
//   str, err := Totc(ctx, my-string)
//   if err != nil {
//       log.Fatalf("Error: %v", err)
//   }
//   fmt.Printf("%v ", str)
func (t Test) Totc(ctx context.Context, value string) (string, error) {
	return "test", nil
}

// Totd is a method that take a ctx of type context.Context, a value of type string
// and returns a string
// It's return an error if fails, otherwise nil.
//
// Example:
//   ctx := context.Background()
//   str, err := Totd(ctx, my-string)
//   if err != nil {
//       log.Fatalf("Error: %v", err)
//   }
//   fmt.Printf("%v ", str)
func Totd(ctx context.Context, value string) (string, error) {
	return "test", nil
}

// Config represents a structure for .
type Config struct {
}

// Adapter represents a structure for .
type Adapter struct {
}

// New creates a new instance of Adapter.
// It initializes the Adapter with the provided config of type Config, logger of type log.Logger.
// It's return an error if the initialization fails, otherwise nil.
func New(config Config, logger log.Logger) (Adapter, error) {

	var5 = "toto"
	Var4 = "tata"
	Var1 = 1
	Var2 = "test"
	Var3 = 1
	u := tata(5)
	log.Printf("%s %s %d", var6, Var7, u)
	return Adapter{}, nil
}

// IsValid is a method to check the valid that take a token of type string
// and returns a bool.
//
// Example:
//   valid := IsValid(my-string)
//   fmt.Printf("%v ", valid)
func IsValid(token string) bool {
	return true
}

// HasPermission is a method to check the permission that take an user of type string
// and returns a bool.
//
// Example:
//   valid := HasPermission(my-string)
//   fmt.Printf("%v ", valid)
func HasPermission(user string) bool {
	return true
}

```

## Installation

Command to generate the binary :

> make build

```shell
$ go install ./
```

## Usage

```text
Usage: gocomments [flags] [path ...]
  -d	display diffs instead of rewriting files
  -l	list files whose formatting differs from goimport's
  -local string
    	put imports beginning with this string after 3rd-party package
  -prefix value
    	relative local prefix to from a new import group (can be given several times)
  -w	write result to (source) file instead of stdout
```

By default, `gocomments` will search the module name from the root go.mod file as local import path.

We can add a `.gocomments` file in an intermediate directory to override the configuration
for all subdirectories of this path. The final configuration to process a file is computed
with all intermediate configuration files.

The file follow the format of this example below:

```yaml
---
# local is the root Go module name. All subpackages of this module
# will be separated from the external packages.
# By default, use the value of the module name in the root go.mod file.
local: ""

# prefixes is a list of relative Go packages from the root package.
# All imports with these prefixes will be separated from each other.
prefixes:
  - "common"

# Signature is the personal signature affixed at the end of the comment in order to allow users and
# other iterations to find the autogenerated comments and to be able to evolve
# them with future updates of the script.
# If empty, the automatically added prefix is "auto".
signature: "AutoComBOT"
# Allows you to know if you update the tagged comments each time the script is executed.
update-comments: false
# Allows to add examples in comments.
active-examples: false
# Do we use OPENAI to generate function comments
openai-active: false
# If openai-active is active, set your OPENAI api key. 
openai-api_key: ""
# If openai-active is active, set your OPENAI api URL. 
openai-url: "https://api.openai.com/v1/chat/completions"
# Do we use Anthropic Claude to generate function comments
anthropic-active: false
```

## In progress

#### A) Generate comments using [Open AI](https://openai.com/index/introducing-chatgpt-and-whisper-apis/)
#### B) Generate comments using [AWS AI](https://aws.amazon.com/fr/ai/)
#### C) Generate comments using our own private AI

- Generate the model with this command :

> make generate-model

## Resources

- This notebook demonstrates how to generate python code based on the natural language description using TensorFlow, HuggingFace and Mostly Basic Python Programming Benchmark from Google Research. [Text-to-code Generation with TensorFlow, 🤗 & MBPP](https://www.kaggle.com/code/rhtsingh/text-to-code-generation-with-tensorflow-mbpp)
- Common Crawl maintains a free, open repository of web crawl data that can be used by anyone. [commoncrawl](https://commoncrawl.org/)
- [tensorflow > Catalog > c4](https://www.tensorflow.org/datasets/catalog/c4?hl=fr) 
- CodeSearchNet is a collection of datasets and benchmarks that explore the problem of code retrieval using natural language. This research is a continuation of some ideas presented in this blog post and is a joint collaboration between GitHub and the Deep Program Understanding group at Microsoft Research - Cambridge. [CodeSearchNet](https://github.com/github/CodeSearchNet)

