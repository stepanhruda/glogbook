package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/garyburd/redigo/redis"
	"github.com/unrolled/render"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	r := render.New()

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		rooms := []Room{
			Room{Uuid: uuid.NewRandom(), LastEventTimestamp: time.Now(), LastEventState: "open", Name: "Foo"},
			Room{Uuid: uuid.NewRandom(), LastEventTimestamp: time.Now(), LastEventState: "open", Name: "Bar"},
		}
		r.HTML(w, http.StatusOK, "dashboard", rooms)
	})

	mux.HandleFunc("/door_events", func(w http.ResponseWriter, req *http.Request) {
		roomUuid := req.FormValue("room_uuid")
		timestamp := req.FormValue("timestamp")
		time, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			http.Error(w, "Timestamp in wrong format.", 400)
			return
		}
		doorState := DoorState(req.FormValue("door_state"))
		doorEvent := createDoorEvent(roomUuid, time, doorState)
		saveEvent(doorEvent)
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

	eventKey := fmt.Sprintf("doorEvent:%s", doorEvent.uuid)
	roomKey := fmt.Sprintf("room:%s", doorEvent.roomUuid)
	c.Send("HMSET", eventKey, "roomUuid", doorEvent.roomUuid, "timestamp", doorEvent.timestamp, "state", doorEvent.state)
	c.Send("HMSET", roomKey, "lastEventTimestamp", doorEvent.timestamp, "lastEventState", doorEvent.state)
	c.Flush()
	c.Receive()
	c.Receive()

	return
}
