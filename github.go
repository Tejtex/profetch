package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
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

	result := make([]string, 4);
	result[1] = ColorText("Stars: ", colorCode) + stars
	result[2] = ColorText("Forks: ", colorCode) + forks
	if data["license"] != nil {
			result[3] = ColorText("License: ", colorCode) + data["license"].(map[string]interface{})["name"].(string)

	}

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