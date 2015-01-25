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
	uuid      uuid.UUID
	roomSlug  string
	Timestamp time.Time
	state     DoorState
}

func createDoorEvent(roomSlug string, Timestamp time.Time, state DoorState) (event DoorEvent) {
	return DoorEvent{uuid: uuid.NewRandom(), roomSlug: roomSlug, Timestamp: Timestamp, state: state}
}
