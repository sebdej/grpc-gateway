package main

import (
	"os"
	"strconv"

	"github.com/sebdej/grpc-gateway/gateway"
)

func main() {
	backendTls, err := strconv.ParseBool(os.Getenv("BACKEND_TLS"))

	if err != nil {
		backendTls = false
	}

	gateway.Run(os.Getenv("GATEWAY_ADDRESS"), os.Getenv("GATEWAY_TLS_CERT"), os.Getenv("GATEWAY_TLS_KEY"), os.Getenv("BACKEND_URI"), backendTls)
}
