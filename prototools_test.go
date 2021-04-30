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
		{"enum", "supported.ev", "EV_Ok", false},
		{"enum", "supported.ev", "Ok", true},
	}

	for _, test := range tests {
		got, err := FieldAsStr(data, test.field, test.pretty)
		if err != nil {
			t.Errorf("TestFieldAsStr(%s): got unexpected error: %s", test.desc, err)
			continue
		}
		if got != test.want {
			t.Errorf("TestFieldAsStr(%s): got %q, want %q", test.desc, got, test.want)
		}
	}
}

func TestUpdateProtoField(t *testing.T) {
	tests := []struct {
		desc      string
		start     *pb.Supported
		want      *pb.Supported
		fieldName string
		value     interface{}
		err       bool
	}{
		{
			desc:      "int32",
			start:     &pb.Supported{},
			want:      &pb.Supported{Vint32: 32},
			fieldName: "vint32",
			value:     int32(32),
		},
		{
			desc:      "int",
			start:     &pb.Supported{},
			want:      &pb.Supported{Vint64: 64},
			fieldName: "vint64",
			value:     64,
		},
		{
			desc:      "int64",
			start:     &pb.Supported{},
			want:      &pb.Supported{Vint64: 64},
			fieldName: "vint64",
			value:     int64(64),
		},
		{
			desc:      "string",
			start:     &pb.Supported{},
			want:      &pb.Supported{Vstring: "John Doak"},
			fieldName: "vstring",
			value:     "John Doak",
		},
		{
			desc:      "int32 to enum",
			start:     &pb.Supported{},
			want:      &pb.Supported{Ev: pb.EnumValues_EV_Ok},
			fieldName: "ev",
			value:     int32(1),
		},
		{
			desc:      "enum to enum",
			start:     &pb.Supported{},
			want:      &pb.Supported{Ev: pb.EnumValues_EV_Ok},
			fieldName: "ev",
			value:     pb.EnumValues_EV_Ok,
		},
		{
			desc:      "bool",
			start:     &pb.Supported{},
			want:      &pb.Supported{Vbool: true},
			fieldName: "vbool",
			value:     true,
		},
	}
	for _, test := range tests {
		err := UpdateProtoField(test.start, test.fieldName, test.value)
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
		if diff := cmp.Diff(test.want, test.start, protocmp.Transform()); diff != "" {
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
