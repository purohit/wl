package oauth

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/pivotal-golang/lager"
	"github.com/robdimsdale/wundergo"
)

// oauthClient is an implementation of wundergo.Client.
type oauthClient struct {
	apiURL      string
	accessToken string
	clientID    string
	logger      lager.Logger
}

// NewClient is a utility method to simplify initialization
// of a new oauthClient.
func NewClient(
	accessToken string,
	clientID string,
	apiURL string,
	logger lager.Logger,
) wundergo.Client {
	return &oauthClient{
		apiURL:      apiURL,
		accessToken: accessToken,
		clientID:    clientID,
		logger:      logger,
	}
}

func (c oauthClient) validateRecurrence(recurrenceType string, recurrenceCount uint) error {
	if recurrenceType == "" && recurrenceCount > 0 {
		return errors.New("recurrenceCount must be zero if provided recurrenceType is not provided")
	}

	if recurrenceCount == 0 && recurrenceType != "" {
		return errors.New("recurrenceType must be valid if provided recurrenceCount is non-zero")
	}

	return nil
}

func (c oauthClient) addAuthHeaders(req *http.Request) {
	req.Header.Add("X-Access-Token", c.accessToken)
	req.Header.Add("X-Client-ID", c.clientID)
}

func (c oauthClient) newGetRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	c.addAuthHeaders(req)
	return req, nil
}

func (c oauthClient) newPostRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addBody(req, body)

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c oauthClient) addBody(req *http.Request, body []byte) {
	if body != nil && len(body) > 0 {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
}

func (c oauthClient) newPutRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addBody(req, body)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (c oauthClient) newPatchRequest(url string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return nil, err
	}

	c.addAuthHeaders(req)
	c.addBody(req, body)

	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (c oauthClient) newDeleteRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	c.addAuthHeaders(req)
	return req, nil
}