package radicle

import (
	"github.com/franela/goblin"
	"go.woodpecker-ci.org/woodpecker/v2/server/forge/radicle/internal"
	"go.woodpecker-ci.org/woodpecker/v2/server/model"
	"testing"
)

func Test_convert(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Radicle converter", func() {

		project := &internal.Repository{
			ID:            "the_radicle_id",
			Name:          "hello_world",
			DefaultBranch: "default_branch",
			Head:          "head_commit",
		}
		user := &model.User{
			ForgeRemoteID: "remote_user_id",
			Login:         "user_login",
			Token:         "user_token",
			Avatar:        "user_avatar",
			Admin:         true,
		}
		rad := &radicle{
			url:          "http://some.url",
			nodeID:       "the_nid",
			sessionToken: "the_token",
		}

		g.It("Should convert user with", func() {
			nodeInfo := &internal.NodeInfo{
				ID: "node_id",
				Config: internal.Node{
					Alias: "my_alias",
				},
			}
			to := convertUser(nodeInfo)
			g.Assert(to.ForgeRemoteID).Equal(model.ForgeRemoteID("node_id"))
			g.Assert(to.Login).Equal("my_alias")
		})

		g.It("Should convert repository with", func() {
			to := convertProject(project, user, rad)
			g.Assert(to.ForgeRemoteID).Equal(model.ForgeRemoteID("the_radicle_id"))
			g.Assert(to.FullName).Equal("user_login/hello_world")
			g.Assert(to.Owner).Equal("user_login")
			g.Assert(to.Name).Equal("hello_world")
			g.Assert(to.Branch).Equal("default_branch")
			g.Assert(to.ForgeURL).Equal("http://some.url/the_radicle_id")
			g.Assert(to.Clone).Equal("http://some.url/the_radicle_id.git")
			g.Assert(to.CloneSSH).Equal("")
			g.Assert(to.Perm.Push).IsTrue()
			g.Assert(to.Perm.Admin).IsTrue()
			g.Assert(to.Perm.Pull).IsTrue()
		})

		g.It("Should convert project file to content with", func() {
			projectFile := &internal.ProjectFile{
				Content: "some unicode content καλημέρα!",
			}
			to, err := convertProjectFileToContent(projectFile)
			g.Assert(err).IsNil()
			g.Assert(to).Equal([]byte("some unicode content καλημέρα!"))
		})

		g.It("Should convert file content with", func() {
			fileTreeEntries := internal.FileTreeEntries{
				Path: "/path/to/file",
				Name: "file name",
			}
			fileContent := []byte("some unicode content καλημέρα!")
			to := convertFileContent(fileTreeEntries, fileContent)
			g.Assert(to.Name).Equal("/path/to/file")
			g.Assert(to.Data).Equal([]byte("some unicode content καλημέρα!"))
		})

		g.It("Should convert project patch with", func() {
			patch := &internal.Patch{
				ID:    "patch_id",
				Title: "Patch title",
				State: internal.State{
					Status: "open",
				},
			}
			to := convertProjectPatch(patch)
			g.Assert(to.Title).Equal("Patch title")
			g.Assert(to.Index).Equal("patch_id")
		})

	})
}
