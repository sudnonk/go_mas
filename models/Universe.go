package models

import "github.com/sudnonk/go_mas/config"

type Universe struct {
	Agents map[int64]Agent
}

func (u *Universe) Init() {
	var i int64
	for i = 0; i < config.MaxAgents; i++ {
		u.Agents[i] = NewAgent(i)
	}
}

func (u *Universe) MakeNetwork() {
	//todo: 実装
}
