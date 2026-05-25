# crypto-ticker

Консольный тикер цены Bitcoin в реальном времени. Данные берутся с Bybit API и обновляются каждые 30 секунд.

## Вывод

- Текущая цена BTC в USD
- Изменение цены относительно предыдущего тика со стрелкой и цветом
- Алерты при пересечении заданных порогов цены

## Установка

```bash
git clone https://github.com/malayegor9-dotcom/crypto-ticker.git
cd crypto-ticker
```

## Запуск

```bash
go run .
```

С алертами:
```bash
go run . -btc-above=80000 -btc-below=75000
```

## Флаги

| Флаг | Описание | Пример |
|---|---|---|
| `-btc-above` | Алерт когда BTC выше порога | `-btc-above=80000` |
| `-btc-below` | Алерт когда BTC ниже порога | `-btc-below=75000` |

## Требования

- Go 1.21+