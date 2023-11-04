package radicle

import (
	"fmt"
	"github.com/woodpecker-ci/woodpecker/server/forge/radicle/internal"
	forge_types "github.com/woodpecker-ci/woodpecker/server/forge/types"
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

// convertProjectFileToContent is a helper function used to convert a Radicle Project File content
// to the Woodpecker file content structure.
func convertProjectFileToContent(projectFile *internal.ProjectFile) ([]byte, error) {
	return projectFile.Content, nil
}

// convertFileContent is a helper function used to convert a Radicle Project file Contents
// to the Woodpecker File Meta structure.
func convertFileContent(fileContentEntries internal.FileTreeEntries, fileContent []byte) *forge_types.FileMeta {
	return &forge_types.FileMeta{
		Name: fileContentEntries.Path, // Might need .Name
		Data: fileContent,
	}
}
