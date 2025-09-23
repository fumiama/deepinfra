package model

import (
	"encoding/base64"
	"strings"
	"testing"
)

const contentTextSmallImage = "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mNk+A8AAQUBAScY42YAAAAASUVORK5CYII="

func TestContentImageDataBytes(t *testing.T) {
	// 从 smallimg 中提取 base64 数据
	parts := strings.Split(contentTextSmallImage, ",")
	if len(parts) != 2 {
		t.Fatal("Invalid data URL format")
	}

	// 解码 base64 数据
	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	s, err := NewContentImageDataBase64URL(data)
	if err != nil {
		t.Fatalf("NewContentImageDataBytes failed: %v", err)
	}

	// 比较结果
	if s != contentTextSmallImage {
		t.Errorf("Expected %s, got %s", contentTextSmallImage, s)
	}
}
