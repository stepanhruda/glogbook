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
	Uuid      uuid.UUID `json:"id"`
	RoomSlug  string    `json:"room_slug"`
	Timestamp time.Time `json:"timestamp"`
	DoorState DoorState `json:"door_state"`
}

func createDoorEvent(RoomSlug string, Timestamp time.Time, DoorState DoorState) (event DoorEvent) {
	return DoorEvent{Uuid: uuid.NewRandom(), RoomSlug: RoomSlug, Timestamp: Timestamp, DoorState: DoorState}
}
