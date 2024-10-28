package main

import (
	"fmt"
	"os"

	"github.com/erobx/golingo/duolingo"
	"github.com/joho/godotenv"
)

const (
	baseUrl = "https://www.duolingo.com"
	abbr    = "es"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}

	username := os.Getenv("DUO_USERNAME")
	token := os.Getenv("DUO_JWT")
	duo := duolingo.NewDuolingo(username, token, baseUrl)

	vocab := duo.GetKnownWords(abbr)

	_, ok := vocab["farmacia"]
	if !ok {
		fmt.Println("working")
	}

	//fmt.Println("You've learned:", len(vocab), "vocab words!")
}
