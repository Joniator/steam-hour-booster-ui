//go:generate npm run build
package web

import (
	"embed"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/Joniator/steam-hour-booster-ui/internal"
)

//go:embed static
var static embed.FS

//go:embed templates
var templates embed.FS

type WebServer struct {
	auth          Auth
	boosterConfig *internal.BoosterConfig
	dockerClient  *internal.DockerClient
}

func CreateWebServer(boosterConfig *internal.BoosterConfig, dockerClient *internal.DockerClient, auth Auth) *WebServer {
	return &WebServer{
		boosterConfig: boosterConfig,
		dockerClient:  dockerClient,
		auth:          auth,
	}
}

func (webServer *WebServer) Serve(port int) {
	http.HandleFunc("/delete/", webServer.deleteHandler)
	http.HandleFunc("/add", webServer.addHandler)
	http.HandleFunc("/setUser", webServer.setUserHandler)
	http.HandleFunc("/docker", webServer.dockerHandler)
	http.HandleFunc("/", webServer.getIndex)
	http.Handle("/static/", http.FileServer(http.FS(static)))

	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func (webServer *WebServer) getUserFromCookie(r *http.Request) string {
	userCookie, err := r.Cookie("shb-user")
	if err != nil {
		log.Printf("Failed to get username from request %v", err)
	} else if !webServer.boosterConfig.UserExists(userCookie.Value) {
		log.Printf("User %s does not exist", userCookie.Value)
	} else {
		return userCookie.Value
	}
	return webServer.boosterConfig.GetDefaultUser()
}

func (webServer *WebServer) isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	if webServer.auth.IsEnabled() {
		authHeader := r.Header.Get("Authorization")
		providedCredentials := strings.Trim(authHeader, "Basic ")

		username, password, err := func() (string, string, error) {
			cred, err := base64.StdEncoding.DecodeString(providedCredentials)
			if err != nil {
				return "", "", err
			}

			authSplit := strings.SplitN(string(cred), ":", 2)
			if len(authSplit) != 2 {
				return "", "", errors.New("Split failed")
			}
			return authSplit[0], authSplit[1], nil
		}()

		if err != nil || webServer.auth.Username != username || webServer.auth.Password != password {
			log.Printf("Unauthorized: %s", providedCredentials)
			w.Header().Add("WWW-Authenticate", "Basic")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return false
		}
	}
	return true
}

func (webServer *WebServer) getIndex(w http.ResponseWriter, r *http.Request) {
	if !webServer.isAuthorized(w, r) {
		return
	}
	t, _ := template.ParseFS(templates, "templates/index.html")
	user := webServer.getUserFromCookie(r)
	library, err := webServer.boosterConfig.ResolveSteamLibrary(user)

	type Context struct {
		Users             []string
		Games             []internal.Game
		User              string
		IsDockerAvailable bool
		DockerName        string
		DockerStatus      string
		DockerLogs        []string
	}
	context := Context{
		Users:             []string{"joniator", "meludoge"},
		Games:             library.Games,
		User:              library.User,
		IsDockerAvailable: webServer.dockerClient.IsAvailable(),
		DockerName:        webServer.dockerClient.ContainerName,
		DockerStatus:      webServer.dockerClient.GetStatus(),
		DockerLogs:        webServer.dockerClient.GetLogs(),
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("shb-user=%s", user))
	err = t.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Panic(err)
	}
}

func (webServer *WebServer) setUserHandler(w http.ResponseWriter, r *http.Request) {
	if !webServer.isAuthorized(w, r) {
		return
	}

	func() {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse form: %v", err)
			return
		}
		user := r.Form.Get("users")
		w.Header().Add("Set-Cookie", fmt.Sprintf("shb-user=%s", user))
	}()
	http.Redirect(w, r, "/", 301)
}

func (webServer *WebServer) addHandler(w http.ResponseWriter, r *http.Request) {
	if !webServer.isAuthorized(w, r) {
		return
	}
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
		user := webServer.getUserFromCookie(r)
		webServer.boosterConfig.AddGame(user, parsedAppId)
		if err := webServer.boosterConfig.Save(); err != nil {
			log.Print("Failed to save config")
		}

	}()
	http.Redirect(w, r, "/", 301)
}

func (webServer *WebServer) deleteHandler(w http.ResponseWriter, r *http.Request) {
	if !webServer.isAuthorized(w, r) {
		return
	}
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
		user := webServer.getUserFromCookie(r)
		webServer.boosterConfig.RemoveGame(user, parsedAppId)
		err = webServer.boosterConfig.Save()
		if err != nil {
			log.Print("Failed to save config")
			return
		}
	}()
	w.Header().Add("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", 301)
}

func (webServer *WebServer) dockerHandler(w http.ResponseWriter, r *http.Request) {
	if !webServer.isAuthorized(w, r) {
		return
	}
	func() {
		r.ParseForm()

		switch r.Form.Get("action") {
		case "restart":
			webServer.dockerClient.Restart()
			break
		case "start":
			webServer.dockerClient.Start()
			break
		case "stop":
			webServer.dockerClient.Stop()
			break
		}

	}()

	w.Header().Add("Cache-Control", "no-cache")
	http.Redirect(w, r, "/", 301)
}

func (webServer *WebServer) startHandler(w http.ResponseWriter, r *http.Request) {
	if !webServer.isAuthorized(w, r) {
		return
	}
	if !webServer.dockerClient.IsAvailable() {
		http.Error(w, "Docker not configured", 500)
		return
	}
}
