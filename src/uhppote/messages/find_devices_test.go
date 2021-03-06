package messages

import (
	"net"
	"reflect"
	"testing"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
)

func TestMarshalFindDevicesRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := FindDevicesRequest{}

	m, err := codec.Marshal(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Invalid byte array for uhppote.Marshal(%s):\nExpected:\n%s\nReturned:\n%s", "FindDevicesRequest", dump(expected, ""), dump(m, ""))
		return
	}
}

func TestFactoryUnmarshalFindDevicesRequest(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request, err := UnmarshalRequest(message)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if request == nil {
		t.Fatalf("Unexpected request: %v\n", request)
	}

	rq, ok := request.(*FindDevicesRequest)
	if !ok {
		t.Fatalf("Invalid request type - expected:%T, got: %T\n", &FindDevicesRequest{}, request)
	}

	if rq.MsgType != 0x94 {
		t.Errorf("Incorrect 'message type' from valid message: %02x\n", rq.MsgType)
	}
}

func TestUnmarshalFindDevicesResponse(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := FindDevicesResponse{}

	err := codec.Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Incorrect 'message type' from valid message: %02x\n", reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Incorrect 'serial number' from valid message: %v\n", reply.SerialNumber)
	}

	if !reflect.DeepEqual(reply.IpAddress, net.IPv4(192, 168, 0, 0)) {
		t.Errorf("Incorrect 'IP address' from valid message: %v\n", reply.IpAddress)
	}

	if !reflect.DeepEqual(reply.SubnetMask, net.IPv4(255, 255, 255, 0)) {
		t.Errorf("Incorrect 'subnet mask' from valid message: %v\n", reply.SubnetMask)
	}

	if !reflect.DeepEqual(reply.Gateway, net.IPv4(0, 0, 0, 0)) {
		t.Errorf("Incorrect 'gateway' from valid message: %v\n", reply.Gateway)
	}

	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	if !reflect.DeepEqual(reply.MacAddress, types.MacAddress(MAC)) {
		t.Errorf("Incorrect 'MAC address' - expected:'%v', got:'%v'\n", types.MacAddress(MAC), reply.MacAddress)
	}

	if reply.Version != 0x0892 {
		t.Errorf("Incorrect 'version' from valid message: %v\n", reply.Version)
	}

	date, _ := time.ParseInLocation("20060102", "20180816", time.Local)
	if reply.Date != types.Date(date) {
		t.Errorf("Incorrect 'date' from valid message: %v\n", reply.Date)
	}
}

func TestFactoryUnmarshalFindDevicesResponse(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	response, err := UnmarshalResponse(message)

	if err != nil {
		t.Fatalf("Unexpected error: %v\n", err)
	}

	if response == nil {
		t.Fatalf("Unexpected response: %v\n", response)
	}

	reply, ok := response.(*FindDevicesResponse)
	if !ok {
		t.Fatalf("Invalid response type - expected:%T, got: %T\n", &FindDevicesResponse{}, response)
	}

	if reply.MsgType != 0x94 {
		t.Errorf("Incorrect 'message type' from valid message: %02x\n", reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Incorrect 'serial number' from valid message: %v\n", reply.SerialNumber)
	}

	if !reflect.DeepEqual(reply.IpAddress, net.IPv4(192, 168, 0, 0)) {
		t.Errorf("Incorrect 'IP address' from valid message: %v\n", reply.IpAddress)
	}

	if !reflect.DeepEqual(reply.SubnetMask, net.IPv4(255, 255, 255, 0)) {
		t.Errorf("Incorrect 'subnet mask' from valid message: %v\n", reply.SubnetMask)
	}

	if !reflect.DeepEqual(reply.Gateway, net.IPv4(0, 0, 0, 0)) {
		t.Errorf("Incorrect 'gateway' from valid message: %v\n", reply.Gateway)
	}

	MAC, _ := net.ParseMAC("00:66:19:39:55:2d")
	if !reflect.DeepEqual(reply.MacAddress, types.MacAddress(MAC)) {
		t.Errorf("Incorrect 'MAC address' - expected:'%v', got:'%v'\n", types.MacAddress(MAC), reply.MacAddress)
	}

	if reply.Version != 0x0892 {
		t.Errorf("Incorrect 'version' from valid message: %v\n", reply.Version)
	}

	date, _ := time.ParseInLocation("20060102", "20180816", time.Local)
	if reply.Date != types.Date(date) {
		t.Errorf("Incorrect 'date' from valid message: %v\n", reply.Date)
	}
}

func TestUnmarshalFindDevicesResponseWithInvalidMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x92, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0xc0, 0xa8, 0x00, 0x00, 0xff, 0xff, 0xff, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x66, 0x19, 0x39, 0x55, 0x2d, 0x08, 0x92, 0x20, 0x18, 0x08, 0x16,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	reply := FindDevicesResponse{}

	err := codec.Unmarshal(message, &reply)

	if err == nil {
		t.Errorf("Expected error: '%v'", "Invalid value in message - expected 0x94, received 0x92")
		return
	}
}
