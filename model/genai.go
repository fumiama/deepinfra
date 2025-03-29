package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	ModelGemini15Flash = "models/gemini-1.5-flash"
)

type Text struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Text `json:"parts"`
	Role  string `json:"role,omitempty"`
}

func (c *Content) String() string {
	sb := strings.Builder{}
	for _, p := range c.Parts {
		sb.WriteString(p.Text)
	}
	return sb.String()
}

type Candidate struct {
	Content      Content `json:"content"`
	FinishReason string  `json:"finishReason"`
	Index        int     `json:"index"`
}

// GenAI is Goole API format
type GenAI struct {
	model    string `json:"-"`
	Protocol `json:"-"`
	// request only
	Contents          []Content `json:"contents,omitempty"`
	SystemInstruction *Content  `json:"systemInstruction,omitempty"`
	GenerationConfig  struct {
		Temperature      float32 `json:"temperature,omitempty"`
		ResponseMimeType string  `json:"responseMimeType,omitempty"`
		TopP             float32 `json:"topP,omitempty"`
		MaxOutputTokens  int     `json:"maxOutputTokens,omitempty"`
	} `json:"generationConfig"`
	// callback only
	Candidates []Candidate `json:"candidates,omitempty"`
}

// NewGenAI use temp 0.7, topp 0.9, maxn 4096 if you don't know the meaning.
func NewGenAI(model string, temp, topp float32, maxn uint) *GenAI {
	opai := new(GenAI)
	opai.model = model
	opai.GenerationConfig.Temperature = temp
	opai.GenerationConfig.ResponseMimeType = "text/plain"
	opai.GenerationConfig.TopP = topp
	opai.GenerationConfig.MaxOutputTokens = int(maxn)
	return opai
}

func (opai *GenAI) API(api, key string) string {
	return fmt.Sprintf("%s/%s:generateContent?key=%s", api, opai.model, key)
}

func (*GenAI) Header(_ string, h http.Header) {
	h.Add("Content-Type", "application/json")
}

func (opai *GenAI) Body() *bytes.Buffer {
	w := bytes.NewBuffer(make([]byte, 0, 8192))
	err := json.NewEncoder(w).Encode(opai)
	if err != nil {
		panic(err)
	}
	return w
}

func (opai *GenAI) Parse(body io.Reader) error {
	return json.NewDecoder(body).Decode(&opai)
}

func (opai *GenAI) Output() string {
	if len(opai.Candidates) == 0 {
		return ""
	}
	return opai.Candidates[0].Content.String()
}

func (opai *GenAI) OutputRaw() string {
	return opai.Output()
}

func (opai *GenAI) System(prompt string) Protocol {
	opai.SystemInstruction = &Content{
		Parts: []Text{{prompt}},
	}
	return opai
}

func (opai *GenAI) User(prompt string) Protocol {
	opai.Contents = append(opai.Contents, Content{
		Parts: []Text{{prompt}},
		Role:  "user",
	})
	return opai
}

func (opai *GenAI) Assistant(prompt string) Protocol {
	opai.Contents = append(opai.Contents, Content{
		Parts: []Text{{prompt}},
		Role:  "model",
	})
	return opai
}
