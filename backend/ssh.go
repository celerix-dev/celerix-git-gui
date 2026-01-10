package backend

// SshKeyInfo information about SSH keys for Git authentication.
type SshKeyInfo struct {
	PublicKey string `json:"public_key"`
	HasKey    bool   `json:"has_key"`
	Path      string `json:"path"`
}

// GetSshKeyInfo returns information about the user's SSH key.
func (a *App) GetSshKeyInfo() (SshKeyInfo, error) {
	// Dummy implementation for now
	return SshKeyInfo{HasKey: false}, nil
}

// GenerateSshKey generates a new SSH key.
func (a *App) GenerateSshKey() (SshKeyInfo, error) {
	// Dummy implementation for now
	return SshKeyInfo{HasKey: false}, nil
}
