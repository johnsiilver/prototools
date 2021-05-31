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
	tests := []struct {
		desc    string
		name    string
		want    string
		options []ReadableOption
	}{
		{
			desc: "Base test",
			name: "this_is_my_field_name_32",
			want: "This Is My Field Name 32",
		},
		{
			desc:    "RemovePrefix test",
			name:    "rcategory_Unknown",
			want:    "Unknown",
			options: []ReadableOption{RemovePrefix()},
		},
	}

	for _, test := range tests {
		got := ReadableProto(test.name, test.options...)
		if got != test.want {
			t.Errorf("TestReadableProto(%s): got %q, want %q", test.name, got, test.want)
		}
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

func TestFieldValue(t *testing.T) {
	myPretty := pretty.Config{
		IncludeUnexported: false,
		TrackCycles:       true,
	}

	const (
		unknown = 0
		stdList = 1
		stdMsg  = 2
		stdVal  = 3
		listMsg = 4
	)
	tests := []struct {
		desc        string
		msg         proto.Message
		field       string
		err         bool
		want        FieldValue
		compareType int
	}{
		{
			desc:  "Bool",
			msg:   &pb.Supported{Vbool: true},
			field: "vbool",
			want: FieldValue{
				Value:     true,
				Kind:      protoreflect.BoolKind,
				FieldDesc: (&pb.Supported{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("vbool")),
			},
			compareType: stdVal,
		},
		{
			desc:  "Int32",
			msg:   &pb.Supported{Vint32: 1},
			field: "vint32",
			want: FieldValue{
				Value:     1,
				Kind:      protoreflect.Int32Kind,
				FieldDesc: (&pb.Supported{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("vint32")),
			},
			compareType: stdVal,
		},
		{
			desc:  "Int64",
			msg:   &pb.Supported{Vint64: 1},
			field: "vint64",
			want: FieldValue{
				Value:     1,
				Kind:      protoreflect.Int64Kind,
				FieldDesc: (&pb.Supported{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("vint64")),
			},
			compareType: stdVal,
		},
		{
			desc:  "String",
			msg:   &pb.Supported{Vstring: "hello"},
			field: "vstring",
			want: FieldValue{
				Value:     "hello",
				Kind:      protoreflect.StringKind,
				FieldDesc: (&pb.Supported{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("vstring")),
			},
			compareType: stdVal,
		},
		{
			desc:  "Float32",
			msg:   &pb.Supported{Vfloat: 1.1},
			field: "vfloat",
			want: FieldValue{
				Value:     float32(1.1),
				Kind:      protoreflect.FloatKind,
				FieldDesc: (&pb.Supported{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("vfloat")),
			},
			compareType: stdVal,
		},
		{
			desc:  "Float64",
			msg:   &pb.Supported{Vdouble: 1.1},
			field: "vdouble",
			want: FieldValue{
				Value:     1.1,
				Kind:      protoreflect.DoubleKind,
				FieldDesc: (&pb.Supported{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("vdouble")),
			},
			compareType: stdVal,
		},
		{
			desc:  "[]bool",
			msg:   &pb.BunchOTypes{LBool: []bool{true}},
			field: "l_bool",
			want: FieldValue{
				Value:     []bool{true},
				Kind:      protoreflect.BoolKind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_bool")),
			},
			compareType: stdList,
		},
		{
			desc:  "[]int32",
			msg:   &pb.BunchOTypes{LInt32: []int32{1}},
			field: "l_int32",
			want: FieldValue{
				Value:     []int32{1},
				Kind:      protoreflect.Int32Kind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_int32")),
			},
			compareType: stdList,
		},
		{
			desc:  "[]int64",
			msg:   &pb.BunchOTypes{LInt64: []int64{1}},
			field: "l_int64",
			want: FieldValue{
				Value:     []int64{1},
				Kind:      protoreflect.Int64Kind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_int64")),
			},
			compareType: stdList,
		},
		{
			desc:  "[]float32",
			msg:   &pb.BunchOTypes{LFloat: []float32{1.1}},
			field: "l_float",
			want: FieldValue{
				Value:     []float32{1.1},
				Kind:      protoreflect.FloatKind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_float")),
			},
			compareType: stdList,
		},
		{
			desc:  "[]float64",
			msg:   &pb.BunchOTypes{LDouble: []float64{1.1}},
			field: "l_double",
			want: FieldValue{
				Value:     []float64{1.1},
				Kind:      protoreflect.DoubleKind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_double")),
			},
			compareType: stdList,
		},
		{
			desc:  "[]string",
			msg:   &pb.BunchOTypes{LString: []string{"hello"}},
			field: "l_string",
			want: FieldValue{
				Value:     []string{"hello"},
				Kind:      protoreflect.StringKind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_string")),
			},
			compareType: stdList,
		},
		{
			desc:  "[]enum",
			msg:   &pb.BunchOTypes{LEv: []pb.EnumValues{pb.EnumValues_EV_Ok}},
			field: "l_ev",
			want: FieldValue{
				Value:     []protoreflect.EnumNumber{protoreflect.EnumNumber(1)},
				Kind:      protoreflect.EnumKind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_ev")),
			},
			compareType: stdList,
		},
		{
			desc:  "[]Message",
			msg:   &pb.BunchOTypes{LMessage: []*pb.Supported{&pb.Supported{Vstring: "hello"}}},
			field: "l_message",
			want: FieldValue{
				Value: []protoreflect.Message{
					(&pb.Supported{Vstring: "hello"}).ProtoReflect(),
				},
				Kind:      protoreflect.MessageKind,
				IsList:    true,
				FieldDesc: (&pb.BunchOTypes{}).ProtoReflect().Descriptor().Fields().ByName(protoreflect.Name("l_message")),
			},
			compareType: listMsg,
		},
	}

	for _, test := range tests {
		got, err := fieldValue(test.msg, test.field)
		switch {
		case err == nil && test.err:
			t.Errorf("TestFieldValue(%s): got err == nil, want err != nil", test.desc)
			continue
		case err != nil && !test.err:
			t.Errorf("TestFieldValue(%s): got err == %s, want err == nil", test.desc, err)
			continue
		case err != nil:
			continue
		}

		switch test.compareType {
		case stdList, stdVal:
			if diff := myPretty.Compare(test.want, got); diff != "" {
				t.Errorf("TestFieldValue(%s): -want/+got:\n%s", test.desc, diff)
			}
			continue
		case listMsg:
			lWant := test.want.Value.([]protoreflect.Message)
			lGot := got.Value.([]protoreflect.Message)
			if len(lWant) != len(lGot) {
				t.Errorf("TestFieldValue(%s): got %d messages, want %d messages", test.desc, len(lGot), len(lWant))
				continue
			}
			for i := 0; i < len(lWant); i++ {
				wp := lWant[i].(protoreflect.Message)
				vp := lGot[i].(protoreflect.Message)
				if diff := Equal(wp.Interface(), vp.Interface()); diff != "" {
					t.Errorf("TestFieldValue(%s) item %d: Val: -want/+got:\n%s", test.desc, i, diff)
				}
			}
		default:
			t.Fatalf("TestFieldValue: broken test parameter: compareType not set")
		}

	}
}
