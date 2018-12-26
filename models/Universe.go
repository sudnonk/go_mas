package models

import (
	"github.com/sudnonk/go_mas/config"
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
	//todo: 実装
}

func (u *Universe) step() {
	for _, a := range u.Agents {
		a.Step(u.Agents)
	}
}
