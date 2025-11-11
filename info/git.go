package info

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/tejtex/profetch/utils"
)

func FetchGitInfo(root string, colorCode int) []string {
	commits, lastCommit, contributors, version, age, err := fetchGit(root);
	if !err {
		return make([]string, 0)
	}


	res := make([]string, 6)
	
	res[1] = utils.Format("Commits", commits, colorCode)
	res[2] = utils.Format("Last commit", lastCommit, colorCode)
	res[3] = utils.Format("Contributors", contributors, colorCode)
	res[4] = utils.Format("Repo age", utils.FormatDuration(age), colorCode)
	res[5] = utils.Format("Version", version, colorCode)
	return res
}

func fetchGit(root string) (int, string, int, string, time.Duration, bool) {
	path, _ := filepath.Abs(filepath.Join(root, ".git"))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, "", 0, "", time.Millisecond ,false 
	}

	git := func(command ...string) (string, error) {
		cmd := exec.Command("git", append([]string{"--no-pager"}, command...)...)
		cmd.Dir = root;
		res, err := cmd.CombinedOutput()
		return string(res), err
	}
	_, err := exec.Command("git", "rev-parse", "--verify", "HEAD").CombinedOutput()
	if err != nil {
		return 0, "", 0, "", time.Millisecond, false
	}
	commits, err := git("rev-list", "--count", "HEAD")
	if err != nil {
		commits = "0"
	}
	lastCommit, err := git("log", "-1", "--format=%s")
	if err != nil {
		lastCommit = "-"
	}
	output, err := git("log", "--format=%an") // returns string

	if err != nil {
		output = ""
	}

	version, err := git("describe", "--tags")
	if err != nil {
		version = "-"
	}

	allCommits, err := git("log", "--reverse", "--format=%ct")
	if err != nil {
		allCommits = ""
	}

	lines := strings.Split(strings.TrimSpace(allCommits), "\n")

    ts, _ := strconv.Atoi(lines[0])
    

    created := time.Unix(int64(ts), 0)
    age := time.Since(created)

	lines = strings.Split(strings.TrimSpace(output), "\n")
	seen := make(map[string]struct{})
	var result []string

	for _, line := range lines {
		if _, ok := seen[line]; !ok {
			seen[line] = struct{}{}
			result = append(result, line)
		}
	}
	contributors := len(result)
	numCommits, err := strconv.Atoi(strings.TrimSpace(commits))
	if err != nil {
		numCommits = 0;
	}
	return numCommits, strings.TrimSpace(lastCommit), contributors, strings.TrimSpace(version) ,age, true
}