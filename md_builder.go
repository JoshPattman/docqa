package docqa

import (
	"fmt"
	"strings"
)

type mdBuilder struct {
	lines []string
}

func (md *mdBuilder) Header(level int, header string) {
	md.lines = append(md.lines, fmt.Sprintf("%s %s", strings.Repeat("#", level), strings.Trim(header, " \n\r\t")))
}

func (md *mdBuilder) Headerf(level int, txt string, args ...any) {
	header := fmt.Sprintf(txt, args...)
	md.Header(level, header)
}

func (md *mdBuilder) Bullet(indent int, bullet string) {
	bullet = strings.Trim(bullet, " \n\r\t.")
	md.lines = append(md.lines, fmt.Sprintf("%s- %s.", strings.Repeat(" ", 2*indent), bullet))
}

func (md *mdBuilder) Bulletf(indent int, txt string, args ...any) {
	bullet := fmt.Sprintf(txt, args...)
	md.Bullet(indent, bullet)
}

func (md *mdBuilder) Break(size int) {
	for range size {
		md.lines = append(md.lines, "")
	}
}

func (ms *mdBuilder) Build() string {
	return strings.Join(ms.lines, "\n")
}
