package telegram

import (
	"crypto-ticker/functions"
	"crypto-ticker/structure"
	"fmt"
	"log"
	"net/http"
	"time"

	tele "gopkg.in/telebot.v3"
)

func Start(token string) {
	bot, err := tele.NewBot(tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	url := "https://api.bybit.com/v5/market/tickers?category=spot&symbol=BTCUSDT"

	bot.Handle("/start", func(c tele.Context) error {
		go runTicker(bot, c.Chat(), client, url, structure.AlertConfig{})
		return c.Send("Тикер запущен! Обновление каждые 30 секунд.")
	})

	bot.Start()
}

func runTicker(bot *tele.Bot, chat *tele.Chat, client *http.Client, url string, config structure.AlertConfig) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	var prevData *structure.CoinData

	for {
		data, err := functions.FetchPrices(client, url)
		if err != nil {
			log.Println(err)
			<-ticker.C
			continue
		}

		msg := formatMessage(data, prevData, config)
		bot.Send(chat, msg)

		prevData = data
		<-ticker.C
	}
}

func formatMessage(data, prev *structure.CoinData, config structure.AlertConfig) string {
	msg := fmt.Sprintf("🕐 %s | BTC: $%.2f\n",
		time.Now().Format("15:04:05"),
		data.BTC,
	)

	if prev != nil {
		alerts := functions.CheckAlerts(data, prev, config)
		for _, alert := range alerts {
			msg += fmt.Sprintf("⚠️ %s\n", alert)
		}

		diffBTC := data.BTC - prev.BTC

		arrow := "▲"
		if diffBTC < 0 {
			arrow = "▼"
		} else if diffBTC > -0.01 && diffBTC < 0.01 {
			arrow = "→"
		}

		msg += fmt.Sprintf("Изменение: %s %+.2f\n", arrow, diffBTC)
	}

	return msg
}
