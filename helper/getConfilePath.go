package helper

import (
	"fmt"
	"os"
)

func GetConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Failed to get home directory path")
	}
	gatorConfigPath := home + "/.gatorconfig.json"
	return gatorConfigPath, err
}
