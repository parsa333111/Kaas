package main

import (
	"log"
	"os"

	"github.com/skye-tan/KaaS/database"
	"github.com/skye-tan/KaaS/endpoints"
	"github.com/skye-tan/KaaS/k8s_client"
	"github.com/skye-tan/KaaS/monitoring"
)

func main() {
	logFile, err := os.OpenFile("api.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	k8s_client.Initialize()
	monitoring.Initalize()
	database.Initalize()

	startupCompleteFile, err := os.Create("/tmp/startup-complete")
	if err != nil {
		panic(err)
	}
	defer startupCompleteFile.Close()

	endpoints.Start(listenAddress)
}
