package server

import (
	"encoding/json"

	"github.com/gocraft/dbr"
	"github.com/pkg/errors"
)

type UserInput struct {
	Age  dbr.NullString
	Name dbr.NullString
}
type UserOutput struct {
	CreatedAt dbr.NullString
	Name      dbr.NullString
	Age       dbr.NullString
}

func NewUserInput(b []byte) (UserInput, error) {
	in := UserInput{}
	if err := json.Unmarshal(b, &in); err != nil {
		return in, errors.Wrap(err, "invalid JSON format")
	}
	if err := in.validate(); err != nil {
		return in, errors.Wrap(err, "validation error")
	}
	return in, nil
}
func (in UserInput) validate() error {
	return nil
}
