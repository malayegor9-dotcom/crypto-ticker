package structure

type BybitTicker struct {
	LastPrice string `json:"lastPrice"`
}

type BybitResult struct {
	List []BybitTicker `json:"list"`
}

type BybitResponse struct {
	Result BybitResult `json:"result"`
}

type CoinData struct {
	BTC float64
}

type AlertConfig struct {
	BTCAbove float64
	BTCBelow float64
}

//Описание структуры для алертов
