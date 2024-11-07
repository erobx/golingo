package duolingo

//learning_language
//languages
//language_data
//learning_language_string
//dict_base_url
//fullname
//username
//id

type UserData struct {
	LearningLanguage    string     `json:"learning_language"`
	Languages           []Language `json:"languages"`
	LanguageData        EsData     `json:"language_data"`
	LearningLanguageStr string     `json:"learning_language_string"`
	DictBaseUrl         string     `json:"dict_base_url"`
	Fullname            string     `json:"fullname"`
	Username            string     `json:"username"`
	ID                  int        `json:"id"`
}

//bio
//tutorConfig
//emailUniversalPractice
//useUniversalSmartReminderTime
//streak
//pushStreakFreezeUsed
//gems
//enableSpeaker
//timezoneOffset
//lastStreak
//shouldPreventMonetizationForSchoolsUser
//insiteImmersionLingots
//enableSoundEffects
//observedClassroomIds
//goalRewards
//acquisitionSurveyReason
//xpChallenges
//sessionCount
//hasGoogleId
//emailSchoolsPromotion
//emailSchoolsNewsletter
//emailEditSuggested
//coachOutfit
//emailStreakSaver
//subscriptionConfigs
//blockerUserIds
//shopItems
//emailAssignment
//hasPlus
//unconsumedGiftIds
//adsEnabled
//emailComment
//rewardBundles
//health
//blockedUserIds
//trackingProperties
//pushClubs
//autoUpdatePreloadedCourses
//lastResurrectionTimestamp
//pushEarlyBird
//pushLiveUpdates
//stateNeedsTOS
//learningLanguage
//hasFacebookId
//emailStreakFreezeUsed
//paymentMethods
//currentCourseId
//pushPassed
//pushResurrectRewards
//requiresParentalConsent
//timerBoostConfig
//optionalFeatures
//_achievements
//currentCourse
//canUseModerationTools
//whatsappAll
//globalAmbassadorStatus
//animationEnabled
//pushNightOwl
//emailClubs
//plusStatus
//pushStreamPost
//pushAnnouncement
//motivation
//emailWeeklyReport
//emailWordOfTheDay
//autoFacebookPost
//pushStreakSaver
//deactivated
//webNotificationIds
//wechatAll
//emailEventsDigest
//emailSchoolsProductUpdate
//classroomLeaderboardsEnabled
//practiceReminderSettings
//currencyRewardBundles
//emailStreamPost
//googleId
//snoozeExpirationTime
//pushSchoolsAssignment
//joinedClassroomIds
//pushHappyHour
//literacyAdGroup
//referralInfo
//smsAll
//emailVerified
//pushEditSuggested
//emailFollow
//shakeToReportEnabled
//hasPhoneNumber
//emailAnnouncement
//weeklyXp
//xpConfig
//zapsRemaining
//emailResearch
//emailSchoolsAnnouncement
//xpGains
//pushComment
//xpGoalMetToday
//privacySettings
//lingots
//inviteURL
//id
//experiments
//pushPromotion
//emailPromotion
//email
//emailBetaFeedback
//pushUniversalPractice
//lssEnabled
//emailClassroomLeave
//picture
//monthlyXp
//longestStreak
//courses
//name
//hasRecentActivity15
//emailPass
//enableMicrophone
//adsConfig
//universalPracticeNotifyTime
//emailClassroomJoin
//creationDate
//pushFollow
//totalXp
//feedbackProperties
//betaStatus
//timezone
//persistentNotifications
//emailDetMarketing
//shouldForceConnectPhoneNumber
//emailWeeklyProgressReport
//gemsConfig
//xpGoal
//roles
//plusDiscounts
//streakData
//subscriberLevel
//username
//achievements
//emailAssignmentComplete
//fromLanguage
//chinaUserModerationRecords
//pushLeaderboards
//insiteSentenceEdited

type UserIdData struct {
	CurrentCourse CurrentCourse `json:"currentCourse"`
}

//numberOfWords
//title
//wordsLearned
//sections
//pathSectioned
//skills

type CurrentCourse struct {
	PathsSectioned []PathSectioned `json:"pathSectioned"`
}

//title
//completedUnits
//units

type PathSectioned struct {
	Title          string `json:"title"`
	CompletedUnits int    `json:"completedUnits"`
	Units          []Unit `json:"units"`
}

//levels

type Unit struct {
	Levels []Level `json:"levels"`
}

//type
//pathLevelClientData
//finishedSessions

type Level struct {
	LevelType           string              `json:"type"`
	PathLevelClientData PathLevelClientData `json:"pathLevelClientData"`
	FinishedSessions    int                 `json:"finishedSessions"`
}

// skillId
// skillIds
// teachingObjective

type PathLevelClientData struct {
	SkillId  string   `json:"skillId"`
	SkillIds []string `json:"skillIds"`
}

// current_learning bool
// language string
// language_string string
// learning bool

type Language struct {
	CurrentLearning bool   `json:"current_learning"`
	Language        string `json:"language"`
	LanguageString  string `json:"language_string"`
	Learning        bool   `json:"learning"`
}

type LangData interface {
}

type EsData struct {
	Abbr AbbrData `json:"es"`
}

//num_skills_learned
//language_strength
//skills

type AbbrData struct {
	NumSkillsLearned int     `json:"num_skills_learned"`
	LanguageStrength float64 `json:"language_strength"`
	Skills           []Skill `json:"skills"`
}

//practice_recommended
//learned
//num_lexemes
//name
//words

type Skill struct {
	PracticeRec bool     `json:"practice_recommended"`
	Learned     bool     `json:"learned"`
	NumLexems   int      `json:"num_lexemes"`
	Name        string   `json:"name"`
	Words       []string `json:"words"`
}

type Overview struct {
	LearnedLexemes []interface{} `json:"learnedLexemes"`
	Pagination     Pagination    `json:"pagination"`
}

type Pagination struct {
	TotalLexemes   int `json:"totalLexemes"`
	NextStartIndex int `json:"nextStartIndex"`
}

type VocabRequest struct {
	LastTotalLexemeCount int           `json:"lastTotalLexemeCount"`
	ProgressedSkills     []interface{} `json:"progressedSkills"`
}
