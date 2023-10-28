package radicle

import (
	"context"
	"github.com/woodpecker-ci/woodpecker/server/forge"
	forge_types "github.com/woodpecker-ci/woodpecker/server/forge/types"
	"github.com/woodpecker-ci/woodpecker/server/model"
	"net/http"
)

// Opts defines configuration options.
type Opts struct {
	URL         string // Radicle node url.
	SecretToken string // Radicle secret token.
}

// radicle implements "Forge" interface
type radicle struct {
	url         string
	secretToken string
}

// New returns a new forge Configuration for integrating with the Radicle
// repository hosting service at https://radicle.xyz
func New(opts Opts) (forge.Forge, error) {
	return &radicle{
		url:         opts.URL,
		secretToken: opts.SecretToken,
	}, nil
}

// Name returns the string name of this driver
func (rad radicle) Name() string {
	return "radicle"
}

// URL returns the root url of a configured forge
func (rad radicle) URL() string {
	return rad.url
}

// Login authenticates the session and returns the
// forge user details.
func (rad radicle) Login(ctx context.Context, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

// Auth authenticates the session and returns the forge user
// login for the given token and secret
func (rad radicle) Auth(ctx context.Context, token, secret string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// Teams fetches a list of team memberships from the forge.
func (rad radicle) Teams(ctx context.Context, u *model.User) ([]*model.Team, error) {
	//TODO implement me
	panic("implement me")
}

// Repo fetches the repository from the forge, preferred is using the ID, fallback is owner/name.
func (rad radicle) Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error) {
	//TODO implement me
	panic("implement me")
}

// Repos fetches a list of repos from the forge.
func (rad radicle) Repos(ctx context.Context, u *model.User) ([]*model.Repo, error) {
	//TODO implement me
	panic("implement me")
}

// File fetches a file from the forge repository and returns it in string
// format.
func (rad radicle) File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, f string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

// Dir fetches a folder from the forge repository
func (rad radicle) Dir(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, f string) ([]*forge_types.FileMeta, error) {
	//TODO implement me
	panic("implement me")
}

// Status sends the commit status to the forge.
// An example would be the GitHub pull request status.
func (rad radicle) Status(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, p *model.Workflow) error {
	//TODO implement me
	panic("implement me")
}

// Netrc returns a .netrc file that can be used to clone
// private repositories from a forge.
func (rad radicle) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	//TODO implement me
	panic("implement me")
}

// Activate activates a repository by creating the post-commit hook.
func (rad radicle) Activate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	//TODO implement me
	panic("implement me")
}

// Deactivate deactivates a repository by removing all previously created
// post-commit hooks matching the given link.
func (rad radicle) Deactivate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	//TODO implement me
	panic("implement me")
}

// Branches returns the names of all branches for the named repository.
func (rad radicle) Branches(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

// BranchHead returns the sha of the head (latest commit) of the specified branch
func (rad radicle) BranchHead(ctx context.Context, u *model.User, r *model.Repo, branch string) (string, error) {
	//TODO implement me
	panic("implement me")
}

// PullRequests returns all pull requests for the named repository.
func (rad radicle) PullRequests(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]*model.PullRequest, error) {
	//TODO implement me
	panic("implement me")
}

// Hook parses the post-commit hook from the Request body and returns the
// required data in a standard format.
func (rad radicle) Hook(ctx context.Context, r *http.Request) (repo *model.Repo, pipeline *model.Pipeline, err error) {
	//TODO implement me
	panic("implement me")
}

// OrgMembership returns if user is member of organization and if user
// is admin/owner in that organization.
func (rad radicle) OrgMembership(ctx context.Context, u *model.User, org string) (*model.OrgPerm, error) {
	//TODO implement me
	panic("implement me")
}

// Org fetches the organization from the forge by name. If the name is a user an org with type user is returned.
func (rad radicle) Org(ctx context.Context, u *model.User, org string) (*model.Org, error) {
	//TODO implement me
	panic("implement me")
}
