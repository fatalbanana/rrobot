package main

type Milter struct {
	RemoveHeaders map[string]int8            `json:"remove_headers"`
	AddHeaders    map[string]MilterAddHeader `json:"add_headers"`
}

type MilterAddHeader struct {
	Order string `json:"order"`
	Value string `json:"value"`
}

type Symbol struct {
	Name    string   `json:"name"`
	Options []string `json:"options"`
	Score   float32  `json:"score"`
}

type RspamdURL struct {
	Host       string `json:"host"`
	Phished    bool   `json:"phished"`
	Redirected bool   `json:"redirected"`
	TLD        string `json:"tld"`
	URL        string `json:"url"`
}

type RspamdResult struct {
	Action        string            `json:"action"`
	DKIMSignature string            `json:"dkim-signature"`
	Emails        []string          `json:"emails"`
	MessageId     string            `json:"message-id"`
	Messages      map[string]string `json:"messages"`
	Milter        Milter            `json:"milter"`
	RequiredScore float32           `json:"required_score"`
	Score         float32           `json:"score"`
	Skipped       bool              `json:"is_skipped"`
	Symbols       map[string]Symbol `json:"symbols"`
	URLs          []RspamdURL       `json:"urls"`
}
