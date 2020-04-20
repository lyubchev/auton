package ibm

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/toneanalyzerv3"
)

func AnalyzeTone() {
	authenticator := &core.IamAuthenticator{
		ApiKey: "{apiKey}",
	}

	options := &toneanalyzerv3.ToneAnalyzerV3Options{
		Version:       "2017-09-21",
		Authenticator: authenticator,
	}

	toneAnalyzer, toneAnalyzerErr := toneanalyzerv3.NewToneAnalyzerV3(options)

	if toneAnalyzerErr != nil {
		panic(toneAnalyzerErr)
	}

	err := toneAnalyzer.SetServiceURL("https://api.eu-gb.tone-analyzer.watson.cloud.ibm.com")
	if err != nil {
		panic(err)
	}

	text := "Team, I know that times are tough! Product sales have been disappointing for the past three quarters. We have a competitive product, but we need to do a better job of selling it!"
	result, _, responseErr := toneAnalyzer.Tone(
		&toneanalyzerv3.ToneOptions{
			ToneInput: &toneanalyzerv3.ToneInput{
				Text: &text,
			},
			ContentType: core.StringPtr("application/json"),
		},
	)
	if responseErr != nil {
		panic(responseErr)
	}

	b, _ := json.MarshalIndent(result, "", "  ")
	fmt.Println(string(b))
}
