package server

import (
	"encoding/json"
	"fmt"

	"github.com/gocraft/dbr"
	"github.com/pkg/errors"
)

type UserInput struct {
	Name dbr.NullString
	Age  dbr.NullInt64
}

type UserOutput struct {
	Name      dbr.NullString
	Age       dbr.NullInt64
	CreatedAt dbr.NullString
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
	{
		if !in.Name.Valid {
			return errors.New("name is not found")
		}
	}
	{
		if !in.Age.Valid {
			return errors.New("age is not found")
		}
	}
	return nil
}

func (out UserOutput) Request() error {
	fmt.Printf("%+v", out)
	return nil
}
