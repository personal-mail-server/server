package pushtrigger

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"
)

type fakeRunner struct {
	responses map[string][]fakeResponse
	calls     []string
}

type fakeResponse struct {
	out string
	err error
}

func normalizeKey(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func (f *fakeRunner) Run(_ context.Context, name string, args ...string) (string, error) {
	key := normalizeKey(name + " " + strings.Join(args, " "))
	f.calls = append(f.calls, key)
	responses, ok := f.responses[key]
	if !ok || len(responses) == 0 {
		return "", fmt.Errorf("unexpected command: %s", key)
	}
	f.responses[key] = responses[1:]
	return responses[0].out, responses[0].err
}

func TestExecuteRejectsProtectedBranch(t *testing.T) {
	runner := &fakeRunner{responses: map[string][]fakeResponse{
		normalizeKey("git branch --show-current"): {{out: "main\n"}},
	}}
	service := NewService(runner, nil, Config{BaseBranch: "main"})

	_, err := service.Execute(context.Background())
	if err == nil || !errors.Is(err, ErrProtectedBranch) {
		t.Fatalf("expected protected branch error, got %v", err)
	}
}

func TestExecuteRejectsDirtyWorktree(t *testing.T) {
	runner := &fakeRunner{responses: map[string][]fakeResponse{
		normalizeKey("git branch --show-current"): {{out: "feature/test\n"}},
		normalizeKey("git status --short"):        {{out: " M README.md\n"}},
	}}
	service := NewService(runner, nil, Config{BaseBranch: "main"})

	_, err := service.Execute(context.Background())
	if err == nil || !errors.Is(err, ErrDirtyWorktree) {
		t.Fatalf("expected dirty worktree error, got %v", err)
	}
}

func TestExecuteCreatesPRAndEnablesAutoMerge(t *testing.T) {
	createCommand := normalizeKey("gh pr create --base main --head feature/test --title feat: add push trigger --body ## 요약\n- feat: add push trigger\n\n## 기준 정보\n- 생성 방식: 예약 트리거 `PUSH`\n- 대상 브랜치: `main`\n- 작업 브랜치: `feature/test`\n")
	runner := &fakeRunner{responses: map[string][]fakeResponse{
		normalizeKey("git branch --show-current"):       {{out: "feature/test\n"}},
		normalizeKey("git status --short"):              {{out: ""}},
		normalizeKey("git remote get-url origin"):       {{out: "git@github.com:owner/repo.git\n"}},
		normalizeKey("gh --version"):                    {{out: "gh version 2.0.0\n"}},
		normalizeKey("gh auth status"):                  {{out: "logged in\n"}},
		normalizeKey("git push -u origin feature/test"): {{out: "pushed\n"}},
		normalizeKey("gh pr view feature/test --json number,url,baseRefName,state"): {
			{err: errors.New("not found")},
			{out: `{"number":12,"url":"https://github.com/owner/repo/pull/12","baseRefName":"main","state":"OPEN"}`},
		},
		normalizeKey("git log -1 --pretty=%s"):         {{out: "feat: add push trigger\n"}},
		normalizeKey("git log --format=%s main..HEAD"): {{out: "feat: add push trigger\n"}},
		createCommand:                                 {{out: "https://github.com/owner/repo/pull/12\n"}},
		normalizeKey("gh pr checks 12"):               {{out: "CI\tpass\n"}},
		normalizeKey("gh pr merge 12 --auto --merge"): {{out: "auto-merge enabled\n"}},
	}}
	service := NewService(runner, func(time.Duration) {}, Config{BaseBranch: "main", CheckRetries: 1})

	pr, err := service.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if pr.Number != 12 {
		t.Fatalf("expected pr number 12, got %d", pr.Number)
	}
	if got := runner.calls[len(runner.calls)-1]; got != normalizeKey("gh pr merge 12 --auto --merge") {
		t.Fatalf("unexpected final command: %s", got)
	}
}

func TestExecuteReusesExistingPR(t *testing.T) {
	runner := &fakeRunner{responses: map[string][]fakeResponse{
		normalizeKey("git branch --show-current"):                                   {{out: "feature/test\n"}},
		normalizeKey("git status --short"):                                          {{out: ""}},
		normalizeKey("git remote get-url origin"):                                   {{out: "git@github.com:owner/repo.git\n"}},
		normalizeKey("gh --version"):                                                {{out: "gh version 2.0.0\n"}},
		normalizeKey("gh auth status"):                                              {{out: "logged in\n"}},
		normalizeKey("git push -u origin feature/test"):                             {{out: "pushed\n"}},
		normalizeKey("gh pr view feature/test --json number,url,baseRefName,state"): {{out: `{"number":7,"url":"https://github.com/owner/repo/pull/7","baseRefName":"main","state":"OPEN"}`}},
		normalizeKey("gh pr checks 7"):                                              {{out: "CI\tpending\n"}},
		normalizeKey("gh pr merge 7 --auto --merge"):                                {{out: "auto-merge enabled\n"}},
	}}
	service := NewService(runner, func(time.Duration) {}, Config{BaseBranch: "main", CheckRetries: 1})

	pr, err := service.Execute(context.Background())
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if pr.Number != 7 {
		t.Fatalf("expected existing pr number 7, got %d", pr.Number)
	}
	for _, call := range runner.calls {
		if strings.HasPrefix(call, normalizeKey("gh pr create")) {
			t.Fatalf("did not expect PR creation when existing PR is reusable")
		}
	}
}
