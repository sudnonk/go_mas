package models

import (
	"github.com/sudnonk/go_mas/config"
	"math/rand"
)

type Universe struct {
	Id      int64
	Agents  map[int64]Agent
	StepNum int64
}

func (u *Universe) Init(id int64) {
	u.Id = id
	u.Agents = map[int64]Agent{}
	u.StepNum = 0

	c := make(chan Agent)
	for i := int64(0); i < config.MaxAgents; i++ {
		NewAgent(i, c)
	}
	for i := int64(0); i < config.MaxAgents; i++ {
		u.Agents[i] = <-c
	}

	u.makeNetwork()
}

func (u *Universe) makeNetwork() {
	//todo: より良いネットワーク
	for _, a := range u.Agents {
		b := map[int64]bool{a.Id: true}
		for i := int64(0); i < rand.Int63n(config.MaxAgents); i++ {
			id := rand.Int63n(config.MaxAgents)
			if _, ok := b[id]; ok {
				continue
			}
			b[id] = true

			a.Following = append(a.Following, id)
		}
	}
}

func (u *Universe) Step() {
	u.StepNum++

	for _, a := range u.Agents {
		a.Step(u.Agents)
	}
}

func (u *Universe) End() {

}
