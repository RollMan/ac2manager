package confjson

import (
	"github.com/RollMan/ac2manager/app/models"
)

func ReadDefaultConfigs() (AssistRules, Settings, Event, Configuration, EventRules) {
	var assistRules AssistRules
	var settings Settings
	var event Event
	var configuration Configuration
	var eventRules EventRules

	assistRules.ParseConfig("assistRules.json")
	settings.ParseConfig("settings.json")
	event.ParseConfig("event.json")
	configuration.ParseConfig("configuration.json")
	eventRules.ParseConfig("eventRules.json")

	return assistRules, settings, event, configuration, eventRules
}

func SetConfigs(src models.Event, assistRules *AssistRules, settings *Settings, event *Event, conf *Configuration, eventRules *EventRules) {
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

	eventRules.PitWindowLengthSec = src.PitWindowLengthSec
	eventRules.IsRefuellingAllowedInRace = src.IsRefuellingAllowedInRace
	eventRules.MandatoryPitstopCount = src.MandatoryPitstopCount
	eventRules.IsMandatoryPitstopRefuellingRequired = src.IsMandatoryPitstopRefuellingRequired
	eventRules.IsMandatoryPitstopTyreChangeRequired = src.IsMandatoryPitstopTyreChangeRequired
	eventRules.IsMandatoryPitstopSwapDriverRequired = src.IsMandatoryPitstopSwapDriverRequired
	eventRules.TyreSetCount = src.TyreSetCount

}
