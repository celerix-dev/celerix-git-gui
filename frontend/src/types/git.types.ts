export interface GitRemote {
    name: string;
    url: string;
}

export interface GitStatusFile {
    path: string;
    status: string;
    is_staged: boolean;
}

export interface SshKeyInfo {
    public_key: string;
    has_key: boolean;
    path: string;
}

export interface RepoTab {
    id: string;
    name: string;
    path: string;
    remotes: GitRemote[];
    statusFiles: GitStatusFile[];
    activeVerticalTab: 'info' | 'local-changes' | 'commit' | 'placeholder1' | 'placeholder2';
    loading: boolean;
    error: string | null;
}
