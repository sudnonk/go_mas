package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sakura-internet/go-rison"
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {

	app := cli.NewApp()

	app.Name = "parser"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "filename",
			Value: nil,
		},
		cli.StringFlag{
			Name:  "prefix, p",
			Usage: "prefix for output file.(e,g, output dir path)",
			Value: nil,
		},
		cli.Int64Flag{
			Name:  "target",
			Usage: "target agent id.",
			Value: nil,
		},
		cli.StringFlag{
			Name:  "type,t",
			Usage: "Parse type.",
			Value: nil,
		},
	}

	app.Action = func(ctx *cli.Context) (err error) {
		err = nil

		fn, t, p := ctx.String("filename"), ctx.String("type"), ctx.String("prefix")
		//--targetが未指定の時0になるので必然的にtargetはAgent.ID=0
		target := ctx.Int64("target")

		if !checkArgs(fn, t, p) {
			os.Exit(1)
		}

		f, r, err := openFile(fn)
		if err != nil {
			return err
		}

		defer func(f *os.File) {
			err = closeFile(f)
		}(f)

		switch t {
		case "fanatic":
			err = fanatic(r, p)
		case "hp":
			err = hp(r, target, p)
		case "ideology":
			err = ideology(r, target, p)
		default:
			os.Exit(1)
		}

		return err
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func checkArgs(args ...interface{}) bool {
	for _, arg := range args {
		if arg == nil {
			return false
		}
	}

	return true
}

func openFile(fname string) (f *os.File, r *bufio.Reader, err error) {
	f, err = os.Open(fname)
	if err != nil {
		return nil, nil, err
	}

	return f, bufio.NewReader(f), nil
}

func writeFile(ofn string, d *string) error {
	o, err := os.OpenFile(ofn, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	if _, err := o.WriteString(*d); err != nil {
		return err
	}
	if err := o.Close(); err != nil {
		return err
	}
	return nil
}

func closeFile(f *os.File) error {
	return f.Close()
}

func parseLineAll(line *[]byte) (*[]models.Agent, error) {
	delimiter := []byte("=")
	as := bytes.Split(*line, delimiter)

	agents := make([]models.Agent, config.MaxAgents)
	for i, a := range as {
		if len(a) < 5 {
			continue
		}

		var ag models.Agent
		err := rison.Unmarshal(a, &ag, rison.ORison)
		if err != nil {
			log.Println(err)
			continue
		}
		agents[i] = ag
	}

	return &agents, nil
}

func parseLineByID(line *[]byte, target int64) (*models.Agent, error) {
	delimiter := []byte("=")
	as := bytes.Split(*line, delimiter)

	for _, a := range as {
		if len(a) < 5 {
			continue
		}

		var ag models.Agent
		err := rison.Unmarshal(a, &ag, rison.ORison)
		if err != nil {
			log.Println(err)
			continue
		}
		if ag.Id == target {
			return &ag, nil
		}
	}

	return nil, fmt.Errorf("could not find Agent.Id =  %d ", target)
}

func fanatic(r *bufio.Reader, prefix string) error {
	is := make(map[int64]int64, config.MaxAgents)

	for step := 1; ; step++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		as, err := parseLineAll(&line)
		if err != nil {
			return err
		}

		for _, a := range *as {
			if _, ok := is[a.Ideology]; !ok {
				is[a.Ideology] = int64(0)
			}
			is[a.Ideology]++
		}

		d := ""
		for i, n := range is {
			d += strconv.FormatInt(i, 10) + "," + strconv.FormatInt(n, 10) + "\n"
		}

		if err := writeFile(fmt.Sprintf("%s_ideology_step%03d.csv", prefix, step), &d); err != nil {
			return err
		}
	}

	return nil
}

func hp(r *bufio.Reader, target int64, prefix string) error {
	d := ""

	for step := int64(1); ; step++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		ag, err := parseLineByID(&line, target)
		if err != nil {
			return err
		}

		d += strconv.FormatInt(step, 10) + "," + strconv.FormatInt(ag.HP, 10) + "\n"
	}

	if err := writeFile(fmt.Sprintf("%s_%d_hp.csv", prefix, target), &d); err != nil {
		return err
	}

	return nil
}

func ideology(r *bufio.Reader, target int64, prefix string) error {
	d := ""

	for step := int64(1); ; step++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		ag, err := parseLineByID(&line, target)
		if err != nil {
			return err
		}

		d += strconv.FormatInt(step, 10) + "," + strconv.FormatInt(ag.Ideology, 10) + "\n"
	}

	if err := writeFile(fmt.Sprintf("%s_%d_ideology.csv", prefix, target), &d); err != nil {
		return err
	}

	return nil
}
