package generator

import (
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/lestrrat/go-jshschema"
	"github.com/pkg/errors"
)

const (
	filePath = "./schema.yml"
)

type Schema struct {
	Schema *hschema.HyperSchema
}

func NewSchema(filePath string) (Schema, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Schema{}, errors.Wrap(err, "open file")
	}

	m := map[string]interface{}{}
	if err := yaml.Unmarshal(b, &m); err != nil {
		return Schema{}, errors.Wrap(err, "convert to yaml")
	}

	h := hschema.New()
	if err := h.Extract(m); err != nil {
		return Schema{}, errors.Wrap(err, "convert to JSON Schema")
	}

	return Schema{
		Schema: h,
	}, nil
}

func (s Schema) Generate() error {
	if err := s.GenerateServer(); err != nil {
		return err
	}
	return nil
}

func (s Schema) GenerateServer() error {
	i, err := s.NewInterface()
	if err != nil {
		return errors.Wrap(err, "create new interface")
	}
	if err := i.Generate(); err != nil {
		return errors.Wrap(err, "generating server interface")
	}

	r, err := s.NewRouting()
	if err != nil {
		return errors.Wrap(err, "create new route")
	}
	if err := r.Generate(); err != nil {
		return errors.Wrap(err, "generating server routing")
	}

	if err := i.Save(); err != nil {
		return errors.Wrap(err, "save generated interface code")
	}

	if err := r.Save(); err != nil {
		return errors.Wrap(err, "save generated routing code")
	}
	return nil
}

func Generate() error {
	s, err := NewSchema(filePath)
	if err != nil {
		return errors.Wrap(err, "create schema")
	}

	if err := s.Generate(); err != nil {
		return errors.Wrap(err, "generate codes from schema")
	}
	return nil
}
