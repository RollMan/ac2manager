package confjson

import (
  "io/ioutil"
  "encoding/json"
  "os"
  "log"
)

type AssistRules struct {
  disableIdealLine         int `json:"disableIdealLine"`
  disableAutosteer         int `json:"disableAutosteer"`
  stabilityControlLevelMax int `json:"stabilityControlLevelMax"`
  disableAutoPitLimiter    int `json:"disableAutoPitLimiter"`
  disableAutoGear          int `json:"disableAutoGear"`
  disableAutoClutch        int `json:"disableAutoClutch"`
  disableAutoEngineStart   int `json:"disableAutoEngineStart"`
  disableAutoWiper         int `json:"disableAutoWiper"`
  disableAutoLights        int `json:"disableAutoLights"`
}

func (*t AssistRules) ParseConfig(filename string){
  jsonFile, err := os.Open(filename)
  if err != nil {
    log.Fatalln(err)
  }
  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  json.Unmarshal([]byte(byteValue), t)
}


type Settings struct {
    serverName                 string     `json:"serverName"`
    adminPassword              string     `json:"adminPassword"`
    carGroup                   string     `json:"carGroup"`
    trackMedalsRequirement     int        `json:"trackMedalsRequirement"`
    safetyRatingRequirement    int        `json:"safetyRatingRequirement"`
    racecraftRatingRequirement int        `json:"racecraftRatingRequirement"`
    password                   string     `json:"password"`
    maxCarSlots                int        `json:"maxCarSlots"`
    spectatorPassword          string     `json:"spectatorPassword"`
    configVersion              int        `json:"configVersion"`
}

func (*t Settings) ParseConfig(filename string){
  jsonFile, err := os.Open(filename)
  if err != nil {
    log.Fatalln(err)
  }
  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  json.Unmarshal([]byte(byteValue), t)
}

type Session struct {
  hourOfDay               int   `json:"hourOfDay"`
  dayOfWeekend            int   `json:"dayOfWeekend"`
  timeMultiplier          int   `json:"timeMultiplier"`
  sessionType             int   `json:"sessionType"`
  sessionDurationMinutes  int   `json:"sessionDurationMinutes"`
}

func (*t Session) ParseConfig(filename string){
  jsonFile, err := os.Open(filename)
  if err != nil {
    log.Fatalln(err)
  }
  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  json.Unmarshal([]byte(byteValue), t)
}

type Event struct {
          track                                              string
          preRaceWaitingTimeSeconds                          int
          sessionOverTimeSeconds                             int
          ambientTemp                                        int
          cloudLevel                                         float32
          rain                                               float32
          weatherRandomness                                  int
          sessions                                           []session
          configVersion                                      int
}

func (*t Event) ParseConfig(filename string){
  jsonFile, err := os.Open(filename)
  if err != nil {
    log.Fatalln(err)
  }
  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  json.Unmarshal([]byte(byteValue), t)
}

type Configuration struct {
    udpPort         int `json:"udpPort"`
    tcpPort         int `json:"tcpPort"`
    maxConnections  int `json:"maxConnections"`
    configVersion   int `json:"configVersion"`
}

func (*t Configuration) ParseConfig(filename string){
  jsonFile, err := os.Open(filename)
  if err != nil {
    log.Fatalln(err)
  }
  defer jsonFile.Close()

  byteValue, _ := ioutil.ReadAll(jsonFile)

  json.Unmarshal([]byte(byteValue), t)
}
