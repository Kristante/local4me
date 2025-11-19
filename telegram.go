package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
	"time"
)

func CreateBot(osToken string) (*tele.Bot, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла или переменных окружения: %v", err)
		return nil, err
	}

	token := os.Getenv(osToken)
	if token == "" {
		log.Fatalf("%v не найден в файле .env или переменных окружения", osToken)
		return nil, err
	}

	settings := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(settings)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return bot, nil
}

func SendMessageForChat(bot *tele.Bot, chatID int64, message string) error {
	_, err := bot.Send(&tele.Chat{ID: chatID}, message)
	if err != nil {
		return errors.New(fmt.Sprintf("не смог отправить ошибку в телеграмм. Ошибка: %v", err))
	}
	return nil
}
