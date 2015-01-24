package main

import (
	"code.google.com/p/go-uuid/uuid"
	"time"
)

type Room struct {
	Uuid               uuid.UUID
	LastEventTimestamp time.Time
	LastEventState     DoorState
	Name               string
}
