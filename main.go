package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gizak/termui"
	"github.com/nsf/termbox-go"
)

func main() {

	err := termui.Init()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, os.Interrupt, os.Kill)
		<-c
		termui.Close()
	}()

	defer close(c)

	termui.UseTheme("helloworld")

	strs := []string{
		"[0] github.com/gizak/termui",
		"[1] editbox.go",
		"[2] iterrupt.go",
		"[3] keyboard.go",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go"}

	ls := termui.NewList()
	ls.Items = strs
	ls.ItemFgColor = termui.ColorYellow
	ls.Border.Label = "List"
	ls.Height = 7
	ls.Width = 25
	ls.Y = 0

	termui.Render(ls)

	termbox.PollEvent()
}
