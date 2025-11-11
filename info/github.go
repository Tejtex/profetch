package info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/tejtex/profetch/utils"
)

func FetchGithubInfo(root string, colorCode int) []string{
	if !isGithubRepo(root) {
		return make([]string, 0)
	} 


	repoURL := getGithubRepo(root);
	repo := strings.TrimSpace(strings.TrimPrefix(repoURL, "https://github.com/"))
	repo = strings.TrimPrefix(repo, "git@github.com:")
	repo = strings.TrimSuffix(repo, ".git")
	res, err := http.Get("https://api.github.com/repos/" + repo);
	if err != nil {
		return make([]string, 0)
	}

	defer res.Body.Close()

	var data map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
        return make([]string, 0)

    }
	var stars string
	if starsVal, ok := data["stargazers_count"].(float64); ok {
		stars = fmt.Sprintf("%.0f", starsVal)
	} else {
		stars = "0" // or "-"
	}
	var forks string
	if forksVal, ok := data["forks_count"].(float64); ok {
		forks = fmt.Sprintf("%.0f", forksVal)
	} else {
		forks = "0" // or "-"
	}
	var openIssues string
	if openIssuesVal, ok := data["open_issues_count"].(float64); ok {
		openIssues = fmt.Sprintf("%.0f", openIssuesVal)
	} else {
		openIssues = "0"
	}

	result := make([]string, 6);
	result[1] = utils.Format("Stars", stars, colorCode)
	result[2] = utils.Format("Forks", forks, colorCode)
	result[3] = utils.Format("Issues", openIssues, colorCode)
	if data["license"] != nil {
		result[4] = utils.Format("License", data["license"].(map[string]interface{})["name"].(string), colorCode)
	} else {
		result[4] = utils.Format("License", "-", colorCode)
	}
	result[5] = utils.Format("URL", strings.TrimSpace(repoURL), colorCode)

	return result
}

func isGithubRepo(root string)bool {
	git := func(command ...string) (string, error) {
		cmd := exec.Command("git", append([]string{"--no-pager"}, command...)...)
		cmd.Dir = root;
		res, err := cmd.CombinedOutput()
		return string(res), err
	}
	output, err := git("remote", "get-url", "origin")
	if err != nil {
		return false
	}

	return strings.Contains(output, "github.com")
}

func getGithubRepo(root string) string {
	git := func(command ...string) (string, error) {
		cmd := exec.Command("git", append([]string{"--no-pager"}, command...)...)
		cmd.Dir = root;
		res, err := cmd.CombinedOutput()
		return string(res), err
	}
	output, _ := git("remote", "get-url", "origin")

	return output
}