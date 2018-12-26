package main

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"log"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	for i := 0; i < config.MaxUniverse; i++ {
		go func() {
			wg.Add(1)
			world()
			wg.Done()
		}()
	}
	wg.Wait()
}

func world() {
	log.Println("called")
	var u models.Universe
	u.Init()
	for i := 0; i < config.MaxSteps; i++ {
		log.Println("step" + string(i))
		u.Step()
	}
	u.End()
	log.Println("end")
}
