package info

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/tejtex/profetch/internal/utils"
)

func FetchGithubInfo(root string, colorCode int) []string{
	if !isGithubRepo(root) {
		return make([]string, 0)
	} 


	repo := getGithubRepo(root);
	repoWithout := strings.TrimSpace(strings.TrimPrefix(repo, "https://github.com/"))
	res, err := http.Get("https://api.github.com/repos/" + repoWithout);
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

	result := make([]string, 3);
	result[1] = utils.ColorText("Stars: ", colorCode) + stars
	result[2] = utils.ColorText("Forks: ", colorCode) + forks

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

	return strings.HasPrefix(output, "https://github.com")
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