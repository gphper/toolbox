export namespace main {
	
	export class TotpRes {
	    code: string;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new TotpRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.url = source["url"];
	    }
	}

}

