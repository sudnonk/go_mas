package main

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"flag"
	"github.com/sakura-internet/go-rison"
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"github.com/sudnonk/go_mas/utils"
	"io"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	flag.Parse()
	outDir := flag.Arg(0)

	initPath := flag.Arg(1)

	wg := new(sync.WaitGroup)
	m := new(sync.Mutex)
	log.Println("start")

	for i := int64(0); i < config.MaxUniverse; i++ {
		wg.Add(1)
		go func(i int64) {
			time.Sleep(time.Duration(i*10) * time.Millisecond)
			world(i, m, outDir, initPath)
			wg.Done()
		}(i)
	}
	wg.Wait()
	log.Println("end")
}

func world(id int64, m *sync.Mutex, outDir string, initPath string) {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	ra := rand.New(rand.NewSource(seed.Int64()))

	var initData map[int64]*models.Agent
	if initPath != "" {
		initData, _ = parseAgent(initPath)
	} else {
		initData = genRandomAgent(ra)
	}

	var u models.Universe
	u.Init(id, initData)

	ids := strconv.FormatInt(id, 10)
	prefix := outDir + ids + "_step"
	var fname string
	for i := 0; i < config.MaxSteps; i++ {
		if i%100 == 0 {
			fname = prefix + strconv.Itoa(i) + ".csv"
		}
		log.Println(ids + " step:" + strconv.Itoa(i))
		models.LogStep(&u, fname, m)
		u.Step(ra)
	}
	models.LogStep(&u, fname, m)
}

func parseAgent(initPath string) (map[int64]*models.Agent, error) {
	f, r, e := openFile(initPath)
	if e != nil {
		return nil, e
	}
	line, err := r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}
	as, err := parseLineAll(&line)
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	return as, nil
}

func openFile(fname string) (f *os.File, r *bufio.Reader, err error) {
	f, err = os.Open(fname)
	if err != nil {
		return nil, nil, err
	}

	return f, bufio.NewReader(f), nil
}

func parseLineAll(line *[]byte) (map[int64]*models.Agent, error) {
	delimiter := []byte("=")
	as := bytes.Split(*line, delimiter)

	agents := make(map[int64]*models.Agent, config.MaxAgents)
	for _, a := range as {
		if len(a) < 5 {
			continue
		}

		var ag models.Agent
		err := rison.Unmarshal(a, &ag, rison.ORison)
		if err != nil {
			log.Println(err)
			continue
		}
		agents[ag.Id] = &ag
	}

	return agents, nil
}

func genRandomAgent(ra *rand.Rand) map[int64]*models.Agent {
	Ags := make(map[int64]*models.Agent, config.MaxAgents)

	for i := int64(0); i < config.MaxAgents; i++ {
		Ags[i] = models.NewAgent(i, ra, config.IsNorm)
	}

	MakeNetwork(Ags, ra)
	return Ags
}

func MakeNetwork(as map[int64]*models.Agent, ra *rand.Rand) {
	//todo: より良いネットワーク
	for _, a := range as {
		a.Following = utils.RandIntSlice(config.MaxAgents, int64(config.InitMaxFollowing), a.Id, ra)
	}
}
