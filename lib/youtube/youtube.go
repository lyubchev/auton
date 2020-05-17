package youtube

import (
	"context"
	"errors"
	"log"

	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

type Order string

const (
	OrderRelevance Order = "relevance"
	OrderTime      Order = "time"
)

type Client struct {
	APIKey string
}

func New(apiKey string) *Client {
	c := &Client{APIKey: apiKey}
	return c
}

func (c *Client) GetComments(videoID string, order Order, maxComments int) ([]string, error) {
	comments := []string{}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(c.APIKey))
	if err != nil {
		cancel()
		return nil, err
	}

	commentThreadsService := youtube.NewCommentThreadsService(youtubeService)

	// Get the instance with which we can do api calls to the youtube api
	apiCall := commentThreadsService.List("snippet")

	apiCall.TextFormat("plainText")
	apiCall.VideoId(videoID)
	apiCall.Order(string(order))
	apiCall.MaxResults(100)

	err = apiCall.Pages(ctx, func(resp *youtube.CommentThreadListResponse) error {
		for _, item := range resp.Items {
			c := item.Snippet.TopLevelComment.Snippet.TextDisplay
			comments = append(comments, c)

			lenComments := len(comments)
			log.Printf("%d/%d comments fetched!", lenComments, maxComments)

			if lenComments == maxComments {
				cancel()
				break
			}
		}

		return nil
	})

	if err != nil && !errors.Is(err, context.Canceled) {
		return nil, err
	}

	return comments, nil
}
