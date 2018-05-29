package looker

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
)

func LookAtInterfaces(pkgPath string, is []reflect.Type) *Package {
	pkg := Package{
		PkgPath:    pkgPath,
		Interfaces: make(Interfaces, 0, len(is)),
	}
	for _, it := range is {
		intf := LookAtInterface(it)
		pkg.Interfaces = append(pkg.Interfaces, intf)
	}

	return &pkg
}

func LookAtInterface(typ reflect.Type) *Interface {
	intf := &Interface{
		Name:    typ.Name(),
		Methods: make(Methods, 0, typ.NumMethod()),
	}

	for i := 0; i < typ.NumMethod(); i++ {
		mt := typ.Method(i)
		m := Method{
			Name: mt.Name,
		}
		in, out := LookAtFuncParameters(typ.Method(i).Type)
		m.In = in
		m.Out = out

		intf.Methods = append(intf.Methods, &m)
	}
	return intf
}

func LookAtFuncParameters(mt reflect.Type) (Parameters, Parameters) {
	var in = make([]*Parameter, 0)
	for i := 0; i < mt.NumIn(); i++ {
		in = append(in, LookAtParameter(mt.In(i)))
	}

	var out = make([]*Parameter, 0)
	for i := 0; i < mt.NumOut(); i++ {
		out = append(out, LookAtParameter(mt.Out(i)))
	}

	return in, out
}

func LookAtParameter(at reflect.Type) *Parameter {
	var pointer bool
	if at.Kind() == reflect.Ptr {
		at = at.Elem()
		pointer = true
	}
	prm := Parameter{
		PkgPath:  at.PkgPath(),
		PkgName:  strings.Split(at.String(), ".")[0],
		BaseType: at.Kind().String(),
		UserType: at.Name(),
		Pointer:  pointer,
	}

	if prm.BaseType == reflect.Struct.String() {
		prm.Fields = LookAtFields(at)
	}

	return &prm
}

func LookAtFields(st reflect.Type) Fields {
	fields := make(Fields, 0, st.NumField())
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		field := &Field{
			Name:      ft.Name,
			PkgPath:   ft.Type.PkgPath(),
			PkgName:   strings.Split(ft.Type.String(), ".")[0],
			BaseType:  ft.Type.Kind().String(),
			UserType:  ft.Type.Name(),
			Anonymous: ft.Anonymous,
		}
		fields = append(fields, field)
	}
	return fields
}

type Package struct {
	PkgPath    string
	Interfaces Interfaces
}

type Interface struct {
	Name    string
	Methods Methods
}

type Interfaces []*Interface

type Method struct {
	Name string
	In   Parameters
	Out  Parameters
}

type Methods []*Method

type Parameter struct {
	PkgPath  string
	PkgName  string
	BaseType string
	UserType string
	Pointer  bool
	Fields   Fields
}

type Field struct {
	Name      string
	PkgPath   string
	PkgName   string
	BaseType  string
	UserType  string
	Anonymous bool
}

type Fields []*Field

type Parameters []*Parameter

func EncodeGob(output string, pkg *Package) error {
	outfile := os.Stdout

	if len(output) != 0 {
		var err error
		if outfile, err = os.Create(output); err != nil {
			return fmt.Errorf("failed to open output file %q: %s", output, err)
		}
		defer func() {
			if err := outfile.Close(); err != nil {
				fmt.Errorf("failed to close output file %q: %s", output, err)
			}
		}()
	}

	if err := gob.NewEncoder(outfile).Encode(pkg); err != nil {
		fmt.Errorf("gob encode: %s", err)
	}

	return nil
}

func p(kv ...interface{}) {
	log.Print(kv...)
}

func pf(s string, kv ...interface{}) {
	log.Printf(s, kv...)
}
