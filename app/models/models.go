package models

import (
	"database/sql/driver"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	guuid "github.com/google/uuid"
	"github.com/mholt/binding"
	"net/http"
	"time"
)

type NextRaceData struct {
	Event
	ServerStatusIcon      string "WIP(SERVER STATUS ICON)"
	ServerStatusStatement string "WIP(SERVER STATUS STATEMENT)"
}

type User struct {
	UserID    string `json:"userid" db:"userid, primarykey"`
	PWHash    string `json:"pwhash" db:"pwhash"`
	Attribute int    `json:"attribute" db:"attribute"`
}

type TokenClaims struct {
	Attribute int
	jwt.StandardClaims
}

type UUID guuid.UUID

func (u *UUID) Scan(src interface{}) error {
	var uuid_slice []byte
	var ok bool
	if uuid_slice, ok = src.([]byte); !ok {
		return fmt.Errorf("DB returned not []byte value for UUID.\n")
	}

	if len(uuid_slice) != 16 {
		return fmt.Errorf("uuid_slice must size of 16 but %d\n", len(uuid_slice))
	}

	copy(u[:], uuid_slice[0:16])
	return nil
}

func (u UUID) Value() (driver.Value, error) {
	return u[:], nil
}

func (u UUID) MarshalJSON() ([]byte, error) {
	uuid_str := (guuid.UUID(u)).String()
	uuid_str = "\"" + uuid_str + "\""
	res := []byte(uuid_str)
	return res, nil
}

type Event struct {
	Id                                   UUID      `json:"id" db:"id, primarykey"`
	Startdate                            time.Time `json:"startdate" db:"startdate"`
	Track                                string    `json:"track" db:"track"`
	WeatherRandomness                    int       `json:"weather_randomness" db:"weatherRandomness"`
	P_hourOfDay                          int       `json:"P_hourOfDay" db:"P_hourOfDay"`
	P_timeMultiplier                     int       `json:"P_timeMultiplier" db:"P_timeMultiplier"`
	P_sessionDurationMinute              int       `json:"P_sessionDurationMinute" db:"P_sessionDurationMinute"`
	Q_hourOfDay                          int       `json:"Q_hourOfDay" db:"Q_hourOfDay"`
	Q_timeMultiplier                     int       `json:"Q_timeMultiplier" db:"Q_timeMultiplier"`
	Q_sessionDurationMinute              int       `json:"Q_sessionDurationMinute" db:"Q_sessionDurationMinute"`
	R_hourOfDay                          int       `json:"R_hourOfDay" db:"R_hourOfDay"`
	R_timeMultiplier                     int       `json:"R_timeMultiplier" db:"R_timeMultiplier"`
	R_sessionDurationMinute              int       `json:"R_sessionDurationMinute" db:"R_sessionDurationMinute"`
	PitWindowLengthSec                   int       `json:"pit_window_length_sec" db:"pitWindowLengthSec"`
	IsRefuellingAllowedInRace            bool      `json:"is_refuelling_allowed_in_race" db:"isRefuellingAllowedInRace"`
	MandatoryPitstopCount                int       `json:"mandatory_pitstop_count" db:"mandatoryPitstopCount"`
	IsMandatoryPitstopRefuellingRequired bool      `json:"is_mandatory_pitstop_refuelling_required" db:"isMandatoryPitstopRefuellingRequired"`
	IsMandatoryPitstopTyreChangeRequired bool      `json:"is_mandatory_pitstop_tyre_change_required" db:"isMandatoryPitstopTyreChangeRequired"`
	IsMandatoryPitstopSwapDriverRequired bool      `json:"is_mandatory_pitstop_swap_driver_required" db:"isMandatoryPitstopSwapDriverRequired"`
	TyreSetCount                         int       `json:"tyre_set_count" db:"tyreSetCount"`
}

func (e *Event) FieldMap(r *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&e.Id: binding.Field{
			Form: "id",
			Binder: func(fieldName string, formVals []string, errors binding.Errors) binding.Errors {
				uuid, err := guuid.Parse(formVals[0])
				if err != nil {
					errors = append(errors, err.(binding.Error))
				}
				e.Id = [16]byte(uuid)
				return errors
			},
		},
		&e.Startdate:                            "startdate",
		&e.Track:                                "track",
		&e.WeatherRandomness:                    "weather_randomness",
		&e.P_hourOfDay:                          "P_hourOfDay",
		&e.P_timeMultiplier:                     "P_timeMultiplier",
		&e.P_sessionDurationMinute:              "P_sessionDurationMinute",
		&e.Q_hourOfDay:                          "Q_hourOfDay",
		&e.Q_timeMultiplier:                     "Q_timeMultiplier",
		&e.Q_sessionDurationMinute:              "Q_sessionDurationMinute",
		&e.R_hourOfDay:                          "R_hourOfDay",
		&e.R_timeMultiplier:                     "R_timeMultiplier",
		&e.R_sessionDurationMinute:              "R_sessionDurationMinute",
		&e.PitWindowLengthSec:                   "pit_window_length_sec",
		&e.IsRefuellingAllowedInRace:            "is_refuelling_allowed_in_race",
		&e.MandatoryPitstopCount:                "mandatory_pitstop_count",
		&e.IsMandatoryPitstopRefuellingRequired: "is_mandatory_pitstop_refuelling_required",
		&e.IsMandatoryPitstopTyreChangeRequired: "is_mandatory_pitstop_tyre_change_required",
		&e.IsMandatoryPitstopSwapDriverRequired: "is_mandatory_pitstop_swap_driver_required",
		&e.TyreSetCount:                         "tyre_set_count",
	}
}

type StartStop struct {
	Id       UUID      `db:"id, primarykey"`
	EventID  UUID      `db:"eventid"`
	Op       string    `db:"op"`
	Datetime time.Time `db:"datetime"`
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
