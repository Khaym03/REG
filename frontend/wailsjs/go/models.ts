export namespace auth {
	
	export class User {
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}

}

export namespace config {
	
	export class BrowserConfig {
	    LoggerOut: any;
	    headless?: boolean;
	    trace?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new BrowserConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.LoggerOut = source["LoggerOut"];
	        this.headless = source["headless"];
	        this.trace = source["trace"];
	    }
	}

}

export namespace domain {
	
	export class DateRange {
	    // Go type: time
	    from: any;
	    // Go type: time
	    to: any;
	
	    static createFrom(source: any = {}) {
	        return new DateRange(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = this.convertValues(source["from"], null);
	        this.to = this.convertValues(source["to"], null);
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

}

export namespace event {
	
	export class Topics {
	    stats_result: string;
	    workflow_started_topic: string;
	    loggin_topic: string;
	    guides_gather_topic: string;
	    inventory_sync_topic: string;
	    reception_topic: string;
	    logout_topic: string;
	    workflow_finished_topic: string;
	
	    static createFrom(source: any = {}) {
	        return new Topics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stats_result = source["stats_result"];
	        this.workflow_started_topic = source["workflow_started_topic"];
	        this.loggin_topic = source["loggin_topic"];
	        this.guides_gather_topic = source["guides_gather_topic"];
	        this.inventory_sync_topic = source["inventory_sync_topic"];
	        this.reception_topic = source["reception_topic"];
	        this.logout_topic = source["logout_topic"];
	        this.workflow_finished_topic = source["workflow_finished_topic"];
	    }
	}

}

export namespace stats {
	
	export class Stats {
	    outstanding_debt: number;
	    intransit_guides: number;
	    expired_guides: number;
	    pending_procedures: number;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.outstanding_debt = source["outstanding_debt"];
	        this.intransit_guides = source["intransit_guides"];
	        this.expired_guides = source["expired_guides"];
	        this.pending_procedures = source["pending_procedures"];
	    }
	}

}

export namespace workflow {
	
	export class WorkFlowInput {
	    user: auth.User;
	    date: domain.DateRange;
	    receive_guides_in_transit?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new WorkFlowInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user = this.convertValues(source["user"], auth.User);
	        this.date = this.convertValues(source["date"], domain.DateRange);
	        this.receive_guides_in_transit = source["receive_guides_in_transit"];
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

}

