package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/cheekybits/jobs"
	"gopkg.in/mgo.v2"
)

func main() {

	log.Println("Connecting to localhost MongoDB...")
	dbsession, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatalln(err)
	}
	defer dbsession.Close()

	r := jobs.NewRunner("runner-1", dbsession.DB("jobs-example").C("jobs"), "notifications", func(j *jobs.J) error {
		log.Println("Notification:", j.Data["message"])
		return nil
	})

	if err := r.Start(); err != nil {
		log.Fatalln(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c // wait for ctrl+C
	log.Println("Stopping...")
	r.Stop()
	<-r.StopChan() // wait for runner to finish

	log.Println("Stopped.")

}
