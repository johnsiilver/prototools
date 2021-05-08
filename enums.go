package prototools

import (
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Rec is a record of the value and names for an enumeration value.
type Rec struct {
	// EnumName is the name of the enumerator this value belongs to.
	EnumName string
	// Int32 is the enumerators numeric value.
	Int32 int32
	// ProtoName is the name of the enumerator in the format of the proto file.
	ProtoName string
	// JSONName is the name of the enumerator in JSON format.
	JSONName string
	// Titled is the name of the enumerator in sentence structure, without the leading [char]_ and
	// with each word titled.
	TitledName string
}

// ForwardLookup provides a map of varying spellings of an enumerator to its int32 value. These spellings include:
// ProtoName, JSONName, all lowercase name, minus leading [chars]_ titled name with spaces (aka PR_WHATEVER_VALUE is "Whatever Value").
type ForwardLookup map[string]Rec

// Find finds the record associated with that enumerator value name.
func (f ForwardLookup) Find(name string) (Rec, bool) {
	r := f[name]
	if r.ProtoName == "" {
		return r, false
	}
	return r, true
}

// ReverseLookup provides a lookup of enumerator values to string if you know the enumerators name in proto form.
type ReverseLookup map[string]map[int32]Rec

// Find returns the name of the enumerator value in proto string format. Empty string if not found.
func (r ReverseLookup) Find(enum string, value int32) (Rec, bool) {
	m, ok := r[enum]
	if !ok {
		return Rec{}, false
	}
	return m[value], true
}

// EnumLookup generates forward and reverse lookup maps for enumerators that are in msg or any child messages.
// Given the proto name for the enum values as seen in the proto file (not the generated Go files), an enum can
// be lookup in with:
// 	the proto name: PR_WHATEVER_YOU_PUT
//	the json name: prWhateverYouPut
// 	name minus(-) leading [char]_ as a sentence titled: PR_WHATEVER_YOU_PUT == "Whatever You Put"
//
// There are some caveats. This can cause name collision, so you should be careful how you name duplicates.
// Things like PR_UNKNOWN will be able to be lookup up via "unknown" or "Unknown". To prevent this, anything that
// ends in _unknown will be able to be looked up except by proto name and json name.
func EnumLookup(msgs []proto.Message) (ForwardLookup, ReverseLookup) {
	msgsParsed := map[string]bool{}
	enumsParsed := map[string]bool{}
	forward := ForwardLookup{}
	reverse := ReverseLookup{}

	for _, msg := range msgs {
		ref := msg.ProtoReflect().Descriptor()
		parseMsg(msgsParsed, enumsParsed, ref, forward, reverse)
	}

	return forward, reverse
}

func parseMsg(msgsParsed, enumsParsed map[string]bool, ref protoreflect.MessageDescriptor, forward ForwardLookup, reverse ReverseLookup) {
	if msgsParsed[string(ref.Name())] {
		return
	}

	for i := 0; i < ref.Fields().Len(); i++ {
		field := ref.Fields().Get(i)
		switch field.Kind() {
		case protoreflect.EnumKind:
			enum := field.Enum()
			enumName := enum.Name()
			if enumsParsed[string(enum.FullName())] {
				continue
			}
			enumsParsed[string(enum.FullName())] = true
			for x := 0; x < enum.Values().Len(); x++ {
				v := enum.Values().Get(x)
				num := v.Number() // int32 wrapper
				vName := v.Name()
				rec, ok := popForward(forward, string(enumName), string(vName), int32(num))
				if !ok {
					continue
				}
				popReverse(reverse, rec)
			}
		case protoreflect.MessageKind:
			parseMsg(msgsParsed, enumsParsed, field.Message(), forward, reverse)
		}
	}
}

func popForward(forward ForwardLookup, enumName, vName string, num int32) (Rec, bool) {
	rec := Rec{
		EnumName:   enumName,
		Int32:      num,
		ProtoName:  vName,
		JSONName:   JSONName(vName),
		TitledName: protoToTitled(vName),
	}

	if _, ok := forward.Find(rec.ProtoName); ok {
		return rec, false
	}

	forward[rec.ProtoName] = rec
	forward[rec.JSONName] = rec
	if notUnknown(rec.ProtoName) {
		forward[rec.TitledName] = rec
	}
	return rec, true
}

func popReverse(reverse ReverseLookup, rec Rec) {
	iToA, ok := reverse[rec.EnumName]
	if !ok {
		reverse[rec.EnumName] = map[int32]Rec{
			rec.Int32: rec,
		}
		return
	}
	iToA[rec.Int32] = rec
}

func notUnknown(s string) bool {
	sp := strings.Split(s, "_")
	if strings.ToLower(sp[len(sp)-1]) == "unknown" {
		return false
	}
	return true
}
