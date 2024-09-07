package model

type InsertRequest struct {
	Identifier string `json:"identifier"`
	Key        string `json:"key"`
	Value      string `json:"value"`
}
