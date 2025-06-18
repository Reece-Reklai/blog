package main

import (
	"fmt"
	"os"
)

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory path")
	}
	gatorConfigPath := home + "/.gatorconfig.json"
	return gatorConfigPath, err
}
