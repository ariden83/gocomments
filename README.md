# gocomments

![languague: Go](https://img.shields.io/badge/language-go-007d9c)

`gocomments` is a tool to update to automatically adding comments to go functions and variables to code such as

```go
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

```go
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

> make install

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

example : 

```shell
gocomments -l -w  ./test/.
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

# Do we use local AI to generate function comments
localai-active: false
# If local-active is active, set your local api URL. 
localai-url: "http://tokenizer_container:5000"
# defined the version of the model to use. 
localai-api_model_version: 10

# Do we use OPENAI to generate function comments
openai-active: false
# If openai-active is active, set your OPENAI api key. 
openai-api_key: ""
# If openai-active is active, set your OPENAI api URL. 
openai-url: "https://api.openai.com/v1/chat/completions"

# Do we use Anthropic Claude to generate function comments
anthropic-active: false
# If anthropic-active is active, set the anthropic api URL. 
anthropic-url: "https://api.anthropic.com/v1/complete"
```

## In progress

#### A) Generate comments using [Open AI](https://openai.com/index/introducing-chatgpt-and-whisper-apis/)
#### B) Generate comments using [AWS AI](https://aws.amazon.com/fr/ai/)
#### C) Generate comments using our own private AI

##### 1) Generate dataset file :

define your env variables with creating a .env file in dataset folder such as : 

```
LOCAL_REPO_PATH=/home/ariden/go/src/github.com/ariden/blockchain
```

and generate your dataset file

> make generate-dataset

with generate a dataset file such as : 

```[
    {
        "name": "func NewSetParcelSenderCustomShippingFunc(parcelAdapter parceladapter.Adapter, addTrackingEventByParcelIDFunc AddTrackingEventByParcelIDFunc, logger logging.Logger, monitorer monitor.Monitorer) (SetParcelSenderCustomShippingFunc) {\n    // Function implementation\n\t}",
        "comment": "NewSetParcelSenderCustomShippingFunc creates a new SetParcelSenderCustomShippingFunc.\n"
    },
    {
        "name": "func NewGetParcelWithEventsFunc(parcelAdapter parceladapter.Adapter, pickupClient pickupclient.Client, trackingURLBuilderFunc TrackingURLBuilderFunc, logger logging.Logger, monitorer monitor.Monitorer) (GetParcelWithEventsFunc) {\n    // Function implementation\n\t}",
        "comment": "NewGetParcelWithEventsFunc creates a new GetParcelWithEventsFunc.\n"
    }]
```

##### 2) Generate a new model

Generate a new model from your previous dataset file.

> make generate-model

##### 3) Test models

It is possible to test the final model or the different versions of temporary models created.

> make generate-test

Example of automatically generated comment

```
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func NewColor(c unknown) (Color)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | NewColor compute the ANSI color string according to the given list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16819 ms) with version 1 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its category
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 17100 ms) with version 2 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 19073 ms) with version 3 of model
go-tf-app_1  | NewColor compute the ANSI color string according to an optional list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 17545 ms) with version 4 of model
go-tf-app_1  | NewColor compute the ANSI color string according to a list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 17401 ms) with version 5 of model
go-tf-app_1  | NewColor compute the ANSI color string according to list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16261 ms) with version 6 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its display name.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16341 ms) with version 7 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16968 ms) with version 8 of model
go-tf-app_1  | NewColor compute the ANSI color string according to given list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 18506 ms) with version 9 of model
go-tf-app_1  | NewColor compute the ANSI color string according to provided list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 17943 ms) with version 10 of model
go-tf-app_1  | NewColor compute the ANSI color string according to given list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func (c Color) String(s string) (string)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14601 ms) with version 1 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14316 ms) with version 2 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 12927 ms) with version 3 of model
go-tf-app_1  | String returns the given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14214 ms) with version 4 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13832 ms) with version 5 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14624 ms) with version 6 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13525 ms) with version 7 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 15110 ms) with version 8 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13433 ms) with version 9 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14504 ms) with version 10 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func (c Color) Stringf(format string, a unknown) (string)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 12685 ms) with version 1 of model
go-tf-app_1  | Stringf is the Sprintf method.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13094 ms) with version 2 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14865 ms) with version 3 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14729 ms) with version 4 of model
go-tf-app_1  | Stringf is the colorized version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14379 ms) with version 5 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14649 ms) with version 6 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14015 ms) with version 7 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13884 ms) with version 8 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14316 ms) with version 9 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14110 ms) with version 10 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func NewColor(c unknown) (Color)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | NewColor compute the ANSI color string according to the given list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 15172 ms) with version 1 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its category
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16847 ms) with version 2 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 17914 ms) with version 3 of model
go-tf-app_1  | NewColor compute the ANSI color string according to an optional list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 17731 ms) with version 4 of model
go-tf-app_1  | NewColor compute the ANSI color string according to a list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 19750 ms) with version 5 of model
go-tf-app_1  | NewColor compute the ANSI color string according to list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 18197 ms) with version 6 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its display name.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16940 ms) with version 7 of model
go-tf-app_1  | NewColor compute the ANSI color string according to its display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 18295 ms) with version 8 of model
go-tf-app_1  | NewColor compute the ANSI color string according to given list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 17988 ms) with version 9 of model
go-tf-app_1  | NewColor compute the ANSI color string according to provided list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16900 ms) with version 10 of model
go-tf-app_1  | NewColor compute the ANSI color string according to given list of display attributes.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func (c Color) String(s string) (string)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13104 ms) with version 1 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13916 ms) with version 2 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13803 ms) with version 3 of model
go-tf-app_1  | String returns the given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13788 ms) with version 4 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14484 ms) with version 5 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13716 ms) with version 6 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14357 ms) with version 7 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14453 ms) with version 8 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13670 ms) with version 9 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14892 ms) with version 10 of model
go-tf-app_1  | String returned the colorized given string.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func (c Color) Stringf(format string, a unknown) (string)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13168 ms) with version 1 of model
go-tf-app_1  | Stringf is the Sprintf method.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14263 ms) with version 2 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 14317 ms) with version 3 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 18980 ms) with version 4 of model
go-tf-app_1  | Stringf is the colorized version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16932 ms) with version 5 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 15432 ms) with version 6 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16442 ms) with version 7 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16261 ms) with version 8 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 15585 ms) with version 9 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16498 ms) with version 10 of model
go-tf-app_1  | Stringf is the Sprintf version of String.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func ReadTargetStream(r io.Reader, processTargetFunc func) (error)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | ReadTargetStream reads the given proto stream and call the given callback for each target.
go-tf-app_1  | Stop and return the first occurred error.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 18890 ms) with version 1 of model
go-tf-app_1  | ReadTargetStream reads the io.Reader and returns an error if one is not found
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 18363 ms) with version 2 of model
go-tf-app_1  | ReadTargetStream reads the given io.reader and call bdd-proc
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 20142 ms) with version 3 of model
go-tf-app_1  | ReadTargetStream reads the given reader to process all target produced by this handler.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 19615 ms) with version 4 of model
go-tf-app_1  | ReadTargetStream reads the given stream reader and call a process to read data from it.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 15467 ms) with version 5 of model
go-tf-app_1  | WriteTargetStream calls the given function for each target.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 76219 ms) with version 6 of model
go-tf-app_1  | ReadTargetStream reads the given stream and call
go-tf-app_1  | the process for each target of a specific file. If no error occurs, returns nil otherwise to allow execution that will return an err where errors on read is not valid or if there's any occurred in this function can be used by another handler instead it would execute at init with status "read" method from here until end date has already been done source (mockedor) use stderr as temporary data does nothing right now but ignore some events were added into main object then do something else we should stop processing through more than one second since last tag could't go further even though you
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 15736 ms) with version 7 of model
go-tf-app_1  | ReadTargetStream reads the given stream.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 61143 ms) with version 8 of model
go-tf-app_1  | ReadTargetStream reads target stream and call the given callback for eachtarget.
go-tf-app_1  | Stop will return an error in case of failure to reach a transaction from ewallet-streamer service as might loop can then stop if there's any problem starting here until no one occurred during all process, we cannot wait io has ended or is not yet closed (either by end) processing go routine which may be called at least once per user creation time between runsners serveing them beforehand inside this client fetcher otherwise it creates
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 70297 ms) with version 9 of model
go-tf-app_1  | ReadTargetStream reads the given proto stream and call a target function for eachtarget of it.
go-tf-app_1  | Stopand return an error if one occurred while reading or unzipped by any issues, otherwise returns nil even when there was no more errors in case we can't process all files are blocked using io on incoming connections instead to serve http server right now but not handled yet resolved file content at user end... last modification found inside this topic: It's used with either `StartWithStaticMemory()` optionor` method which trigger only handles first encounteredError object here (no-premium id). If
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 16349 ms) with version 10 of model
go-tf-app_1  | ReadTargetStream reads the given proto stream and call a function
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | QUERY: func (f *FilesetEntry_SymlinkBehavior) UnmarshalJSON(b []byte) (error)
go-tf-app_1  | #########################
go-tf-app_1  | ORIGINAL:
go-tf-app_1  | Deprecated: Do not use.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 15920 ms) with version 1 of model
go-tf-app_1  | Deprecated: Use FilesetEntry_SymlinkBehavior.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13351 ms) with version 2 of model
go-tf-app_1  | Deprecated: Do not use.
go-tf-app_1  | 
go-tf-app_1  | #########################
go-tf-app_1  | GENERATED (took 13671 ms) with version 3 of model
go-tf-app_1  | Deprecated: Do not use.
go-tf-app_1  | 
```

##### 4) Launch API

Launches the API to use the created models to generate a comment.

> make generate-api

The API will be accessible via an HTTP call on the following URL

> http://localhost:5000

such as

> http://localhost:5000/ping


##### 5) Generate your first comments via your generated AI 

Create an .env file such as in your repo in the file for which you want to generate comments:

example in your **./test** folder

```shell
prefixes:
  - "common"
signature: ""
update-comments: false

localai:
    active: true
    url: "http://:5000"
    api_model_version: 10

openai:
    active: false
    api_key: ""
    url: ""
```

and run 

> gocomments -l -w  ./test/. 

you will have : 

```go
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

// InitDatabase  creates a new database for all ATS ads.
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

// GetEventStatus retrieves the status of a specific event.
func GetEventStatus(eventID string, status string) error {

	return nil
}

// SetEventStatus set the given event ID and status.
func SetEventStatus(eventID string, status string) error {

	return nil
}

// myLongRunningProcess will return the error produced by my-exit_processor.
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

// Var4 is a variable of type string..
var Var4 string

// var5 is a private variable of type string..
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

// Quota returns the difference of all
// remaining files in a test. It is set by calling "time" package, so no one can continue even if invalid query params are given (fallback strategy)).
func (t Test) Tota()	{}

// Totbczzsd is a method to define the technicals around geopoint.
func Totbczzsd()	{}

// Totc returns a token from test.
func (t Test) Totc(ctx context.Context, value string) (string, error) {
	return "test", nil
}

// Totd handle GET request to talentview.
func Totd(ctx context.Context, value string) (string, error) {
	return "test", nil
}

// Config represents a structure for .
type Config struct {
}

// Adapter represents a structure for .
type Adapter struct {
}

// New creates a new AWS S3 adapter.
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

// IsValid verify if the token is valid.
func IsValid(token string) bool {
	return true
}

// HasPermission checks if the user has a permission.
func HasPermission(user string) bool {
	return true
}
```



## Resources

- This notebook demonstrates how to generate python code based on the natural language description using TensorFlow, HuggingFace and Mostly Basic Python Programming Benchmark from Google Research. [Text-to-code Generation with TensorFlow, 🤗 & MBPP](https://www.kaggle.com/code/rhtsingh/text-to-code-generation-with-tensorflow-mbpp)
- Common Crawl maintains a free, open repository of web crawl data that can be used by anyone. [commoncrawl](https://commoncrawl.org/)
- [tensorflow > Catalog > c4](https://www.tensorflow.org/datasets/catalog/c4?hl=fr) 
- CodeSearchNet is a collection of datasets and benchmarks that explore the problem of code retrieval using natural language. This research is a continuation of some ideas presented in this blog post and is a joint collaboration between GitHub and the Deep Program Understanding group at Microsoft Research - Cambridge. [CodeSearchNet](https://github.com/github/CodeSearchNet)
