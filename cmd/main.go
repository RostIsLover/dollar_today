package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"dollar_today/pkg/clients/bank"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	FIRST_MESSAGE = "Выберите валюту или введите сумму в рублях, чтобы перевести её в USD и EURO."
	START         = "/start"
)

var (
	numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bank.USD),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(bank.EURO),
		),
	)
)

func main() {
	token := mustToken()
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		var msg tgbotapi.MessageConfig
		if update.Message == nil {
			continue
		}

		switch update.Message.Text {
		case START:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, FIRST_MESSAGE)
			msg.ReplyMarkup = numericKeyboard
		case bank.USD:
			myValutes := bank.GetDailyRates()
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%.2f", myValutes[bank.USD]))
		case bank.EURO:
			myValutes := bank.GetDailyRates()
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%.2f", myValutes[bank.EURO]))
		default:
			sum, err := strconv.Atoi(update.Message.Text)
			if err != nil {
				continue
			}
			myValutes := bank.GetDailyRates()
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("USD: %.2f\nEURO: %.2f", float64(sum)/myValutes[bank.USD], float64(sum)/myValutes[bank.EURO]))
		}

		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for access to telegam bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is empty")
	}

	return *token
}
