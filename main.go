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

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

func main() {
	ticker := time.NewTicker(30 * time.Second)
	//Создаем канал с переодическим откликом
	defer ticker.Stop()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("%s=== Запуск скрипта. Обновление каждые 30 секунд. Нажмите Ctrl+C для выхода. ===%s\n", colorCyan, colorReset)
	fmt.Println()

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

		fmt.Printf("%s%s%s | BTC: %s$%.2f%s\n",
			colorCyan, time.Now().Format("15:04:05"), colorReset,
			colorBold, data.BTC.USD, colorReset,
		)
		//Выводим текущую цену и время

		if prevData != nil {
			alerts := functions.CheckAlerts(data, prevData, config)
			for _, alert := range alerts {
				fmt.Printf("%s%s%s\n", colorYellow, alert, colorReset)
			}
			diffBTC := data.BTC.USD - prevData.BTC.USD

			arrow := "▲"
			diffColor := colorGreen

			if diffBTC < 0 {
				arrow = "▼"
				diffColor = colorRed
			} else if diffBTC > -0.01 && diffBTC < 0.01 {
				arrow = "→"
				diffColor = colorCyan
			}
			fmt.Printf("Изменение: %s%s %+.2f%s\n", diffColor, arrow, diffBTC, colorReset)
		}

		fmt.Println("─────────────────────────")

		prevData = data
		//Обновляем данные
		<-ticker.C
		//Бескронечный цикл блокируется до следующего тика
	}
}
