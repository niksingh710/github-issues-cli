package main

import (
	"flag"
	"fmt"
	"github-issues-cli/editor"
	"github-issues-cli/github"
	"log"
	"os"
	"strings"
)

func main() {
	action := flag.String("action", "", "[create,update,read,close,repoen]")
	owner := flag.String("owner", "", "Owner Username")
	repo := flag.String("repo", "", "Repo Name")
	issueNumber := flag.Int("issueNumber", 0, "Integer value pointing to specific issue")
	flag.Parse()
	switch *action {
	case "create":
		create(*owner, *repo)
	case "close":
		gclose(*owner, *repo, *issueNumber)
	case "open":
		gopen(*owner, *repo, *issueNumber)
	case "update":
		update(*owner, *repo, *issueNumber)
	default:
		fmt.Fprint(os.Stderr, "Make sure provide an Action")
	}
}

func gclose(owner, repo string, issueNumber int) {
	err := github.CloseIssue(owner, repo, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
}

func gopen(owner, repo string, issueNumber int) {
	err := github.OpenIssue(owner, repo, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
}

func create(owner, repo string) {
	input, err := editor.Edit("")
	if err != nil {
		log.Fatal(err)
	}
	title, body := parseText(input)
	err = github.CreateIssue(owner, repo, title, body)
	if err != nil {
		log.Fatal(err)
	}
}

func parseText(input []byte) (string, string) {
	content := strings.Split(string(input), "\n")
	title := content[0]
	body := strings.TrimSpace(strings.Join(content[1:], "\n"))
	return title, body
}

func update(owner, repo string, issueNumber int) {
	issue, err := github.GetIssue(owner, repo, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
	content := fmt.Sprintf("%s\n\n%s", string(issue.Title), string(issue.Body))
	input, err := editor.Edit(content)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(input))
	title, body := parseText(input)
	err = github.UpdateIssue(owner, repo, title, body, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
}
