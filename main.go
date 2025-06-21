package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Reece-Reklai/blog/helper"
	"github.com/Reece-Reklai/blog/internal/command"
	"github.com/Reece-Reklai/blog/internal/config/gator"
	"github.com/Reece-Reklai/blog/internal/config/state"
	"io"
	"os"
)

func main() {
	const configFileName = "/.gatorconfig.json"
	var conf gator.Conf
	var states state.State
	var cmd command.Command
	var cli command.CLI
	cli.Cmds = make(map[string]func(*state.State, command.Command) error)
	cli.Register("login", handlerLogin)
	filePath, err := helper.GetConfigFilePath()
	if err != nil {
		fmt.Println("Failed to get file path")
	}
	// Open File
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println("Failed to open file")
	}

	// Read File
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Failed to read file")
	}

	// File Content decoding json to byte to json
	err = json.Unmarshal(byteValue, &conf)
	if err != nil {
		fmt.Println("Failed to unmarshal json")
	}

	// Keep track of the state from json file
	states.Current = &conf

	// Handle User input
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("Require Two Arguments")
	} else {
		cmd.Name = args[0]
		cmd.Arguments = append(cmd.Arguments, args[1])
		err = cli.Run(&states, cmd)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = writeFile(conf)
	if err != nil {
		fmt.Print(err)
	}
	return
}

func handlerLogin(states *state.State, cmd command.Command) error {
	if cmd.Name == "" || cmd.Arguments == nil {
		return errors.New("Empty slice")
	}
	if cmd.Name != "login" {
		return errors.New("Failed To Login")
	}
	err := states.Current.SetUser(cmd.Arguments[0])
	if err != nil {
		return err
	}
	fmt.Println("Login successfully")
	return nil
}

func writeFile(wConfig gator.Conf) error {
	filePath, err := helper.GetConfigFilePath()
	if err != nil {
		fmt.Println("Failed to get file path")
	}
	jsonToByte, err := json.Marshal(wConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, jsonToByte, 0666)
	if err != nil {
		return err
	}
	return nil
}
