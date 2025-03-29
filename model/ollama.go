package model

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// OLLaMA as an specified example.
type OLLaMA struct {
	sep      string
	Protocol `json:"-"`
	// callback only
	ID       string    `json:"id,omitempty"`
	Object   string    `json:"object,omitempty"`
	Created  int       `json:"created,omitempty"`
	Messages []Message `json:"messages"`
	// callback/request
	Model       string   `json:"model"`
	Message     *Message `json:"message,omitempty"`
	Temperature float32  `json:"temperature"` // Temperature 0.7
	TopP        float32  `json:"top_p"`       // TopP 0.9
	MaxTokens   int      `json:"max_tokens"`  // MaxTokens 4096
	Stream      bool     `json:"stream"`
}

// NewOLLaMA use temp 0.7, topp 0.9, maxn 4096 if you don't know the meaning.
func NewOLLaMA(model, sep string, temp, topp float32, maxn uint) *OLLaMA {
	opai := new(OLLaMA)
	opai.sep = sep
	opai.Model = model
	opai.Temperature = temp
	opai.TopP = topp
	opai.MaxTokens = int(maxn)
	return opai
}

func (*OLLaMA) API(api, _ string) string {
	return api
}

func (*OLLaMA) Header(key string, h http.Header) {
	h.Add("Content-Type", "application/json")
	h.Add("Authorization", "Bearer "+key)
}

func (ollm *OLLaMA) Body() *bytes.Buffer {
	w := bytes.NewBuffer(make([]byte, 0, 8192))
	err := json.NewEncoder(w).Encode(ollm)
	if err != nil {
		panic(err)
	}
	return w
}

func (ollm *OLLaMA) Parse(body io.Reader) error {
	return json.NewDecoder(body).Decode(&ollm)
}

func (ollm *OLLaMA) Output() string {
	if ollm.Message == nil {
		return ""
	}
	return CutLast(ollm.Message.Content, ollm.sep)
}

func (ollm *OLLaMA) OutputRaw() string {
	if ollm.Message == nil {
		return ""
	}
	return ollm.Message.Content
}

func (ollm *OLLaMA) System(prompt string) Protocol {
	ollm.Messages = make([]Message, 1, 8)
	ollm.Messages[0] = Message{
		Role:    "system",
		Content: prompt,
	}
	return ollm
}

func (ollm *OLLaMA) User(prompt string) Protocol {
	ollm.Messages = append(ollm.Messages, Message{
		Role:    "user",
		Content: prompt,
	})
	return ollm
}

func (ollm *OLLaMA) Assistant(prompt string) Protocol {
	ollm.Messages = append(ollm.Messages, Message{
		Role:    "assistant",
		Content: prompt,
	})
	return ollm
}
