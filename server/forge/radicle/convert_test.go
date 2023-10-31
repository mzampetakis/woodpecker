package radicle

import (
	"github.com/franela/goblin"
	"github.com/woodpecker-ci/woodpecker/server/forge/radicle/internal"
	"github.com/woodpecker-ci/woodpecker/server/model"
	"testing"
)

func Test_helper(t *testing.T) {
	g := goblin.Goblin(t)
	g.Describe("Radicle converter", func() {

		g.It("Should convert repository with", func() {
			from := &internal.Project{
				ID:            "the_radicle_id",
				Name:          "hello_world",
				DefaultBranch: "default_branch",
				Head:          "head_commit",
			}
			rad := &radicle{
				url:         "http://some.url",
				nodeID:      "the_nid",
				alias:       "node_alias",
				secretToken: "the_token",
			}
			to := convertProject(from, rad)
			g.Assert(to.ForgeRemoteID).Equal(model.ForgeRemoteID("the_radicle_id"))
			g.Assert(to.FullName).Equal("node_alias/hello_world")
			g.Assert(to.Owner).Equal("the_nid")
			g.Assert(to.Name).Equal("the_radicle_id")
			g.Assert(to.Branch).Equal("default_branch")
			g.Assert(to.Link).Equal("http://some.url/the_radicle_id")
			g.Assert(to.Clone).Equal("http://some.url/the_radicle_id hello_world")
			g.Assert(to.Perm.Push).IsTrue()
			g.Assert(to.Perm.Admin).IsTrue()
			g.Assert(to.Perm.Pull).IsTrue()
		})
	})
}
