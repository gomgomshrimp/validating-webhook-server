package main

import (
	"crypto/tls"
	"github.com/gofiber/fiber/v2"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gomgomshrimp/validating-webhook-server/api"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const idleTimeout = 5 * time.Second

func main() {
	// Set Logger
	logger := setLogger()
	// Create AdmissionWebhookAPIs
	admissionWebhookApi := api.NewApi(logger)
	// Configure go/fiber
	app := fiber.New(fiber.Config{
		IdleTimeout:             idleTimeout,
		EnableTrustedProxyCheck: true,
	})
	app.Use(fiberLogger.New())
	app.Post("/validate", admissionWebhookApi.Validate)

	// Create tls certificate
	cer, err := tls.LoadX509KeyPair("certs/server.crt", "certs/server.key")
	if err != nil {
		logger.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	// Create tls listener
	listener, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		logger.Fatal(err)
	}
	// Start server from goroutine with https/ssl -- https://localhost:443
	go func() {
		if err = app.Listener(listener); err != nil {
			logger.Fatal(err)
		}
	}()

	// Check shutdown signals(interrupt or termination)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	_ = <-sigChan
	logger.Info("Gracefully shutting down go/fiber...")
	_ = app.Shutdown()

	// App cleaning tasks
	logger.Info("go/fiber was successful shutdown.")
}

func setLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logger.SetLevel(logrus.DebugLevel)
	return logger
}
