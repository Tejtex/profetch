package info

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	ignore "github.com/sabhiram/go-gitignore"
	"github.com/tejtex/profetch/internal/utils"
)

var allowedExts = map[string]string{
	".go":   "Go",
	".py":   "Python",
	".js":   "JavaScript",
	".ts":   "TypeScript",
	".c":    "C",
	".cpp":  "C++",
	".java": "Java",
	".rs":   "Rust",
}


func FetchInfo(root string, colorCode int) ([]string, error) {
	res := make([]string, 2)
	projName, err := fetchProjName(root, colorCode);
	if err != nil {
		return nil, err
	}
	res = append(res, projName...);

	files, err := addFilesAndLines(root, colorCode)
	if err != nil {
		return nil, err
	}
	res = append(res, files...)
	size, err := countSize(root, colorCode)
	if err != nil {
		return nil, err
	}
	res = append(res, size...);

	langs, err := fetchLangs(root, colorCode)
	if err != nil {
		return nil, err
	}
	res = append(res, langs...)

	git := FetchGitInfo(root, colorCode);
	res = append(res, git...);

	github := FetchGithubInfo(root, colorCode)
	res = append(res, github...);

	return res, nil
}

func fetchProjName(root string, colorCode int ) ([]string, error) {
	absPath, _ := filepath.Abs(root)
	dirName := filepath.Base(absPath)
	res := []string{utils.ColorText(utils.ColorText(utils.ColorText(dirName, 1), 4), colorCode), ""}
	return res, nil
}
func fetchLangs(root string, colorCode int) ([]string, error) {
	langCount := make(map[string]int)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if langName, ok := allowedExts[ext]; ok {
			langCount[langName]++
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Format output
	var res []string
	res = append(res, "")
	res = append(res, utils.ColorText("Languages:", colorCode))
	total := 0
	for _, count := range langCount {
		total += count
	}
	for lang, count := range langCount {
		line := "  " + utils.ColorText(lang+": ", colorCode) + strconv.Itoa(count) + " file(s), " + fmt.Sprintf("%d", int64(float64(count) / float64(total) * 100)) + "%"
		res = append(res, line)
	}
	return res, nil
}

func addFilesAndLines(root string, colorCode int) ([]string, error) {
	fileCount, lineCount, err := countFilesAndLines(root)
	if err != nil {
		return make([]string, 0), err
	}
	res := make([]string, 2)
	res[0] = utils.ColorText("Number of lines: ", colorCode) + strconv.Itoa(lineCount)
	res[1] = utils.ColorText("Number of files: ", colorCode) + strconv.Itoa(fileCount)
	return res, nil
}

func countSize(root string, colorCode int) ([]string, error ) {
	var size int64 = 0

	file, err := os.Open(filepath.Join(root, ".gitignore"))
	var parser *ignore.GitIgnore
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
			
		}
	} else {
		parser, err = ignore.CompileIgnoreFile(filepath.Join(root, ".gitignore"));
		if err != nil {
			return nil, err
		}
	}
	
	defer file.Close()
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if parser != nil && parser.MatchesPath(path) {
			return nil
		}
		stats, _ := os.Stat(path)
		size += stats.Size()
		return nil
	});
	if err != nil {
		return make([]string, 0), err
	}
	return []string{utils.ColorText("Size: ", colorCode) + strconv.Itoa(int(size / 1024)) + "K" }, nil
}

func countFilesAndLines(root string) (int, int, error) {
	var fileCount, lineCount int

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			return nil
		}
		ext := filepath.Ext(path)
		if _, ok := allowedExts[ext]; !ok {
			return nil
		}

		fileCount++
		lines, err := countLines(path)
		if err != nil {
			return nil
		}
		lineCount += lines
		return nil
	})

	return fileCount, lineCount, err
}
func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines, scanner.Err()
}