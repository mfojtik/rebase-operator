package fork

import (
	"context"
	"fmt"

	"github.com/google/go-github/v32/github"
)

func PrepareBranch(ctx context.Context, client *github.Client, tagName, branchName string) error {
	/*
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := github.NewClient(tc)

	*/

	tag, _, err := client.Git.GetRef(context.TODO(), "kubernetes", "kubernetes", fmt.Sprintf("tags/%s", tagName))
	if err != nil {
		return fmt.Errorf("failed to fetch tag %q: %v", tagName, err)
	}

	tagRef, _, err := client.Git.GetTag(ctx, "kubernetes", "kubernetes", tag.GetObject().GetSHA())
	if err != nil {
		return fmt.Errorf("failed to fetch tag %q reference %q: %v", tagName, tag.GetObject().GetSHA(), err)
	}

	target := fmt.Sprintf("refs/heads/%s", branchName)

	if _, _, err := client.Git.CreateRef(context.TODO(), "mfojtik", "kubernetes", &github.Reference{
		Ref:    &target,
		Object: tagRef.Object,
	}); err != nil {
		return fmt.Errorf("failed to create branch %q based on tag %q: %v", branchName, tagName, err)
	}
	return nil
}
