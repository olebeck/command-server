module command-server

go 1.20

require (
	github.com/gofiber/fiber/v2 v2.42.0
	github.com/gofiber/websocket/v2 v2.1.4
	github.com/sandertv/mcwss v1.2.0
	github.com/sirupsen/logrus v1.9.0
)

replace github.com/sandertv/mcwss => github.com/olebeck/mcwss v1.2.0-1

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/fasthttp/websocket v1.5.1 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20220530130905-52f3993e8d6d // indirect
	github.com/sergi/go-diff v1.0.0 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.44.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/yudai/gojsondiff v1.0.0 // indirect
	github.com/yudai/golcs v0.0.0-20170316035057-ecda9a501e82 // indirect
	golang.org/x/sys v0.0.0-20220811171246-fbc7d0a398ab // indirect
)
