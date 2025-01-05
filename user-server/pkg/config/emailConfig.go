package config

import (
	"gopkg.in/yaml.v3"
	"micro-services/pkg/utils"
	"os"
	"path/filepath"
)

type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Sender   string `yaml:"sender"`
	Password string `yaml:"password"`
	Protocol string `yaml:"protocol"`
}

type TitleEmail struct {
	Email EmailConfig `yaml:"internal"`
}

var EmailSender *TitleEmail

func InitEmailConfig() error {
	rootPath := utils.GetCurrentPath(2)
	configPath := filepath.Join(rootPath, "../pkg/config", "config.yml")
	//fmt.Println("path:", configPath)
	EmailSender = &TitleEmail{}
	data, err := os.ReadFile(configPath)
	if err != nil {
		//fmt.Println("errorFromConfig.go: ", err)
		return err
	}
	err = yaml.Unmarshal(data, EmailSender)
	if err != nil {
		return err
	}
	return nil
}
