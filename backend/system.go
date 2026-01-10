package backend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// SelectDirectory opens a directory selection dialog.
func (a *App) SelectDirectory(title string) (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
	})
}

// GetHomeDir returns the user's home directory.
func (a *App) GetHomeDir() (string, error) {
	return os.UserHomeDir()
}

// OpenInFileManager opens the given path in the system's default file manager.
func (a *App) OpenInFileManager(path string) error {
	path = filepath.Clean(path)

	// Ensure the path is absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify the path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", absPath)
	}

	var cmd *exec.Cmd

	switch goruntime.GOOS {
	case "darwin":
		// On macOS, 'open' starts Finder. Use -- to prevent flag injection.
		cmd = exec.Command("open", "--", absPath)
	case "windows":
		// On Windows, 'explorer' opens the folder.
		// explorer.exe doesn't support -- but it's less prone to flag injection via path.
		cmd = exec.Command("explorer", absPath)
	case "linux":
		// On Linux, xdg-open uses the default file manager. Use -- to prevent flag injection.
		cmd = exec.Command("xdg-open", "--", absPath)
	default:
		return fmt.Errorf("unsupported platform: %s", goruntime.GOOS)
	}

	return cmd.Start()
}

func (a *App) OpenInBrowser(url string) {
	runtime.BrowserOpenURL(a.ctx, url)
}
