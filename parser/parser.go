package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sakura-internet/go-rison"
	"github.com/sudnonk/go_mas/config"
	"github.com/sudnonk/go_mas/models"
	"github.com/sudnonk/go_mas/utils"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
)

func main() {

	app := cli.NewApp()

	app.Name = "parser"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filename, f",
			Usage: "filename you want to parse",
		},
		cli.StringFlag{
			Name:  "outdir, o",
			Usage: "Full path of outdir. Ensure last letter is DIRECTORY_SEPARATOR",
		},
		cli.Int64Flag{
			Name:  "step,s",
			Usage: "if the file has step 0-99, set 0. 100-199, set 100",
		},
		cli.Int64SliceFlag{
			Name:  "target,t",
			Usage: "Target agent ids",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "list, hp, fanatic,ideology,range",
		},
		cli.Int64Flag{
			Name:  "world,w",
			Usage: "World number",
		},
	}

	app.Action = func(ctx *cli.Context) (err error) {
		err = nil

		fn, o, s, w, t, ts := ctx.String("filename"), ctx.String("outdir"), ctx.Int64("step"), ctx.Int64("world"), ctx.String("type"), ctx.Int64Slice("target")

		if !checkArgs(fn, o, t) {
			return fmt.Errorf("-f and -t and -s and -o is required")
		}

		f, r, err := openFile(fn)
		if err != nil {
			return err
		}

		if err = ensureDir(o); err != nil {
			return err
		}

		defer func(f *os.File) {
			err = closeFile(f)
		}(f)

		switch t {
		case "fanatic":
			err = fanatic(r, o, w, s)
		case "hp":
			err = hp(r, o, w, s, ts)
		case "ideology":
			err = ideology(r, o, w, s, ts)
		case "range":
			err = ideologyRange(r, o, w, s, ts)
		case "list":
			err = list(r, o, w, s)
		case "diversity":
			err = diversity(r, o, w, s)
		default:
			return fmt.Errorf("type must be (fanatic | hp | ideology | range | list | diversity)")
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
		if arg == "" {
			return false
		}
	}

	return true
}

func ensureDir(o string) error {
	//もしoutdirが無かったら
	if _, err := os.Stat(o); os.IsNotExist(err) {
		//作る
		err = os.Mkdir(o, 0777)
		if err != nil {
			return err
		}
	}

	return nil
}

func openFile(fname string) (f *os.File, r *bufio.Reader, err error) {
	f, err = os.Open(fname)
	if err != nil {
		return nil, nil, err
	}

	return f, bufio.NewReader(f), nil
}

func createFile(fname string) (err error) {
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = closeFile(f)
	}(f)

	err = f.Truncate(0)
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	return nil
}

func writeFile(ofn string, d string) error {
	f, err := os.OpenFile(ofn, os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = closeFile(f)
	}(f)

	if _, err := f.WriteString(d); err != nil {
		return err
	}
	return nil
}

func closeFile(f *os.File) error {
	return f.Close()
}

func parseLineAll(line *[]byte) (map[int64]*models.Agent, error) {
	delimiter := []byte("=")
	as := bytes.Split(*line, delimiter)

	agents := make(map[int64]*models.Agent, config.MaxAgents)
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
		agents[ag.Id] = &ag
	}

	return agents, nil
}

func parseLineByID(line *[]byte, target []int64) (map[int64]*models.Agent, error) {
	delimiter := []byte("=")
	as := bytes.Split(*line, delimiter)

	ts := make(map[int64]*models.Agent, len(target))

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
		if utils.InArray(ag.Id, target) {
			ts[ag.Id] = &ag
		}
		if len(ts) == len(target) {
			return ts, nil
		}
	}

	return nil, fmt.Errorf("could not find enough agents")
}

func list(r *bufio.Reader, outdir string, world int64, step int64) error {
	ds := string(os.PathSeparator)

	//{hp,ideology,fanatic,range}Writers
	/*hpw := make(map[int64]string, config.MaxAgents)
	idw := make(map[int64]string, config.MaxAgents)
	raw := make(map[int64]string, config.MaxAgents)
	var fw string*/
	var err error
	var list string

	log.Println("Creating Files...")
	td := fmt.Sprintf("%slist%s", outdir, ds)
	err = ensureDir(td)
	list = fmt.Sprintf("%s%d_list_%03d.csv", td, world, step)
	err = createFile(list)
	if err != nil {
		return err
	}

	/*fw = fmt.Sprintf("%sfanatic%s%d_step_%03d.csv", outdir, ds, world, step)
	err = createFile(fw)
	if err != nil {
		return err
	}*/

	/*for i := int64(0); i < config.MaxAgents; i++ {
		// /path/to/ourdir/{hp,ideology,fanatic,range}/AgentId/
		td := fmt.Sprintf("%shp%s%d%s", outdir, ds, i, ds)
		err = ensureDir(td)
		hpw[i] = fmt.Sprintf("%s%d_step%03d.csv", td, world, step)
		err = createFile(hpw[i])
		if err != nil {
			return err
		}

		td = fmt.Sprintf("%sideology%s%d%s", outdir, ds, i, ds)
		err = ensureDir(td)
		idw[i] = fmt.Sprintf("%s%d_step%03d.csv", td, world, step)
		err = createFile(idw[i])
		if err != nil {
			return err
		}

		td = fmt.Sprintf("%srange%s%d%s", outdir, ds, i, ds)
		err = ensureDir(td)
		raw[i] = fmt.Sprintf("%s%d_step%03d.csv", td, world, step)
		err = createFile(raw[i])
		if err != nil {
			return err
		}
	}*/

	log.Println("Parsing,,,")
	//err = writeFile(fw, fmt.Sprintf("# step, Ideology, Fanatic\n"))
	err = writeFile(list, fmt.Sprintf("ID, Receptivity,Ideoloigy,len(Following), HP, Recovery\n"))
	for s := int64(0); ; s++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		ags, err := parseLineAll(&line)
		if err != nil {
			return err
		}

		is := make(map[int64]int64)
		for i := int64(0); i <= config.MaxIdeology; i++ {
			is[i] = 0
		}
		for _, ag := range ags {
			if s == int64(0) {
				/*err = writeFile(hpw[ag.Id], fmt.Sprintf("# id: %d Receptivity: %f\n# step, HP\n", ag.Id, ag.Receptivity))
				err = writeFile(idw[ag.Id], fmt.Sprintf("# id: %d Receptivity: %f\n# step, Ideology\n", ag.Id, ag.Receptivity))
				err = writeFile(raw[ag.Id], fmt.Sprintf("# id: %d Receptivity: %f\n# step, Following_Ideology\n", ag.Id, ag.Receptivity))*/
				err = writeFile(list, fmt.Sprintf("%d,%f,%d,%d,%d,%d\n", ag.Id, ag.Receptivity, ag.Ideology, len(ag.Following), ag.HP, ag.Recovery))
			}
			/*
				err = writeFile(hpw[ag.Id], fmt.Sprintf("%d,%d\n", s+step, ag.HP))
				err = writeFile(idw[ag.Id], fmt.Sprintf("%d,%d\n", s+step, ag.Ideology))

				is[ag.Ideology]++

				d := ""
				for _, agf := range ag.Following {
					d += fmt.Sprintf("%d,%d\n", s+step, ags[agf].Ideology)
				}
				err = writeFile(raw[ag.Id], d)
			*/
			if err != nil {
				return err
			}
		}
		break

		/*		for i, n := range is {
				err = writeFile(fw, fmt.Sprintf("%d,%d,%d\n", s+step, i, n))

				if err != nil {
					return err
				}
			}*/
	}
	log.Println("Parse end.")

	return nil
}

func fanatic(r *bufio.Reader, outdir string, world int64, step int64) error {
	ds := string(os.PathSeparator)

	log.Println("fanatic start.")
	log.Println("Creating Files...")
	td := fmt.Sprintf("%sfanatic%s", outdir, ds)
	err := ensureDir(td)
	file := fmt.Sprintf("%s%d_step_%03d.csv", td, world, step)
	err = createFile(file)
	if err != nil {
		return err
	}
	err = writeFile(file, fmt.Sprintf("# step, Ideology, Fanatic\n"))

	log.Println("Parsing...")
	for s := int64(0); ; s++ {
		is := make(map[int64]int64, config.MaxIdeology+1)
		for i := int64(0); i < config.MaxIdeology; i++ {
			is[i] = int64(0)
		}

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

		for _, a := range as {
			is[a.Ideology]++
		}

		d := ""
		for i, n := range is {
			d += fmt.Sprintf("%d,%d,%d\n", s+step, i, n)
		}

		if err := writeFile(file, d); err != nil {
			return err
		}
	}

	return nil
}

func hp(r *bufio.Reader, outdir string, world int64, step int64, targets []int64) error {
	ds := string(os.PathSeparator)

	log.Println("hp start.")

	log.Println("Creating Files...")
	td := fmt.Sprintf("%shp%s", outdir, ds)
	err := ensureDir(td)
	file := make(map[int64]string)

	for _, target := range targets {
		file[target] = fmt.Sprintf("%s%d_agent_%04d_step_%03d.csv", td, world, target, step)
		err = createFile(file[target])
		if err != nil {
			return err
		}
	}

	log.Println("Parsing...")
	for s := int64(0); ; s++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		ags, err := parseLineByID(&line, targets)
		if err != nil {
			return err
		}

		for _, ag := range ags {
			if s == int64(0) {
				err := writeFile(file[ag.Id], fmt.Sprintf("#Receptivity: %f\n", ag.Receptivity))
				if err != nil {
					return err
				}
			}

			err := writeFile(file[ag.Id], fmt.Sprintf("%d,%d\n", s+step, ag.HP))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ideology(r *bufio.Reader, outdir string, world int64, step int64, targets []int64) error {
	ds := string(os.PathSeparator)

	log.Println("ideology start.")

	log.Println("Creating Files...")
	td := fmt.Sprintf("%sideology%s", outdir, ds)
	err := ensureDir(td)
	file := make(map[int64]string)

	for _, target := range targets {
		file[target] = fmt.Sprintf("%s%d_agent_%04d_step_%03d.csv", td, world, target, step)
		err = createFile(file[target])
		if err != nil {
			return err
		}
	}

	log.Println("Parsing...")
	for s := int64(0); ; s++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		ags, err := parseLineByID(&line, targets)
		if err != nil {
			return err
		}

		for _, ag := range ags {
			if s == int64(0) {
				err := writeFile(file[ag.Id], fmt.Sprintf("#Receptivity: %f\n", ag.Receptivity))
				if err != nil {
					return err
				}
			}

			err := writeFile(file[ag.Id], fmt.Sprintf("%d,%d\n", s+step, ag.Ideology))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ideologyRange(r *bufio.Reader, outdir string, world int64, step int64, targets []int64) error {
	ds := string(os.PathSeparator)

	log.Println("ideologyRange start.")

	log.Println("Creating Files...")
	td := fmt.Sprintf("%srange%s", outdir, ds)
	err := ensureDir(td)
	file := make(map[int64]string)

	for _, target := range targets {
		file[target] = fmt.Sprintf("%s%d_agent_%04d_step_%03d.csv", td, world, target, step)
		err = createFile(file[target])
		if err != nil {
			return err
		}
	}

	log.Println("Parsing...")
	for s := int64(0); ; s++ {
		line, err := r.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF && len(line) == 0 {
			break
		}

		ags, err := parseLineAll(&line)
		if err != nil {
			return err
		}

		for _, aId := range targets {
			ag := ags[aId]
			//min, max := getMinMaxIdeology(ag, ags)

			if s == int64(0) {
				err := writeFile(file[aId], fmt.Sprintf("# Receptivity: %v\n# step, following_ideology\n", ag.Receptivity))
				if err != nil {
					return err
				}
			}

			for _, a := range ag.Following {
				agf := ags[a]
				err := writeFile(file[aId], fmt.Sprintf("%d,%d\n", s+step, agf.Ideology))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func diversity(r *bufio.Reader, outdir string, world int64, step int64) error {
	ds := string(os.PathSeparator)

	log.Println("diversity start.")
	log.Println("Creating Files...")
	td := fmt.Sprintf("%sdiversity%s", outdir, ds)
	err := ensureDir(td)
	file := fmt.Sprintf("%s%d_step_%03d.csv", td, world, step)
	err = createFile(file)
	if err != nil {
		return err
	}
	err = writeFile(file, fmt.Sprintf("# step, Number of Ideology\n"))

	log.Println("Parsing...")
	for s := int64(0); ; s++ {
		is := make(map[int64]struct{}, config.MaxIdeology+1)

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

		for _, a := range as {
			is[a.Ideology] = struct{}{}
		}

		d := fmt.Sprintf("%d,%d\n", s+step, len(is))

		if err := writeFile(file, d); err != nil {
			return err
		}
	}

	return nil
}
