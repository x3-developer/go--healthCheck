package main

import (
	"HealthCheck/config"
	"HealthCheck/internal/checker"
	"HealthCheck/internal/notifier"
	"github.com/joho/godotenv"
	"log"
	"sync"
	"time"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.LoadConfig()
	httpChecker := checker.NewHTTPChecker(30 * time.Second)
	tgNotifier, tgError := notifier.NewTelegramNotifier(cfg.TelegramToken, cfg.TelegramChatID)

	if tgError != nil {
		log.Fatal(tgNotifier)
	}

	for {
		var wg sync.WaitGroup

		for i, site := range cfg.Sites {
			wg.Add(1)

			go func(i int, s config.Site) {
				defer wg.Done()

				isAvailable, err := httpChecker.Check(s.Url)
				if err != nil {
					log.Printf("ошибка проверки сайта %s: %v", s.Name, err)
					tgNotifier.Notify(notifier.Alert, s, err)
					cfg.Sites[i].IsActive = false
					return
				}

				if s.IsActive && !isAvailable {
					cfg.Sites[i].IsActive = false
					tgNotifier.Notify(notifier.Alert, s, nil)
				} else if !s.IsActive && isAvailable {
					cfg.Sites[i].IsActive = true
					tgNotifier.Notify(notifier.Calm, s, nil)
				}

			}(i, site)
		}

		wg.Wait()

		if saveSitesErr := config.SaveSites(cfg); saveSitesErr != nil {
			log.Printf("Ошибка сохранения файла сайтов: %v", saveSitesErr)
		}

		time.Sleep(time.Duration(cfg.CheckInterval) * time.Second)
	}
}
