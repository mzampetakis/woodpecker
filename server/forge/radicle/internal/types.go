package internal

type ListOpts struct {
	Page    int
	PageLen int
}

type NodeInfo struct {
	ID     string     `json:"id"`
	Config NodeConfig `json:"config"`
}

type NodeConfig struct {
	Alias string `json:"alias"`
}

type Project struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	DefaultBranch string `json:"defaultBranch"`
	Head          string `json:"head"`
}

type Error struct {
	Status int
	Body   struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e Error) Error() string {
	return e.Body.Message
}
