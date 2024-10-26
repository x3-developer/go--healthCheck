package notifier

import (
	"HealthCheck/config"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type TelegramNotifier struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramNotifier(token string, chatID int64) (*TelegramNotifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, fmt.Errorf("ошибка при создании бота Telegram Bot: %s", err)
	}

	return &TelegramNotifier{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (n *TelegramNotifier) Notify(nt NotificationType, site config.Site, err error) {
	var tgMsg string

	switch nt {
	case Alert:
		tgMsg = `⚠️ Сайт "` + site.Name + `" упал!`

		if err != nil {
			tgMsg += "\nПричина: " + err.Error()
		}
	case Calm:
		tgMsg = `✅ Сайт "` + site.Name + `" поднялся!`
	default:
		tgMsg = `Неизвестный тип уведомления для сайта "` + site.Name + `"`
	}

	msg := tgbotapi.NewMessage(n.chatID, tgMsg)

	if _, err := n.bot.Send(msg); err != nil {
		log.Println("ошибка при отправке сообщения в Telegram:", err)
	}
}
