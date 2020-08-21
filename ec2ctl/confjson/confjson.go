package confjson

import (
	"encoding/json"
	"github.com/RollMan/ac2manager/app/models"
	"io/ioutil"
	"log"
	"os"
)

func ReadDefaultConfigs() (AssistRules, Settings, Event, Configuration) {
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

func SetConfigs(src models.Event, assistRules *AssistRules, settings *Settings, event *Event, conf *Configuration) {
	event.Track = src.Track
	event.WeatherRandomness = src.WeatherRandomness

	sessions := make([]Session, 3)
	sessions[0] = Session{
		HourOfDay:              src.P_hourOfDay,
		DayOfWeekend:           1,
		TimeMultiplier:         src.P_timeMultiplier,
		SessionType:            "P",
		SessionDurationMinutes: src.P_sessionDurationMinute,
	}
	sessions[1] = Session{
		HourOfDay:              src.Q_hourOfDay,
		DayOfWeekend:           1,
		TimeMultiplier:         src.Q_timeMultiplier,
		SessionType:            "Q",
		SessionDurationMinutes: src.Q_sessionDurationMinute,
	}
	sessions[2] = Session{
		HourOfDay:              src.R_hourOfDay,
		DayOfWeekend:           1,
		TimeMultiplier:         src.R_timeMultiplier,
		SessionType:            "R",
		SessionDurationMinutes: src.R_sessionDurationMinute,
	}

	event.Sessions = sessions
}
