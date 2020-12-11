package main

import (
	"fmt"
	"github.com/cryptopunkscc/lore/comm/client"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const adminURL = "http://localhost:10768/"

type App struct {
	client *client.Client
}

// Add adds a local file to shared files
func (app *App) Add(path string) {
	var err error

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln(err)
	}

	id, err := app.client.Admin().Add(absPath)
	if err != nil {
		log.Fatalln("error adding file:", err)
	}

	log.Println("added file:", id)
}

// List show a list of all shared files
func (app *App) List() {
	list, err := app.client.Admin().List()
	if err != nil {
		log.Fatalln("api error:", err)
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

// Search searches shared files by name
func (app *App) Search(query string) {
	list, err := app.client.Admin().Search(query)
	if err != nil {
		log.Fatalln("api error:", err)
	}
	for _, item := range list {
		fmt.Println(item)
	}
}

// AddSource adds an address to the sources list
func (app *App) AddSource(address string) {
	err := app.client.Admin().AddSource(address)
	if err != nil {
		log.Fatalln("error adding source:", err)
		return
	}
	log.Println("source added")
}

// RemoveSource removes an address from the sources list
func (app *App) RemoveSource(address string) {
	err := app.client.Admin().RemoveSource(address)
	if err != nil {
		log.Fatalln("error removing source:", err)
		return
	}
	log.Println("source removed")
}

// ListSources fetches the sources list
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

// PlayByID searches for a file and plays it locally using ffplay
func (app *App) PlayByID(id string) {
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

// Play searches and plays all files matching name
func (app *App) Play(name string) {
	list, err := app.client.Admin().Search(name)
	if err != nil {
		log.Fatalln("api error:", err)
	}
	for _, i := range list {
		app.PlayByID(i)
	}
}

// Run executes the command provided by the user
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
	case "search":
		app.Search(os.Args[2])
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
