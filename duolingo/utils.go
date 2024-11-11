package duolingo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (d *Duolingo) getProgressedSkills() []ProgressedSkill {
	currentCourses := d.UserIdData.CurrentCourse.PathsSectioned
	progressedSkills := make([]ProgressedSkill, 0)

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
					newObj := ProgressedSkill{
						1,
						finishedSessions,
						SkillId{
							skillId,
						},
					}

					progressedSkills = append(progressedSkills, newObj)
				} else if pathLevelClientData.SkillIds != nil {
					skillIds := pathLevelClientData.SkillIds
					for _, skillId := range skillIds {
						newObj := ProgressedSkill{
							1,
							finishedSessions,
							SkillId{
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
