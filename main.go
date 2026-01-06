package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	paths := os.Args[1:]
	add, err := gitAdd(paths...)
	if err != nil {
		fmt.Println("Error adding changes:", err)
		os.Exit(1)
	}
	// fmt.Println(add)
	_ = add

	diff, err := getGitDiff()
	if err != nil {
		fmt.Println("Error getting git diff:", err)
		os.Exit(1)
	}
	// fmt.Println(diff)

	commitMsg, err := generateCommitMessage(diff)
	if err != nil {
		fmt.Println("Error generating commit message:", err)
		os.Exit(1)
	}
	fmt.Println("Generated Commit Message:", commitMsg)

	commit, err := gitCommit(commitMsg)
	if err != nil {
		fmt.Println("Error committing changes:", err)
		os.Exit(1)
	}
	fmt.Println(commit)
	_ = commit

	push, err := gitPush()
	if err != nil {
		fmt.Println("Error pushing changes:", err)
		os.Exit(1)
	}
	fmt.Println(push)
	_ = push

	// head := exec.Command("git", "diff")
	// fmt.Println(head.Output())
}

func gitAdd(paths ...string) (string, error) {
	args := append([]string{"add"}, paths...)
	cmd := exec.Command("git", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to add changes: %w", err)
	}
	return string(output), nil
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
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

func gitPush() (string, error) {
	branch := exec.Command("git", "branch", "--show-current")
	branchOutput, err := branch.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	branchName := strings.TrimSpace(string(branchOutput))
	cmd := exec.Command("git", "push", "origin", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to push changes: %w", err)
	}
	return string(output), nil
}

func generateCommitMessage(diff string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	baseURL := os.Getenv("BASE_URL")
	modelName := os.Getenv("MODEL_NAME")
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)
	if apiKey == "" || baseURL == "" || modelName == "" {
		fmt.Println("Error: At least one environment variable is not set")
		os.Exit(1)
	}
	prompt := fmt.Sprintf(`
	# Task
	Please generate a concise git commit message based on the following git diff Head.
	# git diff HEAD result
	%s
	`, diff)
	ctx := context.Background()
	chatCompletion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a helpful assistant that generates git commit messages based on git diff HEAD results."),
			openai.UserMessage(prompt),
		},
		Model: modelName,
	})
	if err != nil {
		return "", fmt.Errorf("Error generating commit message: %w", err)
	}
	return chatCompletion.Choices[0].Message.Content, nil
}
