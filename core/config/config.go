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

func LoadConfig(path string) (*Config, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	config, _ := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var c Config
	json.Unmarshal(config, &c)
	log.Printf("Config loaded: %+v", c)
	slices.Sort(c.AppIds)

	return &c, nil
}

func (c *Config) Save(path string) error {
	configFile, err := os.Create(path)
	if err != nil {
		return err
	}

	b, err := json.Marshal(c)
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

func (c *Config) Add(appId int) {
	for _, id := range c.AppIds {
		if id == appId {
			log.Printf("Game %d already in config", appId)
			return
		}
	}
	log.Printf("Adding %d to config", appId)
	c.AppIds = append(c.AppIds, appId)
	slices.Sort(c.AppIds)
}

func (c *Config) Remove(appId int) {
	log.Printf("Deleting %d from config", appId)
	var ids []int
	for _, id := range c.AppIds {
		if id != appId {
			ids = append(ids, id)
		}
	}
	c.AppIds = ids
}
