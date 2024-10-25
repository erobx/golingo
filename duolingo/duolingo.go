package duolingo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"github.com/golang-jwt/jwt/v5"
)

type Duolingo struct {
	BaseUrl   string
	Session   *http.Client
	UserAgent string
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewDuolingo(url string) Duolingo {
	return Duolingo{
		BaseUrl:   url,
		Session:   &http.Client{},
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
	}
}

func NewLogin(username, password string) *Login {
	return &Login{
		Username: username,
		Password: password,
	}
}

func (d Duolingo) Login(username, password string) {
	url := fmt.Sprintf(d.BaseUrl, "/login")

	out, err := json.Marshal(NewLogin(username, password))
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	if err != nil {
		log.Fatal(err)
	}

	res, err := d.Session.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		log.Fatal(res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successful login")
	fmt.Println(string(body))
	fmt.Println(res.Header["jwt_token"])
}

func (d Duolingo) GetVocab() {

}
