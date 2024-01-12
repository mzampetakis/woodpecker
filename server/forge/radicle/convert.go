package radicle

import (
	"fmt"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge/radicle/internal"
	forge_types "go.woodpecker-ci.org/woodpecker/v2/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v2/server/model"
	"strings"
)

// As radicle does not support user avatars, use radicle's logo as user avatar
const RADICLE_IMAGE = "data:image/png;base64," +
	"iVBORw0KGgoAAAANSUhEUgAAACwAAAAsBAMAAADsqkcyAAAAElBMVEUAAAAzM91VVf/09PT/Vf////+iehdrAAAAAXRSTlMAQObYZgAAAAFiS0dEBfhv6ccAAABTSURBVCjPY2AAAiUlBjhAYlNFGMFBEaSasKAgBFNbWAkFUFMYxDQGAgRNbWEXFwRNLWFBrIA6wgxg8yGRC8JQoSEhDEuUaGmKSsKQZIOSpigXBgAOHTr5ND3M6gAAAABJRU5ErkJggg=="

// convertUser is a helper function used to convert a Radicle Node Info structure
// to the Woodpecker User structure.
func convertUser(nodeInfo *internal.NodeInfo) *model.User {
	return &model.User{
		ForgeRemoteID: model.ForgeRemoteID(nodeInfo.ID),
		Login:         nodeInfo.Config.Alias,
		Avatar:        RADICLE_IMAGE,
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
		Name:          fmt.Sprintf("%s (%s)", project.Name, project.ID),
		FullName:      fmt.Sprintf("%s (%s)", project.Name, project.ID),
		ForgeURL:      fmt.Sprintf("%s/%s", rad.URL(), projectID),
		Clone:         fmt.Sprintf("%s/%s.git", rad.URL(), projectID),
		CloneSSH:      "",
		Branch:        project.DefaultBranch,
		Perm:          &perm,
		Owner:         user.Login,
		SCMKind:       model.RepoGit,
		PREnabled:     true,
	}
}

// convertProjectFileToContent is a helper function used to convert a Radicle Project File content
// to the Woodpecker file content structure.
func convertProjectFileToContent(projectFile *internal.ProjectFile) ([]byte, error) {
	return []byte(projectFile.Content), nil
}

// convertFileContent is a helper function used to convert a Radicle Project file Contents
// to the Woodpecker File Meta structure.
func convertFileContent(fileContentEntries internal.FileTreeEntries, fileContent []byte) *forge_types.FileMeta {
	return &forge_types.FileMeta{
		Name: fileContentEntries.Path, // Might need .Name
		Data: fileContent,
	}
}

// convertProjectPatch is a helper function used to convert a Radicle patch
// to the Woodpecker Pull Request structure.
func convertProjectPatch(patch *internal.Patch) *model.PullRequest {
	return &model.PullRequest{
		//Index: patch.ID,
		Title: patch.Title,
	}
}
