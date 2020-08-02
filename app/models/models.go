package models

import (
  "time"
)

type NextRaceData struct {
  Event
  ServerStatusIcon string "WIP(SERVER STATUS ICON)"
  ServerStatusStatement string "WIP(SERVER STATUS STATEMENT)"
}

type User struct {
  UserID    []byte `json:"userid" db:"userid"`
  PWHash    []byte `json:"pwhash" db:"pwhash"`
  Attribute int    `json:"attribute" db:"attribute"`
}

type Event struct {
    Id                                     uint        `json:"id" db:"id"`
    Startdate                              time.Time   `json:"startdate" db:"startdate"`
    Track                                  string      `json:"track" db:"track"`
    WeatherRandomness                      int         `json:"weatherrandomness" db:"weatherRandomness"`
    P_hourOfDay                            int         `json:"P_hourOfDay" db:"P_hourOfDay"`
    P_timeMultiplier                       int         `json:"P_timeMultiplier" db:"P_timeMultiplier"`
    P_sessionDurationMinute                int         `json:"P_sessionDurationMinute" db:"P_sessionDurationMinute"`
    Q_hourOfDay                            int         `json:"Q_hourOfDay" db:"Q_hourOfDay"`
    Q_timeMultiplier                       int         `json:"Q_timeMultiplier" db:"Q_timeMultiplier"`
    Q_sessionDurationMinute                int         `json:"Q_sessionDurationMinute" db:"Q_sessionDurationMinute"`
    R_hourOfDay                            int         `json:"R_hourOfDay" db:"R_hourOfDay"`
    R_timeMultiplier                       int         `json:"R_timeMultiplier" db:"R_timeMultiplier"`
    R_sessionDurationMinute                int         `json:"R_sessionDurationMinute" db:"R_sessionDurationMinute"`
    PitWindowLengthSec                     int         `json:"pitWindowLengthSec" db:"pitWindowLengthSec"`
    IsRefuellingAllowedInRace              bool        `json:"isRefuellingAllowedInRace" db:"isRefuellingAllowedInRace"`
    MandatoryPitstopCount                  int         `json:"mandatoryPitstopCount" db:"mandatoryPitstopCount"`
    IsMandatoryPitstopRefuellingRequired   bool        `json:"isMandatoryPitstopRefuellingRequired" db:"isMandatoryPitstopRefuellingRequired"`
    IsMandatoryPitstopTyreChangeRequired   bool        `json:"isMandatoryPitstopTyreChangeRequired" db:"isMandatoryPitstopTyreChangeRequired"`
    IsMandatoryPitstopSwapDriverRequired   bool        `json:"isMandatoryPitstopSwapDriverRequired" db:"isMandatoryPitstopSwapDriverRequired"`
    TyreSetCount                           int         `json:"tyreSetCount" db:"tyreSetCount"`
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
