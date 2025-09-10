package model

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	ModelDeepDeek = "deepseek-ai/DeepSeek-R1"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// OpenAI as an specified example.
type OpenAI struct {
	sep      string
	Protocol `json:"-"`
	// callback only
	ID      string   `json:"id,omitempty"`
	Object  string   `json:"object,omitempty"`
	Created int      `json:"created,omitempty"`
	Choices []Choice `json:"choices,omitempty"`
	// callback/request
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"` // Temperature 0.7
	TopP        float32   `json:"top_p,omitempty"`       // TopP 0.9
	MaxTokens   int       `json:"max_tokens,omitempty"`  // MaxTokens 4096
	// extra body
	EnableThinking bool `json:"enable_thinking"`
}

// NewOpenAI use temp 0.7, topp 0.9, maxn 4096 if you don't know the meaning.
func NewOpenAI(model, sep string, temp, topp float32, maxn uint) *OpenAI {
	opai := new(OpenAI)
	opai.sep = sep
	opai.Model = model
	opai.Temperature = temp
	opai.TopP = topp
	opai.MaxTokens = int(maxn)
	return opai
}

func (*OpenAI) API(api, _ string) string {
	return api
}

func (*OpenAI) Header(key string, h http.Header) {
	h.Add("Content-Type", "application/json")
	h.Add("Authorization", "Bearer "+key)
}

func (opai *OpenAI) Body() *bytes.Buffer {
	w := bytes.NewBuffer(make([]byte, 0, 8192))
	err := json.NewEncoder(w).Encode(opai)
	if err != nil {
		panic(err)
	}
	return w
}

func (opai *OpenAI) Parse(body io.Reader) error {
	return json.NewDecoder(body).Decode(&opai)
}

func (opai *OpenAI) Output() string {
	if len(opai.Choices) == 0 {
		return ""
	}
	return CutLast(opai.Choices[len(opai.Choices)-1].Message.Content, opai.sep)
}

func (opai *OpenAI) OutputRaw() string {
	if len(opai.Choices) == 0 {
		return ""
	}
	return opai.Choices[len(opai.Choices)-1].Message.Content
}

func (opai *OpenAI) System(prompt string) Protocol {
	opai.Messages = make([]Message, 1, 8)
	opai.Messages[0] = Message{
		Role:    "system",
		Content: prompt,
	}
	return opai
}

func (opai *OpenAI) User(prompt string) Protocol {
	opai.Messages = append(opai.Messages, Message{
		Role:    "user",
		Content: prompt,
	})
	return opai
}

func (opai *OpenAI) Assistant(prompt string) Protocol {
	opai.Messages = append(opai.Messages, Message{
		Role:    "assistant",
		Content: prompt,
	})
	return opai
}
