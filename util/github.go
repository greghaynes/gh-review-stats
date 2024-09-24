package util

import (
	"context"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// NewGithubClient creates a client for communicating with the GitHub
// API using the provided token.
func NewGithubClient(ctx context.Context, token, enterpriseBaseURL, enterpriseUploadURL string) (*github.Client, error) {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	oauthClient := oauth2.NewClient(ctx, tokenSource)

	if enterpriseBaseURL != "" {
		// If upload url is unset then we should set it to base URL
		if enterpriseUploadURL == "" {
			enterpriseUploadURL = enterpriseBaseURL
		}
		return github.NewEnterpriseClient(enterpriseBaseURL, enterpriseUploadURL, oauthClient)
	}

	return github.NewClient(oauthClient), nil
}
