package radicle

import (
	"encoding/json"
	"fmt"
	types "go.woodpecker-ci.org/woodpecker/v2/server/forge/radicle/hooks"
	"go.woodpecker-ci.org/woodpecker/v2/server/model"
)

func (rad *radicle) parsePushHook(payload []byte) (*model.Repo, *model.Pipeline, error) {
	hook := types.PushPayload{}
	if err := json.Unmarshal(payload, &hook); err != nil {
		return nil, nil, err
	}

	perm := model.Perm{
		Pull:  true,
		Push:  true,
		Admin: true,
	}
	repo := model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(hook.Repository.ID),
		Owner:         "",
		Name:          fmt.Sprintf("%s (%s)", hook.Repository.Name, hook.Repository.ID),
		FullName:      fmt.Sprintf("%s (%s)", hook.Repository.Name, hook.Repository.ID),
		Avatar:        "",
		ForgeURL:      hook.Repository.URL,
		Clone:         hook.Repository.CloneURL,
		CloneSSH:      "",
		Branch:        hook.Repository.DefaultBranch,
		Hash:          hook.Repository.ID,
		Perm:          &perm,
	}

	lastCommit := hook.Commits[len(hook.Commits)-1]
	changedFiles := []string{}
	for _, commit := range hook.Commits {
		changedFiles = append(changedFiles, commit.Modified...)
		changedFiles = append(changedFiles, commit.Added...)
		changedFiles = append(changedFiles, commit.Removed...)
	}
	pipeline := model.Pipeline{
		Author:   hook.Author.ID,
		Event:    model.EventPush,
		Commit:   hook.After,
		Branch:   hook.After,
		Ref:      fmt.Sprintf("refs/heads/%s", hook.After),
		Refspec:  "",
		CloneURL: "",
		Title:    lastCommit.Title,
		Message:  lastCommit.Message,
		//Timestamp:           lastCommit.Timestamp,
		Sender:       lastCommit.Author.Name,
		Email:        lastCommit.Author.Email,
		ForgeURL:     lastCommit.URL,
		ChangedFiles: changedFiles,
	}

	return &repo, &pipeline, nil
}

func (rad *radicle) parsePatchHook(payload []byte) (*model.Repo, *model.Pipeline, error) {
	hook := types.PatchPayload{}
	if err := json.Unmarshal(payload, &hook); err != nil {
		return nil, nil, err
	}

	perm := model.Perm{
		Pull:  true,
		Push:  true,
		Admin: true,
	}
	repo := model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(hook.Repository.ID),
		Owner:         "",
		Name:          fmt.Sprintf("%s (%s)", hook.Repository.Name, hook.Repository.ID),
		FullName:      fmt.Sprintf("%s (%s)", hook.Repository.Name, hook.Repository.ID),
		Avatar:        "",
		ForgeURL:      hook.Repository.URL,
		Clone:         hook.Repository.CloneURL,
		CloneSSH:      "",
		Branch:        hook.Repository.DefaultBranch,
		Hash:          hook.Repository.ID,
		Perm:          &perm,
	}

	lastRevision := hook.Patch.Revisions[len(hook.Patch.Revisions)-1]
	pipeline := model.Pipeline{
		Author:   hook.Patch.Author.ID,
		Event:    model.EventPull,
		Commit:   hook.Patch.After,
		Branch:   lastRevision.ID,
		Ref:      fmt.Sprintf("refs/heads/%s", hook.Patch.After),
		Refspec:  fmt.Sprintf("%s:%s", hook.Patch.After, hook.Patch.Target),
		CloneURL: "",
		Title:    hook.Patch.Title,
		Message:  "",
		//Timestamp:           lastRevision.Timestamp,
		Sender:       hook.Patch.Author.ID,
		Email:        "",
		ForgeURL:     hook.Patch.URL,
		ChangedFiles: nil,
	}

	return &repo, &pipeline, nil
}
