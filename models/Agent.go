package models

import (
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/utils"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"
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

func (a *Agent) Step(as map[int64]*Agent, ra *rand.Rand) {
	s := time.Now().Nanosecond()
	for _, aID := range a.Following {
		if a.HP <= 0 {
			break
		}
		diff := float64(a.Ideology - as[aID].Ideology)
		a.HP -= utils.Round(math.Abs(diff) * (1 - a.Receptivity)) //受容性が低いほど多くのHPを消費して
		a.mix(diff, ra)                                           //思想が混ざる
	}

	//1秒以上かかったら
	if time.Now().Nanosecond()-s > 1000000000 {
		log.Println("HP process took > 1s")
	}

	s = time.Now().Nanosecond()
	//体力がなくなると違うイデオロギーのフォローを外す
	if a.HP <= 0 {
		a.unfollowDifferentIdeology(as)
	}

	//1秒以上かかったら
	if time.Now().Nanosecond()-s > 1000000000 {
		log.Println("Unfollow process took > 1s")
	}

	//受容性が高い人ほど高い値が出る
	followCriteria := a.Receptivity * utils.RandNormDecimal(ra)

	s = time.Now().Nanosecond()

	if followCriteria > 0.7 {
		a.followDifferentIdeology(as, ra)
	} else if followCriteria > 0.3 {
		a.followInfluencer(as)
	} else {
		a.followNearIdeology(as, ra)
	}

	//1秒以上かかったら
	if time.Now().Nanosecond()-s > 1000000000 {
		log.Println("Follow process took > 1s")
	}

	a.recover() //毎ターンの回復
}

//二つのエージェントのイデオロギーの交流
func (a *Agent) mix(diff float64, ra *rand.Rand) {
	//混ざり具合
	mixture := ra.Float64()

	if diff > 0 {
		//a:100,0.7 a2:0 -> a:30
		a.Ideology -= utils.Round(diff * a.Receptivity * mixture)
	} else if diff == 0 {
		//a.Receptivity -= a.Receptivity * mixture * config.Pride
	} else {
		//a:0,0.7 a2:100 -> a:70
		a.Ideology += utils.Round(math.Abs(diff) * a.Receptivity * mixture)
	}
}

//毎ターンの回復
func (a *Agent) recover() {
	a.HP += a.Recovery
}

func (a *Agent) unfollowDifferentIdeology(as map[int64]*Agent) {
	//Ideologyが小さい順に並び替え
	sort.Slice(a.Following, func(i, j int) bool {
		return as[a.Following[i]].Ideology < as[a.Following[j]].Ideology
	})
	//Ideologyが小さい方までの距離が大きい方までの距離より大きければ
	if utils.Abs(a.Ideology-as[a.Following[0]].Ideology) > utils.Abs(a.Ideology-as[a.Following[len(a.Following)-1]].Ideology) {
		//一番小さい人を外す
		a.Following = a.Following[1:]
	} else {
		//一番大きい人を外す
		a.Following = a.Following[:len(a.Following)-1]
	}
}

//近い意見の人をフォローする（自分とのイデオロギー差が最小の人をフォローする）
func (a *Agent) followNearIdeology(as map[int64]*Agent, ra *rand.Rand) {
	minI := int64(config.MaxIdeology + 1)

	d := int64(len(as) + 1)
	minR := d

	for r := int64(0); r < config.MaxAgents; r++ {
		if utils.Abs(as[r].Ideology-a.Ideology) < minI {
			minI = utils.Abs(as[r].Ideology - a.Ideology)
			minR = r
		}
	}

	if minR == d {
		//todo: フォローすべき相手が見つからなかった場合
	} else {
		a.Following = append(a.Following, minR)
	}
}

//違う意見の人をフォローする（自分とのイデオロギー差が最大の人をフォローする）
func (a *Agent) followDifferentIdeology(as map[int64]*Agent, ra *rand.Rand) {
	maxI := int64(config.MaxIdeology + 1)

	d := int64(len(as) + 1)
	maxR := d

	for r := int64(0); r < config.MaxAgents; r++ {
		if utils.Abs(as[r].Ideology-a.Ideology) > maxI {
			maxI = as[r].Ideology
			maxR = r
		}
	}

	if maxR == d {
		//todo: フォローすべき相手が見つからなかった場合
	} else {
		a.Following = append(a.Following, maxR)
	}
}

//フォローしてる人が多い人をフォローする
func (a *Agent) followInfluencer(as map[int64]*Agent) {
	maxA := a.Id
	maxL := 0

	for a2 := range as {
		if !utils.InArray(a2, a.Following) {
			a3 := as[a2]
			if l := len(a3.Following); l > maxL {
				maxA = a3.Id
				maxL = l
			}
		}
	}

	a.Following = append(a.Following, maxA)
}

func NewAgent(id int64, ra *rand.Rand, isNorm bool) *Agent {
	hp := ra.Int63n(config.MaxHP)
	var i int64
	var r float64
	if isNorm {
		i = utils.Round(utils.RandNormDecimal(ra) * config.MaxIdeology)
		r = utils.RandNormDecimal(ra)
	} else {
		i = utils.Round(ra.Float64() * config.MaxIdeology)
		r = ra.Float64()
	}
	return &Agent{
		Id:          id,
		Following:   []int64{},
		HP:          hp,
		Ideology:    i,
		Receptivity: r,
		Recovery:    utils.Round(float64(hp) * config.RecoveryRate),
	}
}
