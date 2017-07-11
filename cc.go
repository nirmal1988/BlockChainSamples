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

	"github.com/hyperledger/fabric/core/chaincode/shim"

)

const   PRODUCER1 = "SUPPLIER1"
const 	PRODUCER2 = "SUPPLIER2"
const 	DEALER1 = "DEALER1"
const 	DEALER2 = "DEALER2"
const 	SERVICECENTER = "SERVICECENTER"
const   SHIPPING = "SHIPPINGCO"
const   RETAILER = "RETAILER"
const 	CONSUMER = "CONSUMER"
const 	CERTIFIER = "CERTIFIER"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Part struct {
	PartId 				string 	`json:"partId"`
	ProductCode 		string  `json:"productCode"`
	DateOfManufacture	string  `json:"dateOfManufacture"`
	DateOfDelivery		string	`json:"dateOfDelivery"`
	DateOfInstallation	string	`json:"dateOfInstallation"`
	VehicleId			string	`json:"vehicleId"`
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
	err = stub.PutState("parts", jsonAsBytes)
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
	if function == "init" {													//initialize the chaincode state
		return t.Init(stub, "init", args)
	} else if function == "createPart" {											//create a batch
		return t.createPart(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub  shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 1 { return nil, errors.New("Incorrect number of arguments passed") }

	if function != "getPart" {
		return nil, errors.New("Invalid query function name.")
	}

	if function == "getPart" { return t.getPart(stub, args[0]) }

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


func (t *SimpleChaincode) createPart(stub  shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running createPart")

	if len(args) != 3 {
		fmt.Println("Incorrect number of arguments. Expecting 3")
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	
	//if (args[2] != PRODUCER1)&&(args[2] != PRODUCER2)  {
	//	fmt.Println("You are not allowed to create a new Part")
	//	return nil, errors.New("You are not allowed to create a new Part")
	//}

	//////// TODO
	var bt Part
	bt.PartId 			= args[0]
	bt.ProductCode			= args[1]
	bt.DateOfManufacture		= args[2]

	//Commit part to ledger
	fmt.Println("createPart Commit Part To Ledger");
	btAsBytes, _ := json.Marshal(bt)
	err = stub.PutState(bt.PartId, btAsBytes)
	if err != nil {
		return nil, err
	}

	//Update All Batches Array
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


// ============================================================================================================================
// Main function
// ============================================================================================================================

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
