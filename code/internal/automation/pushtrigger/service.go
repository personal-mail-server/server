package pushtrigger

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrDirtyWorktree     = errors.New("working tree has uncommitted changes")
	ErrProtectedBranch   = errors.New("refusing to run on protected branch")
	ErrMissingOrigin     = errors.New("git remote 'origin' is not configured")
	ErrMissingChecks     = errors.New("pull request checks did not appear in time")
	ErrUnsupportedPR     = errors.New("existing pull request is not reusable")
	ErrMissingBranchName = errors.New("current branch name is empty")
)

type Runner interface {
	Run(ctx context.Context, name string, args ...string) (string, error)
}

type Sleeper func(time.Duration)

type Config struct {
	BaseBranch    string
	CheckRetries  int
	RetryInterval time.Duration
}

type PRInfo struct {
	Number int
	URL    string
}

type Service struct {
	runner  Runner
	sleeper Sleeper
	config  Config
}

func NewService(runner Runner, sleeper Sleeper, config Config) *Service {
	if config.BaseBranch == "" {
		config.BaseBranch = "main"
	}
	if config.CheckRetries <= 0 {
		config.CheckRetries = 12
	}
	if config.RetryInterval <= 0 {
		config.RetryInterval = 5 * time.Second
	}
	if sleeper == nil {
		sleeper = time.Sleep
	}
	return &Service{runner: runner, sleeper: sleeper, config: config}
}

func (s *Service) Execute(ctx context.Context) (*PRInfo, error) {
	branch, err := s.currentBranch(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.ensureSafeBranch(branch); err != nil {
		return nil, err
	}
	if err := s.ensureCleanWorktree(ctx); err != nil {
		return nil, err
	}
	if err := s.ensureOrigin(ctx); err != nil {
		return nil, err
	}
	if err := s.ensureGitHubCLI(ctx); err != nil {
		return nil, err
	}
	if err := s.ensureGitHubAuth(ctx); err != nil {
		return nil, err
	}

	if _, err := s.runner.Run(ctx, "git", "push", "-u", "origin", branch); err != nil {
		return nil, fmt.Errorf("push branch: %w", err)
	}

	pr, err := s.reuseOrCreatePR(ctx, branch)
	if err != nil {
		return nil, err
	}

	if err := s.waitForChecks(ctx, pr.Number); err != nil {
		return nil, err
	}

	if _, err := s.runner.Run(ctx, "gh", "pr", "merge", fmt.Sprintf("%d", pr.Number), "--auto", "--merge"); err != nil {
		return nil, fmt.Errorf("enable auto-merge: %w", err)
	}

	return pr, nil
}

func (s *Service) currentBranch(ctx context.Context) (string, error) {
	out, err := s.runner.Run(ctx, "git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("read current branch: %w", err)
	}
	branch := strings.TrimSpace(out)
	if branch == "" {
		return "", ErrMissingBranchName
	}
	return branch, nil
}

func (s *Service) ensureSafeBranch(branch string) error {
	if branch == "main" || branch == "master" {
		return fmt.Errorf("%w: %s", ErrProtectedBranch, branch)
	}
	return nil
}

func (s *Service) ensureCleanWorktree(ctx context.Context) error {
	out, err := s.runner.Run(ctx, "git", "status", "--short")
	if err != nil {
		return fmt.Errorf("check worktree: %w", err)
	}
	if strings.TrimSpace(out) != "" {
		return ErrDirtyWorktree
	}
	return nil
}

func (s *Service) ensureOrigin(ctx context.Context) error {
	out, err := s.runner.Run(ctx, "git", "remote", "get-url", "origin")
	if err != nil {
		return fmt.Errorf("check origin remote: %w", err)
	}
	if strings.TrimSpace(out) == "" {
		return ErrMissingOrigin
	}
	return nil
}

func (s *Service) ensureGitHubCLI(ctx context.Context) error {
	if _, err := s.runner.Run(ctx, "gh", "--version"); err != nil {
		return fmt.Errorf("check gh cli: %w", err)
	}
	return nil
}

func (s *Service) ensureGitHubAuth(ctx context.Context) error {
	if _, err := s.runner.Run(ctx, "gh", "auth", "status"); err != nil {
		return fmt.Errorf("check gh auth: %w", err)
	}
	return nil
}

func (s *Service) reuseOrCreatePR(ctx context.Context, branch string) (*PRInfo, error) {
	viewOut, viewErr := s.runner.Run(ctx, "gh", "pr", "view", branch, "--json", "number,url,baseRefName,state")
	if viewErr == nil {
		pr, err := parsePRView(viewOut)
		if err != nil {
			return nil, err
		}
		if pr.BaseRefName != s.config.BaseBranch || pr.State != "OPEN" {
			return nil, ErrUnsupportedPR
		}
		return &PRInfo{Number: pr.Number, URL: pr.URL}, nil
	}

	title, body, err := s.buildPRContent(ctx, branch)
	if err != nil {
		return nil, err
	}

	createOut, err := s.runner.Run(ctx, "gh", "pr", "create", "--base", s.config.BaseBranch, "--head", branch, "--title", title, "--body", body)
	if err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}
	if strings.TrimSpace(createOut) == "" {
		return nil, errors.New("pull request url is empty")
	}

	viewOut, err = s.runner.Run(ctx, "gh", "pr", "view", branch, "--json", "number,url,baseRefName,state")
	if err != nil {
		return nil, fmt.Errorf("reload pull request: %w", err)
	}
	pr, err := parsePRView(viewOut)
	if err != nil {
		return nil, err
	}
	return &PRInfo{Number: pr.Number, URL: pr.URL}, nil
}

func (s *Service) buildPRContent(ctx context.Context, branch string) (string, string, error) {
	subject, err := s.runner.Run(ctx, "git", "log", "-1", "--pretty=%s")
	if err != nil {
		return "", "", fmt.Errorf("read latest commit subject: %w", err)
	}
	title := strings.TrimSpace(subject)
	if title == "" {
		title = fmt.Sprintf("%s 변경 사항", branch)
	}

	logOut, err := s.runner.Run(ctx, "git", "log", "--format=%s", fmt.Sprintf("%s..HEAD", s.config.BaseBranch))
	if err != nil {
		logOut = title
	}
	lines := splitNonEmptyLines(logOut)
	if len(lines) == 0 {
		lines = []string{title}
	}

	var builder strings.Builder
	builder.WriteString("## 요약\n")
	for _, line := range lines {
		builder.WriteString("- ")
		builder.WriteString(line)
		builder.WriteString("\n")
	}
	builder.WriteString("\n## 기준 정보\n")
	builder.WriteString("- 생성 방식: 예약 트리거 `PUSH`\n")
	builder.WriteString("- 대상 브랜치: `")
	builder.WriteString(s.config.BaseBranch)
	builder.WriteString("`\n")
	builder.WriteString("- 작업 브랜치: `")
	builder.WriteString(branch)
	builder.WriteString("`\n")

	return title, builder.String(), nil
}

func (s *Service) waitForChecks(ctx context.Context, number int) error {
	for attempt := 0; attempt < s.config.CheckRetries; attempt++ {
		out, err := s.runner.Run(ctx, "gh", "pr", "checks", fmt.Sprintf("%d", number))
		if err == nil && hasChecks(out) {
			return nil
		}
		if attempt < s.config.CheckRetries-1 {
			s.sleeper(s.config.RetryInterval)
		}
	}
	return ErrMissingChecks
}

func splitNonEmptyLines(s string) []string {
	parts := strings.Split(strings.ReplaceAll(s, "\r\n", "\n"), "\n")
	lines := make([]string, 0, len(parts))
	for _, part := range parts {
		line := strings.TrimSpace(part)
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func hasChecks(out string) bool {
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return false
	}
	if strings.Contains(strings.ToLower(trimmed), "no checks") {
		return false
	}
	return true
}
