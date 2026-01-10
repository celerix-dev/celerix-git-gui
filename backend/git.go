package backend

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type RepoStats struct {
	RepoName      string    `json:"repoName"`
	RemoteURL     string    `json:"remoteUrl"`
	SizeMB        float64   `json:"sizeMb"`
	CommitCount   int       `json:"commitCount"`
	LastCommit    time.Time `json:"lastCommit"`
	FirstCommit   time.Time `json:"firstCommit"`
	IsClean       bool      `json:"isClean"`
	ModifiedFiles []string  `json:"modifiedFiles"`
}

func (a *App) GitInit(path string) error {
	_, err := git.PlainInit(path, false)
	return err
}

func (a *App) IsGitRepo(path string) (bool, error) {
	_, err := git.PlainOpen(path)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (a *App) GetRepoStats(path string) (*RepoStats, error) {
	// 1. Open the repository
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	stats := &RepoStats{
		RepoName: filepath.Base(path),
	}

	indexPath := filepath.Join(path, ".git", "index")
	f, err := os.Open(indexPath)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	// 2. Use go-git's index decoder
	idx := &index.Index{}
	decoder := index.NewDecoder(f)
	if err := decoder.Decode(idx); err != nil {
		return nil, err
	}

	// 3. Sum the sizes (this is just reading metadata, no decompression!)
	var totalBytes int64
	for _, entry := range idx.Entries {
		totalBytes += int64(entry.Size)
	}

	stats.SizeMB = float64(totalBytes) / 1024 / 1024

	// 3. Get Remote URL
	if remotes, err := r.Remotes(); err == nil && len(remotes) > 0 {
		stats.RemoteURL = remotes[0].Config().URLs[0]
	}

	// 4. History Stats (Count, First, Last)
	cIter, err := r.Log(&git.LogOptions{})
	if err == nil {
		var count int
		_ = cIter.ForEach(func(c *object.Commit) error {
			if count == 0 {
				stats.LastCommit = c.Author.When
			}
			stats.FirstCommit = c.Author.When
			count++
			return nil
		})
		stats.CommitCount = count
	}

	// 5. Worktree Status (Uncommitted changes)
	w, err := r.Worktree()
	if err == nil {
		status, err := w.Status()
		if err == nil {
			stats.IsClean = status.IsClean()
			for file, s := range status {
				// We only add files that aren't "Unmodified"
				if s.Worktree != git.Unmodified || s.Staging != git.Unmodified {
					stats.ModifiedFiles = append(stats.ModifiedFiles, file)
				}
			}
		}
	}

	return stats, nil
}

func (a *App) GetRepoReadme(path string) (string, error) {
	readmeFiles := []string{"README.md", "readme.md", "README", "readme"}
	var content []byte

	for _, name := range readmeFiles {
		p := filepath.Join(path, name)
		if _, err := os.Stat(p); err == nil {
			c, err := os.ReadFile(p)
			if err == nil {
				content = c
				break
			}
		}
	}

	if content == nil {
		return "", nil // No readme found
	}

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)

	var buf strings.Builder
	if err := md.Convert(content, &buf); err != nil {
		return "", fmt.Errorf("failed to convert markdown: %w", err)
	}

	return buf.String(), nil
}
