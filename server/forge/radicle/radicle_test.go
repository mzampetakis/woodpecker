package radicle

import (
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
)

func Test_radicle(t *testing.T) {
	gin.SetMode(gin.TestMode)

	g := goblin.Goblin(t)
	g.Describe("Radicle client", func() {

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
	})
}
