package pocket

import (
	"errors"
	"strings"
)

const ()

type (
	requestTokenRequest struct {
		ConsumerKey string `json:"consumer_key"`
		RedirectURI string `json:"redirect_uri"`
	}

	authorizeRequest struct {
		ConsumerKey string `json:"consumer_key"`
		Code        string `json:"code"`
	}

	AuthorizeResponse struct {
		AccessToken string `json:"access_token"`
		Username    string `json:"access_token"`
	}

	addRequest struct {
		URL         string `json:"url"`
		Title       string `json:"title,omitempty"`
		Tags        string `json:"tags,omitempty"`
		AccessToken string `json:"access_token"`
		ConsumerKey string `json:"consumer_key"`
	}

	AddInput struct {
		URL         string
		Title       string
		Tags        []string
		AccessToken string
	}
)

func (i AddInput) validate() error {
	if i.URL == "" {
		return errors.New("required URL value is empty")
	}

	if i.AccessToken == "" {
		return errors.New("access token is empty")
	}

	return nil
}

func (i AddInput) generateRequest(consumerKey string) addRequest {
	return addRequest{
		URL:         i.URL,
		Title:       i.Title,
		Tags:        strings.Join(i.Tags, ","),
		AccessToken: i.AccessToken,
		ConsumerKey: consumerKey,
	}
}
