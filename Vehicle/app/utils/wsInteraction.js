/*eslint-env node */
// ==================================
// incoming messages, look for type
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
	else if(data.type == "getVehicle"){
		console.log("Get Part", data.vehicleId);
		chaincode.query.getVehicle([data.vehicleId], cb_got_vehicle);
	}
	else if(data.type == "getAllVehicles"){
		console.log("Get All Vehicles", owner);
		chaincode.query.getAllVehicles([""], cb_got_allvehicles);
	}
	else if(data.type == "createVehicle"){
		console.log("Create Vehicle ", data, owner);
		if(data.vehicle){			
			chaincode.invoke.createVehicle([data.vehicle.make, data.vehicle.chassisNumber, data.vehicle.vin, owner], cb_invoked_createVehicle);				//create a new Vehicle
		}
	}
	else if(data.type == "updateVehicle"){
		console.log("Update Vehicle ", data, owner);
		if(data.vehicle){			
			chaincode.invoke.updateVehicle([data.vehicle.vehicleId, data.vehicle.ttype, data.vehicle.owner.name, data.vehicle.owner.phoneNumber, data.vehicle.owner.email, data.vehicle.dealer.name, data.vehicle.dealer.phoneNumber, data.vehicle.dealer.email, data.vehicle.licensePlateNumber, data.vehicle.dateofDelivery, data.vehicle.warrantyStartDate, data.vehicle.warrantyEndDate, owner, data.vehicle.parts], cb_invoked_updateVehicle);				//update vehicle
		}
	}
	else if(data.type == "createPart"){
		console.log("Create Part ", data, owner);
		if(data.part){
			console.log('Part manufacture date:'+data.part.dateOfManufacture);
			chaincode.invoke.createPart([data.part.partId, data.part.productCode, data.part.dateOfManufacture, owner], cb_invoked_createpart);				//create a new paper
		}
	}
	else if(data.type == "updatePart"){
		console.log("Update Part ", data, owner);
		if(data.part){
			chaincode.invoke.updatePart([data.part.partId, data.part.vehicleId, data.part.dateOfDelivery, data.part.dateOfInstallation, owner, data.part.warrantyStartDate, data.part.warrantyEndDate, data.part.tranType], cb_invoked_updatepart);	//update part details
		}		
	}
	else if(data.type == "getPart"){
		console.log("Get Part", data.partId);
		chaincode.query.getPart([data.partId], cb_got_part);
	}
	else if(data.type == "getAllParts"){
		console.log("Get All Parts", owner);
		chaincode.query.getAllParts([""], cb_got_allparts);
	}
	else if(data.type == "getAllPartsForUpdateVehicle"){
		console.log("Get All Parts", owner);
		chaincode.query.getAllParts([""], cb_got_allpartsForUpdateVehicle);
	}
	
	function cb_got_part(e, part){
		if(e != null){
			console.log("Get Part error", e);
		}
		else{
			sendMsg({msg: "part", part: JSON.parse(part)});
		}
	}
	
	function cb_got_allparts(e, allParts){
		if(e != null){
			console.log("Get All Parts error", e);
		}
		else{
			sendMsg({msg: "allParts", parts: JSON.parse(allParts).parts});
		}
	}

	function cb_got_allpartsForUpdateVehicle(e, allParts){
		if(e != null){
			console.log("Get All Parts error", e);
		}
		else{
			sendMsg({msg: "allPartsForUpdateVehicle", parts: JSON.parse(allParts).parts});
		}
	}
	
	function cb_got_vehicle(e, vehicle){
		if(e != null){
			console.log("Get Vehicle error", e);
		}
		else{
			sendMsg({msg: "vehicle", vehicle: JSON.parse(vehicle)});
		}
	}

	function cb_got_allvehicles(e, allVehicles){
		if(e != null){
			console.log("Get All Vehicles error", e);
		}
		else{
			sendMsg({msg: "allVehicles", vehicles: JSON.parse(allVehicles).vehicles});
		}
	}

	function cb_invoked_createVehicle(e, a){
		console.log("response: ", e, a);
		if(e != null){
			console.log("Invoked create vehicle error", e);
		}
		else{
			console.log("Vehicle ID #" + data.vehicle.chassisNumber)
			sendMsg({msg: "vehicleCreated", chassisNumber: data.vehicle.chassisNumber});
		}
	}

	function cb_invoked_updateVehicle(e, a){
		console.log("response: ", e, a);
		if(e != null){
			console.log("Invoked update vehicle error", e);
		}
		else{
			console.log("Vehicle ID #" + data.vehicle.chassisNumber)
			sendMsg({msg: "vehicleUpdated", chassisNumber: data.vehicle.chassisNumber});
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
	function cb_invoked_updatepart(e, a){
		console.log("response: ", e, a);
		if(e != null){
			console.log("Invoked update part error", e);
		}
		else{
			console.log("part ID #" + data.part.id)
			sendMsg({msg: "partUpdated", partId: data.part.id});
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
