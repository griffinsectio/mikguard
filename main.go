package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/go-routeros/routeros/v3"
)

func dial(address, username, password string) (*routeros.Client, error) {
	var useTLS = false
	if useTLS {
		return routeros.DialTLS(address, username, password, nil)
	}
	return routeros.Dial(address, username, password)
}

func fatal(log *slog.Logger, message string, err error) {
	log.Error(message, slog.Any("error", err))
	os.Exit(2)
}

func main() {
	var debug = false
	logLevel := slog.LevelInfo
	if debug {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	})

	log := slog.New(handler)

	address, username, password := "127.0.0.1:8728", "admin", ""
	c, err := dial(address, username, password)
	if err != nil {
		fatal(log, "could not connect", err)
	}
	defer c.Close()

	c.SetLogHandler(handler)

	var async = false
	if async {
		c.Async()
	}

	var command = "/interface/wireguard/add =name=wg0 =listen-port=13231"
	r, err := c.RunArgs(strings.Split(command, " "))
	if err != nil {
		fatal(log, "could not run args", err)
	}

	log.Info("received results", slog.Any("results", r))
}
