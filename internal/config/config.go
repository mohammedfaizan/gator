package config

import (
	"encoding/json"
	"log"
	"os"
)






type Config struct {
	DbUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	const configFileName = ".gatorconfig.json"
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fileLoc := homeDir + "/" + configFileName
	return fileLoc, nil
}

func (cfg *Config) Read() (error) {

	

	fileLoc, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	

	file, err := os.ReadFile(fileLoc)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		return err
	}


	return nil 
	
}

//Export a SetUser method on the Config struct that writes the config struct to the JSON file after setting the current_user_name field.

func  (cfg *Config) SetUser(name, db_url string) error {
	
	

	fileLoc, err := getConfigFilePath()
	if  err != nil {
		return err
	}
	
	cfg.CurrentUserName = name
	cfg.DbUrl = db_url

	err = cfg.write(fileLoc)
	if err != nil {
		return err
	}
	
	log.Println("json data successfully written to config.json")

	return nil
}

func (cfg *Config) write(fileLoc string) error {
	jsonData, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileLoc, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

