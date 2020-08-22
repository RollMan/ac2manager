package confjson

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type AssistRules struct {
	DisableIdealLine         int `json:"disableIdealLine"`
	DisableAutosteer         int `json:"disableAutosteer"`
	StabilityControlLevelMax int `json:"stabilityControlLevelMax"`
	DisableAutoPitLimiter    int `json:"disableAutoPitLimiter"`
	DisableAutoGear          int `json:"disableAutoGear"`
	DisableAutoClutch        int `json:"disableAutoClutch"`
	DisableAutoEngineStart   int `json:"disableAutoEngineStart"`
	DisableAutoWiper         int `json:"disableAutoWiper"`
	DisableAutoLights        int `json:"disableAutoLights"`
}

func (t *AssistRules) ParseConfig(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), t)
}

type Settings struct {
	ServerName                 string `json:"serverName"`
	AdminPassword              string `json:"adminPassword"`
	CarGroup                   string `json:"carGroup"`
	TrackMedalsRequirement     int    `json:"trackMedalsRequirement"`
	SafetyRatingRequirement    int    `json:"safetyRatingRequirement"`
	RacecraftRatingRequirement int    `json:"racecraftRatingRequirement"`
	Password                   string `json:"password"`
	MaxCarSlots                int    `json:"maxCarSlots"`
	SpectatorPassword          string `json:"spectatorPassword"`
	ConfigVersion              int    `json:"configVersion"`
}

func (t *Settings) ParseConfig(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), t)
}

type Session struct {
	HourOfDay              int    `json:"hourOfDay"`
	DayOfWeekend           int    `json:"dayOfWeekend"`
	TimeMultiplier         int    `json:"timeMultiplier"`
	SessionType            string `json:"sessionType"`
	SessionDurationMinutes int    `json:"sessionDurationMinutes"`
}

func (t *Session) ParseConfig(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), t)
}

type Event struct {
	Track                     string    `json:"track"`
	PreRaceWaitingTimeSeconds int       `json:"preRaceWaitingTimeSeconds"`
	SessionOverTimeSeconds    int       `json:"sessionOverTimeSeconds"`
	AmbientTemp               int       `json:"ambientTemp"`
	CloudLevel                float32   `json:"cloudLevel"`
	Rain                      float32   `json:"rain"`
	WeatherRandomness         int       `json:"weatherRandomness"`
	Sessions                  []Session `json:"sessions"`
	ConfigVersion             int       `json:"configVersion"`
}

func (t *Event) ParseConfig(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), t)
}

type Configuration struct {
	UdpPort        int `json:"udpPort"`
	TcpPort        int `json:"tcpPort"`
	MaxConnections int `json:"maxConnections"`
	ConfigVersion  int `json:"configVersion"`
}

func (t *Configuration) ParseConfig(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), t)
}

type EventRules struct {
	qualifyStandingType                  int  `json:"qualifyStandingType"`
	pitWindowLengthSec                   int  `json:"pitWindowLengthSec"`
	driverStintTimeSec                   int  `json:"driverStintTimeSec"`
	mandatoryPitstopCount                int  `json:"mandatoryPitstopCount"`
	maxTotalDrivingTime                  int  `json:"maxTotalDrivingTime"`
	maxDriversCount                      int  `json:"maxDriversCount"`
	isRefuellingAllowedInRace            bool `json:"isRefuellingAllowedInRace"`
	isRefuellingTimeFixed                bool `json:"isRefuellingTimeFixed"`
	isMandatoryPitstopRefuellingRequired bool `json:"isMandatoryPitstopRefuellingRequired"`
	isMandatoryPitstopTyreChangeRequired bool `json:"isMandatoryPitstopTyreChangeRequired"`
	isMandatoryPitstopSwapDriverRequired bool `json:"isMandatoryPitstopSwapDriverRequired"`
	tyreSetCount                         int  `json:"tyreSetCount"`
}

func (t *EventRules) ParseConfig(filename string) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), t)
}
