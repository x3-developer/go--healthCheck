package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Site struct {
	Name     string `json:"name"`
	Url      string `json:"url"`
	IsActive bool   `json:"isActive"`
}

type Config struct {
	Sites          []Site `json:"sites"`
	CheckInterval  int
	TelegramToken  string
	TelegramChatID int64
}

func LoadConfig() *Config {
	checkInterval := getEnv("CHECK_INTERVAL", "300")
	interval, intervalErr := strconv.Atoi(checkInterval)

	if intervalErr != nil {
		log.Printf("ошибка парсинга CHECK_INTERVAL: %v", intervalErr)
		interval = 300
	}

	sites, err := getSites()

	if err != nil {
		log.Printf("ошибка получения сайтов: %v", err)
	}

	return &Config{
		Sites:          sites,
		CheckInterval:  interval,
		TelegramToken:  getEnv("TELEGRAM_TOKEN", ""),
		TelegramChatID: getChatID(getEnv("TELEGRAM_CHAT_ID", "0")),
	}
}

func getSites() ([]Site, error) {
	var sites []Site
	sitesFile := getEnv("SITES_FILE", "")

	file, fileErr := os.Open(sitesFile)

	if fileErr != nil {
		return sites, fmt.Errorf("ошибка открытия файла: %s, %v", sitesFile, fileErr)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("ошибка закрытия файла: %s", sitesFile)
		}
	}()

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&sites); err != nil {
		return sites, fmt.Errorf("ошибка декодирования файла: %s, %v", sitesFile, err)
	}

	return sites, nil
}

func SaveSites(cfg *Config) error {
	sitesFile := getEnv("SITES_FILE", "")
	file, fileErr := os.Create(sitesFile)
	if fileErr != nil {
		return fmt.Errorf("не удалось открыть файл на запись: %s, %v", sitesFile, fileErr)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("ошибка закрытия файла: %s", sitesFile)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg.Sites); err != nil {
		return err
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}

	return defaultValue
}

func getChatID(chatIDStr string) int64 {
	chatID, err := strconv.ParseInt(chatIDStr, 10, 64)

	if err != nil {
		log.Printf("ошибка парсинга CHAT_ID: %v", err)
	}

	return chatID
}
