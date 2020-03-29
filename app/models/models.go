package models

type Login struct {
	UserID   string `json:"userid"`
	Password string `json:"pw"`
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
