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
	//このエージェントがフォローしてるエージェントのリスト
	Following []Agent
	//体力
	HP int64
	//思想
	Ideology int64
	//受容しやすさ
	Receptivity float64
	//回復量
	Recovery int64
}

func (a *Agent) Step(as []Agent) {
	for _, a2 := range a.Following {
		diff := math.Abs(float64(a.Ideology) - float64(a2.Ideology))
		a.damage(diff)  //HPを消費して
		a.mix(a2, diff) //思想が混ざる
	}

	//体力がなくなると違うイデオロギーのフォローを外す
	if a.HP <= 0 {
		a.unfollowDifferentIdeology()
	}

	followCriteria := a.Receptivity * utils.RandNormDecimal()

	if followCriteria > 0.7 {
		a.followDifferentIdeology(as)
	} else if followCriteria > 0.3 {
		a.followInfluencer(as)
	} else {
		a.followNearIdeology(as)
	}

	a.recover() //毎ターンの回復
}

//二つのエージェントのイデオロギーの交流
func (a *Agent) mix(a2 Agent, diff float64) {
	//混ざり具合
	mixture := utils.RandDecimal()

	//todo: ここ冗長？
	if diff > 0 {
		//a:100,0.7 a2:0 -> a:30
		a.Ideology = a.Ideology - utils.Round(diff*(1-a.Receptivity)*mixture)
	} else if diff == 0 {

	} else {
		//a:0,0.7 a2:100 -> a:70
		a.Ideology = a.Ideology + utils.Round(diff*(a.Receptivity)*mixture)
	}
}

func (a *Agent) damage(diff float64) {
	a.HP -= utils.Round(diff * a.Receptivity)
}

//毎ターンの回復
func (a *Agent) recover() {
	a.HP += a.Recovery
}

func (a *Agent) unfollowDifferentIdeology() {
	//todo: implement
	for _, a2 := range a.Following {

	}
}

//近い意見の人をフォローする
func (a *Agent) followNearIdeology(as []Agent) {
	maxI := int64(float64(a.Ideology) * (1.0 + config.NearCriteria))
	minI := int64(float64(a.Ideology) * (1.0 - config.NearCriteria))

	checked := map[int64]bool{a.id: true}
	for true {
		r := rand.Int63n(config.MaxAgents)
		if _, ok := checked[r]; ok {
			continue
		}

		if as[r].Ideology < maxI || as[r].Ideology > minI {
			a.Following = append(a.Following, as[r])
			return
		}

		checked[r] = true
	}

	//todo: フォローすべき相手が見つからなかった場合
}

//違う意見の人をフォローする
func (a *Agent) followDifferentIdeology(as []Agent) {
	maxI := int64(float64(a.Ideology) * (1.0 + config.FarCriteria))
	minI := int64(float64(a.Ideology) * (1.0 - config.FarCriteria))

	checked := map[int64]bool{a.id: true}
	for true {
		r := rand.Int63n(config.MaxAgents)
		if _, ok := checked[r]; ok {
			continue
		}

		if as[r].Ideology > maxI || as[r].Ideology < minI {
			a.Following = append(a.Following, as[r])
			return
		}

		checked[r] = true
	}

	//todo: フォローすべき相手が見つからなかった場合
}

//フォロワーの多い人をフォローする
func (a *Agent) followInfluencer(as []Agent) {
	//todo:implement
}

func NewAgent(id int64) Agent {
	return Agent{
		id:          id,
		Following:   []Agent{},
		HP:          rand.Int63n(config.MaxHP),
		Ideology:    rand.Int63n(config.MaxIdeology),
		Receptivity: utils.RandNormDecimal(),
		Recovery:    rand.Int63n(config.MaxRecovery),
	}
}
