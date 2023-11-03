package internal

import (
	"net/url"
	"strconv"
)

type ListOpts struct {
	Page    int
	PerPage int
}

type NodeInfo struct {
	ID     string     `json:"id"`
	Config NodeConfig `json:"config"`
}

type NodeConfig struct {
	Alias string `json:"alias"`
}

type Project struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Delegates     []string `json:"delegates"`
	DefaultBranch string   `json:"defaultBranch"`
	Head          string   `json:"head"`
}

type Commits struct {
	Commits []CommitObject `json:"commits"`
	Stats   CommitStats    `json:"stats"`
}
type CommitStats struct {
	Commits      uint `json:"commits"`
	Branches     uint `json:"branches"`
	Contributors uint `json:"contributors"`
}

type CommitObject struct {
	Commit Commit `json:"commit"`
}

type Commit struct {
	ID      string   `json:"id"`
	Parents []string `json:"parents"`
}

func (o *ListOpts) Encode() string {
	params := url.Values{}
	if o.Page != 0 {
		params.Set("page", strconv.Itoa(o.Page))
	}
	if o.PerPage != 0 {
		params.Set("perPage", strconv.Itoa(o.PerPage))
	}
	return params.Encode()
}

type Error struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e Error) Error() string {
	return e.Body.Message
}
