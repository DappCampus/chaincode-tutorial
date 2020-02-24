package controller

import (
	"fmt"
	"strconv"

	"github.com/erc20/repository"
	"github.com/erc20/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Transfer is invoke function that moves amount token
// from the caller's address to recipient
// params - caller's address, recipient's address, amount of token
func (cc *Controller) Transfer(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 3
	if len(params) != 3 {
		return shim.Error("incorrect number of parameters")
	}

	callerAddress, recipientAddress, transferAmount := params[0], params[1], params[2]

	// check amount is integer & positive
	transferAmountInt, err := util.ConvertToPositive("transferAmount", transferAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	// get caller amount
	callerAmountInt, err := repository.GetBalance(stub, callerAddress, false)
	if err != nil {
		return shim.Error(err.Error())
	}

	// get recipient amount
	recipientAmountInt, err := repository.GetBalance(stub, recipientAddress, true)
	if err != nil {
		return shim.Error(err.Error())
	}

	// calculate amount
	callerResultAmount := *callerAmountInt - *transferAmountInt
	recipientResultAmount := *recipientAmountInt + *transferAmountInt

	// check callerReuslt Amount is positive
	if callerResultAmount < 0 {
		return shim.Error("caller's balance is not sufficient")
	}

	// save the caller's & recipient's amount
	err = repository.SaveBalance(stub, callerAddress, strconv.Itoa(callerResultAmount))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = repository.SaveBalance(stub, recipientAddress, strconv.Itoa(recipientResultAmount))
	if err != nil {
		return shim.Error(err.Error())
	}

	// emit transfer event
	err = repository.EmitTransferEvent(stub, callerAddress, recipientAddress, *transferAmountInt)

	return shim.Success([]byte("transfer Success"))
}

// Approve is invoke function that Sets amount as the allowance
// of spender over the owner tokens
// params - owner's address, spender's address, amount of token
func (cc *Controller) Approve(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 3
	if len(params) != 3 {
		return shim.Error("incorrect number of parameters")
	}

	ownerAddress, spenderAddress, allowanceAmount := params[0], params[1], params[2]

	// check amount is integer & positive
	allowanceAmountInt, err := util.ConvertToPositive("AllowanceAmount", allowanceAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	// save allowance amount
	err = repository.SaveAllowance(stub, ownerAddress, spenderAddress, allowanceAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	// emit approval event
	err = repository.EmitApprovalEvent(stub, ownerAddress, spenderAddress, *allowanceAmountInt)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("approve success"))
}

// TransferFrom is invoke function that Moves amount of tokens from sender(owner) to recipient
// using allowance of spender
// parmas - owner's address, spender's address, recipient's address, amount of token
func (cc *Controller) TransferFrom(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of parmas is 4
	if len(params) != 4 {
		return shim.Error("incorrect number of params")
	}

	ownerAddress, spenderAddress, recipientAddress, transferAmount := params[0], params[1], params[2], params[3]

	// check amount is integer & positive
	transferAmountInt, err := util.ConvertToPositive("TransferAmount", transferAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	// get allowance
	allowanceResponse := cc.Allowance(stub, []string{ownerAddress, spenderAddress})
	if allowanceResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error: " + allowanceResponse.GetMessage())
	}

	// convert allowance response paylaod to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("allowance must be positive")
	}

	// transfer from owner to recipient
	transferResponse := cc.Transfer(stub, []string{ownerAddress, recipientAddress, transferAmount})
	if transferResponse.GetStatus() >= 400 {
		return shim.Error("failed to transfer, error: " + transferResponse.GetMessage())
	}

	// decrease allowance amount
	approveAmountInt := allowanceInt - *transferAmountInt
	approveAmount := strconv.Itoa(approveAmountInt)

	// approve amount of tokens transfered
	approveResponse := cc.Approve(stub, []string{ownerAddress, spenderAddress, approveAmount})
	if approveResponse.GetStatus() >= 400 {
		return shim.Error("failed to approve, error: " + approveResponse.GetMessage())
	}

	return shim.Success([]byte("transferFrom success"))
}

// TransferOtherToken is invoke function that Moves amount other chaincode tokens
// from the caller's address to recipient
// params - chaincode name, caller's address, recipient's address, amount
func (cc *Controller) TransferOtherToken(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of parmas is 4
	if len(params) != 4 {
		return shim.Error("incorrect number of params")
	}

	chaincodeName, callerAddress, recipientAddress, transferAmount := params[0], params[1], params[2], params[3]

	// make arguments
	args := [][]byte{[]byte("transfer"), []byte(callerAddress), []byte(recipientAddress), []byte(transferAmount)}

	// get channel
	channel := stub.GetChannelID()

	// transfer other chaincode token
	transferResponse := stub.InvokeChaincode(chaincodeName, args, channel)
	if transferResponse.GetStatus() >= 400 {
		return shim.Error(fmt.Sprintf("failed to transfer %s, error: %s", chaincodeName, transferResponse.GetMessage()))
	}

	return shim.Success([]byte("transfer other token success"))
}

// IncreaseAllowance is invoke function that increases spender's allowance by owner
// params - owner's address, spender's address, amount of amount
func (cc *Controller) IncreaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of parmas is 3
	if len(params) != 3 {
		return shim.Error("incorrect number of params")
	}

	ownerAddress, spenderAddress, increaseAmount := params[0], params[1], params[2]

	// check amount is integer & positive
	increaseAmountInt, err := util.ConvertToPositive("IncreaseAmount", increaseAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	// get allowance
	allowanceResponse := cc.Allowance(stub, []string{ownerAddress, spenderAddress})
	if allowanceResponse.GetStatus() >= 400 {
		return shim.Error("failed to get allowance, error: " + allowanceResponse.GetMessage())
	}

	// convert allowance response paylaod to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("allowance must be positive")
	}

	// increase allowance
	resultAmountInt := allowanceInt + *increaseAmountInt
	resultAmount := strconv.Itoa(resultAmountInt)

	// call approve
	approveResponse := cc.Approve(stub, []string{ownerAddress, spenderAddress, resultAmount})
	if approveResponse.GetStatus() >= 400 {
		return shim.Error("failed to approve allowance, error: " + approveResponse.GetMessage())
	}

	return shim.Success([]byte("increaseAllowance success"))
}

// DecreaseAllowance is invoke function that decreases spender's allowance by owner
// params - owner's address, spender's address, amount of token
func (cc *Controller) DecreaseAllowance(stub shim.ChaincodeStubInterface, params []string) sc.Response {

	// check the number of params is 3
	if len(params) != 3 {
		return shim.Error("incorrect number of params")
	}

	ownerAddress, spenderAddress, decreaseAmount := params[0], params[1], params[2]

	// check amount is integer & positive
	decreaseAmountInt, err := util.ConvertToPositive("DecreaseAmount", decreaseAmount)
	if err != nil {
		return shim.Error(err.Error())
	}

	// get allowance
	allowanceResponse := cc.Allowance(stub, []string{ownerAddress, spenderAddress})
	if allowanceResponse.Status >= 400 {
		return shim.Error("failed to get allowance, error: " + allowanceResponse.GetMessage())
	}

	// convert allowance response payload to allowance data
	allowanceInt, err := strconv.Atoi(string(allowanceResponse.GetPayload()))
	if err != nil {
		return shim.Error("allowance must be positive")
	}

	// calculate allowance (allowance cannot be negative!!)
	resultAmountInt := allowanceInt - *decreaseAmountInt
	if resultAmountInt < 0 {
		resultAmountInt = 0
	}
	resultAmount := strconv.Itoa(resultAmountInt)

	// call approve
	approveResponse := cc.Approve(stub, []string{ownerAddress, spenderAddress, resultAmount})
	if approveResponse.GetStatus() >= 400 {
		return shim.Error("failed to approve allowance, error: " + approveResponse.GetMessage())
	}

	return shim.Success([]byte("decreaseAllowance success"))
}

func (cc *Controller) Mint(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}

func (cc *Controller) Burn(stub shim.ChaincodeStubInterface, params []string) sc.Response {
	return shim.Success(nil)
}
