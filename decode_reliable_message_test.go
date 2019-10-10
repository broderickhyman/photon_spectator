package photon_spectator

import (
	"reflect"
	"strings"
	"testing"
)

var responses = []struct {
	input  []byte
	output ReliableMessageParamaters
}{
	{
		[]byte{0x00, Int8Type, 0xff},
		ReliableMessageParamaters{0: int8(-1)},
	},
	{
		[]byte{0x00, Float32Type, 0x43, 0x00, 0x20, 0xc5},
		ReliableMessageParamaters{0: float32(128.128)},
	},
	{
		[]byte{0x00, Int32Type, 0x00, 0x00, 0x00, 0x80},
		ReliableMessageParamaters{0: int32(128)},
	},
	{
		[]byte{0x00, Int16Type, 0x00, 0x80},
		ReliableMessageParamaters{0: int16(128)},
	},
	{
		[]byte{0x00, Int64Type, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80},
		ReliableMessageParamaters{0: int64(128)},
	},
	{
		[]byte{0x00, StringType, 0x00, 0x03, 0x61, 0x62, 0x63},
		ReliableMessageParamaters{0: "abc"},
	},
	{
		[]byte{0x00, BooleanType, 0x00},
		ReliableMessageParamaters{0: false},
	},
	{
		[]byte{0x00, BooleanType, 0x01},
		ReliableMessageParamaters{0: true},
	},
	{
		[]byte{0x00, Int8SliceType, 0x00, 0x00, 0x00, 0x01, 0x01},
		ReliableMessageParamaters{0: []int8{1}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, Float32Type, 0x43, 0x00, 0x20, 0xc5},
		ReliableMessageParamaters{0: []float32{128.128}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, Int32Type, 0x00, 0x00, 0x00, 0x80},
		ReliableMessageParamaters{0: []int32{128}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, Int16Type, 0x00, 0x80},
		ReliableMessageParamaters{0: []int16{128}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, Int64Type, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80},
		ReliableMessageParamaters{0: []int64{128}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, StringType, 0x00, 0x03, 0x61, 0x62, 0x63},
		ReliableMessageParamaters{0: []string{"abc"}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, BooleanType, 0x01},
		ReliableMessageParamaters{0: []bool{true}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, BooleanType, 0x00},
		ReliableMessageParamaters{0: []bool{false}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, Int8SliceType, 0x00, 0x00, 0x00, 0x01, 0x01},
		ReliableMessageParamaters{0: [][]int8{[]int8{1}}},
	},
	{
		[]byte{0x00, SliceType, 0x00, 0x01, SliceType, 0x00, 0x01, BooleanType, 0x00},
		ReliableMessageParamaters{0: []interface{}{[]bool{false}}},
	},
}

func TestDecodeReliableMessage(t *testing.T) {
	for _, r := range responses {
		var msg ReliableMessage
		msg.ParamaterCount = 1
		msg.Data = r.input

		actual := DecodeReliableMessage(msg)

		if !reflect.DeepEqual(r.output, actual) {
			t.Errorf("Expected `%#v` but got `%#v`", r.output, actual)
		}
	}
}

func TestDecodeReliableMessage_DefaultError(t *testing.T) {
	var msg ReliableMessage
	msg.ParamaterCount = 1
	msg.Data = []byte{64, 64, 64}

	params := DecodeReliableMessage(msg)

	if !strings.HasPrefix(params[64].(string), "ERROR") {
		t.Fail()
	}
}

func TestDecodeReliableMessage_BooleanError(t *testing.T) {
	var msg ReliableMessage
	msg.ParamaterCount = 1
	msg.Data = []byte{64, BooleanType, 64}

	params := DecodeReliableMessage(msg)

	if !strings.HasPrefix(params[64].(string), "ERROR") {
		t.Fail()
	}
}

func TestDecodeReliableMessage_SliceError(t *testing.T) {
	var msg ReliableMessage
	msg.ParamaterCount = 1
	msg.Data = []byte{0x00, SliceType, 0x00, 0x01, BooleanType, 0xff}

	params := DecodeReliableMessage(msg)

	if !strings.HasPrefix(params[0].(string), "ERROR") {
		t.Fail()
	}
}

func TestDecodeReliableMessage_SliceDefaultError(t *testing.T) {
	var msg ReliableMessage
	msg.ParamaterCount = 1
	msg.Data = []byte{0x00, SliceType, 0x00, 0x01, 64, 0xff}

	params := DecodeReliableMessage(msg)

	if !strings.HasPrefix(params[0].(string), "ERROR") {
		t.Fail()
	}
}

func TestDecodeReliableMessage_SliceNestedError(t *testing.T) {
	var msg ReliableMessage
	msg.ParamaterCount = 1
	msg.Data = []byte{0x00, SliceType, 0x00, 0x01, SliceType, 0x00, 0x01, 64, 0x00}

	params := DecodeReliableMessage(msg)

	if !strings.HasPrefix(params[0].(string), "ERROR") {
		t.Fail()
	}
}
