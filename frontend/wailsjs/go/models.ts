export namespace struct { ExtVersion int; ExtTargetduration int; ExtMediaSequence int; ExtPlaylistType string; ExtList string; ExtKey string; ExtKeyMethod string; ExtKeyUri string; ExtKeyIv string; ExtKeyTrue string; ExtKeyIvTrue string } {
	
	export class  {
	    ExtVersion: number;
	    ExtTargetduration: number;
	    ExtMediaSequence: number;
	    ExtPlaylistType: string;
	    ExtList: string[];
	    ExtKey: string;
	    ExtKeyMethod: string;
	    ExtKeyUri: string;
	    ExtKeyIv: string;
	    ExtKeyTrue: string;
	    ExtKeyIvTrue: string;
	
	    static createFrom(source: any = {}) {
	        return new (source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ExtVersion = source["ExtVersion"];
	        this.ExtTargetduration = source["ExtTargetduration"];
	        this.ExtMediaSequence = source["ExtMediaSequence"];
	        this.ExtPlaylistType = source["ExtPlaylistType"];
	        this.ExtList = source["ExtList"];
	        this.ExtKey = source["ExtKey"];
	        this.ExtKeyMethod = source["ExtKeyMethod"];
	        this.ExtKeyUri = source["ExtKeyUri"];
	        this.ExtKeyIv = source["ExtKeyIv"];
	        this.ExtKeyTrue = source["ExtKeyTrue"];
	        this.ExtKeyIvTrue = source["ExtKeyIvTrue"];
	    }
	}

}

