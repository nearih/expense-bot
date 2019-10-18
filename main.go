package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"

	"expense-bot/config"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func Test(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("success"))
}

func main() {
	http.HandleFunc("/", Test)
	http.HandleFunc("/bot", ExpenseBot)
	fmt.Println("server is ready")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ExpenseBot(w http.ResponseWriter, r *http.Request) {
	bot, err := linebot.New(config.Config.Channelsecret,
		config.Config.Channeltoken,
	)
	events, err := bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	msg := getMessage(events)
	insertToSpreadsheet(msg)
	if _, err := bot.ReplyMessage(events[0].ReplyToken, linebot.NewTextMessage("เซฟให้เเล้วนะ :3")).Do(); err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(200)
}

func loginGoogle() *sheets.Service {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("./creds.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv
}

func getMessage(events []*linebot.Event) string {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				fmt.Println("message.Text", message.Text)
				// if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
				// 	fmt.Println(err)
				// }
				return message.Text
			}
		}
	}
	return ""
}

func insertToSpreadsheet(msg string) {
	srv := loginGoogle()
	insertRange := "sheet1"

	date := time.Now().Format("2006-01-02")

	v := []interface{}{date, msg}
	rb := &sheets.ValueRange{
		MajorDimension: "ROWS",
	}
	rb.Values = append(rb.Values, v)

	res, err := srv.Spreadsheets.Values.Append(config.Config.SpreadsheetID, insertRange, rb).InsertDataOption("INSERT_ROWS").ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("resp %+#v\n", res)
}
