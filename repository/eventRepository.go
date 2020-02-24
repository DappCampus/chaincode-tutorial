package repository

import (
	"encoding/json"

	"github.com/erc20/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const (
	TransferEventKey = "transferEvent"
	ApprovalEventKey = "approvalEvent"
)

func EmitTransferEvent(stub shim.ChaincodeStubInterface, sender, spender string, amount int) error {
	transferEvent := model.NewTransferEvent(sender, spender, amount)
	transferEventBytes, err := json.Marshal(transferEvent)
	if err != nil {
		return model.NewCustomError(model.MarshalErrorType, TransferEventKey, err.Error())
	}
	err = stub.SetEvent(TransferEventKey, transferEventBytes)
	if err != nil {
		return model.NewCustomError(model.SetEventErrorType, TransferEventKey, err.Error())
	}

	return nil
}

func EmitApprovalEvent(stub shim.ChaincodeStubInterface, owner, spender string, allowance int) error {
	approvalEvent := model.NewApproval(owner, spender, allowance)
	approvalBytes, err := json.Marshal(approvalEvent)
	if err != nil {
		return model.NewCustomError(model.MarshalErrorType, ApprovalEventKey, err.Error())
	}

	err = stub.SetEvent(ApprovalEventKey, approvalBytes)
	if err != nil {
		return model.NewCustomError(model.SetEventErrorType, ApprovalEventKey, err.Error())
	}

	return nil
}
