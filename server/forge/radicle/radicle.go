package radicle

import (
	"context"
	"errors"
	"fmt"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge"
	types "go.woodpecker-ci.org/woodpecker/v2/server/forge/radicle/hooks"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge/radicle/internal"
	forge_types "go.woodpecker-ci.org/woodpecker/v2/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v2/server/model"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Opts defines configuration options.
type Opts struct {
	URL         string // Radicle node url.
	NodeID      string // Radicle NID
	SecretToken string // Radicle secret token.
}

// radicle implements "Forge" interface
type radicle struct {
	url         string
	nodeID      string
	secretToken string
}

// New returns a new forge Configuration for integrating with the Radicle
// repository hosting service at https://radicle.xyz
func New(opts Opts) (forge.Forge, error) {
	rad := radicle{
		url:         opts.URL,
		nodeID:      opts.NodeID,
		secretToken: opts.SecretToken,
	}
	rad.url = strings.TrimSuffix(opts.URL, "/")
	if len(rad.url) == 0 {
		return nil, fmt.Errorf("must provide a URL")
	}
	if len(rad.nodeID) == 0 {
		return nil, fmt.Errorf("must provide a NID")
	}
	_, err := url.Parse(rad.url)
	if err != nil {
		return nil, fmt.Errorf("must provide a valid URL: %s", err)
	}
	if len(rad.secretToken) == 0 {
		return nil, fmt.Errorf("must provide a token")
	}
	return &rad, nil
}

// Name returns the string name of this driver
func (rad *radicle) Name() string {
	return "radicle"
}

// URL returns the root url of a configured forge
func (rad *radicle) URL() string {
	return rad.url
}

// NID returns the node ID of the of a configured radicle forge
func (rad *radicle) NID() string {
	return rad.nodeID
}

// Login authenticates the session and returns the
// forge user details.
func (rad *radicle) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	fmt.Println("Called Login")
	client := internal.NewClient(ctx, rad.url, rad.secretToken)
	fmt.Println(fmt.Sprintf("%+v", rad.url))
	sessionInfo, err := client.GetSessionInfo()
	if err != nil {
		return nil, err
	}
	if sessionInfo.Status != internal.AUTHORIZED_SESSION {
		return nil, errors.New("provided secret token is unauthorized")
	}
	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		return nil, err
	}
	return convertUser(nodeInfo), nil
}

// Auth authenticates the session and returns the forge user
// login for the given token and secret
func (rad *radicle) Auth(ctx context.Context, token, secret string) (string, error) {
	fmt.Println("Called Auth")
	// Auth is not used by Radicle as there is no oAuth process
	panic("implement me")
}

// Teams fetches a list of team memberships from the forge.
func (rad *radicle) Teams(ctx context.Context, u *model.User) ([]*model.Team, error) {
	fmt.Println("Called Teams")
	//Radicle does not support teams, workspaces or organizations.
	return nil, nil
}

// Repo fetches the repository from the forge, preferred is using the ID, fallback is owner/name.
func (rad *radicle) Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error) {
	fmt.Println("Called Repo")
	fmt.Println("rad url: " + rad.URL())

	fmt.Println(fmt.Sprintf("%+v", *u))
	fmt.Println("remoteID: " + string(remoteID))
	fmt.Println("owner: " + owner)
	fmt.Println("name: " + name)
	if remoteID.IsValid() {
		name = string(remoteID)
	}
	client := internal.NewClient(ctx, rad.url, rad.secretToken)
	project, err := client.GetProject(name)
	if err != nil {
		return nil, err
	}
	return convertProject(project, u, rad), nil
}

// Repos fetches a list of repos from the forge.
func (rad *radicle) Repos(ctx context.Context, u *model.User) ([]*model.Repo, error) {
	fmt.Println("Called Repos")
	fmt.Println(fmt.Sprintf("%+v", *u))
	client := internal.NewClient(ctx, rad.url, rad.secretToken)
	projects, err := client.GetProjects()
	if err != nil {
		return nil, err
	}
	repos := []*model.Repo{}
	for _, project := range projects {
		repos = append(repos, convertProject(project, u, rad))
	}
	return repos, nil
}

// File fetches a file from the forge repository and returns it in string
// format.
func (rad *radicle) File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, f string) ([]byte, error) {
	fmt.Println("Called File")
	fmt.Println(f)
	client := internal.NewClient(ctx, rad.url, rad.secretToken)
	projectFile, err := client.GetProjectCommitFile(string(r.ForgeRemoteID), b.Commit, f)
	if err != nil {
		return nil, err
	}
	return convertProjectFileToContent(projectFile)
}

// Dir fetches a folder from the forge repository
func (rad *radicle) Dir(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline,
	f string) ([]*forge_types.FileMeta, error) {
	fmt.Println("Called Dir")
	fmt.Println(f)
	client := internal.NewClient(ctx, rad.url, rad.secretToken)
	fileContents, err := client.GetProjectCommitDir(string(r.ForgeRemoteID), b.Commit, f)
	if err != nil {
		return nil, err
	}
	filesMeta := []*forge_types.FileMeta{}
	for _, fileContentEntry := range fileContents.Entries {
		fileContent := []byte{}
		if fileContentEntry.Kind == internal.FileTypeBlob {
			fileContent, err = rad.File(ctx, u, r, b, fileContentEntry.Path)
			if err != nil {
				return nil, err
			}
		}
		filesMeta = append(filesMeta, convertFileContent(fileContentEntry, fileContent))
	}
	return filesMeta, err
}

// Status sends the commit status to the forge.
// An example would be the GitHub pull request status.
func (rad *radicle) Status(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, p *model.Workflow) error {
	fmt.Println("Called Status")
	fmt.Println(b.Status)
	fmt.Println(fmt.Sprintf("%+v", b))
	fmt.Println(fmt.Sprintf("%+v", p))
	//TODO implement me - left with no error for testing purposes
	//will have to add a comment to commit when this is ready
	//comment := internal.CreatePatchComment{
	//	Type:     internal.CreatePatchCommentType,
	//	Body:     "",
	//	Revision: "",
	//}
	//client := internal.NewClient(ctx, rad.url, rad.secretToken)
	//err := client.AddProjectPatchComment()
	//if err != nil {
	//	return err
	//}
	return nil
}

// Netrc returns a .netrc file that can be used to clone
// private repositories from a forge.
func (rad *radicle) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	fmt.Println("Called Netrc")
	//Radicle does not currently support private repos, so there is no need to implement this.
	// Return a dummy Netrc model.
	return &model.Netrc{
		Machine:  rad.URL(),
		Login:    rad.NID(),
		Password: "",
	}, nil
}

// Activate activates a repository by creating the post-commit hook.
func (rad *radicle) Activate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	fmt.Println("Called Activate")
	fmt.Println(fmt.Sprintf("%+v", r))
	fmt.Println("Link: " + link)
	//TODO implement me
	//Added as successful in order to test the rest of the procedure
	return nil
}

// Deactivate deactivates a repository by removing all previously created
// post-commit hooks matching the given link.
func (rad *radicle) Deactivate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	fmt.Println("Called Deactivate")
	fmt.Println(fmt.Sprintf("%+v", r))
	fmt.Println("Link: " + link)

	//TODO implement me
	//Added as successful in order to test the rest of the procedure
	return nil
}

// Branches returns the names of all branches for the named repository.
func (rad *radicle) Branches(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]string, error) {
	// Radicle announces only defaultBranch, so no other branch is globally accessible
	fmt.Println("Called Branches")
	if p.Page > 1 {
		return []string{}, nil
	}
	return []string{r.Branch}, nil
}

// BranchHead returns the sha of the head (latest commit) of the specified branch
func (rad *radicle) BranchHead(ctx context.Context, _ *model.User, r *model.Repo, branch string) (string, error) {
	fmt.Println("Called BranchHead")
	if r.Branch != branch {
		return "", errors.New("branch does not exist")
	}
	listOpts := internal.ListOpts{Page: 0, PerPage: 1}
	client := internal.NewClient(ctx, rad.url, rad.secretToken)
	branchCommits, err := client.GetProjectCommits(string(r.ForgeRemoteID), listOpts)
	if err != nil {
		return "", err
	}
	if branchCommits.Stats.Commits == 0 || len(branchCommits.Commits) == 0 {
		return "", errors.New("branch has no commits")
	}
	return branchCommits.Commits[0].Commit.ID, err
}

// PullRequests returns all pull requests for the named repository.
func (rad *radicle) PullRequests(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]*model.PullRequest, error) {
	fmt.Println("Called PullRequests")
	listOpts := internal.ListOpts{Page: p.Page, PerPage: p.PerPage}
	fmt.Println("Called PullRequests with opts" + fmt.Sprintf("%+v", p))
	client := internal.NewClient(ctx, rad.url, rad.secretToken)
	projectPatches, err := client.GetProjectPatches(string(r.ForgeRemoteID), listOpts)
	if err != nil {
		return nil, err
	}
	pullRequests := []*model.PullRequest{}
	for _, projectPatch := range projectPatches {
		pullRequests = append(pullRequests, convertProjectPatch(projectPatch))

	}
	return pullRequests, err
}

// Hook parses the post-commit hook from the Request body and returns the
// required data in a standard format.
func (rad *radicle) Hook(ctx context.Context, r *http.Request) (repo *model.Repo, pipeline *model.Pipeline, err error) {
	fmt.Println("Called Hook")
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}

	hookType := r.Header.Get(types.EventTypeHeaderKey)
	switch hookType {
	case types.EventTypePush:
		return rad.parsePushHook(payload)
	case types.EventTypePatch:
		return rad.parsePatchHook(payload)
	default:
		return nil, nil, &forge_types.ErrIgnoreEvent{Event: hookType}
	}
}

// OrgMembership returns if user is member of organization and if user
// is admin/owner in that organization.
func (rad *radicle) OrgMembership(_ context.Context, u *model.User, orgName string) (*model.OrgPerm, error) {
	fmt.Println("Called OrgMembership")
	// Radicle does not currently support Orgs, so return membership as org Admin if its user's Org.
	if orgName != u.Login {
		return &model.OrgPerm{
			Member: false,
			Admin:  false,
		}, nil
	}
	return &model.OrgPerm{
		Member: true,
		Admin:  true,
	}, nil
}

// Org fetches the organization from the forge by name. If the name is a user an org with type user is returned.
func (rad *radicle) Org(ctx context.Context, u *model.User, org string) (*model.Org, error) {
	fmt.Println("Called Org")
	// Radicle does not currently support Orgs, so return user as individual org.
	return &model.Org{
		Name:   u.Login,
		IsUser: true,
	}, nil
}
