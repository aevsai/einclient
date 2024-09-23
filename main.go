package main

import (
	"flag"
	"log"
	"os"
	"time"

	"einclient/engine"
	"einclient/loop"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func LogErrorAndCapture(logger zerolog.Logger, err error, msg string) {
	// Log the error using zerolog
	logger.Error().Err(err).Msg(msg)

	// Capture the exception using Sentry
	sentry.CaptureException(err)
}

func LogMessageAndCapture(logger zerolog.Logger, level zerolog.Level, msg string) {
	// Log the message using zerolog
	logger.WithLevel(level).Msg(msg)

	// Capture the message using Sentry
	sentry.CaptureMessage(msg)
}

func main() {
	flag.Parse()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file: %s", err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	defer sentry.Flush(2 * time.Second)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Sample log message
	LogMessageAndCapture(logger, zerolog.InfoLevel, "Hello, World!")

	// Call a function that might produce an error and capture it with Sentry
	scenePath := "./scenes/basic.yml"
	s, err := engine.LoadScene(scenePath)
	if err != nil {
		LogErrorAndCapture(logger, err, "An error occurred")
	}
	l, err := loop.NewLoop(s)
	err = l.Start()
	if err != nil {
		LogErrorAndCapture(logger, err, "An error occurred")
	}
}
