package main

import (
	"bufio"
	"bytes"
	"flag"
	"log"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	ui "github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

type (
	Process struct {
		Pid  int
		User string
		PCPU float64
		Pmem float64
		Comm string
		Row  string
	}
	ByTopMemory []*Process
	ByTopCPU    []*Process
)

func (p ByTopMemory) Len() int           { return len(p) }
func (p ByTopMemory) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ByTopMemory) Less(i, j int) bool { return p[i].Pmem > p[j].Pmem }
func (p ByTopCPU) Len() int              { return len(p) }
func (p ByTopCPU) Swap(i, j int)         { p[i], p[j] = p[j], p[i] }
func (p ByTopCPU) Less(i, j int) bool    { return p[i].PCPU > p[j].PCPU }

var (
	dir = flag.String("dir", ".", "directory to dash")
)

func main() {

	flag.Parse()

	if err := ui.Init(); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.UseTheme("helloworld")

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(6, 0, git_diff_widget(*dir)),
			ui.NewCol(6, 0, process_hogs()),
		),
		// ui.NewRow(
		// 	ui.NewCol(3, 0, widget2),
		// 	ui.NewCol(3, 0, widget30, widget31, widget32),
		// 	ui.NewCol(6, 0, widget4),
		// ),
	)

	// calculate layout
	ui.Body.Align()

	ui.Render(ui.Body)

	// ui.Render(git_diff_widget(*dir))
	termbox.PollEvent()

}

func longest(ss []string) int {
	var l int
	for _, s := range ss {
		if l == 0 || len(s) > l {
			l = len(s)
		}
	}
	return l
}

func git_diff_widget(path string) *ui.List {

	cmd := exec.Command("git", "diff", "--stat", "master@{'1 week ago'}", path)
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	} else if stdout.Len() == 0 {
		log.Fatalf("empty response from git: %q", stderr.String())
	}

	ls := ui.NewList()
	ls.Items = strings.Split(stdout.String(), "\n")
	ls.Items = ls.Items[:len(ls.Items)-1]
	ls.ItemFgColor = ui.ColorDefault
	ls.Border.Label = "Changes < 1 week"
	ls.Height = len(ls.Items) + 1
	ls.Width = longest(ls.Items) + 3

	return ls

}

func process_hogs() *ui.List {

	cmd := exec.Command("ps", "-eo pid,user,pcpu,pmem,comm")
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	} else if stdout.Len() == 0 {
		log.Fatalf("empty response from ps: %q", stderr.String())
	}

	var (
		procs []*Process
		// header []string
	)

	for s := bufio.NewScanner(stdout); s.Scan(); {
		row := strings.Fields(s.Text())
		if len(row) != 5 {
			continue
		} else if strings.TrimSpace(row[0]) == "PID" {
			// header = row
			continue
		}
		pid, _ := strconv.Atoi(strings.TrimSpace(row[0]))
		pcpu, _ := strconv.ParseFloat(row[2], 64)
		pmem, _ := strconv.ParseFloat(row[3], 64)
		procs = append(procs, &Process{
			Pid:  pid,
			User: strings.TrimSpace(row[1]),
			PCPU: pcpu,
			Pmem: pmem,
			Comm: strings.TrimSpace(row[4]),
			Row:  s.Text(),
		})
	}

	sort.Sort(ByTopCPU(procs))
	procs = procs[:10]
	proc_list := []string{}
	for _, p := range procs[:10] {
		proc_list = append(proc_list, p.Row)
	}

	ls := ui.NewList()
	ls.Items = proc_list
	ls.ItemFgColor = ui.ColorDefault
	ls.Border.Label = "Top Processes"
	ls.Height = len(proc_list)
	ls.Width = longest(proc_list)

	return ls

}
