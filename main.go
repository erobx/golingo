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
	duo := duolingo.NewDuolingo(username, token, baseUrl, abbr)

	vocab := duo.GetVocab(token, abbr)
	fmt.Println(vocab)

	fmt.Print()
}
