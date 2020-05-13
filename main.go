package main

import (
	"log"
	"net/http"
	"os"

	"github.com/impzero/auton/lib/ibm"
	"github.com/impzero/auton/lib/youtube"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, reading directly from env variables")
	}

	GoogleAPIKey := os.Getenv("GOOGLE_API_KEY")
	IBMAPIKey := os.Getenv("IBM_API_KEY")

	youtubeClient := youtube.New(GoogleAPIKey)
	ibmClient, err := ibm.New(ibm.Config{APIKey: IBMAPIKey, ServiceURL: "https://api.eu-gb.tone-analyzer.watson.cloud.ibm.com"})
	if err != nil {
		panic(err)
	}

	w := NewWeb(youtubeClient, ibmClient)
	log.Println("Server started, listening on port :8080")
	if err := http.ListenAndServe(":8080", w.Router); err != nil {
		panic(err)
	}
}
