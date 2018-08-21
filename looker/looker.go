package looker

import (
	"path"
	"reflect"
)

type Package struct {
	ImportPath ImportElement
	Interfaces Interfaces
}

// ImportElement represents the imported package.
// Attribute `Alias` represents the optional alias for the package.
//		import foo "github.com/fooooo/baaaar-pppkkkkggg"
type ImportElement struct {
	Path  string
	Alias string
}

func (ie ImportElement) Name() string {
	if ie.Alias != "" {
		return ie.Alias
	}
	if ie.Path == "" {
		return ""
	}

	return path.Base(ie.Path)
}

func LookAtInterfaces(pkgPath string, is []reflect.Type) *Package {
	pkg := Package{
		ImportPath: ImportElement{Path: pkgPath},
		Interfaces: make(Interfaces, 0, len(is)),
	}
	for _, it := range is {
		intf := LookAtInterface(it)
		pkg.Interfaces = append(pkg.Interfaces, intf)
	}

	return &pkg
}

type Interface struct {
	Name    string
	Methods Methods
}

type Interfaces []*Interface

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

type Method struct {
	Name string
	In   Parameters
	Out  Parameters
}

type Methods []*Method

type Parameter interface {
	Kind() string
	Name() string
	Pointer() bool
}

type Parameters []Parameter

func LookAtFuncParameters(mt reflect.Type) (Parameters, Parameters) {
	var in = make(Parameters, 0)
	for i := 0; i < mt.NumIn(); i++ {
		in = append(in, LookAtParameter(mt.In(i)))
	}

	var out = make(Parameters, 0)
	for i := 0; i < mt.NumOut(); i++ {
		out = append(out, LookAtParameter(mt.Out(i)))
	}

	return in, out
}

// Use exported fields because god.Encoder
type StructElement struct {
	ImportPath ImportElement
	UserType   string
	IsPointer  bool
	Fields     Fields
}

func (prm *StructElement) Kind() string {
	return reflect.Struct.String()
}

func (prm *StructElement) Name() string {
	return prm.ImportPath.Name() + "." + prm.UserType
}

func (prm *StructElement) Pointer() bool {
	return prm.IsPointer
}

type SliceElement struct {
	ImportPath ImportElement
	UserType   string
	Item       Parameter
	IsPointer  bool
}

func (prm *SliceElement) Kind() string {
	return reflect.Slice.String()
}

func (prm *SliceElement) Name() string {
	if prm.UserType != "" {
		return prm.ImportPath.Name() + "." + prm.UserType
	}

	var ptr string
	if prm.Item.Pointer() {
		ptr = "*"
	}
	return "[]" + ptr + prm.Item.Name()
}

func (prm *SliceElement) Pointer() bool {
	return prm.IsPointer
}

type InterfaceElement struct {
	ImportPath ImportElement
	UserType   string
}

func (prm *InterfaceElement) Kind() string {
	return reflect.Interface.String()
}

func (prm *InterfaceElement) Name() string {
	if prm.ImportPath.Path == "" {
		return prm.UserType
	}
	return prm.ImportPath.Name() + "." + prm.UserType
}

func (prm *InterfaceElement) Pointer() bool {
	return false
}

type UnsupportedElement struct {
	ImportPath ImportElement
	UserType   string
	BaseType   string
	IsPointer  bool
}

func (prm *UnsupportedElement) Kind() string {
	return prm.BaseType
}

func (prm *UnsupportedElement) Name() string {
	return prm.ImportPath.Name() + "." + prm.UserType
}

func (prm *UnsupportedElement) Pointer() bool {
	return prm.IsPointer
}

func LookAtParameter(at reflect.Type) Parameter {
	var pointer bool
	if at.Kind() == reflect.Ptr {
		at = at.Elem()
		pointer = true
	}
	var prm Parameter

	switch at.Kind() {
	case reflect.Struct:
		prm = &StructElement{
			ImportPath: ImportElement{Path: at.PkgPath()},
			UserType:   at.Name(),
			IsPointer:  pointer,
			Fields:     LookAtFields(at),
		}
	case reflect.Slice:
		prm = &SliceElement{
			ImportPath: ImportElement{Path: at.PkgPath()},
			UserType:   at.Name(),
			IsPointer:  pointer,
			Item:       LookAtParameter(at.Elem()),
		}
	case reflect.Interface:
		prm = &InterfaceElement{
			ImportPath: ImportElement{Path: at.PkgPath()},
			UserType:   at.Name(),
		}
	default:
		prm = &UnsupportedElement{
			ImportPath: ImportElement{Path: at.PkgPath()},
			UserType:   at.Name(),
			BaseType:   at.Kind().String(),
			IsPointer:  pointer,
		}
	}

	return prm
}

type Field struct {
	Name       string
	ImportPath ImportElement
	BaseType   string
	UserType   string
	Anonymous  bool
}

type Fields []Field

func LookAtFields(st reflect.Type) Fields {
	fields := make(Fields, 0, st.NumField())
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		fields = append(fields, LookAtField(ft))
	}
	return fields
}

func LookAtField(ft reflect.StructField) Field {
	return Field{
		Name:       ft.Name,
		ImportPath: ImportElement{Path: ft.Type.PkgPath()},
		BaseType:   ft.Type.Kind().String(),
		UserType:   ft.Type.Name(),
		Anonymous:  ft.Anonymous,
	}
}
