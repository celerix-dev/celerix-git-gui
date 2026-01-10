package backend

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"io"
)

type GitRemoteBranches struct {
	Name     string   `json:"name"`
	Branches []string `json:"branches"`
}

type RepoStats struct {
	RepoName      string              `json:"repoName"`
	RemoteURL     string              `json:"remoteUrl"`
	SizeMB        float64             `json:"sizeMb"`
	CommitCount   int                 `json:"commitCount"`
	LastCommit    time.Time           `json:"lastCommit"`
	FirstCommit   time.Time           `json:"firstCommit"`
	IsClean       bool                `json:"isClean"`
	ModifiedFiles []string            `json:"modifiedFiles"`
	Branches      []string            `json:"branches"`
	Remotes       []GitRemoteBranches `json:"remotes"`
	Tags          []string            `json:"tags"`
	Stashes       []string            `json:"stashes"`
	CurrentBranch string              `json:"currentBranch"`
}

type GitStatusFile struct {
	Path     string `json:"path"`
	Status   string `json:"status"`
	IsStaged bool   `json:"is_staged"`
}

type GitCommit struct {
	Hash         string    `json:"hash"`
	AuthorName   string    `json:"authorName"`
	AuthorEmail  string    `json:"authorEmail"`
	Date         time.Time `json:"date"`
	Subject      string    `json:"subject"`
	Body         string    `json:"body"`
	ParentHashes []string  `json:"parentHashes"`
	Refs         []string  `json:"refs"`
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

	// 6. Branches
	branchIter, err := r.Branches()
	if err == nil {
		_ = branchIter.ForEach(func(ref *plumbing.Reference) error {
			stats.Branches = append(stats.Branches, ref.Name().Short())
			return nil
		})
	}

	// 6.1 Current Branch
	head, err := r.Head()
	if err == nil && head.Name().IsBranch() {
		stats.CurrentBranch = head.Name().Short()
	}

	// 7. Remotes and Remote Branches
	remotes, err := r.Remotes()
	if err == nil {
		remoteRefs, _ := r.References()
		for _, remote := range remotes {
			remoteName := remote.Config().Name
			rb := GitRemoteBranches{
				Name: remoteName,
			}

			// Find branches for this remote
			if remoteRefs != nil {
				_ = remoteRefs.ForEach(func(ref *plumbing.Reference) error {
					if ref.Name().IsRemote() && strings.HasPrefix(ref.Name().Short(), remoteName+"/") {
						rb.Branches = append(rb.Branches, ref.Name().Short())
					}
					return nil
				})
				// Reset iterator for next remote
				remoteRefs, _ = r.References()
			}

			stats.Remotes = append(stats.Remotes, rb)
		}
	}

	// 8. Tags
	tagIter, err := r.Tags()
	if err == nil {
		_ = tagIter.ForEach(func(ref *plumbing.Reference) error {
			stats.Tags = append(stats.Tags, ref.Name().Short())
			return nil
		})
	}

	// 9. Stashes
	stashRef, err := r.Storer.Reference(plumbing.ReferenceName("refs/stash"))
	if err == nil {
		stashCommit, err := r.CommitObject(stashRef.Hash())
		if err == nil {
			// In git, stashes are a stack. refs/stash points to the latest.
			// It has parents: the previous stash and the commit it was made on.
			// This is a bit complex to traverse manually with go-git to get ALL stashes.
			// For now, let's just show the latest one if it exists.
			stats.Stashes = append(stats.Stashes, strings.TrimSpace(stashCommit.Message))
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

func (a *App) GetGitStatus(repoPath string) ([]GitStatusFile, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}
	status, err := w.Status()
	if err != nil {
		return nil, err
	}

	var result []GitStatusFile
	for file, s := range status {
		// Staged changes
		if s.Staging != git.Unmodified && s.Staging != git.Untracked {
			result = append(result, GitStatusFile{
				Path:     file,
				Status:   statusChar(s.Staging),
				IsStaged: true,
			})
		}
		// Unstaged changes (including untracked)
		if s.Worktree != git.Unmodified {
			result = append(result, GitStatusFile{
				Path:     file,
				Status:   statusChar(s.Worktree),
				IsStaged: false,
			})
		}
	}
	return result, nil
}

func statusChar(s git.StatusCode) string {
	switch s {
	case git.Unmodified:
		return " "
	case git.Modified:
		return "M"
	case git.Added:
		return "A"
	case git.Deleted:
		return "D"
	case git.Renamed:
		return "R"
	case git.Copied:
		return "C"
	case git.UpdatedButUnmerged:
		return "U"
	case git.Untracked:
		return "?"
	default:
		return "?"
	}
}

func (a *App) StageFile(repoPath string, filePath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	_, err = w.Add(filePath)
	return err
}

func (a *App) StageAll(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	return w.AddWithOptions(&git.AddOptions{All: true})
}

func (a *App) UnstageFile(repoPath string, filePath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	// For unstaging, we use Reset
	head, err := r.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			// No commits yet, so we probably want to remove from index?
			// go-git doesn't have a direct "unstage" that works perfectly without HEAD
			// But we can try to Reset to an empty tree if needed.
			// Actually w.Remove often deletes the file from worktree too if not careful.
			return fmt.Errorf("cannot unstage: no HEAD reference found")
		}
		return err
	}

	return w.Reset(&git.ResetOptions{
		Commit: head.Hash(),
		Files:  []string{filePath},
	})
}

func (a *App) UnstageAll(repoPath string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	head, err := r.Head()
	if err != nil {
		if errors.Is(err, plumbing.ErrReferenceNotFound) {
			// If no HEAD, maybe we can't reset?
			// In git CLI "git reset" on initial repo works.
			// go-git Reset requires a commit hash usually.
			return fmt.Errorf("cannot unstage all: no HEAD reference found")
		}
		return err
	}

	return w.Reset(&git.ResetOptions{
		Commit: head.Hash(),
		Mode:   git.MixedReset,
	})
}

func (a *App) GetFileDiff(repoPath string, filePath string, staged bool) (string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return "", err
	}

	var fromTree *object.Tree
	head, headErr := r.Head()

	if staged {
		// Staged: Diff between HEAD and Index
		if headErr == nil {
			headCommit, err := r.CommitObject(head.Hash())
			if err != nil {
				return "", err
			}
			fromTree, err = headCommit.Tree()
			if err != nil {
				return "", err
			}
		} else if !errors.Is(headErr, plumbing.ErrReferenceNotFound) {
			return "", headErr
		}

		// For staged, we want to see what's in the index vs what was in HEAD
		// This is actually what 'git diff --cached' does.
		return a.getUnifiedDiff(r, fromTree, filePath, true)
	} else {
		// Unstaged: Diff between Index and Worktree
		// This is what 'git diff' does.
		return a.getUnifiedDiff(r, nil, filePath, false)
	}
}

func (a *App) getUnifiedDiff(r *git.Repository, headTree *object.Tree, filePath string, staged bool) (string, error) {
	w, err := r.Worktree()
	if err != nil {
		return "", err
	}

	status, err := w.Status()
	if err != nil {
		return "", err
	}

	fileStatus := status.File(filePath)

	var oldContent string
	var newContent string

	if staged {
		// Staged: Diff between HEAD and Index
		// Old content is from HEAD
		if headTree != nil {
			entry, err := headTree.File(filePath)
			if err == nil {
				oldContent, _ = entry.Contents()
			}
		}

		// New content is from Index
		if fileStatus.Staging != git.Deleted {
			idx, err := r.Storer.Index()
			if err == nil {
				for _, entry := range idx.Entries {
					if entry.Name == filePath {
						blob, err := r.BlobObject(entry.Hash)
						if err == nil {
							reader, err := blob.Reader()
							if err == nil {
								data, _ := io.ReadAll(reader)
								newContent = string(data)
								_ = reader.Close()
							}
						}
						break
					}
				}
			}
		}
	} else {
		// Unstaged: Diff between Index and Worktree
		// New content is from Worktree (file on disk)
		if fileStatus.Worktree != git.Deleted {
			fullPath := filepath.Join(w.Filesystem.Root(), filePath)
			data, err := os.ReadFile(fullPath)
			if err == nil {
				newContent = string(data)
			}
		}

		// Old content is from Index
		idx, err := r.Storer.Index()
		if err == nil {
			for _, entry := range idx.Entries {
				if entry.Name == filePath {
					blob, err := r.BlobObject(entry.Hash)
					if err == nil {
						reader, err := blob.Reader()
						if err == nil {
							data, _ := io.ReadAll(reader)
							oldContent = string(data)
							_ = reader.Close()
						}
					}
					break
				}
			}
		}

		// If not in index (e.g. newly tracked but not yet in a commit? no, if it's in index it's in index)
		// If it's a new file, oldContent will be empty, which is correct.
	}

	// For now, return a simple diff-like format if we have both contents
	// We'll improve this as we go.
	if oldContent == "" && newContent == "" && fileStatus.Worktree != git.Deleted && fileStatus.Staging != git.Deleted {
		return "Binary file or empty", nil
	}

	return a.generateSimpleDiff(oldContent, newContent, filePath), nil
}

func (a *App) GetBranches(repoPath string) ([]string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	iter, err := r.Branches()
	if err != nil {
		return nil, err
	}

	var branches []string
	err = iter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (a *App) GetCommitHistory(repoPath string, count int) ([]GitCommit, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	// Get all references to map them to commits later
	refMap := make(map[plumbing.Hash][]string)
	refs, _ := r.References()
	if refs != nil {
		_ = refs.ForEach(func(ref *plumbing.Reference) error {
			if ref.Type() == plumbing.HashReference || ref.Type() == plumbing.SymbolicReference {
				hash := ref.Hash()
				// For symbolic refs (like HEAD), resolve to actual hash
				if ref.Type() == plumbing.SymbolicReference {
					resolved, err := r.Reference(ref.Target(), true)
					if err == nil {
						hash = resolved.Hash()
					}
				}
				name := ref.Name().Short()
				refMap[hash] = append(refMap[hash], name)
			}
			return nil
		})
	}

	cIter, err := r.Log(&git.LogOptions{Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}

	var commits []GitCommit
	err = cIter.ForEach(func(c *object.Commit) error {
		if count > 0 && len(commits) >= count {
			return io.EOF
		}

		var parents []string
		for _, ph := range c.ParentHashes {
			parents = append(parents, ph.String())
		}

		subject := strings.Split(c.Message, "\n")[0]
		body := ""
		if strings.Contains(c.Message, "\n") {
			body = strings.TrimSpace(c.Message[strings.Index(c.Message, "\n"):])
		}

		refNames := refMap[c.Hash]
		if refNames == nil {
			refNames = []string{}
		}
		if parents == nil {
			parents = []string{}
		}

		commits = append(commits, GitCommit{
			Hash:         c.Hash.String(),
			AuthorName:   c.Author.Name,
			AuthorEmail:  c.Author.Email,
			Date:         c.Author.When,
			Subject:      subject,
			Body:         body,
			ParentHashes: parents,
			Refs:         refNames,
		})
		return nil
	})

	if err != nil && err != io.EOF {
		return nil, err
	}

	return commits, nil
}

func (a *App) Commit(repoPath string, subject string, body string, amend bool) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}

	msg := subject
	if body != "" {
		msg = subject + "\n\n" + body
	}

	opts := &git.CommitOptions{
		All:   false, // We only commit what is staged
		Amend: amend,
	}

	_, err = w.Commit(msg, opts)
	return err
}

func (a *App) generateSimpleDiff(old, new, path string) string {
	dmp := diffmatchpatch.New()

	// Character-based diff is default, but for unified diff we usually want line-based.
	// diffmatchpatch has a way to do line-based diffing.
	aDiff, bDiff, lineArray := dmp.DiffLinesToChars(old, new)
	diffs := dmp.DiffMain(aDiff, bDiff, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	if len(diffs) == 0 {
		return ""
	}

	var result strings.Builder

	oldLine := 1
	newLine := 1

	// We'll group diffs into hunks. For now, let's just do one big hunk or
	// simple hunks for each change. Standard unified diff uses hunks with context.
	// To keep it simple but include line numbers, we'll generate hunks.

	for i := 0; i < len(diffs); i++ {
		diff := diffs[i]
		if diff.Type == diffmatchpatch.DiffEqual {
			lines := strings.Split(strings.TrimSuffix(diff.Text, "\n"), "\n")
			count := len(lines)
			if diff.Text == "" {
				count = 0
			}
			oldLine += count
			newLine += count
		} else {
			// Start of a hunk
			hunkStartOld := oldLine
			hunkStartNew := newLine

			var hunkLines []string
			hunkOldCount := 0
			hunkNewCount := 0

			// Collect consecutive inserts/deletes
			for j := i; j < len(diffs); j++ {
				d := diffs[j]
				if d.Type == diffmatchpatch.DiffEqual {
					// Check if we should end the hunk or if it's a small gap
					// For simplicity, we end hunk on any Equality for now
					// Real diffs have context (usually 3 lines)
					break
				}

				lines := strings.Split(strings.TrimSuffix(d.Text, "\n"), "\n")
				for _, line := range lines {
					if d.Type == diffmatchpatch.DiffInsert {
						hunkLines = append(hunkLines, "+"+line)
						hunkNewCount++
					} else if d.Type == diffmatchpatch.DiffDelete {
						hunkLines = append(hunkLines, "-"+line)
						hunkOldCount++
					}
				}
				i = j
			}

			// Write hunk header
			result.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", hunkStartOld, hunkOldCount, hunkStartNew, hunkNewCount))
			for _, line := range hunkLines {
				result.WriteString(line + "\n")
			}

			oldLine += hunkOldCount
			newLine += hunkNewCount
		}
	}

	return result.String()
}
