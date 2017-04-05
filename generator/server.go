package generator

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"go/format"

	"io/ioutil"

	hschema "github.com/lestrrat/go-jshschema"
	schema "github.com/lestrrat/go-jsschema"
	"github.com/pkg/errors"
)

const (
	intefaceFilePath = "./server/interface_gen.go"
	routeFilePath    = "./server/route_gen.go"
)

type Interface struct {
	Links  []Link
	output []byte
}

type Link struct {
	Input  Object
	Output Object
}

type Object struct {
	Title      string
	Properties []Property
}

type Property struct {
	Key       string
	Title     string
	Type      string
	Validates []Validate
}

type Validate struct {
	Required bool
}

func (s Schema) NewInterface() (Interface, error) {
	i := Interface{}
	for _, l := range s.Schema.Links {
		link, err := s.NewLink(l)
		if err != nil {
			return i, err
		}
		i.Links = append(i.Links, link)
	}
	return i, nil
}

func (s Schema) NewLink(link *hschema.Link) (Link, error) {
	in, err := s.NewObject(link.Schema)
	if err != nil {
		return Link{}, errors.Wrap(err, "parse input schema")
	}

	out, err := s.NewObject(link.TargetSchema)
	if err != nil {
		return Link{}, errors.Wrap(err, "parse output schema")
	}

	return Link{Input: in, Output: out}, nil
}

func (s Schema) NewObject(schema *schema.Schema) (Object, error) {
	rs, err := resolve(s.Schema, schema)
	if err != nil {
		return Object{}, err
	}
	o := Object{Title: upperCamelCase(rs.Title)}
	for k, v := range rs.Properties {
		p, err := s.NewProperty(k, v)
		if err != nil {
			return o, err
		}
		o.Properties = append(o.Properties, p)
	}
	return o, nil
}

func (s Schema) NewProperty(key string, schema *schema.Schema) (Property, error) {
	rs, err := resolve(s.Schema, schema)
	if err != nil {
		return Property{}, err
	}
	return Property{
		Key:   key,
		Title: upperCamelCase(rs.Title),
		Type:  convertGoType(rs.Type[0]),
	}, nil
}

func resolve(root *hschema.HyperSchema, schema *schema.Schema) (*schema.Schema, error) {
	if root == nil {
		return nil, errors.New("root is null pointer")
	}
	if schema == nil {
		return nil, errors.New("schema is null pointer")
	}

	rs, err := schema.Resolve(root)
	if err != nil {
		return nil, errors.New("resolve href on schema")
	}
	return rs, nil
}

func (i *Interface) Generate() error {
	b := bytes.NewBufferString("package server\n")
	writeImports(b, []string{
		"encoding/json",
		"github.com/gocraft/dbr",
		"github.com/pkg/errors",
	})

	for _, l := range i.Links {
		writeInterface(b, l.Input)
		writeInterface(b, l.Output)
		writeConstructor(b, l.Input)
		writeValidate(b, l.Input)
	}

	fb, err := format.Source(b.Bytes())
	if err != nil {
		return errors.Wrap(err, "go format")
	}
	i.output = fb
	return nil
}

func (i Interface) Save() error {
	return ioutil.WriteFile(intefaceFilePath, i.output, 0777)
}

type Routing struct {
	Routes []Route
	output []byte
}

type Route struct {
	HandlerName string
	Path        string
}

func (s Schema) NewRouting() (Routing, error) {
	r := Routing{}
	for _, link := range s.Schema.Links {
		r.Routes = append(r.Routes, Route{
			Path:        link.Href,
			HandlerName: upperCamelCase(link.Title),
		})

	}
	return r, nil
}

func (r *Routing) Generate() error {
	b := bytes.NewBufferString("package server\n")
	writeImports(b, []string{"net/http"})
	writeHandler(b, r.Routes)

	fb, err := format.Source(b.Bytes())
	if err != nil {
		return errors.Wrap(err, "go format")
	}
	r.output = fb
	return nil
}

func (r Routing) Save() error {
	return ioutil.WriteFile(routeFilePath, r.output, 0777)
}

func writeInterface(b io.Writer, o Object) {
	fmt.Fprintf(b, "type %s struct {\n", o.Title)
	for _, p := range o.Properties {
		fmt.Fprintf(b, "%s %s\n", p.Title, p.Type)
	}
	fmt.Fprintf(b, "}\n")
}

func writeConstructor(b io.Writer, o Object) {
	fmt.Fprintf(b, "func New%s(b []byte) (%s, error) {\n", o.Title, o.Title)
	fmt.Fprintf(b, "	in := %s{}\n", o.Title)
	fmt.Fprintf(b, " 	if err := json.Unmarshal(b, &in); err != nil {\n")
	fmt.Fprintf(b, "   		return in, errors.Wrap(err, \"invalid JSON format\")\n")
	fmt.Fprintf(b, "	 	}\n")
	fmt.Fprintf(b, "		if err := in.validate(); err != nil {\n")
	fmt.Fprintf(b, "	  		return in, errors.Wrap(err, \"validation error\")\n")
	fmt.Fprintf(b, "		}\n")
	fmt.Fprintf(b, "	return in, nil\n")
	fmt.Fprintf(b, "}\n")
}
func writeValidate(b io.Writer, o Object) {
	fmt.Fprintf(b, "func (in %s) validate() error {\n", o.Title)
	for _, p := range o.Properties {
		for _, v := range p.Validates {
			if v.Required {
				writeRequired(b, p)
			}
		}
	}
	fmt.Fprintf(b, "	return nil\n")
	fmt.Fprintf(b, "}\n")
}

func writeRequired(b io.Writer, p Property) {
	fmt.Fprintf(b, "{\n")
	fmt.Fprintf(b, "	if !in.%s.Valid {\n", p.Title)
	fmt.Fprintf(b, "		return errors.New(\"%s is not found\")", p.Title)
	fmt.Fprintf(b, "	}\n")
	fmt.Fprintf(b, "}\n")
}

func writeHandler(b io.Writer, routes []Route) {
	fmt.Fprintf(b, "func Run() error {\n")
	for _, r := range routes {
		writeRouting(b, r.HandlerName, r.Path)
	}
	fmt.Fprintf(b, "return http.ListenAndServe(\":8080\", nil)\n")
	fmt.Fprintf(b, "}\n")
}

func writeRouting(b io.Writer, handlerName, path string) {
	fmt.Fprintf(b, "http.HandleFunc(\"%s\", Hook{handler: %sHandler}.Handler)\n", path, handlerName)
}

func writeImports(b io.Writer, packages []string) {
	fmt.Fprintf(b, "import (\n")
	for _, p := range packages {
		fmt.Fprintf(b, "\"%s\"\n", p)
	}
	fmt.Fprintf(b, ")\n")
}

func upperCamelCase(s string) string {
	return strings.Join(strings.Split(strings.Title(strings.ToLower(s)), " "), "")
}

func convertGoType(s schema.PrimitiveType) string {
	switch s {
	case schema.StringType:
		return "dbr.NullString"
	case schema.IntegerType:
		return "dbr.NullInt64"
	default:
		return ""
	}
}
