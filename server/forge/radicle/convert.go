package radicle

import (
	"fmt"
	"github.com/woodpecker-ci/woodpecker/server/forge/radicle/internal"
	"github.com/woodpecker-ci/woodpecker/server/model"
)

// convertUser is a helper function used to convert a Radicle Node Info structure
// to the Woodpecker User structure.
func convertUser(nodeInfo *internal.NodeInfo) *model.User {
	return &model.User{
		ForgeRemoteID: model.ForgeRemoteID(nodeInfo.ID),
		Login:         nodeInfo.Config.Alias,
	}
}

// convertProject is a helper function used to convert a Radicle Project structure
// to the Woodpecker Repo structure.
func convertProject(project *internal.Project, rad *radicle) *model.Repo {
	perm := model.Perm{
		Pull:  true,
		Push:  true,
		Admin: true,
	}
	return &model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(project.ID),
		Name:          project.Name,
		FullName:      fmt.Sprintf("%s/%s", rad.alias, project.Name),
		Link:          project.ID,
		Clone:         fmt.Sprintf("%s/%s", rad.URL(), project.ID),
		CloneSSH:      "",
		Branch:        project.DefaultBranch,
		Perm:          &perm,
	}
}
