package model

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	modelDeepDeek = "deepseek-ai/DeepSeek-R1"
)

// DeepSeek as an specified example.
type DeepSeek struct {
	Inputer                   `json:"-"`
	Outputer                  `json:"-"`
	MessageBuilder[*DeepSeek] `json:"-"`
	// callback only
	ID      string   `json:"id,omitempty"`
	Object  string   `json:"object,omitempty"`
	Created int      `json:"created,omitempty"`
	Choices []Choice `json:"choices,omitempty"`
	// callback/request
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"` // Temperature 0.7
	TopP        float32   `json:"top_p"`       // TopP 0.9
	MaxTokens   int       `json:"max_tokens"`  // MaxTokens 16384

}

// NewDeepSeek 0.7, 0.9
func NewDeepSeek(temp, topp float32, maxn uint) *DeepSeek {
	ds := new(DeepSeek)
	ds.Model = modelDeepDeek
	ds.Temperature = temp
	ds.TopP = topp
	ds.MaxTokens = int(maxn)
	return ds
}

func (ds *DeepSeek) Body() *bytes.Buffer {
	w := bytes.NewBuffer(make([]byte, 0, 16384))
	err := json.NewEncoder(w).Encode(ds)
	if err != nil {
		panic(err)
	}
	return w
}

func (ds *DeepSeek) Parse(body io.Reader) error {
	return json.NewDecoder(body).Decode(&ds)
}

func (ds *DeepSeek) Output() string {
	if len(ds.Choices) == 0 {
		return ""
	}
	return CutLast(ds.Choices[len(ds.Choices)-1].Message.Content, SeparatorThink)
}

func (ds *DeepSeek) OutputRaw() string {
	if len(ds.Choices) == 0 {
		return ""
	}
	return ds.Choices[len(ds.Choices)-1].Message.Content
}

func (ds *DeepSeek) System(prompt string) *DeepSeek {
	ds.Messages = make([]Message, 1, 8)
	ds.Messages[0] = Message{
		Role:    "system",
		Content: prompt,
	}
	return ds
}

func (ds *DeepSeek) User(prompt string) *DeepSeek {
	ds.Messages = append(ds.Messages, Message{
		Role:    "user",
		Content: prompt,
	})
	return ds
}

func (ds *DeepSeek) Assistant(prompt string) *DeepSeek {
	ds.Messages = append(ds.Messages, Message{
		Role:    "assistant",
		Content: prompt,
	})
	return ds
}
