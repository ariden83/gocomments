package main

import (
	"context"
	"fmt"
	"github.com/ariden83/auto-add-golang-comments/cmd"
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

type toto int

var (

	// Var1 is a variable of type int which provides .
	Var1 int

	// Var2 is a variable of type string which provides .
	Var2 string

	// Var3 is a variable of type toto which provides .
	Var3 toto
)

type Test struct // test// - Key
// test// - KeyB
// test// - KeyC

{
	Key  string
	KeyB int
	KeyC toto
}

// Tota is a method that belongs to the Test struct.
// It does not take any arguments.
func (t Test) Tota() {}

// Totb is a method .
// It does not take any arguments.
func Totb() {}

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
