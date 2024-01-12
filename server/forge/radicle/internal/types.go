package internal

import (
	"net/url"
	"strconv"
)

const CreatePatchCommentType = "revision.comment"

type ListOpts struct {
	Page    int
	PerPage int
}

type NodeInfo struct {
	ID     string     `json:"id"`
	Config NodeConfig `json:"config"`
}

type SessionInfo struct {
	SessionId string `json:"sessionId"`
	Status    string `json:"status"`
	PublicKey string `json:"publicKey"`
	Alias     string `json:"alias"`
	IssuedAt  int64  `json:"issuedAt"`
	ExpiresAt int64  `json:"expiresAt"`
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

type ProjectFile struct {
	Binary  bool   `json:"binary"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Path    string `json:"path"`
}

type FileTree struct {
	Entries []FileTreeEntries `json:"entries"`
}

type FileTreeEntries struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

type Patch struct {
	ID    string     `json:"id"`
	Title string     `json:"title"`
	State PatchState `json:"state"`
}

type PatchState struct {
	Status string `json:"status"`
}

type CreatePatchComment struct {
	Type     string `json:"type"`
	Body     string `json:"body"`
	Revision string `json:"revision"`
}

func (o *ListOpts) Encode() string {
	params := url.Values{}
	if o.Page > 0 {
		page := o.Page - 1
		params.Set("page", strconv.Itoa(page))
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
