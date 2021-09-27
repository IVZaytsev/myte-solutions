package main

import (
	"fmt"
	"log"
	"time"
	"strings"
	"net/http"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// Структура входных данных
type Ticker struct {
	Symbol         string  `json:"symbol"`
	Price24H       float64 `json:"price_24h"`
	Volume24H      float64 `json:"volume_24h"`
	LastTradePrice float64 `json:"last_trade_price"`
}

// Структура выходных данных
type TickerOut struct {
	Price24H       float64 `json:"price_24h"`
	Volume24H      float64 `json:"volume_24h"`
	LastTradePrice float64 `json:"last_trade_price"`
}

const tickersURL string = "https://api.blockchain.com/v3/exchange/tickers"

var DB *sql.DB							// Обработчик запросов к БД
var tickers = []Ticker{}				// Хранилище входных данных
var tickersOut = map[string]TickerOut{}	// Хранилище выходных данных

// Функция инициализации БД MySQL
// Инициализирует глобальную переменную <DB> для работы с БД
// Возвращает true при успешном завершении
func db_init() bool {
	var err error
	
	if DB, err = sql.Open("mysql", "docker:docker@tcp(db:3306)/test_db"); err != nil {
		log.Fatal(err)
		return false
	} else {
		DB.SetMaxOpenConns(50)
		DB.SetMaxIdleConns(10)
		err = DB.Ping()
		if err != nil {
			log.Fatal("MySQL error : ", err)
			return false
		}
		return true
	}
}

// Функция получения данных с сайта <tickersURL>
// Отправляет GET-запрос и десериализует JSON-содержимое в глобальный массив <tickers>
// Возвращает true при успешном завершении
func get_tickers() bool {
	httpClient := http.Client{
		Timeout: time.Second * 2,
	}
	
	request, err := http.NewRequest(http.MethodGet, tickersURL, nil)
	if err != nil {
		log.Fatal(err)
		return false
	}
	
	response, err := httpClient.Do(request)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&tickers); err != nil {
		log.Fatal("JSON parse error : ", err)
		return false
	}
	
	return true
}

// Функция сохранения в базу данных
// Собирает глобальный массив tickersOut
// Возвращает true при успешном завершении
func db_insert() {
	var symbolArray = []string{}
	var recordArray = []string{}
	
	// Собираем два SQL-запроса
	for _, ticker := range tickers {
		symbolStr := fmt.Sprintf("('%s')", ticker.Symbol)
		recordStr := fmt.Sprintf("((SELECT `id` FROM `symbols` WHERE `symbol`='%s'), %f, %f, %f)", ticker.Symbol, ticker.Price24H, ticker.Volume24H, ticker.LastTradePrice)
		
		symbolArray = append(symbolArray, symbolStr)
		recordArray = append(recordArray, recordStr)
	}

	// Сохраняем новые символы, если таковые появились
	var sqlQuery1 string = fmt.Sprintf("INSERT IGNORE INTO `symbols` (symbol) VALUES %s;", strings.Join(symbolArray, ","))
	if _, err := DB.Exec(sqlQuery1); err != nil {
		log.Fatal("MySQL INSERT error : ", err)
	}
	
	// Сохраняем архив значений символов с новой временной меткой
	var sqlQuery2 string = fmt.Sprintf("INSERT IGNORE INTO `records` (symbol, price_24h, volume_24h, last_trade_price) VALUES %s;", strings.Join(recordArray, ","))
	if _, err := DB.Exec(sqlQuery2); err != nil {
		log.Fatal("MySQL INSERT error : ", err)
	}
}

// Функция выборки записей из базы данных
// Собирает глобальный массив tickersOut
// Возвращает true при успешном завершении
func db_select() bool {
	sqlQuery3 := "SELECT symbols.symbol, price_24h, volume_24h, last_trade_price FROM `records` INNER JOIN `symbols` ON records.symbol = symbols.id WHERE `stamp`=(SELECT MAX(stamp) FROM records);"
	rows, err := DB.Query(sqlQuery3)
	if err != nil {
		log.Fatal("MySQL SELECT error : ", err)
		return false
	}

	tickers := make([]*Ticker, 0)
	for rows.Next() {
		ticker := new(Ticker)
		err := rows.Scan(&ticker.Symbol, &ticker.Price24H, &ticker.Volume24H, &ticker.LastTradePrice)
		if err != nil {
			log.Fatal("MySQL result scan error : ", err)
			return false
		}
		tickers = append(tickers, ticker)
	}

	for _, ticker := range tickers {
		tickersOut[ticker.Symbol] = TickerOut{
			Price24H: ticker.Price24H,
			Volume24H: ticker.Volume24H,
			LastTradePrice: ticker.LastTradePrice,
		}
	}
	
	return true
}

// Горутина сбора и сохранения данных с сайта <tickersURL>
// Бесконечный цикл с паузой в 30 секунд
func start_scraping() {
	for ok := true; ok; ok = true {
		if get_tickers() {
			db_insert()
		}
		time.Sleep(30 * time.Second)
	}
}

// Функция запуска http-сервера и настройки функции обработки входящих запросов
func start_http_server() {
	fmt.Println("Starting server at port 8080")
	http.HandleFunc("/", mainHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
		return
	}
}

// Обработчик входящих запросов http-сервера. Принимает только GET-запросы в корень
func mainHandler(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.Error(writer, http.StatusText(404), http.StatusNotFound)
		return
	}
	
	if request.Method != "GET" {
		http.Error(writer, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	
	if db_select() {
		jsonText, err := json.Marshal(tickersOut)
		if err != nil {
			http.Error(writer, http.StatusText(500), http.StatusInternalServerError)
			return
		}
		
		fmt.Fprintf(writer, string(jsonText))
	}
}

func main() {
	// Если подключение к БД успешно
	// Создаём 2 горутины:
	// go start_scraping()			- сбор данных
	// go start_http_server()		- http-сервер для отображения
	if db_init() {
		defer DB.Close()
		go start_scraping()
		go start_http_server()
	}
	
	// Заглушка для удобного завершения работы приложения
	time.Sleep(5 * time.Second)
	fmt.Print("Press 'Enter' to exit...")
	fmt.Scanln()
}