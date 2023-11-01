package radicle

import (
	"fmt"
	"github.com/woodpecker-ci/woodpecker/server/forge/radicle/internal"
	"github.com/woodpecker-ci/woodpecker/server/model"
	"strings"
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
func convertProject(project *internal.Project, user *model.User, rad *radicle) *model.Repo {
	perm := model.Perm{
		Pull:  true,
		Push:  true,
		Admin: true,
	}
	projectID := strings.TrimPrefix(project.ID, "rad:")
	return &model.Repo{
		ForgeRemoteID: model.ForgeRemoteID(projectID),
		Name:          project.Name,
		FullName:      fmt.Sprintf("%s/%s", user.Login, project.Name),
		Link:          fmt.Sprintf("%s/%s", rad.URL(), projectID),
		Clone:         fmt.Sprintf("%s/%s.git %s", rad.URL(), projectID, project.Name),
		CloneSSH:      "",
		Branch:        project.DefaultBranch,
		Perm:          &perm,
		Owner:         user.Login,
	}
}
