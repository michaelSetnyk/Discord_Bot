package main

import (
	"./bot"
	"./config"
	"fmt"

)

func main() {
	err := config.ReadConfig()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Init()

fmt.Println("inted")
}