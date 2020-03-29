package models

import (
  "time"
)

type Login struct {
  UserID   string `json:"userid"`
  Password string `json:"pw"`
}

type Event struct {
  Id                                    uint
  Startdate                             time.Time
  Track                                 string
  WeatherRandomness                     int
  P_hourOfDay                           int
  P_timeMultiplier                      int
  P_sessionDurationMinute               int
  Q_hourOfDay                           int
  Q_timeMultiplier                      int
  Q_sessionDurationMinute               int
  R_hourOfDay                           int
  R_timeMultiplier                      int
  R_sessionDurationMinute               int
  PitWindowLengthSec                    int
  IsRefuellingAllowedInRace             bool
  MandatoryPitstopCount                 int
  IsMandatoryPitstopRefuellingRequired  bool
  IsMandatoryPitstopTyreChangeRequired  bool
  IsMandatoryPitstopSwapDriverRequired  bool
  TyreSetCount                          int
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
