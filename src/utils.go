package main

import "fmt"
func ColorText(text string, colorCode int) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", colorCode, text)
}
