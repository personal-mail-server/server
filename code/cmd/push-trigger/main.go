package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"personal-mail-server/internal/automation/pushtrigger"
)

func main() {
	base := flag.String("base", "main", "target base branch")
	retries := flag.Int("check-retries", 12, "number of retries while waiting for PR checks")
	interval := flag.Duration("retry-interval", 5*time.Second, "interval between PR check retries")
	flag.Parse()

	service := pushtrigger.NewService(pushtrigger.ExecRunner{}, nil, pushtrigger.Config{
		BaseBranch:    *base,
		CheckRetries:  *retries,
		RetryInterval: *interval,
	})

	pr, err := service.Execute(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "push trigger failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("push trigger completed\n")
	fmt.Printf("pr_number=%d\n", pr.Number)
	fmt.Printf("pr_url=%s\n", pr.URL)
}
