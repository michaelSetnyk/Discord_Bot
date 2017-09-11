package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
    "os"
)

var (
	// Public variables
	Token     string
	BotPrefix string

	// Private variables
	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
    file, err := ioutil.ReadFile(os.Getenv("LOCALAPPDATA")+"/.dough_bot/config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}
