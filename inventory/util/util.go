package util

import (
	"log"
	"os"
)

var ApplicationLog = log.New(os.Stdout, "inventory-service ", log.LstdFlags)
