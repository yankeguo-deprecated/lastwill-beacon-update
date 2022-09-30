package main

import (
	"context"
	"errors"
	"flag"
	"github.com/google/go-github/v47/github"
	"github.com/guoyk93/grace"
	"github.com/guoyk93/grace/gracemain"
	"golang.org/x/oauth2"
	"os"
	"strings"
	"time"
)

func main() {
	var (
		err    error
		ctx, _ = gracemain.WithSignalCancel(context.Background())
	)
	defer gracemain.Exit(&err)
	defer grace.Guard(&err)

	var (
		optOwner  string
		optRepo   string
		optPath   string
		optBranch string
	)

	flag.StringVar(&optOwner, "owner", "guoyk93", "github repository owner")
	flag.StringVar(&optRepo, "repo", "lastwill", "github repository name")
	flag.StringVar(&optPath, "path", "beacon.txt", "path to beacon file")
	flag.StringVar(&optBranch, "branch", "main", "branch name")
	flag.Parse()

	envToken := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))

	if envToken == "" {
		err = errors.New("missing environment variable $GITHUB_TOKEN")
		return
	}

	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: envToken},
	)))

	content, _, _ := grace.Must3(
		client.Repositories.GetContents(
			ctx,
			optOwner, optRepo, optPath, &github.RepositoryContentGetOptions{
				Ref: optBranch,
			},
		),
	)

	newContent := time.Now().Format(time.RFC3339)

	_, _ = grace.Must2(client.Repositories.UpdateFile(ctx, optOwner, optRepo, optPath, &github.RepositoryContentFileOptions{
		Message: grace.Ptr("update beacon"),
		Content: []byte(newContent),
		SHA:     content.SHA,
		Branch:  grace.Ptr(optBranch),
	}))
}
