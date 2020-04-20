package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/impzero/autone/autone/ibm"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, reading directly from env variables")
	}

	// GoogleApiKey := os.Getenv("GOOGLE_API_KEY")
	IbmApiKey := os.Getenv("IBM_API_KEY")

	// youtubeClient := youtube.New(GoogleApiKey)
	// comments, err := youtubeClient.GetComments("xvZqHgFz51I", youtube.OrderRelevance, 100)
	if err != nil {
		panic(err)
	}

	ibmClient, err := ibm.New(ibm.Config{ApiKey: IbmApiKey, ServiceUrl: "https://api.eu-gb.tone-analyzer.watson.cloud.ibm.com"})
	if err != nil {
		panic(err)
	}

	result, err := ibmClient.Do("I fucking hate this. This annoys me so badly and I won't use this thing again. Fucking die!")
	if err != nil {
		panic(err)
	}

	spew.Dump(result)

	// spew.Dump(comments)
}
