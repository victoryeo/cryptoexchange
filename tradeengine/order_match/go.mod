module order_match

go 1.15

replace github.com/victoryeo/cryptoexchange/engine => ../engine

require (
	github.com/go-redis/redis/v8 v8.10.0
	github.com/victoryeo/cryptoexchange/engine v0.0.0-00010101000000-000000000000
)
