package uhppote

import (
	"net"
	"reflect"
	"testing"
	"time"
	"uhppote/encoding"
)

func TestMarshalFindDevicesRequest(t *testing.T) {
	expected := []byte{
		0x17, 0x94, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := FindDevicesRequest{
		MsgType: 0x94,
	}

	m, err := uhppote.Marshal(request)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	if !reflect.DeepEqual(m, expected) {
		t.Errorf("Invalid byte array for uhppote.Marshal(%s):\nExpected:\n%s\nReturned:\n%s", "FindDevicesRequest", print(expected), print(m))
		return
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

	err := uhppote.Unmarshal(message, &reply)

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
	if !reflect.DeepEqual(reply.MacAddress, MAC) {
		t.Errorf("Incorrect 'MAC address' from valid message: %v\n", reply.MacAddress)
	}

	if reply.Version != 0x0892 {
		t.Errorf("Incorrect 'version' from valid message: %v\n", reply.Version)
	}

	date, _ := time.ParseInLocation("20060102", "20180816", time.Local)
	if reply.Date.Date != date {
		t.Errorf("Incorrect 'date' from valid message: %v\n", reply.Date)
	}
}