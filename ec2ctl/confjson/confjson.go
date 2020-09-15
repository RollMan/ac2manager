package confjson

import (
	"github.com/RollMan/ac2manager/app/models"
	"path/filepath"
	"runtime"
)

var (
	_, sourcefilepath, _, _ = runtime.Caller(0)
	sourcedirpath           = filepath.Dir(sourcefilepath)
)

func ReadDefaultConfigs() (*AssistRules, *Settings, *Event, *Configuration, *EventRules) {
	var assistRules *AssistRules = &AssistRules{}
	var settings *Settings = &Settings{}
	var event *Event = &Event{}
	var configuration *Configuration = &Configuration{}
	var eventRules *EventRules = &EventRules{}

	// FIXME: specify these paths through command line arguments
	// and set default path if no arguments are passed.
	prefix := sourcedirpath
	assistRules.ParseConfig(prefix + "/" + "assistRules.json")
	settings.ParseConfig(prefix + "/" + "settings.json")
	event.ParseConfig(prefix + "/" + "event.json")
	configuration.ParseConfig(prefix + "/" + "configuration.json")
	eventRules.ParseConfig(prefix + "/" + "eventRules.json")

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
