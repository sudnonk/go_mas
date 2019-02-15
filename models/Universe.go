package models

import (
	"math/rand"
)

type Universe struct {
	Id     int64
	Agents map[int64]*Agent
}

func (u *Universe) Init(id int64, as map[int64]*Agent) {
	u.Id = id
	u.Agents = as
}

func (u *Universe) Step(ra *rand.Rand) {
	for _, a := range u.Agents {
		a.Step(u.Agents, ra)
	}
}
