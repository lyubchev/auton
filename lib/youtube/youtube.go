package youtube

import (
	"context"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

type Order string

const (
	OrderRelevance Order = "relevance"
	OrderTime      Order = "time"
)

const rateLimit = time.Second / 20

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	c := &Client{apiKey: apiKey}
	return c
}

func (c *Client) GetComments(videoId string, order Order, maxComments int) ([]string, error) {
	comments := []string{}
	youtubeService, err := youtube.NewService(context.Background(), option.WithAPIKey(c.apiKey))
	if err != nil {
		return nil, err
	}

	commentThreadsService := youtube.NewCommentThreadsService(youtubeService)

	// Get the instance with which we can do api calls to the youtube api
	apiCaller := commentThreadsService.List("snippet")

	apiCaller.TextFormat("plainText")
	apiCaller.VideoId(videoId)
	apiCaller.Order(string(order))

	var wg sync.WaitGroup
	commentsChan := make(chan string, runtime.NumCPU())
	quitChan := make(chan struct{})

	throttle := time.NewTicker(rateLimit)

	go func() {
	forLoop:
		for {
			select {
			case <-quitChan:
				break forLoop
			default:
				wg.Add(1)

				<-throttle.C
				go func() {
					defer wg.Done()

					resp, err := apiCaller.Do()
					if err != nil {
						log.Println(err)
						return
					}

					pageToken := resp.NextPageToken
					if pageToken == "" {
						quitChan <- struct{}{}
						return
					}
					apiCaller.PageToken(pageToken)

					for _, item := range resp.Items {
						commentsChan <- item.Snippet.TopLevelComment.Snippet.TextDisplay
					}
				}()

			}
		}
		wg.Wait()
		close(commentsChan)
		close(quitChan)
	}()

	for c := range commentsChan {
		comments = append(comments, c)

		spew.Dump(len(comments))
		if len(comments) == maxComments {
			quitChan <- struct{}{}
			return comments, nil
		}
	}

	return comments, nil
}
