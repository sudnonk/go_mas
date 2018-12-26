package main

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
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
	var u models.Universe
	u.Init()
	for i := 0; i < config.MaxSteps; i++ {
		u.Step()
	}
	u.End()
}
