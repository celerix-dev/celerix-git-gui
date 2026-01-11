package backend

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// SshKeyInfo information about SSH keys for Git authentication.
type SshKeyInfo struct {
	PublicKey string `json:"public_key"`
	HasKey    bool   `json:"has_key"`
	Path      string `json:"path"`
}

// GetSshKeyInfo returns information about the user's SSH key.
func (a *App) GetSshKeyInfo() (SshKeyInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return SshKeyInfo{}, err
	}

	sshDir := filepath.Join(home, ".ssh")
	keyPath := filepath.Join(sshDir, "id_ed25519")
	pubKeyPath := keyPath + ".pub"

	info := SshKeyInfo{
		Path:   keyPath,
		HasKey: false,
	}

	if _, err := os.Stat(keyPath); err == nil {
		info.HasKey = true
		if pubData, err := os.ReadFile(pubKeyPath); err == nil {
			info.PublicKey = string(pubData)
		}
	} else {
		// Try RSA as a fallback for info
		rsaPath := filepath.Join(sshDir, "id_rsa")
		if _, err := os.Stat(rsaPath); err == nil {
			info.HasKey = true
			info.Path = rsaPath
			if pubData, err := os.ReadFile(rsaPath + ".pub"); err == nil {
				info.PublicKey = string(pubData)
			}
		}
	}

	return info, nil
}

// GenerateSshKey generates a new SSH key (Ed25519) using ssh-keygen.
func (a *App) GenerateSshKey() (SshKeyInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return SshKeyInfo{}, err
	}

	sshDir := filepath.Join(home, ".ssh")
	if _, err := os.Stat(sshDir); os.IsNotExist(err) {
		err = os.MkdirAll(sshDir, 0700)
		if err != nil {
			return SshKeyInfo{}, err
		}
	}

	keyPath := filepath.Join(sshDir, "id_ed25519")
	if _, err := os.Stat(keyPath); err == nil {
		return SshKeyInfo{}, fmt.Errorf("SSH key already exists at %s", keyPath)
	}

	// Use ssh-keygen to generate the key as it's more reliable than manually marshaling if dependencies are tricky
	cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-f", keyPath, "-N", "", "-C", "celerix-git-gui")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return SshKeyInfo{}, fmt.Errorf("ssh-keygen failed: %v, output: %s", err, string(output))
	}

	return a.GetSshKeyInfo()
}
