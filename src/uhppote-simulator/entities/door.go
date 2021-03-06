package entities

import (
	"encoding/json"
	"time"
)

type Delay time.Duration

type Door struct {
	ControlState uint8 `json:"control"`
	Delay        Delay `json:"delay"`
	open         bool
	openUntil    *time.Time
	button       bool
}

func (delay Delay) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(delay).String())
}

func (delay *Delay) UnmarshalJSON(bytes []byte) error {
	var s string

	err := json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		return err
	}

	*delay = Delay(d)

	return nil
}

func (delay Delay) Seconds() uint8 {
	return uint8(time.Duration(delay).Seconds())
}

func DelayFromSeconds(t uint8) Delay {
	return Delay(time.Duration(int64(t) * 1000000000))
}

func NewDoor(id uint8) *Door {
	door := new(Door)

	door.ControlState = 3
	door.Delay = Delay(5 * 1000000000)
	door.open = false
	door.openUntil = nil
	door.button = false

	return door
}

func (d *Door) Open() bool {
	now := time.Now().UTC()
	closeAt := now.Add(time.Duration(d.Delay))

	d.openUntil = &closeAt

	return true
}

func (d Door) IsOpen() bool {
	if d.openUntil != nil {
		return !time.Now().UTC().After(*d.openUntil)
	}

	return false
}

func (d Door) IsButtonPressed() bool {
	return false
}
