# marketflow

# Description
This system designed to process market data in real-time and a RESTful API for accessing price updates.

# Installation
docker compose up

To run provided programs:
make load
make run

# Usage

## API Endpoints

For live mode use:
Exchanges: exchange1, exchange2, exchange3
Symbols: BTCUSDT, DOGEUSDT, TONUSDT, SOLUSDT, ETHUSDT

For test mode use:
Test exchanges: exchange1_test, exchange2_test, exchange3_test
Symbols: BTCUSDT_test, DOGEUSDT_test, TONUSDT_test, SOLUSDT_test, ETHUSDT_test

Market Data API

GET /prices/latest/{symbol} – Get the latest price for a given symbol.

GET /prices/latest/{exchange}/{symbol} – Get the latest price for a given symbol from a specific exchange.

GET /prices/highest/{symbol} – Get the highest price over a period.

GET /prices/highest/{exchange}/{symbol} – Get the highest price over a period from a specific exchange.

GET /prices/highest/{symbol}?period={duration} – Get the highest price within the last {duration} (e.g., the last 1s, 3s, 5s, 10s, 30s, 1m, 3m, 5m).

GET /prices/highest/{exchange}/{symbol}?period={duration} – Get the highest price within the last {duration} from a specific exchange.

GET /prices/lowest/{symbol} – Get the lowest price over a period.

GET /prices/lowest/{exchange}/{symbol} – Get the lowest price over a period from a specific exchange.

GET /prices/lowest/{symbol}?period={duration} – Get the lowest price within the last {duration}.

GET /prices/lowest/{exchange}/{symbol}?period={duration} – Get the lowest price within the last {duration} from a specific exchange.

GET /prices/average/{symbol} – Get the average price over a period.

GET /prices/average/{exchange}/{symbol} – Get the average price over a period from a specific exchange.

GET /prices/average/{exchange}/{symbol}?period={duration} – Get the average price within the last {duration} from a specific exchange

Data Mode API

POST /mode/test – Switch to Test Mode (use generated data).

POST /mode/live – Switch to Live Mode (fetch data from provided programs).

System Health

GET /health - Returns system status (e.g., connections, Redis availability).