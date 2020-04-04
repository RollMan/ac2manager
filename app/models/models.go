package models

import (
  "time"
  "net/http"
  "github.com/mholt/binding"
)

type Login struct {
  UserID   string `json:"userid"`
  Password string `json:"pw"`
}

type Event struct {
 Id uint `json:"id"`
 Startdate time.Time `json:"startdate"`
 Track string `json:"track"`
 WeatherRandomness int `json:"weatherrandomness"`
 P_hourOfDay int `json:"P_hourOfDay"`
 P_timeMultiplier int `json:"P_timeMultiplier"`
 P_sessionDurationMinute int `json:"P_sessionDurationMinute"`
 Q_hourOfDay int `json:"Q_hourOfDay"`
 Q_timeMultiplier int `json:"Q_timeMultiplier"`
 Q_sessionDurationMinute int `json:"Q_sessionDurationMinute"`
 R_hourOfDay int `json:"R_hourOfDay"`
 R_timeMultiplier int `json:"R_timeMultiplier"`
 R_sessionDurationMinute int `json:"R_sessionDurationMinute"`
 PitWindowLengthSec int `json:"pitWindowLengthSec"`
 IsRefuellingAllowedInRace bool `json:"isRefuellingAllowedInRace"`
 MandatoryPitstopCount int `json:"mandatoryPitstopCount"`
 IsMandatoryPitstopRefuellingRequired bool `json:"isMandatoryPitstopRefuellingRequired"`
 IsMandatoryPitstopTyreChangeRequired bool `json:"isMandatoryPitstopTyreChangeRequired"`
 IsMandatoryPitstopSwapDriverRequired bool `json:"isMandatoryPitstopSwapDriverRequired"`
 TyreSetCount int `json:"tyreSetCount"`
}

func (e *Event) FieldMap(r *http.Request) binding.FieldMap {
  return binding.FieldMap{
    &e.Id                                     :"id",
    &e.Startdate                              :"startdate",
    &e.Track                                  :"track",
    &e.WeatherRandomness                      :"weatherrandomness",
    &e.P_hourOfDay                            :"P_hourOfDay",
    &e.P_timeMultiplier                       :"P_timeMultiplier",
    &e.P_sessionDurationMinute                :"P_sessionDurationMinute",
    &e.Q_hourOfDay                            :"Q_hourOfDay",
    &e.Q_timeMultiplier                       :"Q_timeMultiplier",
    &e.Q_sessionDurationMinute                :"Q_sessionDurationMinute",
    &e.R_hourOfDay                            :"R_hourOfDay",
    &e.R_timeMultiplier                       :"R_timeMultiplier",
    &e.R_sessionDurationMinute                :"R_sessionDurationMinute",
    &e.PitWindowLengthSec                     :"pitWindowLengthSec",
    &e.IsRefuellingAllowedInRace              :"isRefuellingAllowedInRace",
    &e.MandatoryPitstopCount                  :"mandatoryPitstopCount",
    &e.IsMandatoryPitstopRefuellingRequired   :"isMandatoryPitstopRefuellingRequired",
    &e.IsMandatoryPitstopTyreChangeRequired   :"isMandatoryPitstopTyreChangeRequired",
    &e.IsMandatoryPitstopSwapDriverRequired   :"isMandatoryPitstopSwapDriverRequired",
    &e.TyreSetCount                           :"tyreSetCount",
  }
}

type NextRaceData struct {
  Event
  ServerStatusIcon string "WIP(SERVER STATUS ICON)"
  ServerStatusStatement string "WIP(SERVER STATUS STATEMENT)"
}

type User struct {
  UserID    []byte `json:"userid"`
  PWHash    []byte `json:"pwhash"`
  Attribute int    `json:"attribute"`
}

type NoSuchUserError struct{}
type NoMatchingPasswordError struct{}

func (e *NoSuchUserError) Error() string {
  return "No such userid in DB."
}

func (e *NoMatchingPasswordError) Error() string {
  return "Password unmatched."
}

const NoEvent = `<h3>No upcoming events.<h3>
<p>
Come back later or contact administrator.
</p>
`

const EventConfigure = `
<ul>
<h3>Event starts at: %v</h3>
<li>Track: %v</li>
<li>Weather randomness: %v/10</li>
<li>The number of tyre sets: %v</li>
<li>Practice:
<ul class="upcoming_race_rule">
<li>Duration: %v min.</li>
<li>Time multiplier: &times;%v</li>
<li>Hour of day in game: %v</li>
</ul>
</li>
<li>Qualify:
<ul class="upcoming_race_rule">
<li>Duration: %v min.</li>
<li>Time multiplier: &times;%v</li>
<li>Hour of day in game: %v</li>
</ul>
</li>
<li>Race:
<ul class="upcoming_race_rule">
<li>Duration: %v min.</li>
<li>Time multiplier: &times;%v</li>
<li>Hour of day in game: %v</li>
</ul>
<li>Pit rule:
<ul class="upcoming_race_rule">
<li>Mandatory pit count: %v</li>
<li>Pit window length: %v sec.</li>
<li>Refuelling allowed?: %v</li>
<li>Refuelling required?: %v </li>
<li>Tyre change required?: %v </li>
<li>Driver swap required?: %v </li>
</ul>
</li>
</ul>
`
