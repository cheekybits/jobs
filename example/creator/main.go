package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
	fmt.Println("Create notification jobs by typing lines:")
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		job := jobs.New("notifications")
		job.Data["message"] = s.Text()
		if err := jobs.Put(dbsession.DB("jobs-example").C("jobs"), job); err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Finished.")

}
