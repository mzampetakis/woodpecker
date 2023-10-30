package internal

type NodeInfo struct {
	ID     string     `json:"id"`
	Config NodeConfig `json:"config"`
}

type NodeConfig struct {
	Alias string `json:"alias"`
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
