package main

type Response struct {
	Chunk    string     `json:"chunk"`
	Metadata []Metadata `json:"metadata"`
}

type Metadata struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Link    string `json:"link"`
}
