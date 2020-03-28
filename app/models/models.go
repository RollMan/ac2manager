package models

type User struct {
  UserID string `json:"userid"`
  PWHash string `json:"pwhash"`
  Attribute int `json:"attribute"`
}
