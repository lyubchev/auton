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
// returns a map where the key is the tone and the value is the score
func AnalyzeCommentsTone(comments []string, ibmClient *ibm.Client) (map[Tone]float64, error) {
	batches := batchComments(comments)

	// toneComputed stores each tone and because we may have many batches of comments each
	// batch will return us a new result (score) then we will re-calculate the score of the
	// specific tone by averaging it
	//
	// For example:
	// {
	//	"Analytical": [0.75, 0.85, 0.61]
	//  "Anger": [0.98, 0,51, 0,53 ],
	//  }
	//
	// Will be computed to this:
	//  {
	//	"Analytical": [0.73]
	//  "Anger": [0.67],
	//  }
	tc := map[string][]float64{}

	for _, batch := range batches {
		tones, err := ibmClient.Do(batch)
		for k, v := range tones {
			tc[k] = append(tc[k], v)
		}

		if err != nil {
			return nil, err
		}
	}

	result := map[string]float64

	for k, v := range tc {
		avgScore := 0.0

		for _, s := range tc[k] {
			avgScore += s
		}

		avgScore = avgScore / len(tc[k])
		result[k] = avgScore
	}


	return result, nil
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
