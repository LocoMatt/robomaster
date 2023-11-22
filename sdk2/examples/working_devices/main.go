package main

import (
	"log/slog"
	"time"

	"github.com/brunoga/robomaster/sdk2"
	"github.com/brunoga/unitybridge/support/logger"
)

func main() {
	l := logger.New(slog.LevelDebug)

	c, err := sdk2.New(l, 0)
	if err != nil {
		panic(err)
	}

	err = c.Start()
	if err != nil {
		panic(err)
	}
	defer c.Stop()

	time.Sleep(10 * time.Second)
}