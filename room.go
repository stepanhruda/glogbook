package main

import (
	"time"
)

type Room struct {
	LastEventTimestamp time.Time `json:"last_event_timestamp"`
	LastEventState     DoorState `json:"last_event_state"`
	Name               string    `json:"name"`
	Slug               string    `json:"id"`
}
