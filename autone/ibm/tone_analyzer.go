package ibm

import (
	"fmt"

	"github.com/IBM/go-sdk-core/core"
	"github.com/watson-developer-cloud/go-sdk/toneanalyzerv3"
)

type ContentType string

const (
	ContentTypeApplicationJson ContentType = "application/json"
)

type Config struct {
	ApiKey     string
	ServiceUrl string
}

type Client struct {
	toneAnalyzer *toneanalyzerv3.ToneAnalyzerV3
}

func New(config Config) (*Client, error) {
	authenticator := &core.IamAuthenticator{
		ApiKey: config.ApiKey,
	}

	options := &toneanalyzerv3.ToneAnalyzerV3Options{
		Version:       "2017-09-21",
		Authenticator: authenticator,
	}

	toneAnalyzer, err := toneanalyzerv3.NewToneAnalyzerV3(options)
	if err != nil {
		return nil, err
	}

	err = toneAnalyzer.SetServiceURL(config.ServiceUrl)
	if err != nil {
		return nil, err
	}

	ta := &Client{
		toneAnalyzer: toneAnalyzer,
	}

	return ta, nil
}

func (c *Client) Do(text string) ([]string, error) {
	generalTone := []string{}

	result, _, err := c.toneAnalyzer.Tone(
		&toneanalyzerv3.ToneOptions{
			ToneInput: &toneanalyzerv3.ToneInput{
				Text: &text,
			},
			ContentType: core.StringPtr(string(ContentTypeApplicationJson)),
		},
	)
	if err != nil {
		return nil, err
	}

	for _, tone := range result.DocumentTone.Tones {

		scoreInPerc := fmt.Sprintf("%f", *tone.Score*100)[:6]

		generalTone = append(generalTone, *tone.ToneName+": "+scoreInPerc+"%")
	}

	return generalTone, nil
}
