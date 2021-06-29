module github.com/victoryeo/cryptoexchange/orderserver

go 1.15

replace github.com/victoryeo/cryptoexchange/engine => ../engine

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/victoryeo/cryptoexchange/engine v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.38.0
	simple v0.0.0-00010101000000-000000000000
)

replace simple => ../simple
