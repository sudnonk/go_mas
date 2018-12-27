package models

import (
	"github.com/sakura-internet/go-rison"
	"os"
	"strconv"
)

func LogStep(u Universe, file *os.File) {
	if u.StepNum%100 != 0 {
		return
	}

	l := []byte(strconv.FormatInt(u.Id, 10) + ",")
	for _, a := range u.Agents {
		r, _ := rison.Marshal(&a, rison.ORison)
		l = append(l, r...)
	}

	go func() {
		if _, err := file.Write(l); err != nil {

		}
	}()
}

func LogInit(u Universe) {

}

func LogEnd(u Universe) {

}
