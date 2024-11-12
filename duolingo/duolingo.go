package duolingo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Duolingo struct {
	BaseUrl    string
	Session    *http.Client
	Token      string
	UserAgent  string
	UserData   UserData
	UserIdData UserIdData
}

func NewDuolingo(username, token, url, abbr string) *Duolingo {
	duo := &Duolingo{
		BaseUrl:   url,
		Session:   &http.Client{},
		Token:     token,
		UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36",
	}
	duo.setUserData(username)
	duo.setDataFromId()
	return duo
}

func (ud *UserData) SetLang(abbr string) {
	switch abbr {
	case "es":
		var es EsData
		ud.LanguageData = es
	}
}

func (d *Duolingo) GetKnownWords() map[string]struct{} {
	skills := d.UserData.LanguageData.Abbr.Skills
	vocab := map[string]struct{}{}

	for _, skill := range skills {
		if skill.Learned == true {
			for _, word := range skill.Words {
				vocab[word] = struct{}{}
			}
		}
	}

	return vocab
}

func (d *Duolingo) GetVocab(abbr string) []interface{} {
	pSkills := d.getProgressedSkills()

	currIndex := 0
	data := make([]interface{}, 0)

	for {
		overviewUrl := fmt.Sprintf("%s/2017-06-30/users/%d/courses/%s/en/learned-lexemes?sortBy=ALPHABETICAL&startIndex=%d", d.BaseUrl, d.UserData.ID, abbr, currIndex)

		reqBody := VocabRequest{
			LastTotalLexemeCount: 0,
			ProgressedSkills:     pSkills,
		}

		out, err := json.Marshal(reqBody)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", overviewUrl, bytes.NewBuffer(out))
		if err != nil {
			log.Fatal(err)
		}
		req = setAuthHeader(req, d.Token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := d.Session.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var overview Overview
		json.Unmarshal(body, &overview)

		learnedLexemes := overview.LearnedLexemes
		data = append(data, learnedLexemes)

		currIndex = overview.Pagination.NextStartIndex

		if currIndex == 0 {
			break
		}
	}

	return data
}
