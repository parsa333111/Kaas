package main

import (
	"log"
	"os"

	"github.com/skye-tan/KaaS/tools/notepad/database"
	"github.com/skye-tan/KaaS/tools/notepad/endpoints"
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Info: Service started.")

	listenAddress, ok := os.LookupEnv("LISTEN_ADDRESS")
	if !ok {
		log.Println("Warn: Missing enviroment variable LISTEN_ADDRESS.",
			"Using default listen address: [0.0.0.0:8081]")
		listenAddress = "0.0.0.0:8081"
	}

	database.Initalize()

	endpoints.Start(listenAddress)
}
