package main

import (
	"context"
	"fmt"
	"github.com/ariden/auto-add-golang-comments/cmd"
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
	cmd.Execute()

	<-ctx.Done()

	fmt.Println("Processus terminé ou interrompu.")
}

// myLongRunningProcess is a private method that take a ctx of type context.Context.
func myLongRunningProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():

			fmt.Println("Processus annulé.")
			return
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

// Totbczzsd is a method .
// It does not take any arguments.
func Totbczzsd()	{}

// Totc is a method that belongs to the Test struct that take a ctx of type context.Context, a value of type string
// and returns a string and an error.
func (t Test) Totc(ctx context.Context, value string) (string, error) {
	return "test", nil
}

// Totd is a method that take a ctx of type context.Context, a value of type string
// and returns a string and an error.
func Totd(ctx context.Context, value string) (string, error) {
	return "test", nil
}
