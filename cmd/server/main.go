package main

import (
	"log"
	"net/http"
	"storeinspection/api"
	"storeinspection/config"
	"storeinspection/service"
	"storeinspection/store"
	"storeinspection/workflow"
)

func main() {
	c := config.Load()
	db, e := store.Open(c.Database)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	svc := service.New(db, service.FixedClock{})
	svc.Notifier = func(string) error { return nil }
	_ = workflow.Recorder{}
	log.Printf("inspection server on %s", c.Address)
	log.Fatal(http.ListenAndServe(c.Address, api.Handler{Service: svc}))
}
