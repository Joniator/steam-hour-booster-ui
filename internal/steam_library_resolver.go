package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
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

func (boosterConfig *BoosterConfig) ResolveSteamLibrary(user string) (*Library, error) {
	var games []Game
	c, err := boosterConfig.GetUserConfig(user)
	if err != nil {
		return nil, err
	}
	for _, id := range c.AppIds {
		name, err := getNameForAppId(id)
		if err != nil {
			log.Printf("Failed to load appId=%d, err=%s", id, err.Error())
		}
		game := Game{
			Name:  name,
			AppId: id,
		}
		games = append(games, game)
	}
	return &Library{
		Games: games,
		User:  c.Name,
	}, nil
}

func getNameForAppId(appId int) (string, error) {
	cachedName := gameNamesCache[appId]
	if cachedName != "" {
		return cachedName, nil
	}

	response, err := http.Get(fmt.Sprintf("http://store.steampowered.com/api/appdetails?appids=%d", appId))
	if err != nil {
		log.Printf("Failed to load from steam API: %v", err)
		return "Unknown", errors.New("Failed to load")
	}

	rawBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "Unknown", errors.New("Failed to load")
	}

	var gameMap map[string]json.RawMessage

	err = json.Unmarshal(rawBody, &gameMap)
	if err != nil {
		return "Unknown", errors.New("Failed to load")
	}

	err = json.Unmarshal(gameMap[fmt.Sprint(appId)], &gameMap)
	if err != nil {
		return "Unknown", errors.New("Failed to load")
	}

	err = json.Unmarshal(gameMap["data"], &gameMap)
	if err != nil {
		return "Unknown", errors.New("Failed to load")
	}

	var name string
	err = json.Unmarshal(gameMap["name"], &name)
	if err != nil {
		return "Unknown", errors.New("Failed to load")
	}

	gameNamesCache[appId] = name
	return name, nil
}
