export namespace backend {
	
	export class CommitFileChange {
	    path: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new CommitFileChange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.status = source["status"];
	    }
	}
	export class GitCommit {
	    hash: string;
	    authorName: string;
	    authorEmail: string;
	    // Go type: time
	    date: any;
	    subject: string;
	    body: string;
	    parentHashes: string[];
	    refs: string[];
	
	    static createFrom(source: any = {}) {
	        return new GitCommit(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hash = source["hash"];
	        this.authorName = source["authorName"];
	        this.authorEmail = source["authorEmail"];
	        this.date = this.convertValues(source["date"], null);
	        this.subject = source["subject"];
	        this.body = source["body"];
	        this.parentHashes = source["parentHashes"];
	        this.refs = source["refs"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GitRemoteBranches {
	    name: string;
	    branches: string[];
	
	    static createFrom(source: any = {}) {
	        return new GitRemoteBranches(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.branches = source["branches"];
	    }
	}
	export class GitStatusFile {
	    path: string;
	    status: string;
	    is_staged: boolean;
	
	    static createFrom(source: any = {}) {
	        return new GitStatusFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.status = source["status"];
	        this.is_staged = source["is_staged"];
	    }
	}
	export class RepoStats {
	    repoName: string;
	    remoteUrl: string;
	    sizeMb: number;
	    commitCount: number;
	    // Go type: time
	    lastCommit: any;
	    // Go type: time
	    firstCommit: any;
	    isClean: boolean;
	    modifiedFiles: string[];
	    branches: string[];
	    remotes: GitRemoteBranches[];
	    tags: string[];
	    stashes: string[];
	    currentBranch: string;
	
	    static createFrom(source: any = {}) {
	        return new RepoStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.repoName = source["repoName"];
	        this.remoteUrl = source["remoteUrl"];
	        this.sizeMb = source["sizeMb"];
	        this.commitCount = source["commitCount"];
	        this.lastCommit = this.convertValues(source["lastCommit"], null);
	        this.firstCommit = this.convertValues(source["firstCommit"], null);
	        this.isClean = source["isClean"];
	        this.modifiedFiles = source["modifiedFiles"];
	        this.branches = source["branches"];
	        this.remotes = this.convertValues(source["remotes"], GitRemoteBranches);
	        this.tags = source["tags"];
	        this.stashes = source["stashes"];
	        this.currentBranch = source["currentBranch"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SshKeyInfo {
	    public_key: string;
	    has_key: boolean;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new SshKeyInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.public_key = source["public_key"];
	        this.has_key = source["has_key"];
	        this.path = source["path"];
	    }
	}

}

