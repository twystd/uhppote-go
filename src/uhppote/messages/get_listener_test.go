package messages

import (
	"net"
	"reflect"
	"testing"
	codec "uhppote/encoding/UTO311-L0x"
)

func TestMarshalGetListenerRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x92, 0x00, 0x00, 0x2D, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := GetListenerRequest{
		SerialNumber: 423187757,
	}

	m, err := codec.Marshal(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Invalid byte array:\nExpected:\n%s\nReturned:\n%s", dump(expected, ""), dump(m, ""))
		return
	}
}

func TestFactoryUnmarshalGetListenerRequest(t *testing.T) {
	message := []byte{
		0x17, 0x92, 0x00, 0x00, 0x2D, 0x55, 0x39, 0x19, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := GetListenerRequest{
		MsgType:      0x92,
		SerialNumber: 423187757,
	}

	request, err := UnmarshalRequest(message)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	rq, ok := request.(*GetListenerRequest)
	if !ok {
		t.Fatalf("Invalid request type - expected:%T, got: %T\n", &GetListenerRequest{}, request)
	}

	if !reflect.DeepEqual(*rq, expected) {
		t.Errorf("Invalid unmarshalled request:\nexpected:%#v\ngot:     %#v", expected, *rq)
		return
	}
}

func TestUnmarshalGetListenerResponse(t *testing.T) {
	message := []byte{
		0x17, 0x92, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0xe1, 0x92, 0x26, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := GetListenerResponse{}

	err := codec.Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	if reply.MsgType != 0x92 {
		t.Errorf("Incorrect 'message type' - expected:%02X, got:%02x\n", 0x92, reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Incorrect 'serial number' - expected:%v, got:%v\n", 423187757, reply.SerialNumber)
	}

	if !reflect.DeepEqual(reply.Address, net.IPv4(192, 168, 0, 225)) {
		t.Errorf("Incorrect IP address - expected:'%v', got:'%v'\n", net.IPv4(192, 168, 0, 225), reply.Address)
	}

	if reply.Port != 9874 {
		t.Errorf("Incorrect 'port' - expected:%d, got:%v\n", 9874, reply.Port)
	}
}

func TestFactoryUnmarshalGetListenerResponse(t *testing.T) {
	message := []byte{
		0x17, 0x92, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0xe1, 0x92, 0x26, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := GetListenerResponse{
		MsgType:      0x92,
		SerialNumber: 423187757,
		Address:      net.IPv4(192, 168, 0, 225),
		Port:         9874,
	}

	response, err := UnmarshalResponse(message)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	reply, ok := response.(*GetListenerResponse)
	if !ok {
		t.Fatalf("Invalid response type - expected:%T, got: %T\n", &GetListenerResponse{}, reply)
	}

	if !reflect.DeepEqual(*reply, expected) {
		t.Errorf("Invalid unmarshalled response:\nexpected:%#v\ngot:     %#v", expected, *reply)
		return
	}
}

func TestUnmarshalGetListenerResponseWithInvalidMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0xe1, 0x92, 0x26, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := GetListenerResponse{}

	err := codec.Unmarshal(message, &reply)

	if err == nil {
		t.Errorf("Expected error: '%v'", "Invalid value in message - expected 0x92, received 0x94")
		return
	}
}
