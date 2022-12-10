package main

import (
	"log"
	"os"
)

func (e *Emporia) SetToken() {
	e.token = os.Getenv("EMPORIA_TOKEN")
	if e.token == "" {
		log.Panicf("Error: EMPORIA_TOKEN environment variable not set\n")
	}
}

func (e *Emporia) SetDevice() {
	e.device = os.Getenv("EMPORIA_DEVICE")
	if e.device == "" {
		log.Panicf("Error: EMPORIA_DEVICE environment variable not set\n")
	}
}
