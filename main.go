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
	ConfigFilePath string `arg:"--config,-c" help:"Path to the config file" default:"config.json"`
	ContainerName  string `arg:"--container" help:"Name of the container" default:"booster"`
}

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

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
	http.HandleFunc("/docker", dockerHandler)
	http.HandleFunc("/", getIndex)
	http.Handle("/static/", http.FileServer(http.FS(static)))

	log.Print("Listening on port :35888")
	log.Fatal(http.ListenAndServe(":35888", nil))
}

func getIndex(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFS(templates, "templates/index.html")
	user := getUserFromCookie(r)
	g := games.FromConfig(configs, user)

	type Context struct {
		Games             []games.Game
		User              string
		IsDockerAvailable bool
		DockerStatus      string
		DockerLogs        []string
	}
	context := Context{
		Games:             g.Games,
		User:              g.User,
		IsDockerAvailable: dc.IsAvailable(),
		DockerStatus:      dc.GetStatus(),
		DockerLogs:        dc.GetLogs(),
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("shb-user=%s", g.User))
	log.Printf("Context: %=v", context)
	err := t.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Panic(err)
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	func() {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse form: %v", err)
			return
		}
		appId := strings.Trim(r.FormValue("AppId"), " ")
		parsedAppId, err := strconv.Atoi(appId)

		if err != nil {
			log.Printf("Failed to parse appId add '%s'", appId)
			return
		}
		user := getUserFromCookie(r)
		config.Add(configs, user, parsedAppId)
		if err := config.Save(configs, args.ConfigFilePath); err != nil {
			log.Print("Failed to save config")
		}

	}()
	http.Redirect(w, r, "/", 301)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	func() {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse form: %v", err)
			return
		}
		appId := r.Form.Get("item")
		parsedAppId, err := strconv.Atoi(appId)

		if err != nil {
			log.Printf("Failed to parse appId '%s'", appId)
			return
		}
		user := getUserFromCookie(r)
		config.Remove(configs, user, parsedAppId)
		err = config.Save(configs, args.ConfigFilePath)
		if err != nil {
			log.Print("Failed to save config")
			return
		}
	}()
	w.Header().Add("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", 301)
}

func dockerHandler(w http.ResponseWriter, r *http.Request) {
	func() {
		r.ParseForm()

		switch r.Form.Get("action") {
		case "restart":
			dc.Restart()
			break
		case "start":
			dc.Start()
			break
		case "stop":
			dc.Stop()
			break
		}

	}()

	w.Header().Add("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", 301)
}

func getUserFromCookie(r *http.Request) string {
	userCookie, err := r.Cookie("shb-user")
	if err != nil {
		log.Printf("Failed to get username from request %v", err)
		return ""
	}
	return userCookie.Value
}

func startHandler(w http.ResponseWriter, r *http.Request) {
	if !dc.IsAvailable() {
		http.Error(w, "Docker not configured", 500)
		return
	}
}
