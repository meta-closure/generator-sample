package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gocraft/dbr"
	"github.com/pkg/errors"
)

func PostUserHandler(w http.ResponseWriter, r *http.Request) error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "read request body")
	}

	in, err := NewUserInput(b)
	if err != nil {
		return errors.Wrap(err, "reading user input")
	}

	out := UserOutput{
		Name:      in.Name,
		Age:       in.Age,
		CreatedAt: dbr.NewNullString(time.Now()),
	}
	fmt.Println(out)
	return nil
}
