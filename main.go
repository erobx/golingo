package main

import (
	"github.com/erobx/golingo/duolingo"
)

const jwt = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJleHAiOjYzMDcyMDAwMDAsImlhdCI6MCwic3ViIjo1Mzg2MjYzNDl9.LloQTnVPcsmQbboPtUdyT8FVxapdLrH8XDcMZWmmjnY"

func main() {
	duo := duolingo.NewDuolingo("https://www.duolingo.com")

	login, err := duo.Login("", "")
}
