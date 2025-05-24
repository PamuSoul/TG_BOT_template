package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pelletier/go-toml"
)

type Config struct {
	Bot Bot `toml:"bot"`
}

type Bot struct {
	Token string `toml:"token"`
}

func main() {

	botinit := inittoml()
	setCommands(botinit) //設定指令清單

	updateConfig := tgbotapi.NewUpdate(0)           //創建機器人更新配置
	updateConfig.Timeout = 60                       //設定更新時間
	updates := botinit.GetUpdatesChan(updateConfig) //將設定的更新配置傳入機器人

	for alluserdate := range updates { //無線迴圈 監控 更新

		eventHandling(botinit, alluserdate) //判斷何種事件

	}
}

/*
// 初始化
func initialization() *tgbotapi.BotAPI {
	txt, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("讀取 token.txt 失敗: %v", err)
	}
	token := strings.TrimSpace(string(txt))
	botinit, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	return botinit
}
*/
/*
func initenv() *tgbotapi.BotAPI {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("無法載入 .env 檔案")
	}
	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("參數設定錯誤")
	}
	botinit, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	return botinit
}
	12165367573
*/
func inittoml() *tgbotapi.BotAPI {
	tree, err := toml.LoadFile("config.toml")
	if err != nil {
		log.Fatalf("讀取 config.toml 失敗: %v", err)
	}

	// 直接取得 bot.token 字串
	token := tree.Get("bot.token") // 這裡的 token 是 interface{} 型別
	if token == nil {
		log.Fatal("config.toml 裡找不到 bot.token")
	}

	botToken := token.(string) //我相信 INTERFACE{} 是字串 如不是 程式會報錯

	botinit, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	return botinit
}

// 設定 /command 指令清單
func setCommands(botinit *tgbotapi.BotAPI) {
	commands := []tgbotapi.BotCommand{
		{Command: "water", Description: "天氣"},
		{Command: "rain", Description: "下雨"},
		{Command: "bot", Description: "機器人指令"},
	}
	cmdConfig := tgbotapi.NewSetMyCommands(commands...)
	if _, err := botinit.Request(cmdConfig); err != nil {
		log.Fatalf("設定指令失敗: %v", err)
	}
}

// 判斷事件用的 目前有這幾種事件
func eventHandling(botinit *tgbotapi.BotAPI, alluserdate tgbotapi.Update) {
	switch {
	case alluserdate.CallbackQuery != nil:
		handleCallback(botinit, alluserdate.CallbackQuery) //處理按鈕事件
	case alluserdate.Message != nil:
		handleMessage(botinit, alluserdate.Message) //處理訊息事件
	case alluserdate.InlineQuery != nil:
		HandleInlineQuery(botinit, alluserdate.InlineQuery) //處理內聯查詢事件
	}
}

func handleCallback(botinit *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery) {
	// 處理按鈕事件
	var botresp string
	switch callback.Data {
	case "happy":
		botresp = "太棒了！希望你天天開心 😄"
	case "sad":
		botresp = "別難過，希望明天會更好 🌈"
	}
	botinit.Request(tgbotapi.NewCallback(callback.ID, "")) // 回應給API按鈕被點擊 如沒這指令會一直轉圈圈
	msg := tgbotapi.NewMessage(callback.Message.Chat.ID, botresp)
	botinit.Send(msg)
}

// 處理訊息事件
// 這裡可以根據不同的訊息內容來回覆不同的訊息 回答出去的位子
func handleMessage(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	username := msg.Chat.ID
	userchar := msg.Text
	log.Printf("[%s] %s", msg.From.UserName, userchar)

	resp := response(username, userchar)
	if _, err := bot.Send(resp); err != nil {
		log.Println("傳送訊息錯誤:", err)
	}
}

// 這裡可以根據不同的訊息內容來回覆不同的訊息  回應的邏輯地方
func response(username int64, userchar string) tgbotapi.MessageConfig {
	var msg tgbotapi.MessageConfig
	switch userchar {
	case "巴哈姆特":
		msg = tgbotapi.NewMessage(username, "https://www.gamer.com.tw/")
	case "/water":
		msg = tgbotapi.NewMessage(username, "今天天氣真好")
	case "/rain":
		msg = tgbotapi.NewMessage(username, "今天下雨了")
	case "/bot": // 訊息區的機器人表單按鈕設定
		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("好開心 😄", "happy"),
				tgbotapi.NewInlineKeyboardButtonData("不開心 😢", "sad"),
			),
		)
		msg = tgbotapi.NewMessage(username, "你今天心情如何？")
		msg.ReplyMarkup = keyboard
	default: // 用戶端的表單設定
		msg = tgbotapi.NewMessage(username, "請選擇")
		//鍵盤設定
		msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{
			Keyboard: [][]tgbotapi.KeyboardButton{
				{
					tgbotapi.NewKeyboardButton("巴哈姆特"),
					tgbotapi.NewKeyboardButton("water"),
				},
				{
					tgbotapi.NewKeyboardButton("rain"),
				},
			},
			ResizeKeyboard: true,
		}
	}
	return msg
}
