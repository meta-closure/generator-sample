package generator

import (
	"encoding/json"
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
	y, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Schema{}, errors.Wrap(err, "open file")
	}

	j, err := yaml.YAMLToJSON(y)
	if err != nil {
		return Schema{}, errors.Wrap(err, "invalid YAML format")
	}

	h := hschema.New()
	if err := json.Unmarshal(j, h); err != nil {
		return Schema{}, errors.Wrap(err, "invalid JSON Hyper Schema format")
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
