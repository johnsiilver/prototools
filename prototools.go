// Package prototools provide functions for performing reflection based operations on protocol buffers.
// These can be useful when extracting or updating proto fields based on field names.
package prototools

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var (
	ErrFieldNameExist = errors.New("the field with that name does not exist")
)

// JSONName converts the proto name of a field to the JSON equivalent.
// This assumes ASCII names and that the proto name kept best practices
// of name = [lower]_[seperated]_[with]_[underscores] .
func JSONName(protoName string) string {
	if len(protoName) == 0 {
		return protoName
	}

	sp := strings.Split(protoName, "_")

	for i, word := range sp {
		if i == 0 {
			sp[i] = strings.ToLower(word)
		} else {
			sp[i] = strings.Title(strings.ToLower(word))
		}
	}
	return strings.Join(sp, "")
}

// ProtoName converts the JSON name of a field to the proto equivalent.
// This is NOT the Go name, this is the name as seen in the proto file.
// We assume best practices of name = [lower]_[seperated]_[with]_[underscores].
// This is really an alias of strings.Title(jsonName).
func ProtoName(jsonName string) string {
	sp := split(jsonName, false)
	for i, word := range sp {
		sp[i] = strings.ToLower(word)
	}

	return strings.Join(sp, "_")
}

// ReadableJSON splits the JSON field name at capital letters and titles each word.
// This assumes ASCII names and that "s" is a JSON name for a field.
func ReadableJSON(s string) string {
	if len(s) == 0 {
		return s
	}

	words := split(s, true)
	if len(words) == 0 {
		return s
	}
	for i, word := range words {
		words[i] = strings.Title(word)
	}

	return strings.Join(words, " ")
}

// ReadableProto slits the proto field name at "_" and titles each word.
// This assumes ASCII names and following [lower]_[seperated]_[with]_[underscores] .
func ReadableProto(s string) string {
	if len(s) == 0 {
		return s
	}

	words := strings.Split(s, "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

// FieldAsStr returns the content of the field as a string.
func FieldAsStr(msg proto.Message, field string) (string, error) {
	i, k, err := FieldValue(msg, field)
	if err != nil {
		return "", err
	}

	if k == protoreflect.EnumKind {
		ref := msg.ProtoReflect()
		descriptors := ref.Descriptor().Fields()
		fd := descriptors.ByName(protoreflect.Name(field))
		enumDesc := fd.Enum().Values().ByNumber(i.(protoreflect.EnumNumber))
		return string(enumDesc.Name()), nil
	}

	return fmt.Sprintf("%v", i), nil
}

/*
FieldValue returns the value of a field as an interface{}, the kind of protocol field it is and an error if there is one.
The following is the kind to Go type mapping:

	╔════════════╤═════════════════════════════════════╗
	║ Go type    │ Protobuf kind                       ║
	╠════════════╪═════════════════════════════════════╣
	║ bool       │ BoolKind                            ║
	║ int32      │ Int32Kind, Sint32Kind, Sfixed32Kind ║
	║ int64      │ Int64Kind, Sint64Kind, Sfixed64Kind ║
	║ uint32     │ Uint32Kind, Fixed32Kind             ║
	║ uint64     │ Uint64Kind, Fixed64Kind             ║
	║ float32    │ FloatKind                           ║
	║ float64    │ DoubleKind                          ║
	║ string     │ StringKind                          ║
	║ []byte     │ BytesKind                           ║
	║ EnumNumber │ EnumKind                            ║
	║ Message    │ MessageKind, GroupKind              ║
	╚════════════╧═════════════════════════════════════╝

*/
func FieldValue(msg proto.Message, field string) (interface{}, protoreflect.Kind, error) {
	ref := msg.ProtoReflect()
	descriptors := ref.Descriptor().Fields()
	fd := descriptors.ByName(protoreflect.Name(field))
	if fd == nil {
		return nil, 0, fmt.Errorf("could not get field named %q: %w", field, ErrFieldNameExist)
	}

	return ref.Get(fd).Interface(), fd.Kind(), nil
}

type enumDescriptor interface {
	Descriptor() protoreflect.EnumDescriptor
	Number() protoreflect.EnumNumber
}

// UpdateProtoField updates a field in a protocol buffer message with a value.
// The field is assumed to be the proto name format.
// This only supports values of string, int, int32, int64 and bool. An int updates an int64.
func UpdateProtoField(m proto.Message, fieldName string, value interface{}) error {
	v := m.ProtoReflect()
	fd := v.Descriptor().Fields().ByName(protoreflect.Name(fieldName))
	if fd == nil {
		return fmt.Errorf("field %s not found", fieldName)
	}
	switch t := value.(type) {
	case string:
		if fd.Kind() != protoreflect.StringKind {
			return fmt.Errorf("field %s is a %s, you sent a string", fieldName, fd.Kind())
		}
		v.Set(fd, protoreflect.ValueOf(t))
	case int:
		if fd.Kind() != protoreflect.Int64Kind {
			return fmt.Errorf("field %s is a %s, you sent a int64", fieldName, fd.Kind())
		}
		v.Set(fd, protoreflect.ValueOf(int64(t)))
	case int64:
		if fd.Kind() != protoreflect.Int64Kind {
			return fmt.Errorf("field %s is a %s, you sent a int64", fieldName, fd.Kind())
		}
		v.Set(fd, protoreflect.ValueOf(t))
	case int32:
		switch fd.Kind() {
		case protoreflect.Int32Kind:
			v.Set(fd, protoreflect.ValueOf(t))
		case protoreflect.EnumKind:
			n := protoreflect.EnumNumber(t)
			if exists := fd.Enum().Values().ByNumber(n); exists == nil {
				return fmt.Errorf("field %s is an enum and %d is not a valid value", fieldName, t)
			}
			enum := protoreflect.ValueOfEnum(n)
			v.Set(fd, enum)
		default:
			return fmt.Errorf("field %s is a %s, you sent a int32", fieldName, fd.Kind())
		}
	case bool:
		if fd.Kind() != protoreflect.BoolKind {
			return fmt.Errorf("field %s is a %s, you sent a bool", fieldName, fd.Kind())
		}
		v.Set(fd, protoreflect.ValueOf(t))
	case enumDescriptor:
		n := int32(t.Number())
		return UpdateProtoField(m, fieldName, n)
	default:
		return fmt.Errorf("field %s cannot be set to %T, as that type isn't supported", fieldName, value)
	}
	return nil
}

// this code is borrowed and modified faith code.
// TODO(johnsiilver): Add attribution.
func split(src string, splitNum bool) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 1
			if splitNum {
				class = 3
			}
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return
}
