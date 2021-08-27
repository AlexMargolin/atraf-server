package uid

import "github.com/google/uuid"

// UID is an alias for uuid.UUID
type UID = uuid.UUID

// Nil is an alias for uuid.Nil
var Nil = uuid.Nil

// New is a wrapper for uuid.New
func New() UID {
	return uuid.New()
}

// FromString is a wrapper for uuid.Parse
func FromString(string string) (UID, error) {
	return uuid.Parse(string)
}
