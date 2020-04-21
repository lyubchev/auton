package main

import (
	"github.com/impzero/autone/lib/ibm"
	"github.com/impzero/autone/tones"
)

const MaxRequestSize = 128000

// AnalyzeCommentsTone takes all the comments from a youtube video passed
// in array, batches them where each batch is no more than 128kB (the maximum request size IBM accepts)
// returns a map where the key is the tone and the value is the score
func AnalyzeCommentsTone(comments []string, ibmClient *ibm.Client) (map[tones.Tone]float64, error) {
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
	tc := map[tones.Tone][]float64{}

	for _, batch := range batches {
		tones, err := ibmClient.Do(batch)
		for k, v := range tones {
			tc[k] = append(tc[k], v)
		}

		if err != nil {
			return nil, err
		}
	}

	result := map[tones.Tone]float64{}

	for k, _ := range tc {
		avgScore := 0.0

		for _, s := range tc[k] {
			avgScore += s
		}

		avgScore = avgScore / float64(len(tc[k]))
		result[k] = avgScore
	}

	return result, nil
}

func batchComments(comments []string) []string {
	batches := []string{}

	batch := ""
	for _, comment := range comments {
		// We make batches with a maximum size of `MaxRequestSize` measured in chars

		if len(batch+". "+comment) <= MaxRequestSize {
			batch += ". " + comment
		} else {
			batches = append(batches, batch)
			batch = comment
		}
	}

	if len(batches) == 0 {
		batches = append(batches, batch)
	}

	return batches
}
