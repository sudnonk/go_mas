package main

import (
	crand "crypto/rand"
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"log"
	"math"
	"math/big"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func main() {
	wg := new(sync.WaitGroup)
	m := new(sync.Mutex)
	log.Println("start")
	for i := int64(0); i < config.MaxUniverse; i++ {
		wg.Add(1)
		go func(i int64) {
			time.Sleep(time.Duration(i*10) * time.Millisecond)
			world(i, m)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("end")
}

func world(i int64, m *sync.Mutex) {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	ra := rand.New(rand.NewSource(seed.Int64()))

	var u models.Universe
	u.Init(i, ra)

	fname := config.LogPath + strconv.FormatInt(i, 10) + ".csv"

	/*cu := make(chan *models.Universe, 100)
	cf := make(chan *os.File, 100)
	defer close(cu)
	defer close(cf)

	go models.LogStepChan(cu, cf)*/

	for i := 0; i < config.MaxSteps; i++ {
		if i%100 == 0 {
			models.LogStep(&u, fname, m)
		}

		u.Step(ra)
	}
	models.LogStep(&u, fname, m)
}
