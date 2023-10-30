package internal

type NodeInfo struct {
	Version string `json:"version"`
	Node    Node   `json:"node"`
}

type Node struct {
	ID string `json:"id"`
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
