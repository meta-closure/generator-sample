package generater

import (
	"github.com/lestrrat/go-jshschema"
	"github.com/pkg/errors"
)

const (
	filePath = "schema.yml"
)

type Schema struct {
	Schema *hschema.HyperSchema
}

func NewSchema(filePath string) (Schema, error) {
	h, err := hschema.ReadFile(filePath)
	if err != nil {
		return Schema{}, errors.Wrap(err, "reading JSON Schema")
	}

	return Schema{
		Schema: h,
	}, nil
}

func (s Schema) generate() error {
	if err := s.generateServer(); err != nil {
		return err
	}
	return nil
}

func (s Schema) generateServer() error {
	i, err := s.NewInterface()
	if err != nil {
		return errors.Wrap(err, "create new interface")
	}
	if err := i.generate(); err != nil {
		return errors.Wrap(err, "generating server interface")
	}

	r, err := s.NewRouting()
	if err != nil {
		return errors.Wrap(err, "create new route")
	}
	if err := r.generate(); err != nil {
		return errors.Wrap(err, "generating server routing")
	}
	return nil
}

type Interface struct {
}

func (s Schema) NewInterface() (Interface, error) {
	return Interface{}, nil
}

func (i Interface) generate() error {
	return nil
}

type Routing struct {
}

func (s Schema) NewRouting() (Routing, error) {
	return Routing{}, nil
}

func (r Routing) generate() error {
	return nil
}

func Generate() error {
	s, err := NewSchema(filePath)
	if err != nil {
		return errors.Wrap(err, "create schema")
	}

	if err := s.generate(); err != nil {
		return errors.Wrap(err, "generate codes from schema")
	}
	return nil
}
