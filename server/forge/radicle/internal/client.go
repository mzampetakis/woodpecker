package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	get  = "GET"
	post = "POST"
	put  = "PUT"
	del  = "DELETE"
)

const (
	FileTypeBlob      = "blob"
	FileTypeDirectory = "tree"

	AUTHORIZED_SESSION = "authorized"
)

const (
	apiPath                 = "/api"
	apiV1Path               = "/v1"
	pathNode                = "%s/node"
	pathSession             = "%s/sessions/%s"
	pathProject             = "%s/projects/%s"
	pathProjects            = "%s/projects?page=%s&perPage=%s"
	pathProjectCommits      = "%s/projects/%s/commits?%s"
	pathProjectCommitFile   = "%s/projects/%s/blob/%s/%s"
	pathProjectCommitDir    = "%s/projects/%s/tree/%s/%s"
	pathProjectPatches      = "%s/projects/%s/patches?%s"
	pathProjectPatchComment = "%s/projects/%s/patches/%s"
)

type Client struct {
	*http.Client
	base  string
	ctx   context.Context
	token string
}

func NewClient(ctx context.Context, url string, secretToken string) *Client {
	return &Client{
		Client: http.DefaultClient,
		base:   url,
		ctx:    ctx,
		token:  secretToken,
	}
}

func (c *Client) GetNodeInfo() (*NodeInfo, error) {
	out := new(NodeInfo)
	uri := fmt.Sprintf(pathNode, c.base+apiPath+apiV1Path)
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return out, err
}

func (c *Client) GetSessionInfo() (*SessionInfo, error) {
	out := new(SessionInfo)
	uri := fmt.Sprintf(pathSession, c.base+apiPath+apiV1Path, c.token)
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return out, err
}

func (c *Client) GetProject(projectID string) (*Project, error) {
	out := new(Project)
	uri := fmt.Sprintf(pathProject, c.base+apiPath+apiV1Path, projectID)
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return out, err
}

func (c *Client) GetProjects() ([]*Project, error) {
	var projects []*Project
	var err error
	page := 0
	perPage := 100
	for {
		out := new([]*Project)
		uri := fmt.Sprintf(pathProjects, c.base+apiPath+apiV1Path, strconv.Itoa(page), strconv.Itoa(perPage))
		_, err = c.do(uri, get, nil, out)
		if err != nil {
			return nil, err
		}
		if len(*out) == 0 {
			break
		}
		page++
		projects = append(projects, *out...)
	}
	return projects, nil
}

func (c *Client) GetProjectCommits(projectID string, listOpts ListOpts) (*Commits, error) {
	out := new(Commits)
	uri := fmt.Sprintf(pathProjectCommits, c.base+apiPath+apiV1Path, projectID, listOpts.Encode())
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return out, err
}

func (c *Client) GetProjectCommitFile(projectID, commit, file string) (*ProjectFile, error) {
	out := new(ProjectFile)
	uri := fmt.Sprintf(pathProjectCommitFile, c.base+apiPath+apiV1Path, projectID, commit, file)
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return out, err
}

func (c *Client) GetProjectCommitDir(projectID, commit, path string) (FileTree, error) {
	out := new(FileTree)
	uri := fmt.Sprintf(pathProjectCommitDir, c.base+apiPath+apiV1Path, projectID, commit, path)
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return *out, err
}

func (c *Client) GetProjectPatches(projectID string, listOpts ListOpts) ([]*Patch, error) {
	out := new([]*Patch)
	uri := fmt.Sprintf(pathProjectPatches, c.base+apiPath+apiV1Path, projectID, listOpts.Encode())
	fmt.Println(uri)
	_, err := c.do(uri, get, nil, out)
	return *out, err
}

//func (c *Client) AddProjectPatchComment(owner, projectID, revision string, status *PipelineStatus) error {
//	out := new(Project)
//
//	uri := fmt.Sprintf(pathProject, c.base+apiPath+apiV1Path, projectID)
//	fmt.Println(uri)
//	_, err := c.do(uri, get, nil, out)
//	return out, err
//}

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
	if len(c.token) > 0 {
		req.Header.Set("Authorization", "Bearer "+c.token)
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
