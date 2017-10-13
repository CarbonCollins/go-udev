package main

import (
	"bytes"
	"fmt"
	"strings"
)

// See: http://elixir.free-electrons.com/linux/v3.12/source/lib/kobject_uevent.c#L45

const (
	ADD     KObjAction = "add"
	REMOVE  KObjAction = "remove"
	CHANGE  KObjAction = "change"
	MOVE    KObjAction = "move"
	ONLINE  KObjAction = "online"
	OFFLINE KObjAction = "offline"
)

type KObjAction string

func (a KObjAction) String() string {
	return string(a)
}

func ParseKObjAction(raw string) (a KObjAction, err error) {
	a = KObjAction(raw)
	switch a {
	case ADD, REMOVE, CHANGE, MOVE, ONLINE, OFFLINE:
	default:
		err = fmt.Errorf("unknow kobject action (got: %s)", raw)
	}
	return
}

type UEvent struct {
	Action KObjAction
	KObj   string
	Env    map[string]string
}

func ParseUEvent(raw []byte) (e *UEvent, err error) {
	fields := bytes.Split(raw, []byte{0x00}) // 0x00 = end of string

	if len(fields) == 0 {
		err = fmt.Errorf("Wrong uevent format")
		return
	}

	headers := bytes.Split(fields[0], []byte{0x40}) // 0x40 = @
	if len(headers) != 2 {
		err = fmt.Errorf("Wrong uevent header")
		return
	}

	action, err := ParseKObjAction(string(headers[0]))
	if err != nil {
		return
	}

	e = &UEvent{
		Action: action,
		KObj:   string(headers[1]),
		Env:    make(map[string]string, 0),
	}

	for _, envs := range fields[1 : len(fields)-1] {
		// log.Printf("v: %s", envs)
		env := strings.Split(string(envs), "=")
		if len(env) != 2 {
			err = fmt.Errorf("Wrong uevent env")
			return
		}
		e.Env[env[0]] = env[1]
	}
	return
}