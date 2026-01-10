export namespace backend {
	
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

