package repository

import (
	"strconv"

	"github.com/erc20/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const approvalCompositeKey = "approval"

func SaveAllowance(stub shim.ChaincodeStubInterface, owner, spender, allowance string) error {
	// create composite key for allowance - approval/{owner}/{spender}
	approvalKey, err := stub.CreateCompositeKey(approvalCompositeKey, []string{owner, spender})
	if err != nil {
		return model.NewCustomError(model.CreateCompositeKeyErrorType, approvalCompositeKey, err.Error())
	}

	// save allowance amount
	err = stub.PutState(approvalKey, []byte(allowance))
	if err != nil {
		return model.NewCustomError(model.PutStateErrorType, approvalKey, err.Error())
	}

	return nil
}

func GetAllowanceBytes(stub shim.ChaincodeStubInterface, owner, spender string, isZero bool) ([]byte, error) {
	// create composite key
	approvalKey, err := stub.CreateCompositeKey(approvalCompositeKey, []string{owner, spender})
	if err != nil {
		return nil, model.NewCustomError(model.CreateCompositeKeyErrorType, approvalCompositeKey, err.Error())
	}

	allowanceBytes, err := stub.GetState(approvalKey)
	if err != nil {
		return nil, model.NewCustomError(model.GetStateErrorType, approvalKey, err.Error())
	}

	if isZero && allowanceBytes == nil {
		allowanceBytes = []byte("0")
	}

	return allowanceBytes, nil
}

func GetApprovalList(stub shim.ChaincodeStubInterface, owner string) ([]model.Approval, error) {
	// get all approval list (format is iterator)
	approvalIterator, err := stub.GetStateByPartialCompositeKey(approvalCompositeKey, []string{owner})
	if err != nil {
		return nil, model.NewCustomError(model.GetStatePartialCompositeKeyErrorType, approvalCompositeKey, err.Error())
	}

	// make slice for return value
	approvalSlice := []model.Approval{}

	// iterator
	defer approvalIterator.Close()
	if approvalIterator.HasNext() {
		for approvalIterator.HasNext() {
			approvalKV, _ := approvalIterator.Next()

			// get spender address
			_, addresses, err := stub.SplitCompositeKey(approvalKV.GetKey())
			if err != nil {
				return nil, model.NewCustomError(model.SpliteCompositeKeyErrorType, approvalKV.GetKey(), err.Error())
			}
			spenderAddress := addresses[1]

			// get amount
			amountBytes := approvalKV.GetValue()
			amountInt, err := strconv.Atoi(string(amountBytes))
			if err != nil {
				return nil, model.NewCustomError(model.ConvertErrorType, string(amountBytes), err.Error())
			}

			// add approval result
			approval := model.Approval{Owner: owner, Spender: spenderAddress, Allowance: amountInt}
			approvalSlice = append(approvalSlice, approval)
		}
	}

	return approvalSlice, nil
}
