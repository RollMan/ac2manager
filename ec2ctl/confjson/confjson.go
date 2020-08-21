package confjson

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
  "github.com/RollMan/ac2manager/app/models"
)

func ReadDefaultConfigs() (AssistRules, Settings, Event, Configuration){
  var assistRules AssistRules
  var settings Settings
  var event Event
  var configuration Configuration

  assistRules.ParseConfig("assistRules.json")
  settings.ParseConfig("settings.json")
  event.ParseConfig("event.json")
  configuration.ParseConfig("configuration.json")

  return assistRules, settings, event, configuration
}

func SetConfigs(src models.Event, assistRules *AssistRules, settings *Settings, event *Event, conf *Configuration){
}
