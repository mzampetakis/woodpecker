package radicle

import (
	"context"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge/radicle/fixtures"
	"go.woodpecker-ci.org/woodpecker/v2/server/model"
	"net/http/httptest"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func Test_radicle(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := httptest.NewServer(fixtures.Handler())

	g := goblin.Goblin(t)
	g.Describe("Radicle client", func() {
		g.After(func() {
			s.Close()
		})

		// Test New()
		g.Describe("Creating new Forge", func() {
			g.It("should return an error when missing URL", func() {
				opts := Opts{
					URL:         "",
					NodeID:      "NodeID",
					SecretToken: "a_secret_token",
				}
				_, err := New(opts)
				g.Assert(err).IsNotNil()
			})
			g.It("Should return an error when invalid URL", func() {
				opts := Opts{
					URL:         "invalid_%url",
					NodeID:      "NodeID",
					SecretToken: "a_secret_token",
				}
				_, err := New(opts)
				g.Assert(err).IsNotNil()
			})
			g.It("Should return an error when missing token", func() {
				opts := Opts{
					URL:         "http://some.url",
					NodeID:      "NodeID",
					SecretToken: "",
				}
				_, err := New(opts)
				g.Assert(err).IsNotNil("Expected error")
			})
			g.It("Should return an error when missing Node ID", func() {
				opts := Opts{
					URL:         "http://some.url",
					NodeID:      "",
					SecretToken: "a_secret_token",
				}
				_, err := New(opts)
				g.Assert(err).IsNotNil("Expected error")
			})
			g.It("Should return a new Forge with correct data", func() {
				opts := Opts{
					URL:         "http://some.url",
					NodeID:      "NodeID",
					SecretToken: "a_secret_token",
				}
				forge, err := New(opts)
				g.Assert(err).IsNil()
				g.Assert(forge.URL()).Equal("http://some.url")
				g.Assert(forge.Name()).Equal("radicle")
			})
		})

		// Test Login()
		g.Describe("When logging in", func() {
			g.Describe("with non-existing session ID", func() {
				notFoundSessionForgeOpts.URL = s.URL
				forge, _ := New(notFoundSessionForgeOpts)
				g.It("Should fail", func() {
					user, err := forge.Login(context.Background(), nil, nil)
					g.Assert(err).IsNotNil()
					g.Assert(user).IsNil()
				})
			})
			g.Describe("with unauthorized sessions ID", func() {
				unauthorizedSessionForgeOpts.URL = s.URL
				forge, _ := New(unauthorizedSessionForgeOpts)
				g.It("Should fail", func() {
					user, err := forge.Login(context.Background(), nil, nil)
					g.Assert(err).IsNotNil()
					g.Assert(err.Error()).Equal("provided secret token is unauthorized")
					g.Assert(user).IsNil()
				})
			})
			g.Describe("with authorized sessions ID", func() {
				authorizedSessionForgeOpts.URL = s.URL
				forge, _ := New(authorizedSessionForgeOpts)
				g.It("Should succeed", func() {
					user, err := forge.Login(context.Background(), nil, nil)
					g.Assert(err).IsNil()
					g.Assert(user).IsNotNil()
					g.Assert(user.Login).Equal("myalias")
					g.Assert(user.ForgeRemoteID).Equal(model.ForgeRemoteID("someid"))
				})
			})
		})

		authorizedSessionForgeOpts.URL = s.URL
		forge, _ := New(authorizedSessionForgeOpts)

		// Test Repo()
		g.Describe("When requesting a project", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo, err := forge.Repo(context.Background(), user, "not_found", "", "")
					g.Assert(err).IsNotNil()
					g.Assert(repo).IsNil()
				})
			})
			g.Describe("with authorized sessions ID", func() {
				g.It("Should succeed", func() {
					repo, err := forge.Repo(context.Background(), user, "valid_project_id", "", "")
					g.Assert(err).IsNil()
					g.Assert(repo).IsNotNil()
					g.Assert(repo.ForgeRemoteID).Equal(model.ForgeRemoteID("valid_project_id"))
					g.Assert(repo.Name).Equal("a-project")
					g.Assert(repo.FullName).Equal("user_login/a-project")
					g.Assert(repo.ForgeURL).Equal(forge.URL() + "/valid_project_id")
					g.Assert(repo.Clone).Equal(forge.URL() + "/valid_project_id.git")
					g.Assert(repo.CloneSSH).Equal("")
					g.Assert(repo.Branch).Equal("main")
					g.Assert(repo.Perm.Pull).Equal(true)
					g.Assert(repo.Perm.Pull).Equal(true)
					g.Assert(repo.Perm.Admin).Equal(true)
					g.Assert(repo.Owner).Equal("user_login")
				})
			})
		})

		// Test Repos()
		g.Describe("When requesting projects", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			g.Describe("with many projects", func() {
				g.It("Should succeed", func() {
					repos, err := forge.Repos(context.Background(), user)
					g.Assert(err).IsNil()
					g.Assert(len(repos)).Equal(2)
					g.Assert(repos[0].ForgeRemoteID).Equal(model.ForgeRemoteID("valid_project_id"))
					g.Assert(repos[0].Name).Equal("a-project")
					g.Assert(repos[0].FullName).Equal("user_login/a-project")
					g.Assert(repos[0].ForgeURL).Equal(forge.URL() + "/valid_project_id")
					g.Assert(repos[0].Clone).Equal(forge.URL() + "/valid_project_id.git")
					g.Assert(repos[0].CloneSSH).Equal("")
					g.Assert(repos[0].Branch).Equal("main")
					g.Assert(repos[0].Perm.Pull).Equal(true)
					g.Assert(repos[0].Perm.Pull).Equal(true)
					g.Assert(repos[0].Perm.Admin).Equal(true)
					g.Assert(repos[0].Owner).Equal("user_login")

					g.Assert(repos[1].ForgeRemoteID).Equal(model.ForgeRemoteID("another_valid_project_id"))
					g.Assert(repos[1].Name).Equal("b-project")
					g.Assert(repos[1].FullName).Equal("user_login/b-project")
					g.Assert(repos[1].ForgeURL).Equal(forge.URL() + "/another_valid_project_id")
					g.Assert(repos[1].Clone).Equal(forge.URL() + "/another_valid_project_id.git")
					g.Assert(repos[1].CloneSSH).Equal("")
					g.Assert(repos[1].Branch).Equal("master")
					g.Assert(repos[1].Perm.Pull).Equal(true)
					g.Assert(repos[1].Perm.Pull).Equal(true)
					g.Assert(repos[1].Perm.Admin).Equal(true)
					g.Assert(repos[1].Owner).Equal("user_login")
				})
			})
		})

		// Test File()
		g.Describe("When requesting a project file", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			pipeline := &model.Pipeline{
				Commit: "the_commit_id",
			}
			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
					}
					file, err := forge.File(context.Background(), user, repo, pipeline, "file_name.md")
					g.Assert(err).IsNotNil()
					g.Assert(file).IsNil()
				})
			})
			g.Describe("with valid project ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "remote_project_id",
					}
					file, err := forge.File(context.Background(), user, repo, pipeline, "file_name.md")
					g.Assert(err).IsNil()
					g.Assert(file).IsNotNil()
					g.Assert(file).Equal([]byte("file content"))
				})
			})
		})

		// Test Tree()
		g.Describe("When requesting a project dir", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}
			pipeline := &model.Pipeline{
				Commit: "the_commit_id",
			}
			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
					}
					tree, err := forge.Dir(context.Background(), user, repo, pipeline, "")
					g.Assert(err).IsNotNil()
					g.Assert(tree).IsNil()
				})
			})
			g.Describe("with valid project ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "remote_project_id",
					}
					tree, err := forge.Dir(context.Background(), user, repo, pipeline, "")
					g.Assert(err).IsNil()
					g.Assert(tree).IsNotNil()
					g.Assert(len(tree)).Equal(3)
					g.Assert(tree[0].Name).Equal("Readme.md")
					g.Assert(tree[0].Data).Equal([]byte("file content"))
					g.Assert(tree[1].Name).Equal("cargo-checksum.json")
					g.Assert(tree[1].Data).Equal([]byte("file content"))
					g.Assert(tree[2].Name).Equal("debian")
					g.Assert(tree[2].Data).Equal([]byte(""))
				})
			})
		})

		// Test BranchHead()
		g.Describe("When requesting a project's BranchHead", func() {
			user := &model.User{
				ForgeRemoteID: "remote_user_id",
				Login:         "user_login",
			}

			g.Describe("with invalid project ID", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "not_found",
						Branch:        "main",
					}
					branchHead, err := forge.BranchHead(context.Background(), user, repo, "main")
					g.Assert(err).IsNotNil()
					g.Assert(len(branchHead)).Equal(0)
				})
			})
			g.Describe("with valid project ID", func() {
				g.It("Should succeed", func() {
					repo := &model.Repo{
						ForgeRemoteID: "valid_project_id",
						Branch:        "main",
					}
					branchHead, err := forge.BranchHead(context.Background(), user, repo, "main")
					g.Assert(err).IsNil()
					g.Assert(branchHead).IsNotNil()
					g.Assert(branchHead).Equal("00bfa9b18be32001481334126c311c4a327dff2e")
				})
			})
			g.Describe("with valid project ID and invalid branch", func() {
				g.It("Should fail", func() {
					repo := &model.Repo{
						ForgeRemoteID: "valid_project_id",
						Branch:        "invalid_branch",
					}
					branchHead, err := forge.BranchHead(context.Background(), user, repo, "main")
					g.Assert(err).IsNotNil()
					g.Assert(err.Error()).Equal("branch does not exist")
					g.Assert(len(branchHead)).Equal(0)
				})
			})
		})

	})
}

var authorizedSessionForgeOpts = Opts{
	URL:         "http://node.id",
	NodeID:      "NodeID",
	SecretToken: "authorized",
}

var notFoundSessionForgeOpts = Opts{
	URL:         "http://node.id",
	NodeID:      "NodeID",
	SecretToken: "not_found",
}

var unauthorizedSessionForgeOpts = Opts{
	URL:         "http://node.id",
	NodeID:      "NodeID",
	SecretToken: "unauthorized_session",
}
