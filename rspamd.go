package main

type RspamdResult struct {
	Score         float32           `json:"score"`
	RequiredScore float32           `json:"required_score"`
	Action        string            `json:"action"`
	Messages      map[string]string `json:"messages"`
	Milter        struct {
		RemoveHeaders map[string]int8        `json:"remove_headers"`
		AddHeaders    map[string]interface{} `json:"add_headers"`
	} `json:"milter"`
	Symbols map[string]struct {
		Score float32
	} `json:"symbols"`
}
