package main

import (
	"crypto-ticker/functions"
	"crypto-ticker/structure"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	ticker := time.NewTicker(30 * time.Second)
	//Создаем канал с переодическим откликом
	defer ticker.Stop()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Println("=== Запуск скрипта. Обновление каждые 30 секунд. Нажмите Ctrl+C для выхода. ===")

	url := "https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"
	//Нужный url для получения котировок

	var prevData *structure.CoinData
	//Переменная для вывода разницы в цене, изачально пустой (*)

	btcAbove := flag.Float64("btc-above", 0, "Алерт, когда BTC выше")
	btcBelow := flag.Float64("btc-below", 0, "Алерт, когда BTC ниже")
	flag.Parse()

	config := structure.AlertConfig{
		BTCAbove: *btcAbove,
		BTCBelow: *btcBelow,
	}

	for {
		data, err := functions.FetchPrices(client, url)
		if err != nil {
			log.Println(err)
			time.Sleep(30 * time.Second)
			continue
		}
		//Обращаемся к API, ждем следующий цикл если что-то пошло не так

		fmt.Printf("%s | BTC: $%.2f\n",
			time.Now().Format("15:04:05"),
			data.BTC.USD,
		)
		//Выводим текущую цену и время

		if prevData != nil {
			alerts := functions.CheckAlerts(data, prevData, config)
			for _, alert := range alerts {
				fmt.Println(alert)
			}
			diffBTC := data.BTC.USD - prevData.BTC.USD

			fmt.Printf("BTC: %+.2f\n", diffBTC)
		}

		prevData = data
		//Обновляем данные
		<-ticker.C
		//Бескронечный цикл блокируется до следующего тика
	}
}
