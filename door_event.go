package main

import (
	"code.google.com/p/go-uuid/uuid"
	"time"
)

type DoorState string

const (
	Open   DoorState = "open"
	Closed DoorState = "closed"
)

type DoorEvent struct {
	uuid      string
	roomUuid  string
	timestamp time.Time
	state     DoorState
}

func createDoorEvent(roomUuid string, timestamp time.Time, state DoorState) (event DoorEvent) {
	return DoorEvent{uuid: uuid.New(), roomUuid: roomUuid, timestamp: timestamp, state: state}
}
