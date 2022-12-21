package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	collectionsv1 "github.com/sebdej/grpc-gateway/gen/proto/go/collections/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
)

func Run(gatewayAddress string, gatewayCertFile string, gatewayKeyFile string, backendUri string, tlsBackend bool) error {
	log := grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr)
	grpclog.SetLoggerV2(log)

	var transportCredentials credentials.TransportCredentials

	if tlsBackend {
		transportCredentials = credentials.NewTLS(&tls.Config{ServerName: "", InsecureSkipVerify: true})
	} else {
		transportCredentials = insecure.NewCredentials()
	}

	conn, err := grpc.DialContext(
		context.Background(),
		backendUri,
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithBlock(),
	)

	if err != nil {
		return fmt.Errorf("DialContext returned: %w", err)
	}

	defer conn.Close()

	if tlsBackend {
		log.Info("gRPC gateway connected to ", backendUri, " with TLS")
	} else {
		log.Info("gRPC gateway connected to ", backendUri, " in plain text")
	}

	mux := runtime.NewServeMux()

	err = collectionsv1.RegisterCollectionServiceHandler(context.Background(), mux, conn)

	if err != nil {
		return fmt.Errorf("RegisterCollectionServiceHandler failed: %w", err)
	}

	handler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		mux.ServeHTTP(writer, request)
	})

	server := &http.Server{
		Addr:    gatewayAddress,
		Handler: handler,
	}

	if gatewayKeyFile != "" {
		log.Info("gRPC gateway running on: https://", gatewayAddress)

		return fmt.Errorf("ListenAndServeTLS: returned: %w", server.ListenAndServeTLS(gatewayCertFile, gatewayKeyFile))
	} else {
		log.Info("gRPC gateway running on: http://", gatewayAddress)

		return fmt.Errorf("ListenAndServe: returned: %w", server.ListenAndServe())
	}
}
