package main

import (
	"fmt"
	_id "github.com/cryptopunkscc/lore/id"
	"github.com/cryptopunkscc/lore/store/http"
	"io"
	"log"
	"os"
	"os/exec"
)

const storeURL = "http://localhost:10768/store"

type App struct {
	store *http.HTTPStore
}

// List show a list of all shared files
func (app *App) List() {
	list, err := app.store.List()
	if err != nil {
		log.Fatalln("api error:", err)
	}

	list.Each(func(id _id.ID) {
		fmt.Println(id)
	})
}

func (app *App) Read(idStr string) {
	id, err := _id.Parse(idStr)
	if err != nil {
		log.Fatalln("error parsing argument:", err)
	}

	file, err := app.store.Read(id)
	if err != nil {
		log.Fatalln("api error:", err)
	}

	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		log.Fatalln("api error:", err)
	}
}

func (app *App) Create() {
	file, err := app.store.Create()
	if err != nil {
		log.Fatalln("api error:", err)
	}

	_, err = io.Copy(file, os.Stdin)
	if err != nil {
		log.Fatalln("api error:", err)
	}

	id, err := file.Finalize()
	if err != nil {
		log.Fatalln("api error:", err)
	}

	fmt.Println(id)
}

func (app *App) Delete(idStr string) {
	id, err := _id.Parse(idStr)
	if err != nil {
		log.Fatalln("error parsing argument:", err)
	}

	err = app.store.Delete(id)
	if err != nil {
		log.Fatalln("api error:", err)
	}
}

// Play reads a file and plays it locally using ffplay
func (app *App) Play(idStr string) {
	id, err := _id.Parse(idStr)
	if err != nil {
		log.Fatalln("error parsing argument:", err)
	}

	stream, err := app.store.Read(id)
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

// Run executes the command provided by the user
func (app *App) Run(args []string) {
	cmd := args[0]

	switch cmd {
	case "list":
		app.List()
	case "create":
		app.Create()
	case "read":
		app.Read(args[1])
	case "delete":
		app.Delete(args[1])
	case "play":
		app.Play(args[1])
	default:
		fmt.Println("unknown command", cmd)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("Usage: lore-cli <command> [args]")
	}

	app := &App{
		store: http.NewHTTPStore(storeURL),
	}

	app.Run(os.Args[1:])
}
