package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleInlineQuery(bot *tgbotapi.BotAPI, query *tgbotapi.InlineQuery) {
	log.Println("收到 Inline Query:", query.Query) //query.Query 是輸入的文字 判斷使用輸入的訊息 底下可以用switch or if 做回動作

	var results []interface{}

	if query.Query == "影片" {

		// 第一個選項：Markdown 格式影片連結
		text1 := "[FFXIV GMV](https://www.youtube.com/watch?v=L--mZCu5RWU)"
		article1 := tgbotapi.NewInlineQueryResultArticleMarkdown( //tgbotapi.NewInlineQueryResultArticleMarkdown
			"1", ////這裡可以是 html or Markdown 模式
			"影片連結",
			text1,
		)
		article1.Description = "FFXIV GMV" // 這裡可以設定說明文字
		results = append(results, article1)

		// 第二個選項：HTML 格式新聞
		text2 := "[駆け引き](https://www.youtube.com/watch?v=24k8yJQ4W40&list=RD24k8yJQ4W40&start_radio=1&rv=L--mZCu5RWU)"
		article2 := tgbotapi.NewInlineQueryResultArticleHTML(
			"2",
			"駆け引き",
			text2,
		)
		article2.Description = "駆け引き的歌曲"
		results = append(results, article2)

		// 第三個選項：文字範例
		text3 := "隨便打打字不知道要放甚麼文字"
		article3 := tgbotapi.NewInlineQueryResultArticleHTML( //HTML 模式範例
			"3",
			"隨便打打字",
			text3,
		)
		article3.Description = "文字範例"
		results = append(results, article3)

		inlineConf := tgbotapi.InlineConfig{
			InlineQueryID: query.ID, // 這裡是回傳的 ID
			IsPersonal:    true,     // 是否是個人查詢 false 為所有人都能看到查詢結果
			CacheTime:     0,        // 緩存時間
			Results:       results,  // 這裡是回傳的結果
		}

		if _, err := bot.Request(inlineConf); err != nil {
			log.Println("回應 Inline Query 失敗:", err)
		}
	}
}
