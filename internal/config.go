package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const gatorPath = "/.gatorconfig.json"

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (config Config) SetUser(user string) error {
	config.CurrentUserName = user
	err := write(config)
	if err != nil {
		return err
	}
	return nil
}

func Read() (Config, error) {
	var config Config
	filePath, err := getFilePath()
	if err != nil {
		fmt.Println(err)
	}
	open, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	fileContent, err := io.ReadAll(open)
	defer open.Close()
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		fmt.Println(err)
	}
	err = write(config)
	if err != nil {
		fmt.Println("Failed to")
	}
	return config, err
}

func write(config Config) error {
	configByte, err := json.Marshal(config)
	if err != nil {
		return err
	}
	filePath, err := getFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, configByte, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
	}
	filePath := fmt.Sprint(homeDir, gatorPath)
	return filePath, err
}
