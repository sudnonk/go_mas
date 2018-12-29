package main

import (
	"bufio"
	"bytes"
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
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "output file name prefix",
		},
	}

	app.Action = func(ctx *cli.Context) error {
		if err := parse(ctx.String("filename"), ctx.String("output")); err != nil {
			return err
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func parse(fname string, output string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}

	r := bufio.NewReader(f)
	for i := 0; ; i++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		if err := parseLine(&line, output, i); err != nil {
			if err2 := f.Close(); err2 != nil {
				return err2
			}
			return err
		}
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func parseLine(line *[]byte, prefix string, step int) error {
	is := make(map[int64]int64, config.MaxIdeology+1)

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

		if _, ok := is[ag.Ideology]; !ok {
			is[ag.Ideology] = 0
		}
		is[ag.Ideology]++
	}

	d := ""
	for i, n := range is {
		d += strconv.FormatInt(i, 10) + "," + strconv.FormatInt(n, 10) + "\n"
	}

	o, err := os.OpenFile(prefix+"_step"+strconv.Itoa(step*100)+".csv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	if _, err := o.WriteString(d); err != nil {
		return err
	}
	if err := o.Close(); err != nil {
		return err
	}

	return nil
}
