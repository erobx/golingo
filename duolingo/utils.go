package duolingo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func setAuthHeader(req *http.Request, token string) *http.Request {
	reqValue := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", reqValue)
	return req
}

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

	var data UserData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	d.UserData = data
}

func (d *Duolingo) setDataFromId(token string) {
	userUrl := fmt.Sprintf("%s/2017-06-30/users/%d", d.BaseUrl, d.UserData.ID)
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

	var data UserIdData
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	d.UserIdData = data
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

func (d *Duolingo) getProgressedSkills() []interface{} {
	currentCourses := d.UserIdData.CurrentCourse.PathsSectioned
	progressedSkills := make([]interface{}, 0)

	for _, section := range currentCourses {
		completedUnits := section.CompletedUnits
		units := section.Units
		for i := range completedUnits {
			unit := units[i]
			levels := unit.Levels
			for _, l := range levels {
				levelType := l.LevelType
				if levelType == "chest" || levelType == "unit_review" {
					continue
				}
				pathLevelClientData := l.PathLevelClientData
				finishedSessions := l.FinishedSessions
				// maybe check for if skillId in pathLevelClientData
				if pathLevelClientData.SkillId != "" {
					skillId := pathLevelClientData.SkillId
					newObj := struct {
						finishedLevels   int
						finishedSessions int
						skillId          interface{}
					}{
						1,
						finishedSessions,
						struct {
							id string
						}{
							skillId,
						},
					}
					progressedSkills = append(progressedSkills, newObj)
				} else if pathLevelClientData.SkillIds != nil {
					skillIds := pathLevelClientData.SkillIds
					for _, skillId := range skillIds {
						newObj := struct {
							finishedLevels   int
							finishedSessions int
							skillId          interface{}
						}{
							1,
							finishedSessions,
							struct {
								id string
							}{
								skillId,
							},
						}
						progressedSkills = append(progressedSkills, newObj)
					}
				}
			}
		}
	}
	return progressedSkills
}
