package main

import (
	"os"

	"github.com/google/logger"
)

func main() {
	rootlogger := logger.Init("root", false, false, os.Stderr)
	rootlogger.Warning("Testing 123")
}
