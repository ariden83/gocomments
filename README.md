# auto-add-golang-comments

Automatically adding comments to go functions and variables to code such as

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
func GetEventStatus(eventID string, status string) error {

	return nil
}

// SetEventStatus is a method which update the event status that take an eventID of type string, a status of type string
// It's return an error if fails, otherwise nil.
func SetEventStatus(eventID string, status string) error {

	return nil
}

// myLongRunningProcess is a private method which execute my long running process that take a ctx of type context.Context
// It's return an error if fails, otherwise nil.
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
	// Var1 is a variable of type int which provides .
	Var1	int

	// Var2 is a variable of type string which provides .
	Var2	string

	// Var3 is a variable of type toto which provides .
	Var3	toto
)

// Var4 is a variable of type string which provides .
var Var4 string

// var5 is a private variable of type string which provides .
var var5 string

// var6 is a private constant which provides .
const var6 = "tata"

// Var7 is a constant which provides .
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
func (t Test) Totc(ctx context.Context, value string) (string, error) {
	return "test", nil
}

// Totd is a method that take a ctx of type context.Context, a value of type string
// and returns a string
// It's return an error if fails, otherwise nil.
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
```

## Commands

Command to generate the binary :

> make build

Command to displays the specified folder code with added comments

>  ./bin/gocomments ./test/.

Changes the code of the specified folder with added comments

>  ./bin/gocomments -l -w  ./test/.

