package model

import (
	"bytes"
	"encoding/json"
	"io"
)

// Custom as an compatible example.
type Custom struct {
	Inputer                   `json:"-"`
	Outputer                  `json:"-"`
	MessageBuilder[*DeepSeek] `json:"-"`
	sep                       string `json:"-"`
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

func NewCustom(model, sep string, temp, topp float32, maxn uint) *Custom {
	c := new(Custom)
	c.sep = sep
	c.Model = model
	c.Temperature = temp
	c.TopP = topp
	c.MaxTokens = int(maxn)
	return c
}

func (c *Custom) Parse(body io.Reader) error {
	return json.NewDecoder(body).Decode(&c)
}

func (c *Custom) Output() string {
	if len(c.Choices) == 0 {
		return ""
	}
	return CutLast(c.Choices[len(c.Choices)-1].Message.Content, c.sep)
}

func (c *Custom) OutputRaw() string {
	if len(c.Choices) == 0 {
		return ""
	}
	return c.Choices[len(c.Choices)-1].Message.Content
}

func (ds *Custom) System(prompt string) *Custom {
	ds.Messages = make([]Message, 1, 8)
	ds.Messages[0] = Message{
		Role:    "system",
		Content: prompt,
	}
	return ds
}

func (ds *Custom) User(prompt string) *Custom {
	ds.Messages = append(ds.Messages, Message{
		Role:    "user",
		Content: prompt,
	})
	return ds
}

func (ds *Custom) Assistant(prompt string) *Custom {
	ds.Messages = append(ds.Messages, Message{
		Role:    "assistant",
		Content: prompt,
	})
	return ds
}

func (ds *Custom) Body() *bytes.Buffer {
	w := bytes.NewBuffer(make([]byte, 0, 16384))
	err := json.NewEncoder(w).Encode(ds)
	if err != nil {
		panic(err)
	}
	return w
}
