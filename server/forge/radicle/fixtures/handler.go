package fixtures

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returns an http.Handler that is capable of handling a variety of mock
// Radicle requests and returning mock responses.
func Handler() http.Handler {
	gin.SetMode(gin.TestMode)

	e := gin.New()
	e.GET("/api/v1/node", getNodeInfo)
	e.GET("/api/v1/sessions/:session_id", getSession)
	e.GET("/api/v1/projects", getProjects)
	e.GET("/api/v1/projects/:project_id", getProject)
	e.GET("/api/v1/projects/:project_id/blob/:commit/:file", getProjectCommitFile)
	e.GET("/api/v1/projects/:project_id/tree/:commit/", getProjectCommitTree)
	e.GET("/api/v1/projects/:project_id/commits", getProjectCommits)
	return e
}

func getNodeInfo(c *gin.Context) {
	c.String(200, nodePayload)
}

func getSession(c *gin.Context) {
	switch c.Param("session_id") {
	case "not_found":
		c.String(404, notFound)
	case "unauthorized_session":
		c.String(200, sessionUnauthorizedPayload)
	default:
		c.String(200, sessionPayload)
	}
}

func getProjects(c *gin.Context) {
	switch c.Query("page") {
	case "0":
		c.String(200, projectsPayloadPage0)
	case "1":
		c.String(200, projectsPayloadPage1)
	default:
		c.String(200, emptyPayload)
	}
}

func getProject(c *gin.Context) {
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectPayload)
	}
}

func getProjectCommitFile(c *gin.Context) {
	if c.Param("commit") != "the_commit_id" {
		c.String(404, notFound)
		return
	}
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectCommitFilePayload)
	}
}

func getProjectCommitTree(c *gin.Context) {
	if c.Param("commit") != "the_commit_id" {
		c.String(404, notFound)
		return
	}
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectCommitTreePayload)
	}
}

func getProjectCommits(c *gin.Context) {
	switch c.Param("project_id") {
	case "not_found":
		c.String(404, notFound)
	default:
		c.String(200, projectCommitsPayload)
	}
}

const nodePayload = `
{
	"id": "someid",
	"config": {
		"alias": "myalias"
	}
}
`
const sessionPayload = `
{
	"sessionId": "session_id",
	"status": "authorized",
	"publicKey": "a_pub_key",
	"alias": "myalias",
	"issuedAt": 1234567890,
	"expiresAt": 1234567891
}
`

const sessionUnauthorizedPayload = `
{
	"sessionId": "session_id",
	"status": "unauthorized",
	"publicKey": "a_pub_key",
	"alias": "myalias",
	"issuedAt": 1234567890,
	"expiresAt": 1234567891
}
`
const emptyPayload = `[]`

const projectsPayloadPage0 = `
[
	{
		"name": "a-project",
		"description": "a description",
		"defaultBranch": "main",
		"delegates": [
			"did:key:the_key"
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2e",
		"id": "rad:a_project"
	},
	{
		"name": "b-project",
		"description": "b description",
		"defaultBranch": "master",
		"delegates": [
			"did:key:the_other_key"
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2f",
		"id": "rad:b_project"
	}
]
`

const projectsPayloadPage1 = `
[
	{
		"name": "c-project",
		"description": "c description",
		"defaultBranch": "main",
		"delegates": [
			"did:key:the_key"
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2e",
		"id": "rad:c_project"
	},
	{
		"name": "d-project",
		"description": "d description",
		"defaultBranch": "master",
		"delegates": [
			"did:key:the_other_key"
		],
		"head": "00bfa9b18be32001481334126c311c4a327dff2f",
		"id": "rad:d_project"
	}
]
`

const projectPayload = `
{
	"name": "a-project",
	"description": "a description",
	"defaultBranch": "main",
	"delegates": [
		"did:key:the_key"
	],
	"head": "00bfa9b18be32001481334126c311c4a327dff2e",
	"id": "rad:valid_project_id"
}
`

const projectCommitFilePayload = `
{
	"binary": false,
	"name": "file_name.md",
	"content": "file content",
	"path": "file_path/file_name.md"
}
`

const projectCommitTreePayload = `
{
	"entries": [
		{
			"path": "Readme.md",
			"name": "source",
			"kind": "blob"
		},
		{
			"path": "cargo-checksum.json",
			"name": "cargo-checksum.json",
			"kind": "blob"
		},
		{
			"path": "debian",
			"name": "build-deb",
			"kind": "tree"
		}
	]
}
`

const projectCommitsPayload = `
{
	"commits": [
		{
			"commit": {
				"id": "00bfa9b18be32001481334126c311c4a327dff2e",
				"parents": [
					"5bb95551460527ce7c24640683d4c0d5cd55a52e"
				]
			}
		}
	],
	"stats": {
		"commits": 1,
		"branches": 2,
		"contributors": 3
	}
}
`

const notFound = `
{
	error: "Not Found",
	code: 404
}
`
