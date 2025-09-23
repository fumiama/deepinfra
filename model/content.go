package model

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime"
	"strings"

	"github.com/fumiama/imgsz"
	"github.com/pkg/errors"
)

var (
	ErrUnsupportedImageExtension = errors.New("unsupported image extension")
)

type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeImageURL ContentType = "image_url"
)

type Contents []Content

func (cs Contents) String() string {
	switch len(cs) {
	case 0:
		return ""
	case 1:
		if cs[0].Type == ContentTypeText {
			return cs[0].Text
		}
		fallthrough
	default:
		sb := strings.Builder{}
		for _, c := range cs {
			switch c.Type {
			case ContentTypeText:
				sb.WriteString("<text>")
				sb.WriteString(c.Text)
				sb.WriteString("</text>")
			case ContentTypeImageURL:
				if c.ImageURL != nil {
					sb.WriteString("<image_url>")
					sb.WriteString(c.ImageURL.URL)
					sb.WriteString("</image_url>")
				}
			default:
				panic("unsupported ContentType " + c.Type)
			}
		}
		return sb.String()
	}
}

type Content struct {
	Type     ContentType      `json:"type"`
	Text     string           `json:"text,omitempty"`
	ImageURL *ContentImageURL `json:"image_url,omitempty"`
}

func NewContentText(txt string) Content {
	return Content{Type: ContentTypeText, Text: txt}
}

type ContentImageURL struct {
	URL string `json:"url"`
}

func NewContentImageURL(u string) Content {
	return Content{Type: ContentTypeImageURL, ImageURL: &ContentImageURL{URL: u}}
}

func NewContentImageDataBase64URL(data []byte) (string, error) {
	_, ext, err := imgsz.DecodeSize(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	typ := mime.TypeByExtension("." + ext)
	if len(typ) == 0 {
		return "", errors.Wrap(ErrUnsupportedImageExtension, ext)
	}
	sb := strings.Builder{}
	sb.WriteString("data:")
	sb.WriteString(typ)
	sb.WriteString(";base64,")
	enc := base64.NewEncoder(base64.StdEncoding, &sb)
	_, _ = io.Copy(enc, bytes.NewReader(data))
	_ = enc.Close()
	return sb.String(), nil
}
