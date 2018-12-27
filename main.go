package main

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"log"
	"os"
	"strconv"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}
	log.Println("start")
	for i := int64(0); i < config.MaxUniverse; i++ {
		go func() {
			wg.Add(1)
			world(i)
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println("end")
}

func world(i int64) {
	var u models.Universe
	u.Init(i)

	file, err := os.Open(config.LogPath + strconv.FormatInt(i, 10) + ".csv")
	log.Println(config.LogPath + strconv.FormatInt(i, 10) + ".csv")
	if err != nil {
		log.Fatal(err)
	}

	models.LogInit(u)
	for i := 0; i < config.MaxSteps; i++ {
		u.Step()
		models.LogStep(u, file)
	}
	u.End()
	models.LogEnd(u)

	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
}
