package model

type FetchRequest struct {
	Identifier string `json:"identifier"`
	Key        string `json:"key"`
}

type FetchResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
