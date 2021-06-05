/*
Package prototools provide functions for performing reflection based operations on protocol buffers.
These can be useful when extracting or updating proto fields based on field names.

It should be noted that this package makes some assumptions and lacks certain features at the moment.

The big assumption is how fields are named.  Best practice says fields are lower_case_seperated_with_underscore.
Enuerators are CAPITALS_WITH_A_LEADING_WORD_FOR_UNIQUENESS_WITH_UNDERSCORES. Many things here might not work
as expected if you are not following these guidelines.

There are two big things that we mostly ignore, maps and arrays. We just don't introspect them except where noted
as that gets complicated and I don't need the capability at the moment.

Finally, I am mostly ignoring all the "fixed" types, Any and whatever the types were before Any (my brain can't remember
what those were called, I wouldn't even use them when I worked at Google, so not doing it here).

This "may" work with OneOf's, but I haven't tried.

You might ask yourself, why even bother if you don't do these?  Well, most of the time for what this package will get used for,
which is data exchange for web stuff, these things don't matter. Again, fits my purpose for the moment.  If I need more
complicated stuff, I'll add it at a later date.
*/
package prototools

import (
	"errors"
	"fmt"
	"reflect"
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
	// ErrNotMessage indicates that a value in the path is not a message. Commonly this happens when
	// trying to retrieve a value from a repeated message or map. You cannot pull this directly, you
	// must get the repeated messaged and then look through each value.
	ErrNotMessage = 4
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

type readableOpts struct {
	removePrefix bool
}

type ReadableOption func(r *readableOpts)

// RemovePrefix will remove the prefix word on the field name. In proto format,
// this is anything before (and including) the first underscore(_) character.
func RemovePrefix() ReadableOption {
	return func(r *readableOpts) {
		r.removePrefix = true
	}
}

// ReadableProto slits the proto field name at "_" and titles each word.
// This assumes ASCII names and following [lower]_[seperated]_[with]_[underscores] .
func ReadableProto(s string, options ...ReadableOption) string {
	if len(s) == 0 {
		return s
	}
	opts := readableOpts{}
	for _, o := range options {
		o(&opts)
	}

	words := strings.Split(s, "_")

	if opts.removePrefix && len(words) > 1 {
		words = words[1:]
	}

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
func FieldAsStr(msg proto.Message, fqPath string, pretty bool) (string, protoreflect.Kind, error) {
	fv, err := GetField(msg, fqPath)
	if err != nil {
		return "", 0, err
	}

	switch fv.Kind {
	case protoreflect.BoolKind:
		if pretty {
			return strings.Title(fmt.Sprintf("%v", fv.Value)), fv.Kind, nil
		}
		return fmt.Sprintf("%v", fv.Value), fv.Kind, nil
	case protoreflect.StringKind:
		return fv.Value.(string), fv.Kind, nil
	case protoreflect.BytesKind:
		return fmt.Sprintf("[%d]bytes", len(fv.Value.([]byte))), fv.Kind, nil
	case protoreflect.Int32Kind:
		return fmt.Sprintf("%v", fv.Value), fv.Kind, nil
	case protoreflect.Int64Kind:
		if strings.HasSuffix(fqPath, "_time") {
			t := time.Unix(fv.Value.(int64), 0).Truncate(0).UTC()
			return fmt.Sprintf("%v", t), fv.Kind, nil
		}
		return fmt.Sprintf("%v", fv.Value), fv.Kind, nil
	case protoreflect.FloatKind, protoreflect.DoubleKind:
		return fmt.Sprintf("%.2f", fv.Value), fv.Kind, nil
	case protoreflect.EnumKind:
		if pretty {
			return prettyEnum(string(fv.EnumDesc.Name())), fv.Kind, nil
		}
		return string(fv.EnumDesc.Name()), fv.Kind, nil
	case protoreflect.MessageKind:
		b, err := protojson.Marshal(fv.Value.(proto.Message))
		if err != nil {
			return "", fv.Kind, err
		}
		return string(b), fv.Kind, nil
	}
	return "", fv.Kind, fmt.Errorf("type not supported")
}

func protoToTitled(s string) string {
	sp := strings.Split(s, "_")
	sp = sp[1:]
	for i, w := range sp {
		sp[i] = strings.Title(strings.ToLower(w))
	}
	return strings.Join(sp, " ")
}

// TODO(jdoak): This and protoToTitled are the same. Test which is more efficient
// and use that one.
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
	// which is an int32. Be aware that if the Kind is a MessageKind, this value
	// can be a nil value stored in an interface{}. That means Value != nil, but
	// the value inside is nil (yeah, I know). Had to do this for certain reasons.
	// There is a IsNil() method if you want to test if the stored value == nil.
	// If the value is a list and the type is a base type, then it will be []<type>.
	// However, it it is a message, this will be []protoreflect.Message, because we have no way
	// to get to the concrete type. You can use Kind to determine what you need to do.
	// If it is an Enum, we leave it as a protoreflect.EnumNumber, which is really an int32.
	Value interface{}
	// Kind is the proto Kind that was stored. If IsList == true, Kind represents the underlying
	// value stored in the list.
	Kind protoreflect.Kind
	// IsList is set if the Kind == MessageKind, but the message represents a repeated value.
	IsList bool
	// FieldDesc is the field descriptor for this value.
	FieldDesc protoreflect.FieldDescriptor
	// EnumDesc is the enumerator descriptor if the Kind was EnumKind.
	// Usually this is used to call .Name() to get the text string representation
	// or FullName() if you want the package path + name.
	EnumDesc protoreflect.EnumValueDescriptor
	// MsgDesc is the message descriptor if the Kind was MessageKind.
	MsgDesc protoreflect.MessageDescriptor
}

// IsNil determins if the value stored in .Value is nil.
func (f FieldValue) IsNil() bool {
	if f.Value == nil {
		return true
	}
	return reflect.ValueOf(f.Value).IsNil()
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
		var ok bool
		msg, ok = fv.Value.(proto.Message)
		if !ok {
			return FieldValue{}, Errorf(ErrNotMessage, "message field(%s) is not a proto.Message. Usually this means you are trying to retrieve inside a repeated field or map, which cannot be done", strings.Join(fields[0:x], "."))
		}
	}

	fv, err := fieldValue(msg, fields[len(fields)-1])
	if err != nil {
		return FieldValue{}, Errorf(ErrBadFieldName, "field(%s) could not be found", fqPath)
	}
	return fv, nil
}

// getLastMessage will takes the path and returns the len(fqPath) -1 proto.Message. If createmessage is set, this will create the
// message values if they are not set through the entire path except the inital message passed as "msg".
func getLastMessage(msg proto.Message, fqPath []string, createMessages bool) (protoreflect.Message, error) {
	fields := fqPath[0 : len(fqPath)-1]
	for x, field := range fields {
		fv, err := fieldValue(msg, field)
		if err != nil {
			return nil, Errorf(ErrBadFieldName, "field(%s) could not be found", strings.Join(fields[0:x], "."))
		}
		if fv.Kind != protoreflect.MessageKind {
			return nil, Errorf(ErrIntermediateNotMessage, "field(%s) should be a message, was a %s", strings.Join(fields[0:x], "."), fv.Kind)
		}
		if fv.IsNil() {
			if createMessages {
				fieldMsg := fv.Value.(proto.Message).ProtoReflect()
				n := fieldMsg.New()
				msg.ProtoReflect().Set(fv.FieldDesc, protoreflect.ValueOf(n))
				msg = n.Interface()
				continue
			}
			return nil, Errorf(ErrIntermdiateNotSet, "message field(%s) is an empty message", strings.Join(fields[0:x], "."))
		}
		msg = fv.Value.(proto.Message)
	}
	return msg.ProtoReflect(), nil
}

// fieldValue gets a field from msg.
func fieldValue(msg proto.Message, field string) (FieldValue, error) {
	ref := msg.ProtoReflect()
	descriptors := ref.Descriptor().Fields()
	fd := descriptors.ByName(protoreflect.Name(field))
	if fd == nil {
		return FieldValue{}, errors.New("bad field name")
	}

	switch {
	case fd.IsList():
		return listFieldValue(ref, fd)
	case fd.IsMap():
		return FieldValue{}, errors.New("we do not currently support maps")
	case fd.IsExtension():
		return FieldValue{}, errors.New("we do not currently support extensions")
	}

	switch fd.Kind() {
	case protoreflect.MessageKind:
		var v = ref.Get(fd).Message().Interface()
		return FieldValue{
			Value:     v,
			Kind:      protoreflect.MessageKind,
			FieldDesc: fd,
			MsgDesc:   fd.Message(),
		}, nil
	case protoreflect.EnumKind:
		i := ref.Get(fd).Interface()
		enumDesc := fd.Enum().Values().ByNumber(i.(protoreflect.EnumNumber))
		return FieldValue{
			Value:     protoreflect.ValueOfEnum(enumDesc.Number()).Interface(),
			Kind:      protoreflect.EnumKind,
			FieldDesc: fd,
			EnumDesc:  enumDesc,
		}, nil
	}
	return FieldValue{
		Value:     ref.Get(fd).Interface(),
		Kind:      fd.Kind(),
		FieldDesc: fd,
	}, nil
}

func listFieldValue(ref protoreflect.Message, fd protoreflect.FieldDescriptor) (FieldValue, error) {
	// These are currently unsupported.
	/*
		Int32Kind    Kind = 5
		Sint32Kind   Kind = 17
		Uint32Kind   Kind = 13
		Sint64Kind   Kind = 18
		Uint64Kind   Kind = 4
		Sfixed32Kind Kind = 15
		Fixed32Kind  Kind = 7
		Sfixed64Kind Kind = 16
		Fixed64Kind  Kind = 6
		BytesKind    Kind = 12
		GroupKind    Kind = 10
	*/

	fv := FieldValue{
		Kind:      fd.Kind(),
		IsList:    true,
		FieldDesc: fd,
	}
	var l = ref.Get(fd).List()
	switch fd.Kind() {
	case protoreflect.BoolKind:
		v := make([]bool, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = entry.Bool()
		}
		fv.Value = v
	case protoreflect.Int32Kind:
		v := make([]int32, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = int32(entry.Int())
		}
		fv.Value = v
	case protoreflect.Int64Kind:
		v := make([]int64, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = int64(entry.Int())
		}
		fv.Value = v
	case protoreflect.FloatKind:
		v := make([]float32, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = float32(entry.Float())
		}
		fv.Value = v
	case protoreflect.DoubleKind:
		v := make([]float64, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = entry.Float()
		}
		fv.Value = v
	case protoreflect.StringKind:
		v := make([]string, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = entry.String()
		}
		fv.Value = v
	case protoreflect.EnumKind:
		v := make([]protoreflect.EnumNumber, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = entry.Enum()
		}
		fv.Value = v
	case protoreflect.MessageKind:
		v := make([]protoreflect.Message, l.Len())
		for i := 0; i < l.Len(); i++ {
			entry := l.Get(i)
			v[i] = entry.Message()
		}
		fv.Value = v
		fv.MsgDesc = fd.Message()
	default:
		return FieldValue{}, fmt.Errorf("we do not support a list value of this type(%s)", fd.Kind())
	}
	return fv, nil
}

type enumDescriptor interface {
	Descriptor() protoreflect.EnumDescriptor
	Number() protoreflect.EnumNumber
}

// UpdateProtoField updates a field in a protocol buffer message with a value.
// The field is assumed to be the proto name format.
// This only supports values of string, int, int32, int64 and bool. An int updates an int64.
func UpdateProtoField(m proto.Message, fqPath string, value interface{}) error {
	fields := FQPathSplit(fqPath)
	if len(fields) == 0 {
		return fmt.Errorf("cannot send a path(%s) of zero len", fqPath)
	}
	fieldName := fields[len(fields)-1]

	v, err := getLastMessage(m, fields, false)
	if err != nil {
		return err
	}

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
		n := t.Number()
		enum := protoreflect.ValueOfEnum(n)
		v.Set(fd, enum)
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
