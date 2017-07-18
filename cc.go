/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"

)

const   PRODUCER1 = "SUPPLIER1"
const 	PRODUCER2 = "SUPPLIER2"
const 	DEALER = "DEALER"
const 	SERVICE_CENTER = "SERVICE_CENTER"


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Part struct {
	PartId 			string 	`json:"partId"`
	ProductCode 		string  `json:"productCode"`
	Transactions		[]Transaction `json:"transactions"`
}

// PART TRANSACTION HISTORY
type Transaction struct {
	User  			string  `json:"user"`
	DateOfManufacture	string  `json:"dateOfManufacture"`
	DateOfDelivery		string	`json:"dateOfDelivery"`
	DateOfInstallation	string	`json:"dateOfInstallation"`
	VehicleId		string	`json:"vehicleId"`
	WarrantyStartDate	string	`json:"warrantyStartDate"`
	WarrantyEndDate		string	`json:"warrantyEndDate"`
	TType 			string   `json:"ttype"`
}

//==============================================================================================================================
//				Used as an index when querying all parts.
//==============================================================================================================================

type AllParts struct{
	Parts []string `json:"parts"`
}


// ============================================================================================================================
// Init
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub  shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	var parts AllParts
	jsonAsBytes, _ := json.Marshal(parts)
	err = stub.PutState("allParts", jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}



// ============================================================================================================================
// Run - Our entry point
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub  shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke is running " + function)

	// Handle different functions
	if function == "init" {					//initialize the chaincode state
		return t.Init(stub, "init", args)
	} else if function == "createPart" {			//create a part
		return t.createPart(stub, args)
	} else if function == "updatePart" {			//update a part
		return t.updatePart(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)	//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub  shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 { return nil, errors.New("Incorrect number of arguments passed") }

	if function != "getPart" && function != "getAllParts" {
		return nil, errors.New("Invalid query function name.")
	}

	if function == "getPart" { return t.getPart(stub, args[0]) }
	if function == "getAllParts" { return t.getAllParts(stub, args[0]) }

	return nil, nil
}


// ============================================================================================================================
// Get Part Details
// ============================================================================================================================
func (t *SimpleChaincode) getPart(stub  shim.ChaincodeStubInterface, partId string)([]byte, error){

	fmt.Println("Start find Part")
	fmt.Println("Looking for Part #" + partId);

	//get the part index
	bAsBytes, err := stub.GetState(partId)
	if err != nil {
		return nil, errors.New("Failed to get Part #" + partId)
	}

	return bAsBytes, nil

}

// ============================================================================================================================
// Get All Parts
// ============================================================================================================================
func (t *SimpleChaincode) getAllParts(stub  shim.ChaincodeStubInterface, user string)([]byte, error){

	fmt.Println("getAllParts:Looking for All Parts");

	//get the AllParts index
	allBAsBytes, err := stub.GetState("allParts")
	if err != nil {
		return nil, errors.New("Failed to get all Parts")
	}

	var res AllParts
	err = json.Unmarshal(allBAsBytes, &res)
	//fmt.Println(allBAsBytes);
	if err != nil {
		fmt.Println("Printing Unmarshal error:-");
		fmt.Println(err);
		return nil, errors.New("Failed to Unmarshal all Parts")
	}

	var rab AllParts

	for i := range res.Parts{

		sbAsBytes, err := stub.GetState(res.Parts[i])
		if err != nil {
			return nil, errors.New("Failed to get Part")
		}
		var sb Part
		json.Unmarshal(sbAsBytes, &sb)

		// currently we show all parts to the users
		//if(user == DEALER) {
			rab.Parts = append(rab.Parts,sb.PartId);
		//}
		//else{
		//	var _owner = sb.Owner
		//	if (user == _owner){
		//		rab.Parts = append(rab.Parts,sb.PartId);
		//		break;
		//	}
		//}
	}

	rabAsBytes, _ := json.Marshal(rab)

	return rabAsBytes, nil

}


func (t *SimpleChaincode) createPart(stub  shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running createPart")

	if len(args) != 4 {
		fmt.Println("Incorrect number of arguments. Expecting 4 - PartId, Product Code, Manufacture Date, User")
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("Arguments :"+args[0]+","+args[1]+","+args[2]+","+args[3]);
	// currently there is no such validation
	//if (args[2] != PRODUCER1)&&(args[2] != PRODUCER2)  {
	//	fmt.Println("You are not allowed to create a new Part")
	//	return nil, errors.New("You are not allowed to create a new Part")
	//}

	var bt Part
	bt.PartId 			= args[0]
	bt.ProductCode			= args[1]
	var tx Transaction
	tx.DateOfManufacture		= args[2]
	tx.TType 			= "CREATE"
	tx.User 			= args[3]
	bt.Transactions = append(bt.Transactions, tx)

	//Commit part to ledger
	fmt.Println("createPart Commit Part To Ledger");
	btAsBytes, _ := json.Marshal(bt)
	err = stub.PutState(bt.PartId, btAsBytes)
	if err != nil {
		return nil, err
	}

	//Update All Parts Array
	allBAsBytes, err := stub.GetState("allParts")
	if err != nil {
		return nil, errors.New("Failed to get all Parts")
	}
	var allb AllParts
	err = json.Unmarshal(allBAsBytes, &allb)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Parts")
	}
	allb.Parts = append(allb.Parts,bt.PartId)

	allBuAsBytes, _ := json.Marshal(allb)
	err = stub.PutState("allParts", allBuAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) updatePart(stub  shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running updatePart")

	if len(args) != 7 {
		fmt.Println("Incorrect number of arguments. Expecting 7 - PartId, Vehicle Id, Delivery Date, Installation Date, User, Warranty Start Date, Warranty End Date")
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
	}
	fmt.Println("Arguments :"+args[0]+","+args[1]+","+args[2]+","+args[3]+","+args[4]+","+args[5]+","+args[6]);

	//Get and Update Part data
	bAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Part #" + args[0])
	}
	var bch Part
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Part #" + args[0])
	}

	var tx Transaction
	if (strings.Contains(args[4], DEALER)) {
		tx.TType 	= "DELIVERY"
	} else if (strings.Contains(args[4], SERVICE_CENTER)) {
		tx.TType 	= "INSTALLED"
	}

	tx.VehicleId		= args[1]
	tx.DateOfDelivery	= args[2]
	tx.DateOfInstallation	= args[3]
	tx.User  		= args[4]
	tx.WarrantyStartDate	= args[5]
	tx.WarrantyEndDate	= args[6]


	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates part to ledger
	fmt.Println("updatePart Commit Updates To Ledger");
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.PartId, btAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil

}


// ============================================================================================================================
// Main function
// ============================================================================================================================

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
