/*eslint-env node */
// ==================================
// Part 2 - incoming messages, look for type
// ==================================
var ibc = {};
var chaincode = {};
var async = require("async");

module.exports.setup = function(sdk, cc){
	ibc = sdk;
	chaincode = cc;
};

module.exports.process_msg = function(ws, data, owner){
	
	if(data.type == "chainstats"){
		console.log("Chainstats msg");
		ibc.chain_stats(cb_chainstats);
	}
	else if(data.type == "createPart"){
		console.log("Create Part ", data, owner);
		if(data.part){
			chaincode.invoke.createPart([data.part.partId, data.part.productCode, data.part.dateOfManufacture], cb_invoked_createpart);				//create a new paper
		}
	}
	else if(data.type == "getPart"){
		console.log("Get Part", data.partId);
		chaincode.query.getPart([data.partId], cb_got_part);
	}
	/*else if(data.type == "getAllBatches"){
		console.log("Get All Batches", owner);
		chaincode.query.getAllBatches([owner], cb_got_allbatches);
	}*/
	
	function cb_got_part(e, part){
		if(e != null){
			console.log("Get Part error", e);
		}
		else{
			sendMsg({msg: "part", part: JSON.parse(part)});
		}
	}
	
	function cb_got_allbatches(e, allBatches){
		if(e != null){
			console.log("Get All Batches error", e);
		}
		else{
			sendMsg({msg: "allBatches", batches: JSON.parse(allBatches).batches});
		}
	}
	
	function cb_invoked_createpart(e, a){
		console.log("response: ", e, a);
		if(e != null){
			console.log("Invoked create part error", e);
		}
		else{
			console.log("part ID #" + data.part.id)
			sendMsg({msg: "partCreated", partId: data.part.id});
		}
		

	}
	
	//call back for getting the blockchain stats, lets get the block height now
	var chain_stats = {};
	function cb_chainstats(e, stats){
		chain_stats = stats;
		if(stats && stats.height){
			var list = [];
			for(var i = stats.height - 1; i >= 1; i--){										//create a list of heights we need
				list.push(i);
				if(list.length >= 8) break;
			}
			list.reverse();																//flip it so order is correct in UI
			console.log(list);
			async.eachLimit(list, 1, function(key, cb) {								//iter through each one, and send it
				ibc.block_stats(key, function(e, stats){
					if(e == null){
						stats.height = key;
						sendMsg({msg: "chainstats", e: e, chainstats: chain_stats, blockstats: stats});
					}
					cb(null);
				});
			}, function() {
			});
		}
	}

	//call back for getting a block's stats, lets send the chain/block stats
	function cb_blockstats(e, stats){
		if(chain_stats.height) stats.height = chain_stats.height - 1;
		sendMsg({msg: "chainstats", e: e, chainstats: chain_stats, blockstats: stats});
	}
	

	//send a message, socket might be closed...
	function sendMsg(json){
		if(ws){
			try{
				ws.send(JSON.stringify(json));
			}
			catch(e){
				console.log("error ws", e);
			}
		}
	}
};
