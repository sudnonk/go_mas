package main

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"log"
)

func main() {
	cs := make(chan struct{})
	for i := 0; i < config.MaxUniverse; i++ {
		go world(cs)
	}
	<-cs
}

func world(c chan struct{}) {
	log.Println("called")
	var u models.Universe
	u.Init()
	for i := 0; i < config.MaxSteps; i++ {
		log.Println("step" + string(i))
		u.Step()
	}
	u.End()
	log.Println("end")
	c <- struct{}{}
}
