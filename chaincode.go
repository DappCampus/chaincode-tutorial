/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// ERC20Chaincode is the definition of the chaincode structure.
type ERC20Chaincode struct {
}

// ERC20Metadata is the definition of Token Meta Info
type ERC20Metadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Owner       string `json:"owner"`
	TotalSupply uint64 `json:"totalSupply"`
}

// TransferEvent is the event definition of Transfer
type TransferEvent struct {
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Amount    int    `json:"amount"`
}

// Approval is the definition of Approval Event & Data format
type Approval struct {
	Owner     string `json:"owner"`
	Spender   string `json:"spender"`
	Allowance int    `json:"allowance"`
}

// Init is called when the chaincode is instantiated by the blockchain network.
// params - tokenName, symbol, owner(address), amount
func (cc *ERC20Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Init called with params: ", params)
	if len(params) != 4 {
		return shim.Error("incorrect number of parameter")
	}

	tokenName, symbol, owner, amount := params[0], params[1], params[2], params[3]

	// check amount is unsigned int
	amountUint, err := strconv.ParseUint(string(amount), 10, 64)
	if err != nil {
		return shim.Error("amount must be a number or amount cannot be negative")
	}

	// tokenName & symbol & owner cannot be empty
	if len(tokenName) == 0 || len(symbol) == 0 || len(owner) == 0 {
		return shim.Error("tokenName or symbol or owner cannot be emtpy")
	}

	// make metadata
	erc20 := &ERC20Metadata{Name: tokenName, Symbol: symbol, Owner: owner, TotalSupply: amountUint}
	erc20Bytes, err := json.Marshal(erc20)
	if err != nil {
		return shim.Error("failed to Marshal erc20, error: " + err.Error())
	}

	// save token meta data
	err = stub.PutState(tokenName, erc20Bytes)
	if err != nil {
		return shim.Error("failed to PutState, error: " + err.Error())
	}

	// save owner balance
	err = stub.PutState(owner, []byte(amount))
	if err != nil {
		return shim.Error("failed to PutState, error: " + err.Error())
	}

	// response
	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *ERC20Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()

	switch fcn {
	case "totalSupply":
		return cc.totalSupply(stub, params)
	case "balanceOf":
		return cc.balanceOf(stub, params)
	case "transfer":
		return cc.transfer(stub, params)
	case "allowance":
		return cc.allowance(stub, params)
	case "approve":
		return cc.approve(stub, params)
	case "approvalList":
		return cc.approvalList(stub, params)
	case "transferFrom":
		return cc.transferFrom(stub, params)
	case "increaseAllowance":
		return cc.increaseAllowance(stub, params)
	case "decreaseAllowance":
		return cc.decreaseAllowance(stub, params)
	case "mint":
		return cc.mint(stub, params)
	case "burn":
		return cc.burn(stub, params)
	default:
		return sc.Response{Status: 404, Message: "404 Not Found", Payload: nil}
	}
}

// totalSuuply is query function
// params - tokenName
// Returns the amount of token in existence
func (cc *ERC20Chaincode) totalSupply(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is one
	if len(params) != 1 {
		return shim.Error("incorrect number of parameter")
	}

	tokenName := params[0]

	// Get ERC20 Metadata
	erc20 := ERC20Metadata{}
	erc20Bytes, err := stub.GetState(tokenName)
	if err != nil {
		return shim.Error("failed to GetState, error: " + err.Error())
	}
	err = json.Unmarshal(erc20Bytes, &erc20)
	if err != nil {
		return shim.Error("failed to Unmarshal, error: " + err.Error())
	}

	// Convert TotalSupply to Bytes
	totalSupplyBytes, err := json.Marshal(erc20.TotalSupply)
	if err != nil {
		return shim.Error("failed to Marshal totalSupply, error: " + err.Error())
	}
	fmt.Println(tokenName + "'s totalSupply is " + string(totalSupplyBytes))

	return shim.Success(totalSupplyBytes)
}

// balanceOf is query function
// params - address
// Returns the amount of tokens owned by addresss
func (cc *ERC20Chaincode) balanceOf(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is one
	if len(params) != 1 {
		return shim.Error("incorrect number of parameters")
	}

	address := params[0]

	// get Balance
	amountBytes, err := stub.GetState(address)
	if err != nil {
		return shim.Error("failed to GetState, error: " + err.Error())
	}

	fmt.Println(address + "'s balance is " + string(amountBytes))

	if amountBytes == nil {
		return shim.Success([]byte("0"))
	}
	return shim.Success(amountBytes)
}

// transfer is invoke function that moves amount token
// from the caller's address to recipient
// params - caller's address, recipient's address, amount of token
func (cc *ERC20Chaincode) transfer(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 3
	if len(params) != 3 {
		return shim.Error("incorrect number of parameters")
	}

	callerAddress, recipientAddress, transferAmount := params[0], params[1], params[2]

	// check amount is integer & positive
	transferAmountInt, err := strconv.Atoi(transferAmount)
	if err != nil {
		return shim.Error("transfer amount must be integer")
	}
	if transferAmountInt <= 0 {
		return shim.Error("transfer amount must be positive")
	}

	// get caller amount
	callerAmount, err := stub.GetState(callerAddress)
	if err != nil {
		return shim.Error("failed to GetState, error: " + err.Error())
	}
	callerAmountInt, err := strconv.Atoi(string(callerAmount))
	if err != nil {
		return shim.Error("caller amount must be integer")
	}

	// get recipient amount
	recipientAmount, err := stub.GetState(recipientAddress)
	if err != nil {
		return shim.Error("failed to GetState, error: " + err.Error())
	}
	if recipientAmount == nil {
		recipientAmount = []byte("0")
	}
	recipientAmountInt, err := strconv.Atoi(string(recipientAmount))
	if err != nil {
		return shim.Error("caller amount must be integer")
	}

	// calculate amount
	callerResultAmount := callerAmountInt - transferAmountInt
	recipientResultAmount := recipientAmountInt + transferAmountInt

	// check callerReuslt Amount is positive
	if callerResultAmount < 0 {
		return shim.Error("caller's balance is not sufficient")
	}

	// save the caller's & recipient's amount
	err = stub.PutState(callerAddress, []byte(strconv.Itoa(callerResultAmount)))
	if err != nil {
		return shim.Error("failed to PutState of caller, error: " + err.Error())
	}
	err = stub.PutState(recipientAddress, []byte(strconv.Itoa(recipientResultAmount)))
	if err != nil {
		return shim.Error("failed to PutState of caller, error: " + err.Error())
	}

	// emit transfer event
	transferEvent := TransferEvent{Sender: callerAddress, Recipient: recipientAddress, Amount: transferAmountInt}
	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return shim.Error("failed to Marshal transferEvent, error: " + err.Error())
	}
	err = stub.SetEvent("transferEvent", transferEventBytes)
	if err != nil {
		return shim.Error("failed to SetEvent of TransferEvent, error: " + err.Error())
	}

	fmt.Println(callerAddress + " send " + transferAmount + " to " + recipientAddress)

	return shim.Success([]byte("transfer Success"))
}

// allowance is query function
// params - owner's address, spender's address
// Returns the remaining amount of token to invoke {transferFrom}
func (cc *ERC20Chaincode) allowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 2
	if len(params) != 2 {
		return shim.Error("incorrect number of parameters")
	}

	ownerAddress, spenderAddress := params[0], params[1]

	// create composite key
	approvalKey, err := stub.CreateCompositeKey("approval", []string{ownerAddress, spenderAddress})
	if err != nil {
		return shim.Error("failed to CreateCompositeKey for approval")
	}

	// get amount
	amountBytes, err := stub.GetState(approvalKey)
	if err != nil {
		return shim.Error("failed to GetState for amount")
	}
	if amountBytes == nil {
		amountBytes = []byte("0")
	}

	return shim.Success(amountBytes)

}

// approve is invoke function that Sets amount as the allowance
// of spender over the owner tokens
// params - owner's address, spender's address, amount of token
func (cc *ERC20Chaincode) approve(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 3
	if len(params) != 3 {
		return shim.Error("incorrect number of parameters")
	}

	ownerAddress, spenderAddress, allowanceAmount := params[0], params[1], params[2]

	// check amount is integer & positive
	allowanceAmountInt, err := strconv.Atoi(allowanceAmount)
	if err != nil {
		return shim.Error("allowance amount must be integer")
	}
	if allowanceAmountInt < 0 {
		return shim.Error("allowance amount must be positve")
	}

	// create composite key for allowance - approval/{owner}/{spender}
	approvalKey, err := stub.CreateCompositeKey("approval", []string{ownerAddress, spenderAddress})
	if err != nil {
		return shim.Error("failed to CreateCompositeKey for approval")
	}

	// save allowance amount
	err = stub.PutState(approvalKey, []byte(allowanceAmount))
	if err != nil {
		return shim.Error("failed to PutState for approval")
	}

	// emit approval event
	approvalEvent := Approval{Owner: ownerAddress, Spender: spenderAddress, Allowance: allowanceAmountInt}
	approvalBytes, err := json.Marshal(approvalEvent)
	if err != nil {
		return shim.Error("failed to SetEvent for ApprovalEvent")
	}
	err = stub.SetEvent("approvalEvent", approvalBytes)
	if err != nil {
		return shim.Error("failed to SetEvent for ApprovalEvent")
	}

	return shim.Success([]byte("approve success"))
}

// approvalList is query function
// params - owner's address
// Returns the approval list approved by owner
func (cc *ERC20Chaincode) approvalList(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of parmas is 1
	if len(params) != 1 {
		return shim.Error("incorrect number of params")
	}

	ownerAddress := params[0]

	// get all approval list (format is iterator)
	approvalIterator, err := stub.GetStateByPartialCompositeKey("approval", []string{ownerAddress})
	if err != nil {
		return shim.Error("failed to GetStateByPartialCompositeKey for approval iterationm error: " + err.Error())
	}

	// make slice for return value
	approvalSlice := []Approval{}

	// iterator
	defer approvalIterator.Close()
	if approvalIterator.HasNext() {
		for approvalIterator.HasNext() {
			approvalKV, _ := approvalIterator.Next()

			// get spender address
			_, addresses, err := stub.SplitCompositeKey(approvalKV.GetKey())
			if err != nil {
				return shim.Error("failed to SplitCompositeKey, error: " + err.Error())
			}
			spenderAddress := addresses[1]

			// get amount
			amountBytes := approvalKV.GetValue()
			amountInt, err := strconv.Atoi(string(amountBytes))
			if err != nil {
				return shim.Error("failed to get amount, error: " + err.Error())
			}

			// add approval result
			approval := Approval{Owner: ownerAddress, Spender: spenderAddress, Allowance: amountInt}
			approvalSlice = append(approvalSlice, approval)
		}
	}

	// convert approvalSlice to bytes for return
	response, err := json.Marshal(approvalSlice)
	if err != nil {
		return shim.Error("failed to Marshal approvalSlice, error: " + err.Error())
	}

	return shim.Success(response)
}

// transferFrom is invoke function that Moves amount of tokens from sender(owner) to recipient
// using allowance of spender
// parmas - owner's address, spender's address, recipient's address, amount of token
func (cc *ERC20Chaincode) transferFrom(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of parmas is 4
	if len(params) != 4 {
		return shim.Error("incorrect number of params")
	}

	ownerAddress, spenderAddress, recipientAddress, transferAmount := params[0], params[1], params[2], params[3]

	// check amount is integer & positive
	transferAmountInt, err := strconv.Atoi(transferAmount)
	if err != nil {
		return shim.Error("amount must be integer")
	}
	if transferAmountInt <= 0 {
		return shim.Error("amount must be positve")
	}

	// get allowance
	allowanceResponse := cc.allowance(stub, []string{ownerAddress, spenderAddress})
	if allowanceResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error: " + allowanceResponse.GetMessage())
	}

	// convert allowance response paylaod to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("allowance must be positive")
	}

	// transfer from owner to recipient
	transferResponse := cc.transfer(stub, []string{ownerAddress, recipientAddress, transferAmount})
	if transferResponse.GetStatus() >= 400 {
		return shim.Error("failed to transfer, error: " + transferResponse.GetMessage())
	}

	// decrease allowance amount
	approveAmountInt := allowanceInt - transferAmountInt
	approveAmount := strconv.Itoa(approveAmountInt)

	// approve amount of tokens transfered
	approveResponse := cc.approve(stub, []string{ownerAddress, spenderAddress, approveAmount})
	if approveResponse.GetStatus() >= 400 {
		return shim.Error("failed to approve, error: " + approveResponse.GetMessage())
	}

	return shim.Success([]byte("transferFrom success"))
}

func (cc *ERC20Chaincode) increaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) decreaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) mint(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) burn(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
