package models

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/utils"
	"math"
	"math/rand"
)

type Agent struct {
	//一意に特定するもの
	id int64
	//繋がってるエージェントのリスト
	Neighbors []Agent
	//体力
	HP int64
	//思想
	Ideology int64
	//受容しやすさ
	Receptivity float64
	//回復量
	Recovery int64
}

//二つのエージェントのイデオロギーの交流
func (a *Agent) Mix(a2 Agent) {
	if a.HP <= 0 {
		//todo: その人とのネットワークを切る処理
	}

	diff := math.Abs(float64(a.Ideology) - float64(a2.Ideology))
	a.HP -= utils.Round(diff * a.Receptivity)

	//todo: ここ冗長？
	if diff > 0 {
		//a:100,0.7 a2:0 -> a:30
		a.Ideology = a2.Ideology + utils.Round(diff*(1-a.Receptivity))
	} else if diff == 0 {

	} else {
		//a:0,0.7 a2:100 -> a:70
		a.Ideology = a.Ideology + utils.Round(diff*(a.Receptivity))
	}
}

//毎ターンの回復
func (a *Agent) Recover() {
	a.HP += a.Recovery
}

//隣人を追加
func (a *Agent) addNeighbor(as []Agent) {
	a.Neighbors = append(a.Neighbors, as...)
}

func NewAgent(id int64) Agent {
	return Agent{
		id:          id,
		Neighbors:   []Agent{},
		HP:          rand.Int63n(config.MaxHP),
		Ideology:    rand.Int63n(config.MaxIdeology),
		Receptivity: rand.NormFloat64(),
		Recovery:    rand.Int63n(config.MaxRecovery),
	}
}
