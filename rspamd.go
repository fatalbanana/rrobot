package main

type Milter struct {
	AddHeaders    map[string]MilterAddHeader `json:"add_headers"`
	AddRecipients []string                   `json:"add_recipients"`
	ChangeFrom    string                     `json:"change_from"`
	NoAction      bool                       `json:"no_action"`
	Reject        string                     `json:"reject"`
	RemoveHeaders map[string]int8            `json:"remove_headers"`
	SpamHeader    string                     `json:"spam_header"`
}

type MilterAddHeader struct {
	Order string `json:"order"`
	Value string `json:"value"`
}

type Symbol struct {
	Description string   `json:"description"`
	MetricScore float32  `json:"metric_score"`
	Name        string   `json:"name"`
	Options     []string `json:"options"`
	Score       float32  `json:"score"`
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
	Subject       string            `json:"subject"`
	Symbols       map[string]Symbol `json:"symbols"`
	TimeReal      float32           `json:"time_real"`
	URLs          []string          `json:"urls"`
}
