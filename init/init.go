package main

import (
	crand "crypto/rand"
	"flag"
	"github.com/sakura-internet/go-rison"
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"github.com/sudnonk/go_mas/utils"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
)

func main() {
	flag.Parse()
	outFile := flag.Arg(0)
	norm := flag.Arg(1)
	var isNorm bool
	if norm == "true" {
		isNorm = true
	} else if norm == "false" {
		isNorm = false
	} else {
		log.Println("second arg must be true or false.")
		return
	}

	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	ra := rand.New(rand.NewSource(seed.Int64()))

	as := genRandomAgent(ra, isNorm)

	lr := []byte("\n")
	q := []byte("=")
	var l []byte
	for _, a := range as {
		r, _ := rison.Marshal(&a, rison.ORison)
		l = append(l, r...)
		l = append(l, q...)
	}
	l = append(l, lr...)

	writeLog(&l, outFile)
}

func genRandomAgent(ra *rand.Rand, isNorm bool) map[int64]*models.Agent {
	Ags := make(map[int64]*models.Agent, config.MaxAgents())

	for i := int64(0); i < config.MaxAgents(); i++ {
		Ags[i] = models.NewAgent(i, ra, isNorm)
	}

	MakeNetwork(Ags, ra)
	return Ags
}

func MakeNetwork(as map[int64]*models.Agent, ra *rand.Rand) {
	//todo: より良いネットワーク
	for _, a := range as {
		a.Following = utils.RandIntSlice(config.MaxAgents(), int64(config.InitMaxFollowing()), a.Id, ra)
	}
}

func writeLog(data *[]byte, fname string) {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := file.Write(*data); err != nil {
		log.Println(err)
	}

	if err := file.Close(); err != nil {
		log.Println(err)
	}
}
