package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/garyburd/redigo/redis"
	"github.com/unrolled/render"
	"net/http"
	"strings"
	"time"
)

func main() {
	mux := http.NewServeMux()
	r := render.New()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		rooms, _ := loadRooms()
		r.HTML(w, http.StatusOK, "dashboard", rooms)
	})

	mux.HandleFunc("/rooms", func(w http.ResponseWriter, req *http.Request) {
		rooms, _ := loadRooms()
		r.JSON(w, http.StatusOK, map[string][]Room{"rooms": rooms})
	})

	mux.HandleFunc("/door_events", func(w http.ResponseWriter, req *http.Request) {
		// API ENDPOINT: '/door_events'
		// PARAMS:
		// room_slug: string
		// timestamp: rfc3339 timestamp i.e. "2006-01-02T15:04:05Z"
		// door_state: one of the values from the DoorState const

		roomSlug := req.FormValue("room_slug")
		timestamp := req.FormValue("timestamp")
		time, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			http.Error(w, "Timestamp in wrong format.", 400)
			return
		}
		doorState := DoorState(req.FormValue("door_state"))
		doorEvent := createDoorEvent(roomSlug, time, doorState)
		saveEvent(doorEvent)
		r.JSON(w, http.StatusCreated, doorEvent)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func saveEvent(doorEvent DoorEvent) (err error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return err
	}
	defer c.Close()

	eventKey := fmt.Sprintf("doorEvent:%s", doorEvent.Uuid)
	roomKey := fmt.Sprintf("room:%s", doorEvent.RoomSlug)
	c.Send("HMSET", eventKey, "roomSlug", doorEvent.RoomSlug, "timestamp", doorEvent.Timestamp, "doorState", doorEvent.DoorState)
	c.Send("HMSET", roomKey, "lastEventTimestamp", doorEvent.Timestamp, "lastEventState", doorEvent.DoorState)
	c.Flush()
	c.Receive()
	c.Receive()

	return
}

func loadRooms() (rooms []Room, err error) {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		return nil, err
	}
	defer c.Close()

	roomKeysQuery := "room:*"
	roomKeys, err := redis.Strings(c.Do("KEYS", roomKeysQuery))
	for _, key := range roomKeys {
		parts := strings.Split(key, "room:")

		roomSlug := parts[len(parts)-1]
		roomValues, _ := redis.Strings(c.Do("HGETALL", key))
		time, _ := time.Parse(time.RFC3339, roomValues[1])
		room := Room{Slug: roomSlug, LastEventTimestamp: time, Name: roomSlug, LastEventState: DoorState(roomValues[3])}
		rooms = append(rooms, room)
	}
	return
}
