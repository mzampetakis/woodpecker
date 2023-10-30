package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	//shared_utils "github.com/woodpecker-ci/woodpecker/shared/utils"
)

const (
	get  = "GET"
	post = "POST"
	put  = "PUT"
	del  = "DELETE"
)

const (
	apiPath      = "/api"
	apiV1Path    = "/v1"
	pathNodeInfo = "%s/node"
)

type Client struct {
	*http.Client
	base string
	ctx  context.Context
}

func NewClient(ctx context.Context, url string, secretToken string) *Client {
	return &Client{
		Client: http.DefaultClient,
		base:   url,
		ctx:    ctx,
	}
}

func (c *Client) GetNodeInfo() (*NodeInfo, error) {
	out := new(NodeInfo)
	uri := fmt.Sprintf(pathNodeInfo, c.base+apiPath+apiV1Path)
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return out, err
}

func (c *Client) do(rawurl, method string, in, out interface{}) (*string, error) {
	uri, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	// if we are posting or putting data, we need to
	// write it to the body of the request.
	var buf io.ReadWriter
	if in != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(in)
		if err != nil {
			return nil, err
		}
	}

	// creates a new http request to radicle httpd.
	req, err := http.NewRequestWithContext(c.ctx, method, uri.String(), buf)
	if err != nil {
		return nil, err
	}
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// if an error is encountered, parse and return the
	// error response.
	if resp.StatusCode > http.StatusPartialContent {
		err := Error{}
		_ = json.NewDecoder(resp.Body).Decode(&err)
		err.Status = resp.StatusCode

		return nil, err
	}

	// if a json response is expected, parse and return
	// the json response.
	if out != nil {

		return nil, json.NewDecoder(resp.Body).Decode(out)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, err
	}
	bodyString := string(bodyBytes)

	return &bodyString, nil
}
