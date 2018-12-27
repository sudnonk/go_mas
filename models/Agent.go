package models

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/utils"
	"math"
	"math/rand"
)

type Agent struct {
	//一意に特定するもの
	Id int64 `json:"id"`
	//このエージェントがフォローしてるエージェントIDのリスト
	Following []int64 `json:"f"`
	//体力
	HP int64 `json:"h"`
	//思想
	Ideology int64 `json:"i"`
	//受容しやすさ
	Receptivity float64 `json:"rcp"`
	//回復量
	Recovery int64 `json:"rcv"`
}

func (a *Agent) Step(as map[int64]Agent) {
	for _, aID := range a.Following {
		a2 := as[aID]
		diff := math.Abs(float64(a.Ideology) - float64(a2.Ideology))
		a.damage(diff)  //HPを消費して
		a.mix(a2, diff) //思想が混ざる
	}

	//体力がなくなると違うイデオロギーのフォローを外す
	if a.HP <= 0 {
		a.unfollowDifferentIdeology(as)
	}

	//受容性が高い人ほど高い値が出る
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

func (a *Agent) unfollowDifferentIdeology(as map[int64]Agent) {
	maxA := a.Id
	maxD := int64(0)

	for _, aID := range a.Following {
		a2 := as[aID]
		if diff := utils.Abs(a2.Ideology - a.Ideology); diff > maxD {
			maxA = a2.Id
			maxD = diff
		}
	}

	var f []int64
	for _, aID := range a.Following {
		if aID != maxA {
			f = append(f, aID)
		}
	}

	a.Following = f
}

//近い意見の人をフォローする
func (a *Agent) followNearIdeology(as map[int64]Agent) {
	maxI := int64(float64(a.Ideology) * (1.0 + config.NearCriteria))
	minI := int64(float64(a.Ideology) * (1.0 - config.NearCriteria))

	checked := map[int64]bool{a.Id: true}
	for true {
		r := rand.Int63n(config.MaxAgents)
		if _, ok := checked[r]; ok {
			continue
		}

		if as[r].Ideology < maxI || as[r].Ideology > minI {
			a.Following = append(a.Following, r)
			return
		}

		checked[r] = true
	}

	//todo: フォローすべき相手が見つからなかった場合
}

//違う意見の人をフォローする
func (a *Agent) followDifferentIdeology(as map[int64]Agent) {
	maxI := int64(float64(a.Ideology) * (1.0 + config.FarCriteria))
	minI := int64(float64(a.Ideology) * (1.0 - config.FarCriteria))

	checked := map[int64]bool{a.Id: true}
	for true {
		r := rand.Int63n(config.MaxAgents)
		if _, ok := checked[r]; ok {
			continue
		}

		if as[r].Ideology > maxI || as[r].Ideology < minI {
			a.Following = append(a.Following, r)
			return
		}

		checked[r] = true
	}

	//todo: フォローすべき相手が見つからなかった場合
}

//フォローしてる人が多い人をフォローする
func (a *Agent) followInfluencer(as map[int64]Agent) {
	maxA := a.Id
	maxL := 0

	for _, aID := range a.Following {
		a2 := as[aID]
		if l := len(a2.Following); l > maxL {
			maxA = a2.Id
			maxL = l
		}
	}

	a.Following = append(a.Following, maxA)
}

func NewAgent(id int64) Agent {
	return Agent{
		Id:          id,
		Following:   []int64{},
		HP:          rand.Int63n(config.MaxHP),
		Ideology:    rand.Int63n(config.MaxIdeology),
		Receptivity: utils.RandNormDecimal(),
		Recovery:    rand.Int63n(config.MaxRecovery),
	}
}
