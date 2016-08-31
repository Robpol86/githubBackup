package main

import (
	log "github.com/Sirupsen/logrus"
)

func main() {
	log.Debug("Debug.")
	log.Info("Info.")
	log.Warn("Warn.")
	log.Error("Error.")
}
