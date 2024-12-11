package radicle

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge/common"
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
	URL        string // Radicle node url.
	NodeID     string // Radicle NID
	Host       string // Woodpecker hostUrl URL
	HookSecret string // secret used for signature generation
}

// radicle implements "Forge" interface
type radicle struct {
	url          string
	hostUrl      string
	nodeID       string
	sessionToken string
	alias        string
	hookSecret   string
}

// New returns a new forge Configuration for integrating with the Radicle
// repository hosting service at https://radicle.xyz
func New(opts Opts) (forge.Forge, error) {
	fmt.Println("Called New")
	rad := radicle{
		url:        opts.URL,
		hookSecret: opts.HookSecret,
	}
	rad.url = strings.TrimSuffix(opts.URL, "/")
	if len(rad.url) == 0 {
		return nil, fmt.Errorf("must provide a URL")
	}
	_, err := url.Parse(rad.url)
	if err != nil {
		return nil, fmt.Errorf("must provide a valid URL: %s", err)
	}

	rad.hostUrl = strings.TrimSuffix(opts.Host, "/")
	if len(rad.url) == 0 {
		return nil, fmt.Errorf("must provide a hostUrl URL")
	}
	_, err = url.Parse(rad.hostUrl)
	if err != nil {
		return nil, fmt.Errorf("must provide a valid hostUrl URL: %s", err)
	}

	return &rad, nil
}

// Name returns the string name of this driver
func (rad *radicle) Name() string {
	fmt.Println("Called Name")
	return "radicle"
}

// URL returns the root url of a configured forge
func (rad *radicle) URL() string {
	fmt.Println("Called URL")
	return rad.url
}

// NID returns the node ID of the of a configured radicle forge
func (rad *radicle) NID() string {
	fmt.Println("Called NID")
	return rad.nodeID
}

// Login authenticates the session and returns the
// forge user details.
func (rad *radicle) Login(ctx context.Context, _ *forge_types.OAuthRequest) (*model.User, string, error) {
	fmt.Println("Called Login")

	loginURL := internal.GetLoginURL(rad.url, rad.hostUrl)
	rad.sessionToken = ""
	ginCtx, ok := ctx.(*gin.Context)
	if ok {
		if len(ctx.(*gin.Context).Request.FormValue("session_id")) > 0 {
			rad.sessionToken = ginCtx.Request.FormValue("session_id")
		}
	}
	if len(rad.sessionToken) == 0 {
		fmt.Println("Session Token is empty")
		return nil, loginURL, nil
	}
	client := internal.NewClient(ctx, rad.url, rad.sessionToken)
	sessionInfo, err := client.GetSessionInfo()
	if err != nil {
		rad.sessionToken = ""
		return nil, loginURL, err
	}
	if sessionInfo.Status != internal.SessionStatusAuthorized {
		rad.sessionToken = ""
		return nil, loginURL, errors.New("provided secret token is unauthorized")
	}
	nodeInfo, err := client.GetNodeInfo()
	if err != nil {
		rad.sessionToken = ""
		return nil, loginURL, err
	}
	rad.nodeID = nodeInfo.Config.ID
	rad.alias = sessionInfo.Alias
	return convertUser(rad), loginURL, nil
}

// Auth authenticates the session and returns the forge user
// login for the given token and secret
func (rad *radicle) Auth(_ context.Context, _, _ string) (string, error) {
	// Auth is not used by Radicle as there is no oAuth process
	fmt.Println("Called Auth")
	return "", nil
}

// Teams fetches a list of team memberships from the forge.
func (rad *radicle) Teams(_ context.Context, _ *model.User) ([]*model.Team, error) {
	fmt.Println("Called Teams")
	//Radicle does not support teams, workspaces or organizations.
	return nil, nil
}

// Repo fetches the repository from the forge, preferred is using the ID, fallback is owner/name.
func (rad *radicle) Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error) {
	fmt.Println("Called Repo")
	if remoteID.IsValid() {
		name = string(remoteID)
	}
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
	project, err := client.GetProject(name)
	if err != nil {
		return nil, err
	}
	return convertProject(project, u, rad), nil
}

// Repos fetches a list of repos from the forge.
func (rad *radicle) Repos(ctx context.Context, u *model.User) ([]*model.Repo, error) {
	fmt.Println("Called Repos")
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
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
func (rad *radicle) File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, f string) ([]byte,
	error) {
	fmt.Println("Called File")
	fmt.Println(f)
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
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
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
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
func (rad *radicle) Status(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline,
	_ *model.Workflow) error {
	fmt.Println("Called Status")
	//Do not add comment if pipeline in progress
	if b.Status == model.StatusPending || b.Status == model.StatusRunning {
		return nil
	}
	patchID, patchIDExists := b.AdditionalVariables["patch_id"]
	revisionID, revisionIDExists := b.AdditionalVariables["revision_id"]
	if !patchIDExists || !revisionIDExists {
		return errors.New("branch does not contain all required information for adding patch comment")
	}
	statusIcon := "⏳"
	if b.Status == model.StatusFailure || b.Status == model.StatusKilled || b.Status == model.StatusError || b.
		Status == model.StatusBlocked || b.Status == model.StatusDeclined {
		statusIcon = "❌"
	} else if b.Status == model.StatusSuccess {
		statusIcon = "✅"
	} else if b.Status == model.StatusSkipped {
		statusIcon = "↪️"
	}
	radicleComment := internal.CreatePatchComment{
		Type: internal.CreatePatchCommentType,
		Body: fmt.Sprintf("Pipeline #%v completed with result: %s. %s \n - Details: %s", b.ID, b.Status, statusIcon,
			common.GetPipelineStatusURL(r, b, nil)),
		Revision: revisionID,
	}
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
	err := client.AddProjectPatchComment(r.ForgeRemoteID, patchID, radicleComment)
	return err
}

// Netrc returns a .netrc file that can be used to clone
// private repositories from a forge.
func (rad *radicle) Netrc(_ *model.User, _ *model.Repo) (*model.Netrc, error) {
	fmt.Println("Called Netrc")
	//Radicle's private repos should be accessible through the node
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
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
	_, err := client.GetProject(string(r.ForgeRemoteID))
	if err != nil {
		return err
	}
	fmt.Println("Activate Repo: " + r.ForgeRemoteID)
	fmt.Println("Activate Link: " + link)
	webhookOpts := internal.RepoWebhook{
		RepoID:      string(r.ForgeRemoteID),
		URL:         link,
		Secret:      rad.hookSecret,
		ContentType: internal.AppJsonType,
	}
	return client.AddProjectWebhook(r.ForgeRemoteID, webhookOpts)
}

// Deactivate deactivates a repository by removing all previously created
// post-commit hooks matching the given link.
func (rad *radicle) Deactivate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	fmt.Println("Deactivate Repo: " + r.ForgeRemoteID)
	fmt.Println("Deactivate Link: " + link)
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
	hook, _ := client.GetProjectWebhook(r.ForgeRemoteID, link)
	return client.RemoveProjectWebhook(r.ForgeRemoteID, hook)
}

// Branches returns the names of all branches for the named repository.
func (rad *radicle) Branches(_ context.Context, _ *model.User, r *model.Repo, p *model.ListOptions) ([]string, error) {
	fmt.Println("Called Branches")
	// Radicle announces only defaultBranch, so no other branch is globally accessible
	if p.Page > 1 {
		return []string{}, nil
	}
	return []string{r.Branch}, nil
}

// BranchHead returns the sha of the head (latest commit) of the specified branch
func (rad *radicle) BranchHead(ctx context.Context, u *model.User, r *model.Repo, branch string) (*model.Commit,
	error) {
	fmt.Println("Called BranchHead")
	if r.Branch != branch {
		return nil, errors.New("branch does not exist")
	}
	listOpts := internal.ListOpts{Page: 0, PerPage: 1}
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
	branchCommits, err := client.GetProjectCommits(string(r.ForgeRemoteID), listOpts)
	if err != nil {
		return nil, err
	}
	if len(branchCommits) == 0 {
		return nil, errors.New("branch has no commits")
	}
	commit := model.Commit{
		SHA:      branchCommits[0].ID,
		ForgeURL: "",
	}
	return &commit, err
}

// PullRequests returns all pull requests for the named repository.
func (rad *radicle) PullRequests(ctx context.Context, u *model.User, r *model.Repo,
	p *model.ListOptions) ([]*model.PullRequest, error) {
	fmt.Println("Called PullRequests")
	listOpts := internal.ListOpts{Page: p.Page, PerPage: p.PerPage}
	client := internal.NewClient(ctx, rad.url, u.AccessToken)
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
func (rad *radicle) Hook(_ context.Context, r *http.Request) (repo *model.Repo, pipeline *model.Pipeline, err error) {
	fmt.Println("Called Hook")
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, nil, err
	}
	if len(rad.hookSecret) > 0 {
		signatureGot := r.Header.Get(internal.SignatureHashKey)
		signatureGenerated := generateHmacSignature(rad.hookSecret, payload)
		if signatureGot != signatureGenerated {
			return nil, nil, errors.New("invalid hook message signature")
		}
	}

	hookType := r.Header.Get(internal.EventTypeHeaderKey)
	switch hookType {
	case internal.EventTypePush:
		return rad.parsePushHook(payload)
	case internal.EventTypePatch:
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
func (rad *radicle) Org(_ context.Context, _ *model.User, _ string) (*model.Org, error) {
	fmt.Println("Called Org")
	// Radicle does not currently support Orgs, so return user as individual org.
	return &model.Org{
		Name:   rad.Name(),
		IsUser: true,
	}, nil
}
