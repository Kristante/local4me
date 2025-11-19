package main

import (
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v4"
	"log"
	"os"
	"time"
)

const chatID = 1062210573
const BusinessAccount = "Служебная УЗ"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	bot, err := CreateBot("BOT_TOKEN")
	if err != nil {
		log.Fatalf("Произошла ошибка %v", err)
	}

	apiToken := os.Getenv("TOKEN_4ME")
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	pool := CreateDatabasePool()

	fmt.Println("Запуск планировщика...")
	go timing(ticker, bot, apiToken, pool)

	bot.Start()

}

func timing(ticker *time.Ticker, bot *tele.Bot, apiToken string, pool *pgxpool.Pool) {
	for {
		select {
		case <-ticker.C:
			// Выполняем функцию каждые 30 минут
			fmt.Println("Таймер сработал, выполняется получение заявок...")
			err, requests := getAllRequests(apiToken)
			if err != nil {
				log.Fatalf("Произошла ошибка при получении списка всех заявок: %v", err)
			}
			for _, request := range requests {
				if CheckMemberName(request) {
					if SelectRequest(pool, int64(request.ID)) == 0 {
						err, request := getInfoForRequest(request.ID, apiToken)
						if err != nil {
							log.Fatalf("Произошла ошибка при получении информации о заявке %v. Текст ошибки: %v", request.ID, err)
						}

						err = SendMessageForChat(bot, chatID, ConvertInfoForMessageTelegram(*request))
						if err != nil {
							log.Fatalf("Произошла ошибка при отправке телеграмм: %v", err)
						}
						err = AddRequest(pool, int64(request.ID))
						if err != nil {
							log.Fatalf("Произошла ошибка при взаимодействии с базой данных: %v", err)
						}
					}
				} else {
					if request.Status == "assigned" {
						UpdateTime := request.UpdatedAt.Format("2006-01-02T15:04:05-07:00")
						if !CheckNotes(pool, int64(request.ID), UpdateTime) && request.Member.Name != BusinessAccount {
							err, notes := GetNotesForRequest(request.ID, apiToken)
							if err != nil {
								log.Fatalf("Произошла ошибка при получении заметок в заявке: %v", err)
							}
							comments := GetComments(notes)
							if comments != nil {
								err = SendMessageForChat(bot, chatID, ConvertNotesForMessageTelegram(*comments, request.ID))
								if err != nil {
									log.Fatalf("Произошла ошибка при отправке телеграмм: %v", err)
								}
								err = AddNote(pool, int64(request.ID), UpdateTime)
								if err != nil {
									log.Fatalf("Произошла ошибка при взаимодействии с базой данных: %v", err)
								}
							}
						}
					}
				}
			}
		}
	}
}
