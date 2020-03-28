package models

type User struct {
  UserID string `json:"userid"`
  PWHash string `json:"pwhash"`
  Attribute int `json:"attribute"`
}

type NoSuchUserError struct {}
type NoMatchingPasswordError struct {}

func (e *NoSuchUserError) Error() string {
  return "No such userid in DB."
}

func (e *NoMatchingPasswordError) Error() string {
  return "Password unmatched."
}
