package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/telebot.v4"
	"gopkg.in/yaml.v3"
)

var config = "./config.yaml"
var env = "./.env"

type Config struct {
	Token  string `yaml:"token"`
	ChatId string `yaml:"chat_id"`
	DB     struct {
		Url string `yaml:"db_url"`
		Key string `yaml:"db_key"`
	} `yaml:"db"`
}

func LoadConfig() (*Config, error) {
	file, err := os.ReadFile(config)
	if err != nil {
		return nil, fmt.Errorf("cannot read yaml file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	if err := godotenv.Load(env); err != nil {
		log.Printf("Warning: could not load .env file, using default environment variables")
	}

	if token := os.Getenv("TOKEN"); token != "" {
		config.Token = token
	}

	if chatId := os.Getenv("CHAT_ID"); chatId != "" {
		config.ChatId = chatId
	}

	if url := os.Getenv("DB_URL"); url != "" {
		config.DB.Url = url
	}

	if key := os.Getenv("DB_TOKEN"); key != "" {
		config.DB.Key = key
	}

	return &config, nil
}

func BotInit(token string) (*telebot.Bot, error) {
	pref := telebot.Settings{
		Token:     token,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeHTML,
	}

	return telebot.NewBot(pref)
}
