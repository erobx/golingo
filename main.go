package main

import (
	"os"

	"github.com/erobx/golingo/duolingo"
	"github.com/joho/godotenv"
)

const (
	baseUrl = "https://www.duolingo.com"
	abbr    = "ES"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		os.Exit(1)
	}
	username := os.Getenv("DUO_USERNAME")
	token := os.Getenv("DUO_JWT")
	duo := duolingo.NewDuolingo(baseUrl)
	duo.GetVocab(username, token, abbr)
}
