package main

import (
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

// TestHandler is for connetion testing with rest-api
func TestHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("success"))
}

func main() {
	http.HandleFunc("/", TestHandler)
	http.HandleFunc("/bot", ExpenseBot)
	log.Println("server is ready")
	port := config.Config.Port
	if (port==0) {
		port = 8080	
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v",port), nil))
}

// ExpenseBot save line bot message to google spreadsheet
func ExpenseBot(w http.ResponseWriter, r *http.Request) {
	// new linebot client
	bot, err := linebot.New(config.Config.Line.Channelsecret,
		config.Config.Line.Channeltoken,
	)

	// get event from line request
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

	// reply back to line
	if _, err := bot.ReplyMessage(events[0].ReplyToken, linebot.NewTextMessage("เซฟให้เเล้วนะ :3")).Do(); err != nil {
		log.Fatal(err)
	}
	w.WriteHeader(200)
}

// getMessage extracts msg from line event
func getMessage(events []*linebot.Event) string {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				log.Println("message.Text", message.Text)
				return message.Text
			}
		}
	}
	return ""
}

// insertToSpreadsheet insert message to spreadsheet
func insertToSpreadsheet(msg string) {

	// convert msg to int
	msgF, err := strconv.ParseFloat(msg, 64)
	if err != nil {
		log.Println("cannot convert msg", err)
		return
	}

	// login to google spread sheet
	srv := loginGoogle()
	sheetRange := config.Config.SheetRange // this is sheetName or tab name eg: sheet1
	//you can specific column/row if not thing specify it mean an entire sheet

	// get current time
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Println("cannot load location: ", err)
		return
	}
	date := time.Now().In(loc)
	dateS := date.Format("02/01/2006")

	// get last row to verify, add month name if it is new month
	get, err := srv.Spreadsheets.Values.Get(config.Config.SpreadsheetID, sheetRange).Do()
	if err != nil {
		log.Println("srv.Spreadsheets.Values.Get: ", err)
		return
	}

	if len(get.Values) == 0 {
		log.Println("value out of range, data now found")
		return
	}

	lastDate, ok := (get.Values[len(get.Values)-1][0].(string))
	if !ok {
		log.Println("cannot extract data from map")
		return
	}

	// pLastDate = parsed last date
	pLastDate, err := time.Parse("02/01/2006", lastDate)
	if err != nil {
		log.Println("parse date error: ", err)
		return
	}

	if pLastDate.Month() != date.Month() {

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

		_, err := srv.Spreadsheets.Values.Append(config.Config.SpreadsheetID, sheetRange+"!A1:A1000", rb).InsertDataOption("INSERT_ROWS").ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			log.Println("cannot append: ", err)
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
		log.Println("append value error: ", err)
		return
	}

	return
}

// loginGoogle log in to google sheet sdk with credentialfile
func loginGoogle() *sheets.Service {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("./creds.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return nil
	}
	return srv
}
