//go:generate npm run build
package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"steam-hour-booster-ui/core/config"
	"steam-hour-booster-ui/core/docker"
	"steam-hour-booster-ui/core/games"
	"strconv"
	"strings"

	"github.com/alexflint/go-arg"
)

var args struct {
	ConfigFilePath      string `arg:"--config,-c" help:"Path to the config file" default:"config.json"`
	ContainerName string `arg:"--container" help:"Name of the container"`
}

//go:embed static
var static embed.FS
var configs *[]config.Config
var dc docker.DockerClient

func main() {
	arg.MustParse(&args)
	dc = docker.New(args.ContainerName)
	var err error
	configs, err = config.LoadConfig(args.ConfigFilePath)
	if err != nil {
		log.Panic(err)
	}

	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/", getIndex)
	http.Handle("/static/", http.FileServer(http.FS(static)))

	log.Print("Listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	g := games.FromConfig(configs, w.Header().Get("shb-user"))


	type Context struct {
		Games             []games.Game
		User              string
		IsDockerAvailable bool
	}
	context := Context{
		Games:             g.Games,
		User:              g.User,
		IsDockerAvailable: dc.IsAvailable(),
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("user=%s", g.User))
	err := t.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Panic(err)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	appId := strings.Trim(r.FormValue("AppId"), " ")
	parsedAppId, err := strconv.Atoi(appId)

	if err != nil {
		log.Printf("Failed to parse appId %s", appId)
	} else {
		config.Add(configs, w.Header().Get("shb-user"), parsedAppId)
		if err := config.Save(configs, args.ConfigFilePath); err != nil {
			log.Print("Failed to save config")
		}
	}
	http.Redirect(w, r, "/", 301)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	appId := strings.TrimPrefix(r.URL.Path, "/delete/")
	parsedAppId, err := strconv.Atoi(appId)
	log.Printf("%+v", configs)

	if err != nil {
		log.Printf("Failed to parse appId %s", appId)
	} else {
		config.Remove(configs, w.Header().Get("shb-user"), parsedAppId)

		if err := config.Save(configs, args.ConfigFilePath); err != nil {
			log.Print("Failed to save config")
		}
	}

	w.Header().Add("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", 301)
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	if !dc.IsAvailable() {
		http.Error(w, "Docker not configured", 500)
		return
	}
}
