package models

import (
	"github.com/sudnonk/go_mas/config"
	"math/rand"
)

type Universe struct {
	Agents map[int64]Agent
}

func (u *Universe) Init() {
	u.Agents = map[int64]Agent{}

	for i := int64(0); i < config.MaxAgents; i++ {
		u.Agents[i] = NewAgent(i)
	}

	u.makeNetwork()
}

func (u *Universe) makeNetwork() {
	//todo: より良いネットワーク
	for _, a := range u.Agents {
		b := map[int64]bool{a.id: true}
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
	for _, a := range u.Agents {
		a.Step(u.Agents)
	}
}

func (u *Universe) End() {

}
