package http_handler

type RequestObject struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RequestKey struct {
	Key string `json:"key"`
}
