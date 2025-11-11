package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/tejtex/profetch/src/ascii"
	"github.com/tejtex/profetch/src/display"
	"github.com/tejtex/profetch/src/info"
)

func main() {
	rand.Seed(time.Now().UnixNano()) // seed with current timestamp

	if rand.Intn(50) == 0 {
		fmt.Println("sorry, profetch decided to not work this time D:");
		os.Exit(67);
	}

	logo, ok := ascii.FetchLogo(".");
	info, err := info.FetchInfo(".", rand.Intn(37 - 31) + 31);
	if err != nil {
		fmt.Fprintf(os.Stderr, "profetch: %v\n", err);
		os.Exit(67);
	}
	display.Render(info, logo, ok);
}