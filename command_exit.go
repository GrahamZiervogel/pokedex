package main

import (
	"errors"
	"fmt"
	"os"
)

func commandExit(cfg *Config, args ...string) error {
	if len(args) > 0 {
		return errors.New("exit command does not take any arguments")
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
