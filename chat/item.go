package chat

import "strings"

type item struct {
	isatme   bool
	usr, txt string
}

func (item *item) writeToBuilder(sb *strings.Builder, atprefix, namel, namer string) {
	if item.isatme {
		sb.WriteString(atprefix)
	}
	sb.WriteString(namel)
	sb.WriteString(item.usr)
	sb.WriteString(namer)
	sb.WriteString(item.txt)
}
