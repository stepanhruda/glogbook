package main

import (
	"time"
)

type Room struct {
	LastEventTimestamp time.Time
	LastEventState     DoorState
	Name               string
	Slug               string
}
