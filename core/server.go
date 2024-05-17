package core

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logger = logrus.New()

var rootCmd = &cobra.Command{
	Use:   "soft-security",
	Short: "Simple api server",
}

type ProductsResponse struct {
	Products []*Product `json:"products"`
	Total    int        `json:"total"`
}

var srvCmd = &cobra.Command{
	Use:   "server",
	Short: "Starting server",
	Run:   start,
}

func Execute() {
	rootCmd.AddCommand(srvCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func start(command *cobra.Command, args []string) {
	StartServer()
}

func StartServer() {
	err := godotenv.Load(filepath.Join("./", ".env"))
	if err != nil {
		panic("Cannot find .env file")
	}

	port := os.Getenv("APP_PORT")
	store := NewStore()
	server := NewServer(store)

	tls_cert, err := tls.LoadX509KeyPair("./certificate/certificate.crt", "./certificate/private.key")
	if err != nil {
		logger.Fatalf("Failed to load X509 key pair: %v", err)
	}

	tls_config := &tls.Config{
		Certificates: []tls.Certificate{tls_cert},
	}
	httpServer := &http.Server{
		Handler:      server.Router,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		TLSConfig:    tls_config,
	}

	go func() {
		logger.Println("Starting server on port:", port)
		logger.Fatal(httpServer.ListenAndServeTLS("", ""))
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	logger.Println("Shutting down server gracefully in 1 second.")
	time.Sleep(time.Second)
	defer cancel()

	logger.Fatal(httpServer.Shutdown(ctx))
	os.Exit(0)
}
