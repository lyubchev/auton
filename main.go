package main

import (
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/impzero/autone/autone/ibm"
	"github.com/impzero/autone/autone/youtube"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, reading directly from env variables")
	}

	GoogleApiKey := os.Getenv("GOOGLE_API_KEY")

	youtubeClient := youtube.New(GoogleApiKey)
	comments, err := youtubeClient.GetComments("xvZqHgFz51I", youtube.OrderRelevance, 100)

	ibm.AnalyzeTone()
	if err != nil {
		panic(err)
	}

	spew.Dump(comments)
}
