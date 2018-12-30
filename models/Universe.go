package models

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/utils"
	"math/rand"
)

type Universe struct {
	Id     int64
	Agents map[int64]*Agent
}

func (u *Universe) Init(id int64, ra *rand.Rand) {
	u.Id = id
	u.Agents = make(map[int64]*Agent, config.MaxAgents)

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
	for _, a := range u.Agents {
		a.Step(u.Agents, ra)
	}
}
