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
  branches: GitBranch[];
  remoteBranches: string[];
  tags: string[];
  stashes: GitStash[];
  remotes: GitRemote[];
  commits: GitCommit[];
  graphNodes: any[]; 
  statusFiles: GitStatusFile[];
  selectedCommitHash: string | null;
  selectedFilePath: string | null;
  selectedFileStaged: boolean | null;
  currentDiff: string | null;
  sidebarSelection: 'local-changes' | 'all-commits';
  selectedFiles: string[];
  commitSubject: string;
  commitDescription: string;
  commitAmend: boolean;
  loading: boolean;
  operationLoading: boolean;
  operationName: string | null;
  diffLoading: boolean;
  error: string | null;
}
