package functions

import (
	"crypto-ticker/structure"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func FetchPrices(client *http.Client, url string) (*structure.CoinData, error) {
	resp, err := client.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Выполням GET запрос, defer гарантирует закрытие тела обьекта

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//Читаем тело ответа

	var bybit structure.BybitResponse
	if err := json.Unmarshal(body, &bybit); err != nil {
		return nil, err
	}

	if len(bybit.Result.List) == 0 {
		return nil, fmt.Errorf("пустой ответ от Bybit")
	}

	price, err := strconv.ParseFloat(bybit.Result.List[0].LastPrice, 64)
	if err != nil {
		return nil, err
	}
	return &structure.CoinData{BTC: price}, nil
}

func CheckAlerts(current, prev *structure.CoinData, config structure.AlertConfig) []string {

	var alerts []string
	//Создаем слайс для хранения сообщений об алертах

	if config.BTCAbove > 0 {

		prevWasBelow := prev.BTC < config.BTCAbove
		nowIsAbove := current.BTC >= config.BTCAbove
		if prevWasBelow && nowIsAbove {
			msg := fmt.Sprintf("BTC пересёк верхний порог $%.2f",
				config.BTCAbove)
			alerts = append(alerts, msg)
		}
		//Проверяем, был ли предыдущий курс ниже порога и стал ли текущий выше или равен порогу.
		//Если да, добавляем сообщение в алерты
	}

	if config.BTCBelow > 0 {
		prevWasAbove := prev.BTC > config.BTCBelow
		nowIsBelow := current.BTC <= config.BTCBelow
		if prevWasAbove && nowIsBelow {
			msg := fmt.Sprintf("BTC пересёк нижний порог $%.2f",
				config.BTCBelow)
			alerts = append(alerts, msg)
		}
		//Проверяем, был ли предыдущий курс выше порога и стал ли текущий ниже или равен порогу.
		//Если да, добавляем сообщение в алерты
	}

	return alerts
}
