package messages

import (
	"reflect"
	"testing"
	"time"
	codec "uhppote/encoding/UTO311-L0x"
	"uhppote/types"
)

func TestMarshalGetEventRequest(t *testing.T) {
	expected := []byte{
		0x17, 0xb0, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	request := GetEventRequest{
		SerialNumber: 423187757,
		Index:        1,
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

func TestFactoryUnmarshalGetEventRequest(t *testing.T) {
	message := []byte{
		0x17, 0xb0, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	expected := GetEventRequest{
		MsgType:      0xb0,
		SerialNumber: 423187757,
		Index:        1,
	}

	request, err := UnmarshalRequest(message)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	rq, ok := request.(*GetEventRequest)
	if !ok {
		t.Fatalf("Invalid request type - expected:%T, got: %T\n", &GetEventRequest{}, request)
	}

	if !reflect.DeepEqual(*rq, expected) {
		t.Errorf("Invalid unmarshalled request:\nexpected:%#v\ngot:     %#v", expected, *rq)
		return
	}
}

func TestUnmarshalGetEventResponse(t *testing.T) {
	message := []byte{
		0x17, 0xb0, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x00, 0x00, 0x00, 0x02, 0x01, 0x03, 0x01,
		0xad, 0xe8, 0x5d, 0x00, 0x20, 0x19, 0x02, 0x10, 0x07, 0x12, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4a, 0x26, 0x80, 0x39, 0x08, 0x92, 0x00, 0x00,
	}

	reply := GetEventResponse{}

	err := codec.Unmarshal(message, &reply)

	if err != nil {
		t.Errorf("Unexpected error: %v\n", err)
	}

	if reply.MsgType != 0xb0 {
		t.Errorf("Incorrect 'message type' - expected:%02X, got:%02x\n", 0xb0, reply.MsgType)
	}

	if reply.SerialNumber != 423187757 {
		t.Errorf("Incorrect 'serial number' - expected:%d, got: %v\n", 423187757, reply.SerialNumber)
	}

	if reply.Index != 8 {
		t.Errorf("Incorrect 'index' - expected:%d, got:%d\n", 8, reply.Index)
	}

	if reply.Type != 2 {
		t.Errorf("Incorrect 'type' - expected:%d, got:%d\n", 2, reply.Type)
	}

	if reply.Granted != true {
		t.Errorf("Incorrect 'granted' - expected:%v, got:%v\n", true, reply.Granted)
	}

	if reply.Door != 3 {
		t.Errorf("Incorrect 'door' - expected:%d, got:%d\n", 3, reply.Door)
	}

	if reply.DoorOpened != true {
		t.Errorf("Incorrect 'door opened' - expected:%v, got:%v\n", true, reply.DoorOpened)
	}

	if reply.UserID != 6154413 {
		t.Errorf("Incorrect 'user ID' - expected:%d, got: %v\n", 6154413, reply.UserID)
	}

	if reply.Result != 6 {
		t.Errorf("Incorrect 'result' - expected:%d, got: %v\n", 6, reply.Result)
	}

	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-02-10 07:12:01", time.Local)
	if reply.Timestamp != types.DateTime(timestamp) {
		t.Errorf("Incorrect 'timestamp' - expected:%s, got:%s\n", timestamp.Format("2006-01-02 15:04:05"), reply.Timestamp)
	}
}

func TestFactoryUnmarshalGetEventResponse(t *testing.T) {
	message := []byte{
		0x17, 0xb0, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x00, 0x00, 0x00, 0x02, 0x01, 0x03, 0x01,
		0xad, 0xe8, 0x5d, 0x00, 0x20, 0x19, 0x02, 0x10, 0x07, 0x12, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4a, 0x26, 0x80, 0x39, 0x08, 0x92, 0x00, 0x00,
	}

	timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", "2019-02-10 07:12:01", time.Local)
	expected := GetEventResponse{
		MsgType:      0xb0,
		SerialNumber: 423187757,
		Index:        8,
		Type:         2,
		Granted:      true,
		Door:         3,
		DoorOpened:   true,
		UserID:       6154413,
		Timestamp:    types.DateTime(timestamp),
		Result:       0x06,
	}

	response, err := UnmarshalResponse(message)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	reply, ok := response.(*GetEventResponse)
	if !ok {
		t.Fatalf("Invalid response type - expected:%T, got: %T\n", &GetEventResponse{}, reply)
	}

	if !reflect.DeepEqual(*reply, expected) {
		t.Errorf("Invalid unmarshalled response:\nexpected:%#v\ngot:     %#v", expected, *reply)
		return
	}
}

func TestUnmarshalGetEventResponseWithInvalidMsgType(t *testing.T) {
	message := []byte{
		0x17, 0x94, 0x00, 0x00, 0x2d, 0x55, 0x39, 0x19, 0x08, 0x00, 0x00, 0x00, 0x02, 0x01, 0x03, 0x01,
		0xad, 0xe8, 0x5d, 0x00, 0x20, 0x19, 0x02, 0x10, 0x07, 0x12, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x4a, 0x26, 0x80, 0x39, 0x08, 0x92, 0x00, 0x00,
	}

	reply := GetEventResponse{}

	err := codec.Unmarshal(message, &reply)

	if err == nil {
		t.Errorf("Expected error: '%v'", "Invalid value in message - expected 0xb0, received 0x94")
		return
	}
}
