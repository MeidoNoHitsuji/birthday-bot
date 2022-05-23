package client

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"os"
	"strings"
	"time"
)

type Client struct {
	sheets.Service
	Users []map[string]interface{}
}

const (
	layout = "02.01.2006"
)

func New() Client {
	ctx := context.Background()

	srv, err := sheets.NewService(ctx, option.WithCredentialsFile(os.Getenv("CREDENTIAL_FILE")), option.WithScopes(sheets.SpreadsheetsReadonlyScope))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	return Client{*srv, []map[string]interface{}{}}
}

func (s *Client) InitData() {
	spreadsheetId := os.Getenv("SPREADSHEET_ID")
	list := os.Getenv("READ_LIST")
	readRange := ""
	if len(list) != 0 {
		readRange += list + "!"
	}
	readRange += os.Getenv("READ_RANGE")

	resp, err := s.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		for _, row := range resp.Values {
			if len(row) == 0 {
				continue
			}
			var user = make(map[string]interface{})
			// Print columns A and E, which correspond to indices 0 and 4.
			fio := strings.Split((row[0]).(string), " ")

			user["surname"] = fio[0]
			user["name"] = fio[1]
			if len(fio) > 2 {
				user["patronymic"] = fio[2]
			}

			user["telegram"] = strings.Replace(row[6].(string), "@", "", 1)
			user["birthday"], _ = time.Parse(layout, row[7].(string))
			s.Users = append(s.Users, user)
		}
	}
}
