package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Quote struct {
	ID         int    `gorm:"primaryKey" json:"id"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"createDate"`
}

type QuoteResponse map[string]Quote

func main() {
	http.HandleFunc("/cotacao", FindMoneyQuoteHandler)
	http.ListenAndServe(":8080", nil)
}

func FindMoneyQuoteHandler(responseWriter http.ResponseWriter, request *http.Request) {
	db, err := gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Quote{})

	currency := request.URL.Query().Get("currency")
	quote, err := SearchMoneyQuote(currency)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = InsertQuote(db, quote)

	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(quote)
}

func SearchMoneyQuote(currency string) (*Quote, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/"+currency, nil)

	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var quoteResponse QuoteResponse
	err = json.Unmarshal(body, &quoteResponse)

	if err != nil {
		return nil, err
	}

	currency = strings.Replace(currency, "-", "", -1)
	quote := quoteResponse[currency]

	return &quote, nil
}

func InsertQuote(db *gorm.DB, quote *Quote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	return db.WithContext(ctx).Create(&quote).Error
}
