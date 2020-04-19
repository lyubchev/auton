package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"

	youtube "google.golang.org/api/youtube/v3"
)

const MaxComments = 500

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, reading directly from env variables")
	}

	GoogleApiKey := os.Getenv("GOOGLE_API_KEY")

	comments := []string{}

	ctx := context.Background()
	ctx, cancelCtxFunc := context.WithCancel(ctx)

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(GoogleApiKey))
	if err != nil {
		panic(err)
	}

	commentThreadsService := youtube.NewCommentThreadsService(youtubeService)

	// Get the instance with which we can do api calls to the youtube api
	apiCall := commentThreadsService.List("snippet")

	apiCall.TextFormat("plainText")
	apiCall.VideoId("1DfbhdG0tEk")
	apiCall.Order("relevance")

	err = apiCall.Pages(ctx, func(resp *youtube.CommentThreadListResponse) error {
		for _, item := range resp.Items {
			c := item.Snippet.TopLevelComment.Snippet.TextDisplay
			comments = append(comments, c)

			if len(comments) == MaxComments {
				cancelCtxFunc()
				break
			}
		}

		return nil
	})

	if err != nil && !errors.Is(err, context.Canceled) {
		panic(err)
	}

	spew.Dump(comments)
}
