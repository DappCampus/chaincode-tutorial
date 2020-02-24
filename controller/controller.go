package controller

import (
	"strconv"

	"github.com/erc20/repository"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

// Init is called when the chaincode is instantiated by the blockchain network.
// params - tokenName, symbol, owner(address), amount
func (cc *Controller) Init(stub shim.ChaincodeStubInterface, params []string) sc.Response {
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

	// save token meta data
	err = repository.SaveERC20Metadata(stub, tokenName, symbol, owner, amountUint)
	if err != nil {
		return shim.Error(err.Error())
	}

	// save owner balance
	err = repository.SaveBalance(stub, owner, amount)
	if err != nil {
		return shim.Error(err.Error())
	}

	// response
	return shim.Success(nil)
}
