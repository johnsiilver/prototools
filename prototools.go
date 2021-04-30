// Package prototools provide functions for performing reflection based operations on protocol buffers.
// These can be useful when extracting or updating proto fields based on field names.
package prototools

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
)

//go:generate stringer -type=ErrCode

// ErrCode reprsents an error code.
type ErrCode int8

const (
	// ErrUnknown means the code wasn't set.
	ErrUnknown ErrCode = 0
	// ErrIntermediateNotMessage means that one of the intermediate fields
	// was not a proto Message type. Aka, msg1.msg2.field , if msg1 or msg2
	// was not a Message.
	ErrIntermediateNotMessage ErrCode = 1
	// ErrIntermdiateNotSet indicates that one of the intermediates messsages was
	// nil.
	ErrIntermdiateNotSet ErrCode = 2
	// ErrBadFieldName indicates that one the fields did not exist in the message.
	// This is not the same as a message having a nil value, which is ErrIntermdiateNotSet.
	ErrBadFieldName = 3
)

// Error is our internal error types with error codes.
type Error struct {
	Code ErrCode
	Msg  string
}

// Error implements error.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

// Errorf is like fmt.Errorf but with our coded error type.
func Errorf(code ErrCode, msg string, i ...interface{}) Error {
	return Error{Code: code, Msg: fmt.Sprintf(msg, i...)}
}

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

// FieldAsStr returns the content of the field as a string. If pretty is set, it will try to pretty
// an enumerator by chopping off the text before the first "_", replacing the rest with a space, and
// doing a string.Title() on all the words. Aka: TYPE_UNKNOWN_DEVICE become: "Unknown Device". A user
// should not depend on the output of this string, as this may change over time without warning.
// If the field is _time and an int64, it is assumed to be unix time(epoch) in nanoseconds. If the field is
// a message, we protojson.Marshal() it. float or double values are printed out with 2 decimal places rounded up.
// We only support these values: boo, string, int32, int64, float, double, enum and message. We do not supports groups (repeated).
func FieldAsStr(msg proto.Message, fqPath string, pretty bool) (string, error) {
	fv, err := GetField(msg, fqPath)
	if err != nil {
		return "", err
	}

	switch fv.Kind {
	case protoreflect.BoolKind:
		if pretty {
			return strings.Title(fmt.Sprintf("%v", fv.Value)), nil
		}
		return fmt.Sprintf("%v", fv.Value), nil
	case protoreflect.StringKind:
		return fv.Value.(string), nil
	case protoreflect.BytesKind:
		return fmt.Sprintf("[%d]bytes", len(fv.Value.([]byte))), nil
	case protoreflect.Int32Kind:
		return fmt.Sprintf("%v", fv.Value), nil
	case protoreflect.Int64Kind:
		if strings.HasSuffix(fqPath, "_time") {
			t := time.Unix(fv.Value.(int64), 0).Truncate(0).UTC()
			return fmt.Sprintf("%v", t), nil
		}
		return fmt.Sprintf("%v", fv.Value), nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return fmt.Sprintf("%.2f", fv.Value), nil
	case protoreflect.EnumKind:
		if pretty {
			return prettyEnum(string(fv.EnumDesc.Name())), nil
		}
		return string(fv.EnumDesc.Name()), nil
	case protoreflect.MessageKind:
		b, err := protojson.Marshal(fv.Value.(proto.Message))
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", fmt.Errorf("type not supported")
}

func prettyEnum(s string) string {
	sp := strings.Split(s, "_")
	if len(sp) == 1 {
		return strings.Title(strings.ToLower(s))
	}
	sp = sp[1:]
	b := &strings.Builder{}
	for i, word := range sp {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(strings.Title(strings.ToLower(word)))
	}
	return b.String()
}

// FQPathSplit separates fqpath at ".".
func FQPathSplit(fqpath string) []string {
	return strings.Split(fqpath, ".")
}

// FQPathField extracts just the field's name.
func FQPathField(fqpath string) string {
	sp := FQPathSplit(fqpath)
	if len(sp) == 0 {
		return ""
	}
	return sp[len(sp)-1]
}

// FieldValue provides the proto value of a field.
type FieldValue struct {
	// Value is Go value of that field. Enumerators are of type protoreflect.EnumNumber
	// which is an int32.
	Value interface{}
	// Kind is the proto Kind that was stored.
	Kind protoreflect.Kind
	// EnumDesc is the enumerator descriptor if the Kind was EnumKind.
	// Usually this is used to call .Name() to get the text string representation
	// or FullName() if you want the package path + name.
	EnumDesc protoreflect.EnumValueDescriptor
}

/*
GetField searches into a proto to get a field value. It returns the value as
an interface{}, the kind of the field and if the field was found. You use a "."
notation to dive into the proto (field.field.field , where everything but the
last must be a Message type). We use the proto file spelling, not JSON or local
language spellings of the fields. You cannot look into groups (repeated values/array/slice...).

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
func GetField(msg proto.Message, fqPath string) (FieldValue, error) {
	fields := FQPathSplit(fqPath)
	for x, field := range fields[0 : len(fields)-1] {
		fv, err := fieldValue(msg, field)
		if err != nil {
			return FieldValue{}, Errorf(ErrBadFieldName, "field(%s) could not be found", strings.Join(fields[0:x], "."))
		}
		if fv.Kind != protoreflect.MessageKind {
			return FieldValue{}, Errorf(ErrIntermediateNotMessage, "field(%s) should be a message, was a %s", strings.Join(fields[0:x], "."), fv.Kind)
		}
		if fv.Value == nil {
			return FieldValue{}, Errorf(ErrIntermdiateNotSet, "message field(%s) is an empty message", strings.Join(fields[0:x], "."))
		}
		msg = fv.Value.(proto.Message)
	}

	fv, err := fieldValue(msg, fields[len(fields)-1])
	if err != nil {
		return FieldValue{}, Errorf(ErrBadFieldName, "field(%s) could not be found", fqPath)
	}
	return fv, nil
}

// fieldValue gets a field from msg.
func fieldValue(msg proto.Message, field string) (FieldValue, error) {
	ref := msg.ProtoReflect()
	descriptors := ref.Descriptor().Fields()
	fd := descriptors.ByName(protoreflect.Name(field))
	if fd == nil {
		return FieldValue{}, errors.New("bad field name")
	}

	switch fd.Kind() {
	case protoreflect.MessageKind:
		return FieldValue{
			Value: ref.Get(fd).Message().Interface(),
			Kind:  protoreflect.MessageKind,
		}, nil
	case protoreflect.EnumKind:
		i := ref.Get(fd).Interface()
		enumDesc := fd.Enum().Values().ByNumber(i.(protoreflect.EnumNumber))
		return FieldValue{
			Value:    protoreflect.ValueOfEnum(enumDesc.Number()).Interface(),
			Kind:     protoreflect.EnumKind,
			EnumDesc: enumDesc,
		}, nil
	}
	return FieldValue{
		Value: ref.Get(fd).Interface(),
		Kind:  fd.Kind(),
	}, nil
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

// HumanDiff is a wrapper aound go-cmp using the protocmp.Transform. It outputs a string of what changes from a (older) to b (newer).
// Options to pass can be found at: https://pkg.go.dev/google.golang.org/protobuf/testing/protocmp .
func HumanDiff(a, b proto.Message, options ...cmp.Option) string {
	options = append(options, protocmp.Transform())
	return cmp.Diff(a, b, options...)
}

// Equal returns an empty string if the two protos are equal. Otherwise it returns -want/+got.
// This is the same as HumanDiff, but order is reversed.
func Equal(want, got proto.Message, options ...cmp.Option) string {
	options = append(options, protocmp.Transform())
	return cmp.Diff(want, got, options...)
}

/*
// FieldChange details a field that changed in a proto. If the field is a non-message, .From* and .To*
// will contain the value.
type FieldChange struct {
	// Name is the name of the field.
	Name string
	// Path is the path to that field in the proto you diffed. The field will be at:
	// strings.Join(fc.Path, ".") + "." + fc.Name .
	Path []string
	Kind protoreflect.Kind
	From, To interface{}
}

// FQPath will return the fully qualified path to the value.
func (f FieldChange) FQPath() string {
	return strings.Join(f.Path, ".") + "." + f.Name
}

func DiffField(msg proto.Message, fqPath string) (FieldChange, error) {

}
*/

// this code is borrowed and modified faith code(github.com/fatih/camelcase)
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
