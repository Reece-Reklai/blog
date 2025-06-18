package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Reece-Reklai/blog/internal/command"
	"github.com/Reece-Reklai/blog/internal/config"
)

func main() {
	const configFileName = "/.gatorconfig.json"
	var config config.Config
	filePath, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Failed to get file path")
	}
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		fmt.Println("Failed to open file")
	}
	byteValue, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Failed to read file")
	}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		fmt.Println("Failed to unmarshal json")
	}
	config.SetUser()
	err = writeFile(config)
	if err != nil {
		fmt.Println("Failed to write to gatorconfig.json file")
	}
	fmt.Println(config.URL)
	return
}

func writeFile(wConfig config.Config) error {
	jsonToByte, err := json.Marshal(wConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile("gatorconfig.json", jsonToByte, 0666)
	if err != nil {
		return err
	}
	return err
}

func handlerLogin(current *config.State, cmd command.Command) {

}
