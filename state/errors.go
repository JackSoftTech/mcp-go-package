package state

import "errors"

var (
	ErrSchemaLocked  = errors.New("state schema already configured")
	ErrSchemaNotSet  = errors.New("state schema not configured")
	ErrSchemaEmpty   = errors.New("state schema is empty")
	ErrSchemaInvalid = errors.New("state schema is invalid")
	ErrStateInvalid  = errors.New("state does not match schema")
)
