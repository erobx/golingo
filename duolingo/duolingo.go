package duolingo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Duolingo struct {
	BaseUrl   string
	Session   *http.Client
	UserAgent string
}

func NewDuolingo(url string) Duolingo {
	return Duolingo{
		BaseUrl:   url,
		Session:   &http.Client{},
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
	}
}

func setAuthHeader(req *http.Request, token string) *http.Request {
	reqValue := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", reqValue)
	return req
}

func dumpUserData(data []byte) {
	f, _ := os.Create("user.txt")
	defer f.Close()
	_, err := f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

func (d Duolingo) getUserId(username, token string) int {
	userUrl := fmt.Sprintf("%s/users/%s", d.BaseUrl, username)
	req, err := http.NewRequest("GET", userUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req = setAuthHeader(req, token)

	resp, err := d.Session.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var v map[string]interface{}
	err = json.Unmarshal(body, &v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(v["language_data"])

	tmpId, ok := v["id"].(float64)
	if !ok {
		log.Fatal("couldn't get user id")
	}

	userId := int(tmpId)
	return userId
}

func (d Duolingo) GetVocab(username, token, abbr string) {
	userId := d.getUserId(username, token)
	currIndex := 0
	vocabUrl := fmt.Sprintf("%s/2017-06-30/users/%d/courses/%s/en/learned-lexemes?sortBy=ALPHABETICAL&startIndex=%d", d.BaseUrl, userId, abbr, currIndex)
	req, err := http.NewRequest("GET", vocabUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req = setAuthHeader(req, token)

	resp, err := d.Session.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var v map[string]interface{}
	json.Unmarshal(body, &v)

	fmt.Println(string(body))
}
