package model

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	ModelDeepDeek = "deepseek-ai/DeepSeek-V3.1-Terminus"
)

type OpenAIMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"` // Contents or string
}

type OpenAIChoice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// OpenAI as an specified example.
type OpenAI struct {
	sep      string
	Protocol `json:"-"`
	// callback only
	ID      string         `json:"id,omitempty"`
	Object  string         `json:"object,omitempty"`
	Created int            `json:"created,omitempty"`
	Choices []OpenAIChoice `json:"choices,omitempty"`
	// callback/request
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float32         `json:"temperature,omitempty"` // Temperature 0.7
	TopP        float32         `json:"top_p,omitempty"`       // TopP 0.9
	MaxTokens   int             `json:"max_tokens,omitempty"`  // MaxTokens 4096
	// extra body
	EnableThinking bool `json:"enable_thinking"` // EnableThinking is always false in non-stream mode, adapt to 百炼
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

func (opai *OpenAI) Output() Contents {
	if len(opai.Choices) == 0 {
		return nil
	}
	cs := make(Contents, 0, 8)
	c := opai.Choices[len(opai.Choices)-1].Message.Content
	err := json.Unmarshal(c, &cs)
	if err == nil {
		for i := range cs {
			if cs[i].Type == ContentTypeText {
				cs[i].Text = CutLast(cs[i].Text, opai.sep)
			}
		}
		return cs
	}
	s := ""
	err = json.Unmarshal(c, &s)
	if err == nil {
		return Contents{{Type: ContentTypeText, Text: CutLast(s, opai.sep)}}
	}
	return nil
}

func (opai *OpenAI) OutputRaw() Contents {
	if len(opai.Choices) == 0 {
		return nil
	}
	cs := make(Contents, 0, 8)
	c := opai.Choices[len(opai.Choices)-1].Message.Content
	err := json.Unmarshal(c, &cs)
	if err == nil {
		return cs
	}
	s := ""
	err = json.Unmarshal(c, &s)
	if err == nil {
		return Contents{{Type: ContentTypeText, Text: s}}
	}
	return nil
}

func (opai *OpenAI) System(prompt ...Content) Protocol {
	opai.Messages = make([]OpenAIMessage, 1, 8)
	raw, err := json.Marshal(&prompt)
	if err != nil {
		panic(err)
	}
	opai.Messages[0] = OpenAIMessage{
		Role:    "system",
		Content: raw,
	}
	return opai
}

func (opai *OpenAI) User(prompt ...Content) Protocol {
	raw, err := json.Marshal(&prompt)
	if err != nil {
		panic(err)
	}
	opai.Messages = append(opai.Messages, OpenAIMessage{
		Role:    "user",
		Content: raw,
	})
	return opai
}

func (opai *OpenAI) Assistant(prompt ...Content) Protocol {
	raw, err := json.Marshal(&prompt)
	if err != nil {
		panic(err)
	}
	opai.Messages = append(opai.Messages, OpenAIMessage{
		Role:    "assistant",
		Content: raw,
	})
	return opai
}
