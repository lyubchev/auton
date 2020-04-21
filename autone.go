package main

import (
	"github.com/impzero/autone/lib/ibm"
)

type Tone string

const (
	Anger      Tone = "Anger"
	Fear       Tone = "Fear"
	Joy        Tone = "Joy"
	Sadness    Tone = "Sadness"
	Analytical Tone = "Analytical"
	Confident  Tone = "Confident"
	Tentative  Tone = "Tentative"
)

const MaxRequestSize = 128000

// AnalyzeCommentsTone takes all the comments from a youtube video passed
// in array, batches them where each batch is no more than 128kB (the maximum request size IBM accepts)
// returns a map where the key is the tone and the value is the score in percentages
func AnalyzeCommentsTone(comments []string, ibmClient *ibm.Client) (map[Tone]string, error) {
	batches := batchComments(comments)

	for _, batch := range batches {
		tones, err := ibmClient.Do(batch)
		if err != nil {
			return nil, err
		}
	}

}

func batchComments(comments []string) []string {
	batches := []string{}
	batchID := 0

	for _, comment := range comments {
		// We make batches with a maximus sie of `MaxRequestSize` measured in chars
		if len(batches[batchID]+". "+comment) > MaxRequestSize {
			batchID++
		}

		batches[batchID] = batches[batchID] + ". " + comment
	}
}
