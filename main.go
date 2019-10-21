package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	bot, err := linebot.New(config.Config.Line.Channelsecret,
		config.Config.Line.Channeltoken,
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
	// msg := "95.5"
	insertToSpreadsheet(msg)
	if _, err := bot.ReplyMessage(events[0].ReplyToken, linebot.NewTextMessage("เซฟให้เเล้วนะ :3")).Do(); err != nil {
		fmt.Println(err)
	}
	w.WriteHeader(200)
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

	// convert msg to int
	msgF, err := strconv.ParseFloat(msg, 64)
	if err != nil {
		fmt.Println("cannot convert msg", err)
		return
	}

	// login to google spread sheet
	srv := loginGoogle()
	sheetRange := "sheet1" // this mean an entire sheet

	// get current time
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("cannot load location: ", err)
		return
	}
	date := time.Now().In(loc)
	dateS := date.Format("02-01-2006")

	// get last row to verify, add month summary if it is new month
	get, err := srv.Spreadsheets.Values.Get(config.Config.SpreadsheetID, sheetRange).Do()
	if err != nil {
		fmt.Println(err)
	}

	if len(get.Values) == 0 {
		fmt.Println("value out of range")
		return
	}

	lastDate, ok := (get.Values[len(get.Values)-1][0].(string))
	if !ok {
		fmt.Println("not ok")
		return
	}
	pOldDate, err := time.Parse("02-01-2006", lastDate)
	if err != nil {
		fmt.Println("parse date error: ", err)
		return
	}

	if pOldDate.Month() != date.Month() {
		// add value got from line and create new month
		b := []interface{}{}
		m := []interface{}{date.Month().String()}
		v := []interface{}{dateS, msgF}
		rb := &sheets.ValueRange{
			MajorDimension: "ROWS",
		}
		rb.Values = append(rb.Values, b)
		rb.Values = append(rb.Values, m)
		rb.Values = append(rb.Values, v)

		_, err := srv.Spreadsheets.Values.Append(config.Config.SpreadsheetID, sheetRange, rb).InsertDataOption("INSERT_ROWS").ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			fmt.Println("cannot append: ", err)
			return
		}
		return
	}

	// add value got from line
	v := []interface{}{dateS, msgF}
	rb := &sheets.ValueRange{
		MajorDimension: "ROWS",
	}
	rb.Values = append(rb.Values, v)

	_, err = srv.Spreadsheets.Values.Append(config.Config.SpreadsheetID, sheetRange, rb).InsertDataOption("INSERT_ROWS").ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		fmt.Println(err)
	}

	return
}

func loginGoogle() *sheets.Service {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("./creds.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return srv
}
