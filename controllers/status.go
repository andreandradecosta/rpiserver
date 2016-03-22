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
	router.
		Methods("GET").
		Path("/uptime").
		Handler(s.handleAction(s.uptime))
	router.
		Methods("GET").
		Path("/updated").
		Handler(s.handleAction(s.updated))
}

func (s *status) index(w http.ResponseWriter, r *http.Request) error {
	conn := s.redisPool.Get()
	defer conn.Close()
	b, err := redis.Bytes(conn.Do("GET", "status"))
	if err != nil {
		return err
	}
	var d models.Info
	err = json.Unmarshal(b, &d)
	if err != nil {
		return err
	}

	delete(d, "disk_part")
	delete(d, "host")
	s.JSON(w, http.StatusOK, d)
	return nil
}

func (s *status) uptime(w http.ResponseWriter, r *http.Request) error {
	conn := s.redisPool.Get()
	defer conn.Close()
	b, err := redis.Bytes(conn.Do("GET", "status"))
	if err != nil {
		return err
	}

	var d models.Info
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil
	}
	host := d["host"].(map[string]interface{})
	days := host["uptime"].(float64) / (60 * 60 * 24)
	res := map[string]float64{
		"uptime": days,
	}
	s.JSON(w, http.StatusOK, res)
	return nil
}

func (s *status) updated(w http.ResponseWriter, r *http.Request) error {
	conn := s.redisPool.Get()
	defer conn.Close()
	up, err := redis.Int64(conn.Do("GET", "updated"))
	if err != nil {
		return err
	}
	date := time.Unix(up, 0).In(s.hostLocation)
	res := map[string]time.Time{
		"updated": date,
	}
	s.JSON(w, http.StatusOK, res)
	return nil
}
