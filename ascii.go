package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
	"golang.org/x/image/draw"
)

var imageExts = []string{".png", ".jpg", ".jpeg", ".svg"}

func FindLogoFile(root string) (string, error) {
	// Compile .gitignore if exists
	var parser *ignore.GitIgnore
	gitignorePath := filepath.Join(root, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		var err2 error
		parser, err2 = ignore.CompileIgnoreFile(gitignorePath)
		if err2 != nil {
			return "", err2
		}
	}

	var found string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			// Skip unreadable files/dirs
			return nil
		}
		if d.IsDir() {
			return nil
		}

		// Skip ignored files
		if parser != nil && parser.MatchesPath(path) {
			return nil
		}

		// Check if filename contains "logo" (case-insensitive)
		name := strings.ToLower(filepath.Base(path))
		if !strings.Contains(name, "logo") {
			return nil
		}

		// Check if extension is one of the allowed image types
		ext := strings.ToLower(filepath.Ext(path))
		for _, e := range imageExts {
			if ext == e {
				found = path
				return filepath.SkipDir // stop after first match
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	if found == "" {
		return "", os.ErrNotExist
	}
	return found, nil
}

func FetchLogo(root string) ([]string, bool) {
	file, err := FindLogoFile(root)
	if err != nil {
		return nil, false
	}
	res, err := imageToASCII(file, 50, 25);
	if err != nil {
		return nil, false
	}
	return res, true

}


var asciiChars = []rune(" .'`^\",:;Il!i~+_-?][}{1)(|\\/*tfrxnvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$")

func imageToASCII(path string, widthChars, heightChars int) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(path))
	var img image.Image

	switch ext {
	case ".png":
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
	case ".svg":
		icon, err := oksvg.ReadIconStream(file)
		if err != nil {
			return nil, err
		}
		width := int(icon.ViewBox.W)
		height := int(icon.ViewBox.H)
		rgba := image.NewRGBA(image.Rect(0, 0, width, height))
		icon.SetTarget(0, 0, float64(width), float64(height))
		icon.Draw(rasterx.NewDasher(width, height, rasterx.NewScannerGV(width, height, rgba, rgba.Bounds())), 1.0)
		img = rgba
	default:
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	// Resize to target ASCII size
	resized := image.NewRGBA(image.Rect(0, 0, widthChars, heightChars))
	draw.ApproxBiLinear.Scale(resized, resized.Bounds(), img, img.Bounds(), draw.Over, nil)

	rows := make([]string, heightChars)
	for y := 0; y < heightChars; y++ {
		var row strings.Builder
		for x := 0; x < widthChars; x++ {
			c := color.NRGBAModel.Convert(resized.At(x, y)).(color.NRGBA)
			// Compute brightness using perceptual weighting
			brightness := 0.2126*float64(c.R) + 0.7152*float64(c.G) + 0.0722*float64(c.B)
			idx := int((brightness / 255.0) * float64(len(asciiChars)-1))
			char := asciiChars[idx]

			// Add ANSI color escape
			row.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", c.R, c.G, c.B, char))
		}
		rows[y] = row.String()
	}

	return rows, nil
}
