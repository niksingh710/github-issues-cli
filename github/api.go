package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
	UdpatedAt time.Time `json:"udpated_at"`
}

func getIssueUrl(owner, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
}

func getIssueUrlById(owner, repo string, issueNumber int) string {
	return fmt.Sprintf("%s/%d", getIssueUrl(owner, repo), issueNumber)
}

func setAuth(request *http.Request) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN variable is not set in env or contains a blank value")
	}
	request.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
}

func CreateIssue(owner, repo, title, body string) error {
	requestBody, err := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})
	if err != nil {
		log.Fatal(err)
	}
	request, err := http.NewRequest("POST", getIssueUrl(owner, repo), bytes.NewBuffer(requestBody))
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")
	setAuth(request)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		msg := ""
		if response.StatusCode == 401 {
			msg = "Maybe Invalid GITHUB_TOKEN"
		}
		return fmt.Errorf("failed to create issue: %s\n\t%s", response.Status, msg)
	}
	for key, value := range response.Header {
		if key == "Location" {
			fmt.Fprintf(os.Stdout, "Successfully Created issue at: %s\n", value[0])
		}
	}
	return nil
}

func patchIssue(owner, repo string, issueNumber int, values map[string]string) error {
	requestBody, err := json.Marshal(values)
	if err != nil {
		log.Fatal(err)
	}
	request, err := http.NewRequest(
		"PATCH",
		getIssueUrlById(owner, repo, issueNumber),
		bytes.NewBuffer(requestBody),
	)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")
	setAuth(request)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		msg := ""
		if response.StatusCode == 401 {
			msg = "Maybe Invalid GITHUB_TOKEN"
		}
		return fmt.Errorf("failed for issue: %s\n\t%s", response.Status, msg)
	}
	for key, value := range response.Header {
		if key == "Location" {
			fmt.Fprintf(os.Stdout, "check issue at: %s\n", value[0])
		}
	}
	return nil
}

func CloseIssue(owner, repo string, issueNumber int) error {
	values := map[string]string{
		"state": "closed",
	}
	return patchIssue(owner, repo, issueNumber, values)
}

func OpenIssue(owner, repo string, issueNumber int) error {
	values := map[string]string{
		"state": "open",
	}
	return patchIssue(owner, repo, issueNumber, values)
}

func UpdateIssue(owner, repo, title, body string, issueNumber int) error {
	values := map[string]string{
		"title": title,
		"body":  body,
	}
	return patchIssue(owner, repo, issueNumber, values)
}

func GetIssue(owner, repo string, issueNumber int) (*Issue, error) {
	request, err := http.NewRequest("PATCH", getIssueUrlById(owner, repo, issueNumber), nil)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("Content-Type", "application/json")
	setAuth(request)
	if err != nil {
		log.Fatal(err)
	}

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get issue %d: %s", issueNumber, response.Status)
	}

	var issue Issue
	if err := json.NewDecoder(response.Body).Decode(&issue); err != nil {
		return nil, err
	}
	return &issue, nil
}
