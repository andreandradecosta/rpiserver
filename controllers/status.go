package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/andreandradecosta/rpimonitor/models"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"gopkg.in/unrolled/render.v1"
)

type status struct {
	controller
}

//NewStatus creates this controller and register it with router.
func NewStatus(renderer *render.Render, router *mux.Router, redisPool *redis.Pool, hostLocation *time.Location) {
	s := &status{controller{Render: renderer, redisPool: redisPool, hostLocation: hostLocation}}
	router.
		Methods("GET").
		Path("/").
		Name("Index").
		Handler(s.handleAction(s.index))
}

func (s *status) index(w http.ResponseWriter, r *http.Request) error {
	conn := s.redisPool.Get()
	defer conn.Close()
	reply, err := redis.Values(conn.Do("MGET", "status", "updated"))
	if err != nil {
		return err
	}
	var status []byte
	var updated int64
	if _, err = redis.Scan(reply, &status, &updated); err != nil {
		return err
	}

	var info models.Info
	err = json.Unmarshal(status, &info)
	if err != nil {
		return err
	}
	response := make(models.Info)

	host := info["host"].(map[string]interface{})
	days := host["uptime"].(float64) / (60 * 60 * 24)
	response["uptime"] = days

	date := time.Unix(updated, 0).In(s.hostLocation)
	response["updated"] = date

	response["users"] = info["users"]

	s.JSON(w, http.StatusOK, response)
	return nil
}
