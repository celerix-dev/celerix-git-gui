export interface GitBranch {
  name: string;
  is_current: boolean;
}

export interface GitStash {
  index: number;
  message: string;
  branch: string;
}

export interface GitCommitFile {
  path: string;
  status: string;
}

export interface GitCommit {
  hash: string;
  author: string;
  author_email: string;
  message: string;
  body: string;
  date: string;
  parents: string[];
  branches: string[];
  tags: string[];
}

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
  activeVerticalTab: 'info' | 'commit' | 'placeholder1' | 'placeholder2';
  loading: boolean;
  error: string | null;
}
