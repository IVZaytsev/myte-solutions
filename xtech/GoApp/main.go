package main

import (
	"fmt"
	"log"
	"time"
	"strings"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type Ticker struct {
    Symbol	string	`json:"symbol"`
    Price	float64	`json:"price_24h"`
    Volume	float64	`json:"volume_24h"`
	LastTrade	float64	`json:"last_trade_price"`
}

type TickerOut struct {
    Price	float64	`json:"price_24h"`
    Volume	float64	`json:"volume_24h"`
	LastTrade	float64	`json:"last_trade_price"`
}

var URL_Tickers = "https://api.blockchain.com/v3/exchange/tickers"
var Tickers = []Ticker{}

func get_tickers() {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	
	req, err := http.NewRequest(http.MethodGet, URL_Tickers, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	res, err := httpClient.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	jsonErr := json.Unmarshal(body, &Tickers)
	if jsonErr != nil {
		log.Fatal(jsonErr)
		return
	}
	
	db_record()
}

func db_record() {
	DB, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/test_db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer DB.Close()


	symArray := []string{}
	recArray := []string{}

	for _, ticker := range Tickers {
		synStr := fmt.Sprintf("('%s')", ticker.Symbol)
		recStr := fmt.Sprintf("((SELECT `id` FROM `symbols` WHERE `symbol`='%s'), %f, %f, %f)", ticker.Symbol, ticker.Price, ticker.Volume, ticker.LastTrade)
		
		symArray = append(symArray, synStr)
		recArray = append(recArray, recStr)
	}
	
	sqlQuery1 := fmt.Sprintf("INSERT IGNORE INTO `symbols` (symbol) VALUES %s;", strings.Join(symArray, ","))
	_, err = DB.Exec(sqlQuery1)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	sqlQuery2 := fmt.Sprintf("INSERT IGNORE INTO `records` (symbol, price_24h, volume_24h, last_trade_price) VALUES %s;", strings.Join(recArray, ","))
	_, err = DB.Exec(sqlQuery2)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func runHTTPServer() {
	fmt.Printf("Starting server at port 8080\n")
	http.HandleFunc("/", mainHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
		return
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		return
	}
	
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	DB, err := sql.Open("mysql", "docker:docker@tcp(db:3306)/test_db")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer DB.Close()

	sqlQuery3 := "SELECT symbols.symbol, price_24h, volume_24h, last_trade_price FROM `records` INNER JOIN `symbols` ON records.symbol = symbols.id WHERE `stamp`=(SELECT MAX(stamp) FROM records);"
	rows, err := DB.Query(sqlQuery3)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		log.Fatal(err)
		return
	}
	defer rows.Close()

	tickers := make([]*Ticker, 0)
	for rows.Next() {
		ticker := new(Ticker)
		err := rows.Scan(&ticker.Symbol, &ticker.Price, &ticker.Volume, &ticker.LastTrade)
		if err != nil {
			http.Error(w, http.StatusText(500), http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
		tickers = append(tickers, ticker)
	}

	mapTickers := make(map[string]TickerOut)
	for _, ticker := range tickers {
		mapTickers[ticker.Symbol] = TickerOut{
			Price: ticker.Price,
			Volume: ticker.Volume,
			LastTrade: ticker.LastTrade,
		}
	}
	
	jsonText, err := json.Marshal(mapTickers)
	if err != nil {
		log.Fatal(err)
		return
	}	
	
	fmt.Fprintf(w, string(jsonText))
}

func myTask(){
	for ok := true; ok; ok = true {
		get_tickers()
		time.Sleep(30 * time.Second)
	}
}

func main() {
	fmt.Println("Program starting...")
	time.Sleep(5 * time.Second)
	go myTask()
	runHTTPServer()
}