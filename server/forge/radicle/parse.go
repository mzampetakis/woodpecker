package radicle

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge/radicle/internal"
	"go.woodpecker-ci.org/woodpecker/v2/server/model"
)

// parsePushHook parses the push hook payload
func (rad *radicle) parsePushHook(payload []byte) (*model.Repo, *model.Pipeline, error) {
	hook := internal.ΗοοκPushPayload{}
	if err := json.Unmarshal(payload, &hook); err != nil {
		return nil, nil, err
	}
	repo := convertProject(&hook.Repository, nil, rad)
	lastCommit := internal.Commit{}
	if len(hook.Commits) == 0 {
		return nil, nil, errors.New("no commits found in push")
	}
	lastCommit = hook.Commits[len(hook.Commits)-1]
	changedFiles := []string{}
	for _, commit := range hook.Commits {
		changedFiles = append(changedFiles, commit.Modified...)
		changedFiles = append(changedFiles, commit.Added...)
		changedFiles = append(changedFiles, commit.Removed...)
	}
	pipeline := model.Pipeline{
		Author:       hook.Author.ID,
		Event:        model.EventPush,
		Commit:       hook.After,
		Branch:       hook.After,
		Ref:          fmt.Sprintf("refs/heads/%s", hook.After),
		Refspec:      fmt.Sprintf("refs/heads/%s:refs/heads/%s", hook.After, hook.After),
		CloneURL:     "",
		Avatar:       RADICLE_IMAGE,
		Message:      lastCommit.Title,
		Timestamp:    lastCommit.Timestamp,
		Sender:       lastCommit.Author.Name,
		Email:        lastCommit.Author.Email,
		ForgeURL:     lastCommit.URL,
		ChangedFiles: changedFiles,
	}

	return repo, &pipeline, nil
}

// parsePatchHook parses the patch hook payload
func (rad *radicle) parsePatchHook(payload []byte) (*model.Repo, *model.Pipeline, error) {
	hook := internal.ΗοοκPatchPayload{}
	if err := json.Unmarshal(payload, &hook); err != nil {
		return nil, nil, err
	}
	repo := convertProject(&hook.Repository, nil, rad)
	if len(hook.Patch.Revisions) == 0 {
		return nil, nil, errors.New("no revision found in patch")
	}
	lastRevision := hook.Patch.Revisions[len(hook.Patch.Revisions)-1]
	changedFiles := []string{}
	for _, commit := range hook.Patch.Commits {
		changedFiles = append(changedFiles, commit.Modified...)
		changedFiles = append(changedFiles, commit.Added...)
		changedFiles = append(changedFiles, commit.Removed...)
	}
	vars := map[string]string{}
	vars["patch_id"] = hook.Patch.ID
	vars["revision_id"] = lastRevision.ID
	pipeline := model.Pipeline{
		Author:              hook.Patch.Author.Alias,
		Event:               model.EventPull,
		Commit:              hook.Patch.After,
		Branch:              hook.Patch.ID,
		Ref:                 fmt.Sprintf("refs/heads/%s", hook.Patch.After),
		Refspec:             fmt.Sprintf("refs/heads/%s:refs/heads/%s", hook.Patch.After, hook.Patch.After),
		CloneURL:            "",
		Avatar:              RADICLE_IMAGE,
		PullRequestLabels:   hook.Patch.Labels,
		Message:             hook.Patch.Title,
		Timestamp:           lastRevision.Timestamp,
		Sender:              hook.Patch.Author.ID,
		Email:               hook.Patch.Author.Alias,
		ForgeURL:            hook.Patch.URL,
		ChangedFiles:        changedFiles,
		AdditionalVariables: vars,
	}

	return repo, &pipeline, nil
}
