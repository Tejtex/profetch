package src

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tejtex/profetch/src/utils"
)

func FetchGitInfo(root string, colorCode int) []string {
	commits, lastCommit, contributors, version, firstCommit, err := fetchGit(root);
	if !err {
		return make([]string, 0)
	}

	res := make([]string, 6)
	
	res[1] = utils.ColorText("Number of commits: ", colorCode) + strconv.Itoa(commits)
	res[2] = utils.ColorText("Last commit: ", colorCode) + lastCommit
	res[3] = utils.ColorText("Number of contributors: ", colorCode) + strconv.Itoa(contributors) 
	res[4] = utils.ColorText("First commit date: ", colorCode) + firstCommit
	res[5] = utils.ColorText("Version: ", colorCode) + version
	return res
}

func fetchGit(root string) (int, string, int, string, string, bool) {
	path, _ := filepath.Abs(filepath.Join(root, ".git"))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, "", 0, "", "" ,false 
	}

	git := func(command ...string) (string, error) {
		cmd := exec.Command("git", append([]string{"--no-pager"}, command...)...)
		cmd.Dir = root;
		res, err := cmd.CombinedOutput()
		return string(res), err
	}
	_, err := exec.Command("git", "rev-parse", "--verify", "HEAD").CombinedOutput()
	if err != nil {
		return 0, "", 0, "", "", false
	}
	commits, err := git("rev-list", "--count", "HEAD")
	if err != nil {
		return 0, "", 0, "", "", false
	}
	lastCommit, err := git("log", "-1", "--format=%s")
	if err != nil {
		return 0, "", 0, "", "", false
	}
	output, err := git("log", "--format=%an") // returns string

	if err != nil {
		return 0, "", 0 , "", "", false
	}

	version, err := git("describe", "--tags")
	if err != nil {
		return 0, "", 0,"","", false
	}

	allCommits, err := git("log", "--reverse", "--format=%cd")
	if err != nil {
		return 0, "", 0, "", "", false
	}

	firstCommit := strings.Split(strings.TrimSpace(allCommits), "\n")[0]


	lines := strings.Split(strings.TrimSpace(output), "\n")
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
		return 0, "", 0, "", "", false
	}
	return numCommits, strings.TrimSpace(lastCommit), contributors, strings.TrimSpace(version) ,firstCommit, true
}