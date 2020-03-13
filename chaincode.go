/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	"github.com/erc20/controller"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// ERC20Chaincode is the definition of the chaincode structure.
type ERC20Chaincode struct {
	controller *controller.Controller
}

// NewChaincode is constructor function for ERC20Chaincode
func NewChaincode() *ERC20Chaincode {
	controller := controller.NewController()
	return &ERC20Chaincode{controller}
}

// Init is called when the chaincode is instantiated by the blockchain network.
// params - tokenName, symbol, owner(address), amount
func (cc *ERC20Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	_, params := stub.GetFunctionAndParameters()
	fmt.Println("Init called with params: ", params)

	return cc.controller.Init(stub, params)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *ERC20Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()

	switch fcn {
	case "totalSupply":
		return cc.controller.TotalSupply(stub, params)
	case "balanceOf":
		return cc.controller.BalanceOf(stub, params)
	case "transfer":
		return cc.controller.Transfer(stub, params)
	case "allowance":
		return cc.controller.Allowance(stub, params)
	case "approve":
		return cc.controller.Approve(stub, params)
	case "approvalList":
		return cc.controller.ApprovalList(stub, params)
	case "transferFrom":
		return cc.controller.TransferFrom(stub, params)
	case "transferOtherToken":
		return cc.controller.TransferOtherToken(stub, params)
	case "increaseAllowance":
		return cc.controller.IncreaseAllowance(stub, params)
	case "decreaseAllowance":
		return cc.controller.DecreaseAllowance(stub, params)
	case "mint":
		return cc.controller.Mint(stub, params)
	case "burn":
		return cc.controller.Burn(stub, params)
	case "transactionAPI":
		return cc.transactionAPI(stub, params)
	default:
		return sc.Response{Status: 404, Message: "404 Not Found", Payload: nil}
	}
}

// <Transaction API>
//   - GetTxID
//   - GetTxTimestamp()
//   - GetCreator()
//   - GetSignedProposal()
// <State data API>
//   - GetStateByRange()
//   - GetStateByRangeWithPagination()()
//   - GetStateByPartialCompositeKeyWithPagination()
// <Key API>
//   - GetHistoryForKey()

func (cc *ERC20Chaincode) transactionAPI(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// GetTxID
	fmt.Println("==================== TX ID ====================")
	txID := stub.GetTxID()
	fmt.Println(txID)
	fmt.Println()

	// GetTxTimestamp
	fmt.Println("==================== TX Timestamp ====================")
	txTimeStamp, _ := stub.GetTxTimestamp()
	fmt.Println(txTimeStamp.String())
	fmt.Println()

	// GetCreator
	fmt.Println("==================== Creator ====================")
	creator, _ := stub.GetCreator()
	fmt.Println(string(creator))
	fmt.Println()

	// GetSignedProposal
	fmt.Println("==================== Signed Proposal ====================")
	signedProposal, _ := stub.GetSignedProposal()
	fmt.Println(signedProposal.String())
	fmt.Println()

	return shim.Success(nil)
}
