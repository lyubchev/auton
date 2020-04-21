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
	apiKey string
}

func New(apiKey string) *Client {
	c := &Client{apiKey: apiKey}
	return c
}

func (c *Client) GetComments(videoId string, order Order, maxComments int) ([]string, error) {

	comments := []string{}

	ctx := context.Background()
	ctx, cancelCtxFunc := context.WithCancel(ctx)

	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(c.apiKey))
	if err != nil {
		cancelCtxFunc()
		return nil, err
	}

	commentThreadsService := youtube.NewCommentThreadsService(youtubeService)

	// Get the instance with which we can do api calls to the youtube api
	apiCall := commentThreadsService.List("snippet")

	apiCall.TextFormat("plainText")
	apiCall.VideoId(videoId)
	apiCall.Order(string(order))

	err = apiCall.Pages(ctx, func(resp *youtube.CommentThreadListResponse) error {
		for _, item := range resp.Items {
			c := item.Snippet.TopLevelComment.Snippet.TextDisplay
			comments = append(comments, c)

			lenComments := len(comments)

			log.Printf("%d/%d comments fetched!", lenComments, maxComments)
			if lenComments == maxComments {
				cancelCtxFunc()
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
