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
	userData  map[string]interface{}
}

func NewDuolingo(username, token, url string) *Duolingo {
	duo := &Duolingo{
		BaseUrl:   url,
		Session:   &http.Client{},
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
	}
	duo.setUserData(username, token)
	return duo
}

func setAuthHeader(req *http.Request, token string) *http.Request {
	reqValue := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", reqValue)
	return req
}

func dumpUserDataAsString(data map[string]struct{}) {
	f, _ := os.Create("user_vocab.txt")
	defer f.Close()
	//for i, s := range data {
	//	if i == len(data)-1 {
	//		f.WriteString(s)
	//	} else {
	//		s = s + ", "
	//		f.WriteString(s)
	//	}
	//}
	fmt.Println("success writing user vocab")
}

func dumpUserDataAsBytes(data []byte) {
	f, _ := os.Create("user_vocab.txt")
	defer f.Close()
	_, err := f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("success writing user vocab")
}

// keys from main user data request
// learning_language, languages, language_data, id

func (d *Duolingo) setUserData(username, token string) {
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

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	d.userData = data
}

func (d *Duolingo) getUserId() int {
	// do this since json always returns 64-bit
	tmpId, ok := d.userData["id"].(float64)
	if !ok {
		log.Fatal("couldn't get user id")
	}

	userId := int(tmpId)
	return userId
}

func (d *Duolingo) GetKnownWords(abbr string) map[string]struct{} {
	lang_data := d.userData["language_data"].(map[string]interface{})
	lang_data = lang_data[abbr].(map[string]interface{})
	skills := lang_data["skills"].([]interface{})

	vocab := make(map[string]struct{})

	for _, skill := range skills {
		topic := skill.(map[string]interface{})
		if topic["learned"] == true {
			wordSlice := skill.(map[string]interface{})["words"]
			interSlice := wordSlice.([]interface{})
			for _, item := range interSlice {
				word := item.(string)
				vocab[word] = struct{}{}
			}
		}
	}

	return vocab
}

type Overview struct {
	lastTotalLexemeCount int           `json:"lastTotalLexemeCount"`
	progressedSkills     []interface{} `json:"progressedSkills"`
}

func (d *Duolingo) GetVocab(abbr string) {
	//currIndex := 0
	//data := make([]string, 0)
	//overviewUrl := fmt.Sprintf("%s/2017-06-30/users/%d/courses/%s/en/learned-lexemes?sortBy=ALPHABETICAL&startIndex=%d", d.BaseUrl, userId, abbr, currIndex)
	//req, err := http.NewRequest("GET", overviewUrl, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//req = setAuthHeader(req, token)

	//resp, err := d.Session.Do(req)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//body, err := io.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//var v map[string]interface{}
	//json.Unmarshal(body, &v)

	//fmt.Println(string(body))

}
