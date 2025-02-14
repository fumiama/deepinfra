package model

import "strings"

const (
	SeparatorThink = "</think>"
)

func CutLast(txt, sep string) string {
	if sep == "" { // no need to cut
		return txt
	}
	a := strings.LastIndex(txt, sep)
	if a < 0 {
		return ""
	}
	a += len(sep)
	if a >= len(txt) {
		return ""
	}
	return strings.TrimSpace(txt[a:])
}
