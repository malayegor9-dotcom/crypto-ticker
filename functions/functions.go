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
	//Создаем слайс для хранения сообщений об алертах

	if config.BTCAbove > 0 {

		prevWasBelow := prev.BTC.USD < config.BTCAbove
		nowIsAbove := current.BTC.USD >= config.BTCAbove
		if prevWasBelow && nowIsAbove {
			msg := fmt.Sprintf("BTC пересёк верхний порог $%.2f)",
				config.BTCAbove)
			alerts = append(alerts, msg)
		}
		//Проверяем, был ли предыдущий курс ниже порога и стал ли текущий выше или равен порогу.
		//Если да, добавляем сообщение в алерты
	}

	if config.BTCBelow > 0 {
		prevWasAbove := prev.BTC.USD > config.BTCBelow
		nowIsBelow := current.BTC.USD <= config.BTCBelow
		if prevWasAbove && nowIsBelow {
			msg := fmt.Sprintf("BTC пересёк нижний порог $%.2f)",
				config.BTCBelow)
			alerts = append(alerts, msg)
		}
		//Проверяем, был ли предыдущий курс выше порога и стал ли текущий ниже или равен порогу.
		//Если да, добавляем сообщение в алерты
	}

	return alerts
}
