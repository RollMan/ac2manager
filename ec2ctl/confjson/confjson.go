package confjson

import (
  "encoding/json"
  "io/ioutil"
  "log"
  "os"
)

func readConfigs() (AssistRules, Settings, Event, Configuration){
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

