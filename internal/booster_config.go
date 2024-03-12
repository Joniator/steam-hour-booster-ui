package internal

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"slices"
)

type UserConfig struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	AppIds   []int  `json:"games"`
}

type BoosterConfig struct {
	UserConfigs []UserConfig
	configPath  string
}

func LoadBoosterConfig(path string) (*BoosterConfig, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	configRaw, _ := io.ReadAll(configFile)
	if err != nil {
		return nil, err
	}

	var config = BoosterConfig{configPath: path}

	err = json.Unmarshal(configRaw, &config.UserConfigs)
	if err != nil {
		return nil, err
	}
	log.Printf("Config loaded: %+v", config)
	for _, c := range config.UserConfigs {
		slices.Sort(c.AppIds)
	}

	return &config, nil
}

func (boosterConfig *BoosterConfig) Save() error {
	configFile, err := os.Create(boosterConfig.configPath)
	if err != nil {
		return err
	}

	b, err := json.MarshalIndent(boosterConfig.UserConfigs, "", "    ")
	if err != nil {
		return err
	}

	_, err = configFile.Write(b)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf("Config saved: %+v", boosterConfig)
	return nil
}

func (boosterConfig *BoosterConfig) GetDefaultUser() string {
	return boosterConfig.UserConfigs[0].Name
}

func (boosterConfig *BoosterConfig) UserExists(user string) bool {
	_, err := boosterConfig.GetUserConfig(user)
	return err == nil
}

func (boosterConfig *BoosterConfig) GetUserConfig(user string) (*UserConfig, error) {
	userConfigs := boosterConfig.UserConfigs
	for index, _ := range userConfigs {
		if userConfigs[index].Name == user {
			return &userConfigs[index], nil
		}
	}
	return nil, errors.New("User not found")
}

func (boosterConfig *BoosterConfig) AddGame(user string, appId int) error {
	userConfig, err := boosterConfig.GetUserConfig(user)
	if err != nil {
		return err
	}
	for _, id := range userConfig.AppIds {
		if id == appId {
			log.Printf("Game %d already in config", appId)
			return nil
		}
	}
	log.Printf("Adding %d to %s config", appId, user)
	userConfig.AppIds = append(userConfig.AppIds, appId)
	slices.Sort(userConfig.AppIds)
	return nil
}

func (boosterConfig *BoosterConfig) RemoveGame(user string, appId int) error {
	log.Printf("Deleting %d from %s config", appId, user)
	userConfig, err := boosterConfig.GetUserConfig(user)
	if err != nil {
		return err
	}
	var ids []int
	for _, id := range userConfig.AppIds {
		if id != appId {
			ids = append(ids, id)
		}
	}
	userConfig.AppIds = ids
	return nil
}
