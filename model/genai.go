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

type GenAIInlineData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"` // Data is base64 repr
}

type GenAIPart struct {
	Text       string           `json:"text,omitempty"`
	InlineData *GenAIInlineData `json:"inline_data,omitempty"`
}

type GenAIContent struct {
	Parts []GenAIPart `json:"parts"`
	Role  string      `json:"role,omitempty"`
}

func (c *GenAIContent) String() string {
	sb := strings.Builder{}
	for _, p := range c.Parts {
		sb.WriteString(p.Text)
	}
	return sb.String()
}

type GenAICandidate struct {
	Content      GenAIContent `json:"content"`
	FinishReason string       `json:"finishReason"`
	Index        int          `json:"index"`
}

// GenAI is Goole API format
type GenAI struct {
	model    string `json:"-"`
	Protocol `json:"-"`
	// request only
	Contents          []GenAIContent `json:"contents,omitempty"`
	SystemInstruction *GenAIContent  `json:"systemInstruction,omitempty"`
	GenerationConfig  struct {
		Temperature      float32 `json:"temperature,omitempty"`
		ResponseMimeType string  `json:"responseMimeType,omitempty"`
		TopP             float32 `json:"topP,omitempty"`
		MaxOutputTokens  int     `json:"maxOutputTokens,omitempty"`
	} `json:"generationConfig"`
	// callback only
	Candidates []GenAICandidate `json:"candidates,omitempty"`
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

func (opai *GenAI) Output() Contents {
	return opai.OutputRaw()
}

func (opai *GenAI) OutputRaw() Contents {
	if len(opai.Candidates) == 0 {
		return nil
	}
	raw := opai.Candidates[0].Content
	cs := make(Contents, len(raw.Parts))
	for i, c := range raw.Parts {
		switch {
		case c.Text != "":
			cs[i].Type = ContentTypeText
			cs[i].Text = c.Text
		case c.InlineData != nil:
			cs[i].Type = ContentTypeImageURL
			if strings.HasPrefix(c.InlineData.MimeType, "image/") {
				cs[i].ImageURL = &ContentImageURL{
					URL: "data:" + c.InlineData.MimeType + ";base64," + c.InlineData.Data,
				}
			}
		default:
			panic("unsupported genai part")
		}
	}
	return cs
}

func (cs Contents) ToGenAIParts() []GenAIPart {
	ps := make([]GenAIPart, 0, len(cs))
	for _, c := range cs {
		switch c.Type {
		case ContentTypeText:
			ps = append(ps, GenAIPart{Text: c.Text})
		case ContentTypeImageURL:
			if strings.HasPrefix(c.ImageURL.URL, "data:") {
				typ, dat, ok := strings.Cut(strings.TrimPrefix(c.ImageURL.URL, "data:"), ";")
				if !ok {
					continue
				}
				ps = append(ps, GenAIPart{InlineData: &GenAIInlineData{
					MimeType: typ,
					Data:     dat[1:], // skip ;
				}})
				continue
			}
			resp, err := http.Get(c.ImageURL.URL)
			if err != nil {
				continue
			}
			defer resp.Body.Close()
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				continue
			}
			s, err := NewContentImageDataBase64URL(data)
			if err != nil {
				continue
			}
			typ, dat, ok := strings.Cut(strings.TrimPrefix(s, "data:"), ",")
			if !ok {
				continue
			}
			ps = append(ps, GenAIPart{InlineData: &GenAIInlineData{
				MimeType: typ,
				Data:     dat[1:], // skip ,
			}})
		default:
			panic("unsupported ContentType " + c.Type)
		}
	}
	return ps
}

func (opai *GenAI) System(prompt string) Protocol {
	opai.SystemInstruction = &GenAIContent{
		Parts: []GenAIPart{{Text: prompt}},
	}
	return opai
}

func (opai *GenAI) User(prompt ...Content) Protocol {
	opai.Contents = append(opai.Contents, GenAIContent{
		Parts: Contents(prompt).ToGenAIParts(),
		Role:  "user",
	})
	return opai
}

func (opai *GenAI) Assistant(prompt ...Content) Protocol {
	opai.Contents = append(opai.Contents, GenAIContent{
		Parts: Contents(prompt).ToGenAIParts(),
		Role:  "model",
	})
	return opai
}
