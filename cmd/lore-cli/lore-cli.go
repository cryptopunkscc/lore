package main

import (
	"fmt"
	"github.com/cryptopunkscc/lore/comm/client"
	"io"
	"log"
	"os"
	"os/exec"
)

const adminURL = "http://localhost:10768/"

type App struct {
	client *client.Client
}

func (app *App) Add(path string) {
	id, err := app.client.Admin().Add(path)
	if err != nil {
		log.Fatalln("error adding file:", err)
	}
	log.Println("added file:", id)
}

func (app *App) List() {
	list, err := app.client.Admin().List()
	if err != nil {
		log.Fatalln("api error:", err)
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

func (app *App) AddSource(address string) {
	err := app.client.Admin().AddSource(address)
	if err != nil {
		log.Fatalln("error adding source:", err)
		return
	}
	log.Println("source added")
}

func (app *App) RemoveSource(address string) {
	err := app.client.Admin().RemoveSource(address)
	if err != nil {
		log.Fatalln("error removing source:", err)
		return
	}
	log.Println("source removed")
}

func (app *App) ListSources() {
	list, err := app.client.Admin().ListSources()
	if err != nil {
		log.Fatalln("api error:", err)
	}
	fmt.Printf("%d sources on the list:\n", len(list))
	for _, i := range list {
		fmt.Println(i)
	}
}

func (app *App) Play(id string) {
	stream, err := app.client.Local().Stream(id)
	if err != nil {
		log.Fatalln("Error playing:", err)
	}

	// Add streaming player to ffplay
	cmd := exec.Command("ffplay", "-fs", "-")
	in, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	_, _ = io.Copy(in, stream)
	_ = cmd.Wait()
}

func (app *App) Run(args []string) {
	cmd := args[0]

	switch cmd {
	case "add":
		app.Add(os.Args[2])
	case "addsource":
		app.AddSource(os.Args[2])
	case "removesource":
		app.RemoveSource(os.Args[2])
	case "play":
		app.Play(os.Args[2])
	case "list":
		app.List()
	case "listsources":
		app.ListSources()
	default:
		fmt.Println("unknown command")
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: lore-cli <command> [args]")
	}

	c := client.NewClient(adminURL)
	app := &App{client: c}

	app.Run(os.Args[1:])
}
