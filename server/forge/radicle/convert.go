package radicle

import (
	"github.com/woodpecker-ci/woodpecker/server/forge/radicle/internal"
	"github.com/woodpecker-ci/woodpecker/server/model"
)

// convertUser is a helper function used to convert a Radicle Node Info structure
// to the Woodpecker User structure.
func convertUser(from *internal.NodeInfo) *model.User {
	return &model.User{
		Login:         from.Node.ID,
		ForgeRemoteID: model.ForgeRemoteID(from.Node.ID),
	}
}
