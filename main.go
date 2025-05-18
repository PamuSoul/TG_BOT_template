package main

import (
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {

	txt, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("讀取 token.txt 失敗: %v", err)
	}
	token := strings.TrimSpace(string(txt))
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	//設定 Command 區域
	//Commadnd key 使用者前面需要加上 "/"  Description 是使用者看到的指令
	commands := []tgbotapi.BotCommand{
		{Command: "water", Description: "天氣"},
		{Command: "rain", Description: "下雨"},
		{Command: "bot", Description: "機器人指令"},
	}
	setCmdCfg := tgbotapi.NewSetMyCommands(commands...)

	if _, err := bot.Request(setCmdCfg); err != nil {
		log.Fatalf("設定指令失敗: %v", err)
	}

	updateConfig := tgbotapi.NewUpdate(0) //設定初始化
	updateConfig.Timeout = 60             //更新時間間隔

	updates := bot.GetUpdatesChan(updateConfig) //對更新的監控

	for alluserdate := range updates { //處理數據

		//按鈕事件
		//如果不先處理按鈕事件 會導致按鈕事件無法正常運行
		if alluserdate.CallbackQuery != nil {
			callback := alluserdate.CallbackQuery
			var botresp string

			switch callback.Data {
			case "happy":
				botresp = "太棒了！希望你天天開心 😄"
			case "sad":
				botresp = "別難過，希望明天會更好 🌈"
			}

			bot.Request(tgbotapi.NewCallback(callback.ID, ""))
			msg := tgbotapi.NewMessage(callback.Message.Chat.ID, botresp)
			bot.Send(msg)
			continue
		}

		//訊息事件
		if alluserdate.Message != nil {
			log.Printf("[%s] %s", alluserdate.Message.From.UserName, alluserdate.Message.Text)

			userchar := alluserdate.Message.Text
			username := alluserdate.Message.Chat.ID

			if userchar == "/bot" {
				keyboard := tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("好開心 😄", "happy"),
						tgbotapi.NewInlineKeyboardButtonData("不開心 😢", "sad"),
					),
				)
				msg := tgbotapi.NewMessage(username, "你今天心情如何？")
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
				continue
			}

			msg := response(username, userchar)
			if _, err := bot.Send(msg); err != nil {
				log.Println(err)
			}

		}
	}
}

func response(username int64, userchar string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	switch userchar {
	case "巴哈姆特":
		msg = tgbotapi.NewMessage(username, "https://www.gamer.com.tw/")
	case "/water":
		msg = tgbotapi.NewMessage(username, "今天天氣真好")
	case "/rain":
		msg = tgbotapi.NewMessage(username, "今天下雨了")
	default:
		msg = tgbotapi.NewMessage(username, "請選擇")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("巴哈姆特"),
				tgbotapi.NewKeyboardButton("/water"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("/rain"),
			),
		)
	}
	return msg
}
