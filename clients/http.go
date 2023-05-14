package clients

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Header     map[string]string
}

func (b *Client) prepRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	for key, value := range b.Header {
		req.Header.Add(key, value)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	return req, nil
}

func (b *Client) MakeRequest(method string, path string, reqBody io.Reader) ([]byte, int, error) {
	req, err := b.prepRequest(method, path, reqBody)
	if err != nil {
		return nil, 500, errors.Wrap(err, "constructing request")
	}

	resp, err := b.HTTPClient.Do(req)
	if err != nil {
		return nil, 500, err
	}
	defer resp.Body.Close()

	requestStr := fmt.Sprintf("%s %s", method, path)

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 500, errors.Wrapf(err, "reading response from request %q", requestStr)
	}
	return respBody, resp.StatusCode, nil
}

//set the vault client
func NewHTTPClient(baseURL string) *Client {

	header := make(map[string]string)

	return &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    baseURL,
		Header:     header,
	}
}
