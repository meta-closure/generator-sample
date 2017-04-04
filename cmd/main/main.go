package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/meta-closure/sample-generater/generater"
	"github.com/meta-closure/sample-generater/server"
)

type Operation struct {
	command string
}

func NewOperation(args []string) (Operation, error) {
	if len(args) != 2 {
		return Operation{}, errors.New("required argument is 1")
	}

	return Operation{
		command: args[1],
	}, nil
}

func (op Operation) isGenerate() bool {
	return op.command == "generate"
}

func (op Operation) isRun() bool {
	return op.command == "run"
}

func main() {
	if err := _main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func _main() error {
	cmd, err := NewOperation(os.Args)
	if err != nil {
		return err
	}
	if cmd.isGenerate() {
		return generater.Generate()
	}

	if cmd.isRun() {
		return server.Run()
	}

	return errors.New("unknown operation")
}
