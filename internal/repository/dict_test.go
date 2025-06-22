package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Namchee/actions-case-police/internal"
	"github.com/Namchee/actions-case-police/internal/entity"
	"github.com/google/go-github/v43/github"
	"github.com/stretchr/testify/assert"
)

type githubClientMock struct {}

func (c *githubClientMock) GetIssue(ctx context.Context, meta *entity.Meta, number int) (*github.Issue, error) { return nil, nil }
func (c *githubClientMock) EditIssue(ctx context.Context, meta *entity.Meta, number int, issue *entity.IssueData) error { return nil }
func (c *githubClientMock) GetRepositoryContents(ctx context.Context, meta *entity.Meta, path string) (string, error) {
	mappings := map[string]string{
		"packages/case-police/dict/foo.json": `{
			"vscode": "VS Code"
		}`,
		"packages/case-police/dict/bar.json": `{
			"wifi": "Wi-Fi"
		}`,
	}

	if (mappings[path] != "") {
		return mappings[path], nil
	} else {
		return "", errors.New("Failed to get contents of " + path)
	}
}

func TestGetDictionary(t *testing.T) {
	tests := []struct {
		name     string
		files  	[]string
		clientMock internal.GithubClient
		want     map[string]string
	}{
		{
			name: "should parse single dictionary",
			files: []string{"foo"},
			want: map[string]string{
				"vscode": "VS Code",
			},
		},
		{
			name: "should combine dictionaries",
			files: []string{"foo", "bar"},
			want: map[string]string{
				"vscode": "VS Code",
				"wifi":   "Wi-Fi",
			},
		},
		{
			name: "should return empty dictionary on error",
			files: []string{"baz"},
			want: make(map[string]string),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			client := &githubClientMock{}

			got := GetDictionary(context.TODO(), client, tc.files)

			assert.Equal(t, tc.want, got)
		})
	}
}
