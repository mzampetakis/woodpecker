package types

const (
	EventTypeHeaderKey = "x-radicle-event-type"
	EventTypePush      = "push"
	EventTypePatch     = "patch"
)

type PushPayload struct {
	Author     Peer       `json:"author"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Commits    []Commit   `json:"commits"`
	Repository Repository `json:"repository"`
}

type PatchPayload struct {
	Action     string     `json:"action"`
	Patch      Patch      `json:"patch"`
	Repository Repository `json:"repository"`
}

type Repository struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Private       bool     `json:"private"`
	DefaultBranch string   `json:"default_branch"`
	URL           string   `json:"url"`
	CloneURL      string   `json:"clone_url"`
	Delegates     []string `json:"delegates"`
}

type Peer struct {
	ID    string `json:"id"`
	Alias string `json:"alias"`
}

type Commit struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Message   string   `json:"message"`
	Timestamp string   `json:"timestamp"`
	URL       string   `json:"url"`
	Author    Author   `json:"author"`
	Added     []string `json:"added"`
	Modified  []string `json:"modified"`
	Removed   []string `json:"removed"`
}

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Patch struct {
	ID        string          `json:"id"`
	Author    Peer            `json:"author"`
	Title     string          `json:"title"`
	State     State           `json:"state"`
	Before    string          `json:"before"`
	After     string          `json:"after"`
	URL       string          `json:"url"`
	Target    string          `json:"target"`
	Labels    []string        `json:"labels"`
	Assignees []string        `json:"assignees"`
	Revisions []PatchRevision `json:"revisions"`
}

type PatchRevision struct {
	ID          string `json:"id"`
	Author      Peer   `json:"author"`
	Description string `json:"description"`
	Base        string `json:"base"`
	Oid         string `json:"oid"`
	Timestamp   int64  `json:"timestamp"`
}

type Draft struct {
	Status string `json:"status"`
}

type State struct {
	Status    string      `json:"status"`
	Conflicts []Conflicts `json:"conflicts"`
}

type Archived struct {
	Status string `json:"status"`
}

type Merged struct {
	Status   string `json:"status"`
	Revision string `json:"revision"`
	Commit   string `json:"commit"`
}

type Conflicts struct {
	RevisionID string `json:"revision_id"`
	Oid        string `json:"oid"`
}
