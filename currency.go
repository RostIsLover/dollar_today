package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/html/charset"
)

const (
	DOLLAR        = "Доллар США"
	EURO          = "Евро"
	START         = "/start"
	FIRST_MESSAGE = "Выберите валюту"
	BOT_API_TOKEN = "6096941081:AAGe5wWH9HwDrHmPVEo5Jo8l8h9pjuEdyDk"
)

var (
	numericKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(DOLLAR),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(EURO),
		),
	)
	myValutes = make(map[string]float64)
)

func main() {
	currentTime := time.Now().Format("02/01/2006")
	url := "http://www.cbr.ru/scripts/XML_daily.asp?date_req=" + currentTime
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "__ddg1_=bMz7QAI3fDT4y8GS26rJ")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "BatPhone/7.26.8")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := string(body)
	valCurs := new(CBRValCurs)
	r := bytes.NewReader([]byte(data))
	d := xml.NewDecoder(r)
	d.CharsetReader = charset.NewReaderLabel
	err = d.Decode(&valCurs)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}

	addToMyValutes(valCurs.Val, DOLLAR, EURO)

	bot, err := tgbotapi.NewBotAPI(BOT_API_TOKEN)
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
		case DOLLAR:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%f", myValutes[DOLLAR]))
		case EURO:
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("%f", myValutes[EURO]))

		}
		if _, err := bot.Send(msg); err != nil {
			panic(err)
		}
	}
}

func addToMyValutes(vals []Valute, names ...string) {
	for _, val := range vals {
		for _, name := range names {
			if val.Name == name {
				f := val.Value
				f = strings.Replace(f, ",", ".", -1)
				s, err := strconv.ParseFloat(f, 64)
				if err != nil {
					fmt.Println(err.Error())
				}
				myValutes[name] = s
			}
		}
	}

}

type CBRValCurs struct {
	XMLName xml.Name `xml:"ValCurs"`
	Val     []Valute `xml:"Valute"`
}

type Valute struct {
	NumCode  string `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}
