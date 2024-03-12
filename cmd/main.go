package main

import (
	"log"

	"github.com/Joniator/steam-hour-booster-ui/internal"
	"github.com/Joniator/steam-hour-booster-ui/web"
	"github.com/alexflint/go-arg"
)

var args struct {
	ConfigFilePath string `arg:"--config,-c" help:"Path to the config file" default:"config.json"`
	ContainerName  string `arg:"--container" help:"Name of the container" default:"steam_hour_booster"`
	Username       string `arg:"--user,-u" help:"Username for basic auth"`
	Password       string `arg:"--password,-p" help:"Password for basic auth"`
	Port           int    `arg:"--port,-P" default:"80" help:"The port to listen on"`
}

func main() {
	arg.MustParse(&args)
	configs, err := internal.LoadBoosterConfig(args.ConfigFilePath)
	dockerClient := internal.NewDockerClient(args.ContainerName)
	if err != nil {
		log.Panic(err)
	}
	auth := web.Auth{
		Username: args.Username,
		Password: args.Password,
	}
	webServer := web.CreateWebServer(configs, dockerClient, auth)
	webServer.Serve(args.Port)
}
