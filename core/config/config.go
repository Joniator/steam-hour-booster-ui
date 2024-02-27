package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"slices"
)

type Config struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	AppIds   []int  `json:"games"`
}

func LoadConfig(path string) (*[]Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	configRaw, _ := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config []Config
	err = json.Unmarshal(configRaw, &config)
	if err != nil {
		return nil, err
	}
	log.Printf("Config loaded: %+v", config)
	for _, c := range config {
		slices.Sort(c.AppIds)
	}

	return &config, nil
}

func Save(c *[]Config, path string) error {
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}

	_, err = configFile.Write(b)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Config saved: %+v", c)
	return nil
}

func Add(configs *[]Config, user string, appId int) {
	c := GetUserConfig(configs, user)
	for _, id := range c.AppIds {
		if id == appId {
			log.Printf("Game %d already in config", appId)
			return
		}
	}
	log.Printf("Adding %d to %s config", appId, user)
	c.AppIds = append(c.AppIds, appId)
	slices.Sort(c.AppIds)
}

func Remove(configs *[]Config, user string, appId int) {
	log.Printf("Deleting %d from %s config", appId, user)
	c := GetUserConfig(configs, user)
	var ids []int
	for _, id := range c.AppIds {
		if id != appId {
			ids = append(ids, id)
		}
	}
	c.AppIds = ids
}

func GetUserConfig(configs *[]Config, user string) *Config {
	cfgs := (*configs)
	first := &cfgs[0]
	if user == "" {
		return first
	}
	for index, _ := range *configs {
		if cfgs[index].Name == user {
			return &cfgs[index]
		}
	}
	return first
}
