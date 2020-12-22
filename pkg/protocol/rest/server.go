package rest

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "getitqec.com/server/chat/pkg/api/v1"
	// pb "../api/v1"
	"getitqec.com/server/chat/pkg/logger"
	"getitqec.com/server/chat/pkg/protocol/rest/middleware"
)

func CustomMatcher(key string) (string, bool) {
	switch key {
	case "Token":
		return "token", true
	// case "Authorization":
	// 	return "Authorization", true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

// RunServer runs HTTP/REST gateway
func RunServer(ctx context.Context, grpcPort, httpPort, certFilePath string, keyFilePath string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// mux := runtime.NewServeMux()
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
		runtime.WithIncomingHeaderMatcher(CustomMatcher),
		runtime.WithErrorHandler(DefaultHTTPProtoErrorHandler),
		// runtime.WithProtoErrorHandler(DefaultHTTPProtoErrorHandler),
	)
	opts := []grpc.DialOption{}
	if certFilePath != "" && keyFilePath != "" {
		// creds, err := credentials.NewServerTLSFromFile(certFilePath, keyFilePath)
		// creds, err := credentials.NewClientTLSFromFile(certFilePath, "CheeTest")
		// if err != nil {
		// 	log.Fatalf("Failed to generate credentials %v", err)
		// }

		b, _ := ioutil.ReadFile(certFilePath)
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			log.Fatalf("fail to dial: %v", errors.New("credentials: failed to append certificates"))
		}
		config := &tls.Config{
			InsecureSkipVerify: true,
			RootCAs:            cp,
		}
		creds := credentials.NewTLS(config)

		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	if err := pb.RegisterChatServiceHandlerFromEndpoint(ctx, mux, "localhost:"+grpcPort, opts); err != nil {
		// log.Fatalf("failed to start HTTP gateway: %v", err)
		logger.Log.Fatal("failed to start HTTP gateway: %v", zap.String("reason", err.Error()))
	}
	fmt.Println("REST : gRPC client up")

	// cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath)
	// if err != nil {
	// 	log.Println(err)
	// 	return err
	// }

	// tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		Addr: ":" + httpPort,
		// Handler: mux,
		Handler: middleware.AddRequestID(middleware.AddLogger(logger.Log, mux)),
		// TLSConfig: tlsConfig,
	}

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
		}

		_, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		_ = srv.Shutdown(ctx)
	}()

	fmt.Println("starting HTTP/REST gateway...")
	logger.Log.Info("starting HTTP/REST gateway...")
	// return srv.ListenAndServeTLS(certFilePath, keyFilePath)
	return srv.ListenAndServe()
}
