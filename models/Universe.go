package models

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/utils"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type Universe struct {
	Id      int64
	Agents  map[int64]*Agent
	StepNum int64
}

func (u *Universe) Init(id int64, ra *rand.Rand) {
	u.Id = id
	u.Agents = make(map[int64]*Agent, config.MaxAgents)
	u.StepNum = 0

	for i := int64(0); i < config.MaxAgents; i++ {
		u.Agents[i] = NewAgent(i, ra)
	}

	u.MakeNetwork(ra)
}

func (u *Universe) MakeNetwork(ra *rand.Rand) {
	//todo: より良いネットワーク
	for _, a := range u.Agents {
		a.Following = utils.RandIntSlice(config.MaxAgents, int64(config.InitMaxFollowing), a.Id, ra)
	}
}

func (u *Universe) Step(ra *rand.Rand) {
	u.StepNum++
	if u.Id == 0 {
		s := time.Now()

		for _, a := range u.Agents {
			a.Step(u.Agents, ra)
		}

		log.Println(strconv.FormatInt(u.Id, 10) + ": " + "step: " + strconv.FormatInt(u.StepNum, 10) + " " + strconv.FormatInt(time.Since(s).Nanoseconds(), 10) + " ns")
	} else {
		for _, a := range u.Agents {
			a.Step(u.Agents, ra)
		}
	}
}
