package games

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"steam-hour-booster-ui/core/config"
)

type Game struct {
	AppId int
	Name  string
}

type Library struct {
	User  string
	Games []Game
}

var gameNamesCache = make(map[int]string)

func FromConfig(c *config.Config) Library {
	var games []Game
	for _, id := range c.AppIds {
		name := getNameForAppId(id)
		game := Game{
			Name:  name,
			AppId: id,
		}
		games = append(games, game)
	}
	return Library{
		Games: games,
		User:  c.Name,
	}
}

func getNameForAppId(appId int) string {
	cachedName := gameNamesCache[appId]
	if cachedName != "" {
		return cachedName
	}

	response, err := http.Get(fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%d", appId))
	if err != nil {
		return "Unknown"
	}

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "Unknown"
	}

	var gameMap map[string]json.RawMessage

	err = json.Unmarshal(rawBody, &gameMap)
	if err != nil {
		return "Unknown"
	}

	err = json.Unmarshal(gameMap[fmt.Sprint(appId)], &gameMap)
	if err != nil {
		return "Unknown"
	}

	err = json.Unmarshal(gameMap["data"], &gameMap)
	if err != nil {
		return "Unknown"
	}

	var name string
	err = json.Unmarshal(gameMap["name"], &name)
	if err != nil {
		return "Unknown"
	}

	gameNamesCache[appId] = name
	return name
}
