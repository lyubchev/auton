package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/impzero/autone/lib/ibm"
	"github.com/impzero/autone/lib/youtube"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, reading directly from env variables")
	}

	GoogleApiKey := os.Getenv("GOOGLE_API_KEY")
	IbmApiKey := os.Getenv("IBM_API_KEY")

	youtubeClient := youtube.New(GoogleApiKey)
	comments, err := youtubeClient.GetComments("sM_NlKAcGpo", youtube.OrderRelevance, 100)
	if err != nil {
		panic(err)
	}

	ibmClient, err := ibm.New(ibm.Config{ApiKey: IbmApiKey, ServiceUrl: "https://api.eu-gb.tone-analyzer.watson.cloud.ibm.com"})
	if err != nil {
		panic(err)
	}

	tones, err := AnalyzeCommentsTone(comments, ibmClient)
	if err != nil {
		panic(err)
	}

	spew.Dump(tones)
}
