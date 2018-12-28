package models

import (
	"github.com/sakura-internet/go-rison"
	"log"
	"os"
)

func LogStepChan(cu chan *Universe, cf chan *os.File) {
	for u := range cu {
		for file := range cf {
			lr := []byte("\n")
			q := []byte("=")
			var l []byte
			for _, a := range u.Agents {
				r, _ := rison.Marshal(&a, rison.ORison)
				l = append(l, r...)
				l = append(l, q...)
			}
			l = append(l, lr...)

			if _, err := file.Write(l); err != nil {

			}
		}
	}
}

func LogStep(u *Universe, fname string) {
	lr := []byte("\n")
	q := []byte("=")
	var l []byte
	for _, a := range u.Agents {
		r, _ := rison.Marshal(&a, rison.ORison)
		l = append(l, r...)
		l = append(l, q...)
	}
	l = append(l, lr...)

	writeLog(&l, fname)
}

func writeLog(data *[]byte, fname string) {
	file, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := file.Write(*data); err != nil {
		log.Println(err)
	}

	if err := file.Close(); err != nil {
		log.Println(err)
	}
}
