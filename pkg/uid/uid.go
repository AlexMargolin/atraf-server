package uid

import "github.com/google/uuid"

type UID = uuid.UUID

var Nil = uuid.Nil

// New provides a wrapper around Google's UUID package
func New() UID {
	return uuid.New()
}

func FromString(string string) (UID, error) {
	return uuid.Parse(string)
}

func ToString(uid UID) string {
	return uid.String()
}
