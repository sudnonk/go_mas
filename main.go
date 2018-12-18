package main

import (
	"github.com/sudnonk/go_mas/models"
	"log"
	"time"
)

func main() {
	ss := int64(0)
	log.Println("start")
	for i := 0; i < 10000; i++ {
		var u models.Universe
		s := time.Now()
		u.Init()
		ss += time.Since(s).Nanoseconds()
	}

	log.Println("end")
	log.Println(ss / 10000)
}

func world() {
	var u models.Universe
	u.Init()
}
