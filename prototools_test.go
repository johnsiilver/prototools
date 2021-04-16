package prototools

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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
		t.Errorf("TestProtoName(%s): got %q, want %q",  jsonName, got, want)
	}
}

func TestFieldAsStr(t *testing.T) {
	data := &pb.Supported{
		Vint32: 32,
		Vint64: 64,
		Vstring: "string",
		Vbool: true,
		Ev: pb.EnumValues_EV_Ok,
	}

	tests := []struct{
		desc string
		field string
		want string
	}{
		{"int32", "vint32", "32"},
		{"int64", "vint64", "64"},
		{"string", "vstring", "string"},
		{"bool", "vbool", "true"},
		{"enum", "ev", "EV_Ok"},
	}

	for _, test := range tests {
		got, err := FieldAsStr(data, test.field)
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
	tests := []struct{
		desc string
		start *pb.Supported
		want *pb.Supported
		fieldName string
		value interface{}
		err bool
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
