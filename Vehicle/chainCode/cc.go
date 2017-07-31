package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"		
	"crypto/rand"		
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//VehicleId, Description, RegistrationNumber, Make, VIN, DateofRegistration, ChassisNumber, Color, OwnerName, OwnerPhoneNumber, OwnerEmail
type Vehicle struct {
	VehicleId 			string 	`json:"vehicleId"`
	Make 		string  `json:"make"`
	ChassisNumber 		string  `json:"chassisNumber"`
	Vin 		string  `json:"vin"`
	DateOfManufacture 		string  `json:"dateOfManufacture"`	
	Owner Owner `json:"owner"`
	Dealer Dealer `json:"dealer"`
	LicensePlateNumber 		string  `json:"licensePlateNumber"`
	WarrantyStartDate 		string  `json:"warrantyStartDate"`	
	WarrantyEndDate 		string  `json:"warrantyEndDate"`	
	DateofDelivery 		string  `json:"dateofDelivery"`
	ServiceRequestRaised	string 	`json:"serviceRequestRaised"`	
	Parts		[]Part `json:"parts"`
	VehicleTransactions		[]VehicleTransaction `json:"vehicleTransactions"`
}

type VehicleTransaction struct {	
	WarrantyStartDate 		string  `json:"warrantyStartDate"`	
	WarrantyEndDate 		string  `json:"warrantyEndDate"`	
	TType 			string   `json:"ttype"`
	TValue 			string   `json:"tvalue"`
	UpdatedBy  			string  `json:"updatedBy"`
	UpdatedOn  			string  `json:"updatedOn"`
}

type Owner struct {
	Name 		string  `json:"name"`
	PhoneNumber 		string  `json:"phoneNumber"`
	Email 		string  `json:"email"`
}

type Dealer struct {
	Name 		string  `json:"name"`
	PhoneNumber 		string  `json:"phoneNumber"`
	Email 		string  `json:"email"`
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
//				Used as an index when querying all vehicles and parts.
//==============================================================================================================================

type AllVehicles struct{
	Vehicles []string `json:"vehicles"`
}


type AllParts struct{
	Parts []string `json:"parts"`
}


// ============================================================================================================================
// Init --- 
// Local - *shim.ChaincodeStub
// Server - shim.ChaincodeStubInterface
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub  shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error

	var vehicles AllVehicles
	var parts AllParts

	jsonAsBytesVehicles, _ := json.Marshal(vehicles)
	err = stub.PutState("allVehicles", jsonAsBytesVehicles)
	if err != nil {
		return nil, err
	}

	jsonAsBytesParts, _ := json.Marshal(parts)
	err = stub.PutState("allParts", jsonAsBytesParts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub  shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
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
	} else if function == "createVehicle" {			//create a vehicle
		return t.createVehicle(stub, args)	
	} else if function == "updateVehicle" {			//create a vehicle
		return t.updateVehicle(stub, args)
	} else if function == "addPart" {			//create a part
		return t.createPart(stub, args)	
	} else if function == "updatePart" {			//create a part
		return t.updatePart(stub, args)			
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

	if function == "getVehicle" { return t.getVehicle(stub, args[0]) }
	if function == "getPart" { return t.getPart(stub, args[0]) }
	if function == "getAllVehicles" { return t.getAllVehicles(stub, args[0]) }
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
		rab.Parts = append(rab.Parts,sb.PartId);
	}

	rabAsBytes, _ := json.Marshal(rab)

	return rabAsBytes, nil

}

// ============================================================================================================================
// Get Vehicle Details
// ============================================================================================================================
func (t *SimpleChaincode) getVehicle(stub  shim.ChaincodeStubInterface, vehicleId string)([]byte, error){

	fmt.Println("Start find Vehicle")
	fmt.Println("Looking for Vehicle #" + vehicleId);

	//get the vehicle index
	bAsBytes, err := stub.GetState(vehicleId)
	if err != nil {
		return nil, errors.New("Failed to get Vehicle Id #" + vehicleId)
	}

	return bAsBytes, nil

}


// ============================================================================================================================
// Get All Vehicles
// ============================================================================================================================
func (t *SimpleChaincode) getAllVehicles(stub  shim.ChaincodeStubInterface, user string)([]byte, error){

	fmt.Println("getAllVehicles:Looking for All Vehicles");

	//get the AllVehicles index
	allBAsBytes, err := stub.GetState("allVehicles")
	if err != nil {
		return nil, errors.New("Failed to get all Vehicles")
	}

	var res AllVehicles
	err = json.Unmarshal(allBAsBytes, &res)
	//fmt.Println(allBAsBytes);
	if err != nil {
		fmt.Println("Printing Unmarshal error:-");
		fmt.Println(err);
		return nil, errors.New("Failed to Unmarshal all Vehicles")
	}

	var rab AllVehicles

	for i := range res.Vehicles{

		sbAsBytes, err := stub.GetState(res.Vehicles[i])
		if err != nil {
			return nil, errors.New("Failed to get Vehicle")
		}
		var sb Vehicle
		json.Unmarshal(sbAsBytes, &sb)
		
		if user != "" {
			// return only customer vehicles
			if sb.Owner.Name == user {
				rab.Vehicles = append(rab.Vehicles, sb.VehicleId +"-"+ sb.ChassisNumber);
			}
		} else if user == "" {
			// return all vehicles for mfr, dealer, service center user
			rab.Vehicles = append(rab.Vehicles, sb.VehicleId +"-"+ sb.ChassisNumber);
		}
	}

	rabAsBytes, _ := json.Marshal(rab)

	return rabAsBytes, nil

}

// creating new vehicle in blockchain
func (t *SimpleChaincode) createVehicle(stub  shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running createVehicle")

	if len(args) != 4 {
		fmt.Println("Incorrect number of arguments. Expecting 4 - Make, ChassisNumber, Vin, User")
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("Arguments :"+args[0]+","+args[1]+","+args[2]+","+args[3]);

	var bt Vehicle
	bt.VehicleId = NewUniqueId()
	bt.Make			= args[0]
	bt.ChassisNumber = args[1]
	bt.Vin = args[2]
	bt.DateOfManufacture = time.Now().Local().String()
	
	var own Owner
	own.Name = ""
	own.PhoneNumber = ""
	own.Email = ""
	var del Dealer
	del.Name = ""
	del.PhoneNumber = ""
	del.Email = ""
	bt.Owner = own
	bt.Dealer = del
	
	var tx VehicleTransaction 	
	tx.TType 			= "CREATE"
	tx.UpdatedBy 			= args[3]
	tx.UpdatedOn   			= time.Now().Local().String()
	bt.VehicleTransactions = append(bt.VehicleTransactions, tx)

	//Commit vehicle to ledger
	fmt.Println("createVehicle Commit Vehicle To Ledger");
	btAsBytes, _ := json.Marshal(bt)
	err = stub.PutState(bt.VehicleId, btAsBytes)
	if err != nil {
		return nil, err
	}

	//Update All Vehicles Array
	allBAsBytes, err := stub.GetState("allVehicles")
	if err != nil {
		return nil, errors.New("Failed to get all Vehicles")
	}
	var allb AllVehicles
	err = json.Unmarshal(allBAsBytes, &allb)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Vehicles")
	}
	allb.Vehicles = append(allb.Vehicles,bt.VehicleId)

	allBuAsBytes, _ := json.Marshal(allb)
	err = stub.PutState("allVehicles", allBuAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Updating existing vehicle in blockchain
func (t *SimpleChaincode) updateVehicle(stub  shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running updateVehicle")

	fmt.Println("Arguments :"+args[0]+","+args[1]+","+args[2]+","+args[3]+","+args[4]+","+args[5]+","+args[6]+","+args[7]);

	//Get and Update Part data
	bAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Vehicle #" + args[0])
	}
	var bch Vehicle
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Vehicle #" + args[0])
	}	
	
	var updateStr string
	if bch.Owner.Name 	!= args[2] {
		bch.Owner.Name 	= args[2]
		updateStr += ",Owner Name to "+ args[2]
	}

	if bch.Owner.PhoneNumber != args[3] {
		bch.Owner.PhoneNumber 	= args[3]
		updateStr += ",Owner Phone to "+ args[3]
	}

	if bch.Owner.Email != args[4] {
		bch.Owner.Email 	= args[4]
		updateStr += ",Owner Email to "+ args[4]
	}
	
	bch.Dealer.Name 	= args[5]
	bch.Dealer.PhoneNumber 	= args[6]
	bch.Dealer.Email 	= args[7]
	
	if bch.LicensePlateNumber != args[8] {
		bch.LicensePlateNumber=  args[8]
		updateStr += ",License Plate Number"+ args[8]
	}

	if bch.DateofDelivery != args[9] {
		bch.DateofDelivery =  args[9]
		updateStr += ",Date of Delivery"+ args[9]
	}

	if bch.WarrantyStartDate != args[10] {
		bch.WarrantyStartDate =  args[10]
		updateStr += ",Warranty Start Date"+ args[10]
	}

	if bch.WarrantyEndDate != args[11] {
		bch.WarrantyEndDate =  args[11]
		updateStr += ",Warranty End Date"+ args[11]
	}	
	
	var tx VehicleTransaction 
	
	tx.WarrantyStartDate	= args[10]
	tx.WarrantyEndDate	= args[11]		
	tx.UpdatedBy   	= args[12]
	tx.UpdatedOn   	= time.Now().Local().String()
	
	tx.TType 	= args[1] +" |"+ updateStr
	bch.VehicleTransactions = append(bch.VehicleTransactions, tx)

	//Commit updates part to ledger
	fmt.Println("updateVehicle Commit Updates To Ledger");
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.VehicleId, btAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// creating new part in blockchain
func (t *SimpleChaincode) createPart(stub  shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running createPart")

	if len(args) != 4 {
		fmt.Println("Incorrect number of arguments. Expecting 4 - PartId, Product Code, Manufacture Date, User")
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	fmt.Println("Arguments :"+args[0]+","+args[1]+","+args[2]+","+args[3]);

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

// Updating existing part in blockchain
func (t *SimpleChaincode) updatePart(stub  shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running updatePart")

	if len(args) != 8 {
		fmt.Println("Incorrect number of arguments. Expecting 8 - PartId, Vehicle Id, Delivery Date, Installation Date, User, Warranty Start Date, Warranty End Date, Type")
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}
	fmt.Println("Arguments :"+args[0]+","+args[1]+","+args[2]+","+args[3]+","+args[4]+","+args[5]+","+args[6]+","+args[7]);

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
	tx.TType 	= args[7];

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

func NewUniqueId() string{
	n := 10
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }
	s := ""
    s = fmt.Sprintf("%X", b)
	return s    
}

// ============================================================================================================================
// Main function
// ============================================================================================================================

func main() {
	
    fmt.Println(time.Now().Local().String())
	
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
