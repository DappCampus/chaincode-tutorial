/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Chaincode is the definition of the chaincode structure.
type ERC20Chaincode struct {
}

// Init is called when the chaincode is instantiated by the blockchain network.
func (cc *ERC20Chaincode) Init(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()
	fmt.Println("Init()", fcn, params)

	return shim.Success(nil)
}

// Invoke is called as a result of an application request to run the chaincode.
func (cc *ERC20Chaincode) Invoke(stub shim.ChaincodeStubInterface) sc.Response {
	fcn, params := stub.GetFunctionAndParameters()

	// GetArgs
	args := stub.GetArgs()
	fmt.Println("GetArgs(): ")
	for _, arg := range args {
		argStr := string(arg)
		fmt.Printf("%s ", argStr)
	}
	fmt.Println()

	// GetStringArgs()
	stringArgs := stub.GetStringArgs()
	fmt.Println("GetStringArgs(): ", stringArgs)

	// GetArgsSlice()
	argsSlice, _ := stub.GetArgsSlice()
	fmt.Println("GetArgsSlice(): ", string(argsSlice))

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

func (cc *ERC20Chaincode) totalSupply(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success([]byte("totalSupply call"))
}

func (cc *ERC20Chaincode) balanceOf(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success([]byte("balanceOf call !!"))
}

func (cc *ERC20Chaincode) transfer(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success([]byte("transfer call!!!!!"))
}

func (cc *ERC20Chaincode) allowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) approve(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *ERC20Chaincode) transferFrom(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
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
