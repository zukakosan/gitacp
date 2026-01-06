package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// apiKey := os.Getenv("OPENAI_API_KEY")
	// if apiKey == "" {
	// 	fmt.Println("Error: OPENAI_API_KEY environment variable not set")
	// 	os.Exit(1)
	// }
	diff, err := getGitDiff()
	if err != nil {
		fmt.Println("Error getting git diff:", err)
		os.Exit(1)
	}
	fmt.Println(diff)
	add, err := gitAddAll()
	if err != nil {
		fmt.Println("Error adding changes:", err)
		os.Exit(1)
	}
	fmt.Println(add)
	commit, err := gitCommit("test commit message")
	if err != nil {
		fmt.Println("Error committing changes:", err)
		os.Exit(1)
	}
	fmt.Println(commit)
	head := exec.Command("git", "diff")
	fmt.Println(head.Output())
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "HEAD")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	return string(output), nil
}

func gitAddAll() (string, error) {
	cmd := exec.Command("git", "add", ".")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to add changes: %w", err)
	}
	return string(output), nil
}

func gitCommit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to commit changes: %w", err)
	}
	return string(output), nil
}

func gitPuch(branch string) (string, error) {
	cmd := exec.Command("git", "push", "origin", branch)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to push changes: %w", err)
	}
	return string(output), nil
}
