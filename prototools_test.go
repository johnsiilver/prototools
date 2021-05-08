package prototools

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kylelemons/godebug/pretty"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/johnsiilver/prototools/sample"
)

func TestJSONName(t *testing.T) {
	protoName := "this_is_my_field_name32"
	want := "thisIsMyFieldName32"
	got := JSONName(protoName)
	if got != want {
		t.Errorf("TestJSONName(%s): got %q, want %q", protoName, got, want)
	}
}

func TestProtoName(t *testing.T) {
	jsonName := "thisIsMyFieldName32"
	want := "this_is_my_field_name32"
	got := ProtoName(jsonName)
	if got != want {
		t.Errorf("TestProtoName(%s): got %q, want %q", jsonName, got, want)
	}
}

func TestReadableJSON(t *testing.T) {
	name := "thisIsMyFieldName32"
	want := "This Is My Field Name 32"
	got := ReadableJSON(name)
	if got != want {
		t.Errorf("TestReadableJSON(%s): got %q, want %q", name, got, want)
	}
}

func TestReadbleProto(t *testing.T) {
	name := "this_is_my_field_name_32"
	want := "This Is My Field Name 32"
	got := ReadableProto(name)
	if got != want {
		t.Errorf("TestReadableProto(%s): got %q, want %q", name, got, want)
	}
}

func TestFieldAsStr(t *testing.T) {
	data := &pb.Layer1{
		Supported: &pb.Supported{
			Vint32:  32,
			Vint64:  64,
			Vstring: "string",
			Vbool:   true,
			Ev:      pb.EnumValues_EV_Ok,
			VTime:   1619820228,
			Vfloat:  float32(3.4569),
			Vdouble: float64(8.9645),
		},
		Vstring: "Hello",
	}

	tests := []struct {
		desc   string
		field  string
		want   string
		pretty bool
	}{
		{"int32", "supported.vint32", "32", false},
		{"int64", "supported.vint64", "64", false},
		{"string", "supported.vstring", "string", false},
		{"bool", "supported.vbool", "true", false},
		{"bool", "supported.vbool", "True", true},
		{"enum", "supported.ev", "EV_Ok", false},
		{"enum", "supported.ev", "Ok", true},
		{"time", "supported.v_time", "2021-04-30 22:03:48 +0000 UTC", false},
		{"float", "supported.vfloat", "3.46", false},
		{"double", "supported.vdouble", "8.96", false},
	}

	for _, test := range tests {
		got, _, err := FieldAsStr(data, test.field, test.pretty)
		if err != nil {
			t.Errorf("TestFieldAsStr(%s): got unexpected error: %s", test.desc, err)
			continue
		}
		if got != test.want {
			t.Errorf("TestFieldAsStr(%s): got %q, want %q", test.desc, got, test.want)
		}
	}
}

func TestGetLastMessage(t *testing.T) {
	// Note: Almost all cases are tested in UpdateProtoField, except one.
	// That option isn't used there, so here's that test.
	toUpdate := &pb.Layer0{}

	msg, err := getLastMessage(toUpdate, []string{"layer1", "supported", "vstring"}, true)
	if err != nil {
		t.Fatalf("TestGetLastMessage: got err == %s, want err == nil", err)
	}
	msg.Interface().(*pb.Supported).Vstring = "Hello"

	if toUpdate.Layer1.Supported.Vstring != "Hello" {
		t.Fatalf("TestGetLastMessage: createMessages option isn't working")
	}
}

func TestUpdateProtoField(t *testing.T) {
	toUpdate := &pb.Layer1{}

	tests := []struct {
		desc           string
		supportedExist bool
		createMessages bool
		want           *pb.Layer1
		fieldName      string
		value          interface{}
		err            bool
	}{
		{
			desc: "string at top message",
			want: &pb.Layer1{
				Vstring: "hello",
			},
			fieldName: "vstring",
			value:     "hello",
		},

		{
			desc:           "int32",
			supportedExist: true,
			want: &pb.Layer1{
				Supported: &pb.Supported{Vint32: 32},
			},
			fieldName: "supported.vint32",
			value:     int32(32),
		},
		{
			desc:           "int",
			supportedExist: true,
			want: &pb.Layer1{
				Supported: &pb.Supported{Vint64: 64},
			},
			fieldName: "supported.vint64",
			value:     64,
		},
		{
			desc:           "int64",
			supportedExist: true,
			want: &pb.Layer1{
				Supported: &pb.Supported{Vint64: 64},
			},
			fieldName: "supported.vint64",
			value:     int64(64),
		},
		{
			desc:           "string",
			supportedExist: true,
			want: &pb.Layer1{
				Supported: &pb.Supported{Vstring: "John Doak"},
			},
			fieldName: "supported.vstring",
			value:     "John Doak",
		},
		{
			desc:           "int32 to enum",
			supportedExist: true,
			want: &pb.Layer1{
				Supported: &pb.Supported{Ev: pb.EnumValues_EV_Ok},
			},
			fieldName: "supported.ev",
			value:     int32(1),
		},
		{
			desc:           "enum to enum",
			supportedExist: true,
			want: &pb.Layer1{
				Supported: &pb.Supported{Ev: pb.EnumValues_EV_Ok},
			},
			fieldName: "supported.ev",
			value:     pb.EnumValues_EV_Ok,
		},
		{
			desc:           "bool",
			supportedExist: true,
			want: &pb.Layer1{
				Supported: &pb.Supported{Vbool: true},
			},
			fieldName: "supported.vbool",
			value:     true,
		},
		{
			desc:      "error: supported is nil",
			fieldName: "supported.vbool",
			value:     true,
			err:       true,
		},
	}
	for _, test := range tests {
		start := proto.Clone(toUpdate).(*pb.Layer1)
		if test.supportedExist {
			start.Supported = &pb.Supported{}
		}

		err := UpdateProtoField(start, test.fieldName, test.value)
		switch {
		case err == nil && test.err:
			t.Errorf("TestUpdateProtoField(%s): got err == nil, want err != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestUpdateProtoField(%s): got err == %s, want err == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}
		if diff := cmp.Diff(test.want, start, protocmp.Transform()); diff != "" {
			t.Errorf("TestUpdateProtoField(%s): -want/+got:\n%s", test.desc, diff)
		}
	}
}

func TestGetField(t *testing.T) {
	msg := &pb.Layer0{
		Vint32: 3,
		Layer1: &pb.Layer1{
			Vstring: "hello",
			Supported: &pb.Supported{
				Ev:    pb.EnumValues_EV_Ok,
				Vbool: true,
			},
		},
	}

	tests := []struct {
		desc     string
		fqPath   string
		err      bool
		wantVal  interface{}
		wantKind protoreflect.Kind
		protocmp bool
	}{
		{
			desc:   "Non-existent field in first message",
			fqPath: "what",
			err:    true,
		},
		{
			desc:     "Field in first message",
			fqPath:   "vint32",
			err:      false,
			wantVal:  int32(3),
			wantKind: protoreflect.Int32Kind,
		},
		{
			desc:     "Field in the first message that is a message",
			fqPath:   "layer1",
			err:      false,
			wantVal:  msg.Layer1,
			wantKind: protoreflect.MessageKind,
			protocmp: true,
		},
		{
			desc:   "Non-existent field second message",
			fqPath: "layer1.blah",
			err:    true,
		},
		{
			desc:     "Field in second message",
			fqPath:   "layer1.vstring",
			err:      false,
			wantVal:  "hello",
			wantKind: protoreflect.StringKind,
		},
		{
			desc:     "enum in third message",
			fqPath:   "layer1.supported.ev",
			err:      false,
			wantVal:  pb.EnumValues_EV_Ok,
			wantKind: protoreflect.EnumKind,
		},
	}

	for _, test := range tests {
		fv, err := GetField(msg, test.fqPath)
		switch {
		case err == nil && test.err:
			t.Errorf("TestGetField(%s): got err == nil, want err != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestGetField(%s): got err == %s, want err == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		if fv.Kind != test.wantKind {
			t.Errorf("TestGetField(%s): Kind: got %v, want %v", test.desc, fv.Kind, test.wantKind)
			continue
		}

		if test.protocmp {
			wp := test.wantVal.(proto.Message)
			vp := fv.Value.(proto.Message)
			if diff := Equal(wp, vp); diff != "" {
				t.Errorf("TestGetField(%s): Val: -want/+got:\n%s", test.desc, diff)
			}
			continue
		}

		if diff := pretty.Compare(fv.Value, test.wantVal); diff != "" {
			t.Errorf("TestGetField(%s): Val: -want/+got:\n%s", test.desc, diff)
		}
	}
}

func TestEnumLookup(t *testing.T) {
	msg := &pb.Layer0{}

	evUnknown := Rec{
		EnumName:   "EnumValues",
		Int32:      0,
		ProtoName:  "EV_Unknown",
		JSONName:   "evUnknown",
		TitledName: "Unknown",
	}
	evOk := Rec{
		EnumName:   "EnumValues",
		Int32:      1,
		ProtoName:  "EV_Ok",
		JSONName:   "evOk",
		TitledName: "Ok",
	}
	evNotOk := Rec{
		EnumName:   "EnumValues",
		Int32:      2,
		ProtoName:  "EV_Not_Ok",
		JSONName:   "evNotOk",
		TitledName: "Not Ok",
	}
	evEh := Rec{
		EnumName:   "EnumValues",
		Int32:      3,
		ProtoName:  "EV_Eh",
		JSONName:   "evEh",
		TitledName: "Eh",
	}
	eeUnknown := Rec{
		EnumName:   "EnumEmbedded",
		Int32:      0,
		ProtoName:  "EE_UNKNOWN",
		JSONName:   "eeUnknown",
		TitledName: "Unknown",
	}
	eeWhatever := Rec{
		EnumName:   "EnumEmbedded",
		Int32:      1,
		ProtoName:  "EE_WHATEVER",
		JSONName:   "eeWhatever",
		TitledName: "Whatever",
	}

	wantForward := ForwardLookup{ //map[string]Rec
		evUnknown.ProtoName:   evUnknown,
		evUnknown.JSONName:    evUnknown,
		evOk.ProtoName:        evOk,
		evOk.JSONName:         evOk,
		evOk.TitledName:       evOk,
		evNotOk.ProtoName:     evNotOk,
		evNotOk.JSONName:      evNotOk,
		evNotOk.TitledName:    evNotOk,
		evEh.ProtoName:        evEh,
		evEh.JSONName:         evEh,
		evEh.TitledName:       evEh,
		eeUnknown.ProtoName:   eeUnknown,
		eeUnknown.JSONName:    eeUnknown,
		eeWhatever.ProtoName:  eeWhatever,
		eeWhatever.JSONName:   eeWhatever,
		eeWhatever.TitledName: eeWhatever,
	}

	wantReverse := ReverseLookup{ // map[string]map[int32]Rec
		"EnumValues": map[int32]Rec{
			0: evUnknown,
			1: evOk,
			2: evNotOk,
			3: evEh,
		},
		"EnumEmbedded": map[int32]Rec{
			0: eeUnknown,
			1: eeWhatever,
		},
	}

	forward, reverse := EnumLookup([]proto.Message{msg})

	if diff := pretty.Compare(wantForward, forward); diff != "" {
		t.Errorf("TestEnumLookup(forward): -want/+got:\n%s", diff)
	}

	if diff := pretty.Compare(wantReverse, reverse); diff != "" {
		t.Errorf("TestEnumLookup(reverse): -want/+got:\n%s", diff)
	}
}
