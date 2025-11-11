package src

import (
	"fmt"
	"regexp"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)
func visibleLen(s string) int {
	return len(ansiRegex.ReplaceAllString(s, ""))
}

func Render(info []string, logo []string, logoExists bool) {
	if !logoExists {
		for _, line := range info {
			fmt.Println(line)
		}
		return
	}

	logoWidth := 0
	for _, line := range logo {
		if l := visibleLen(line); l > logoWidth {
			logoWidth = l
		}
	}

	maxHeight := len(info)
	if len(logo) > maxHeight {
		maxHeight = len(logo)
	}

	for i := 0; i < maxHeight; i++ {
		var logoLine, infoLine string

		if i < len(logo) {
			logoLine = logo[i]
		}
		if i < len(info) {
			infoLine = info[i]
		}

		// Pad based on visible length, not raw string length
		fmt.Printf("%-*s  %s\n", logoWidth+len(logoLine)-visibleLen(logoLine), logoLine, infoLine)
	}
}
