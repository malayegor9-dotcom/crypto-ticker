package functions

import (
	"crypto-ticker/structure"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func FetchPrices(client *http.Client, url string) (*structure.CoinData, error) {
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Выполням GET запрос, defer гарантирует закрытие тела обьекта

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %s", resp.Status)
	}
	//Обрабатываем ответ сервера на ошибку

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//Читаем тело ответа

	var data structure.CoinData

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	//Отправляем тело в структуру
	return &data, nil
}

func CheckAlerts(current, prev *structure.CoinData, config structure.AlertConfig) []string {

	var alerts []string

	if config.BTCAbove > 0 {

		prevWasBelow := prev.BTC.USD < config.BTCAbove
		nowIsAbove := current.BTC.USD >= config.BTCAbove
		if prevWasBelow && nowIsAbove {
			msg := fmt.Sprintf("BTC пересёк верхний порог $%.2f (текущая $%.2f)",
				config.BTCAbove, current.BTC.USD)
			alerts = append(alerts, msg)
		}
	}

	if config.BTCBelow > 0 {
		prevWasAbove := prev.BTC.USD > config.BTCBelow
		nowIsBelow := current.BTC.USD <= config.BTCBelow
		if prevWasAbove && nowIsBelow {
			msg := fmt.Sprintf("BTC упал ниже порога $%.2f (текущая $%.2f)",
				config.BTCBelow, current.BTC.USD)
			alerts = append(alerts, msg)
		}
	}

	return alerts
}
