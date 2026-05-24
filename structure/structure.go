package structure

type Price struct {
	USD float64 `json:"usd"`
}

type CoinData struct {
	BTC Price `json:"bitcoin"`
	ETH Price `json:"ethereum"`
}

//описание структуры с JSON для котировок

type AlertConfig struct {
    BTCAbove float64
    BTCBelow float64
    ETHAbove float64
    ETHBelow float64
}

//Описание структуры для алертов
