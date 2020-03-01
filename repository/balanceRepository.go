package repository

import (
	"encoding/json"
	"strconv"

	"github.com/erc20/model"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func SaveERC20Metadata(stub shim.ChaincodeStubInterface, tokenName, symbol, owner string, amount uint64) error {
	// make metadata
	erc20 := model.NewERC20MetaData(tokenName, symbol, owner, amount)
	erc20Bytes, err := json.Marshal(erc20)
	if err != nil {
		return model.NewCustomError(model.MarshalErrorType, "erc20", err.Error())
	}

	// save token meta data
	err = stub.PutState(tokenName, erc20Bytes)
	if err != nil {
		return model.NewCustomError(model.PutStateErrorType, "erc20Metadata", err.Error())
	}

	return nil
}

func GetERC20Metadata(stub shim.ChaincodeStubInterface, tokenName string) (*model.ERC20Metadata, error) {
	// Get ERC20 Metadata
	erc20 := model.ERC20Metadata{}
	erc20Bytes, err := stub.GetState(tokenName)
	if err != nil {
		return nil, model.NewCustomError(model.GetStateErrorType, "balance", err.Error())
	}
	err = json.Unmarshal(erc20Bytes, &erc20)
	if err != nil {
		return nil, model.NewCustomError(model.UnMarshalErrorType, "erc20Metadata", err.Error())
	}
	return &erc20, nil
}

func GetERC20TotalSupply(stub shim.ChaincodeStubInterface, tokenName string) (*uint64, error) {
	// Get ERC20 Metadata
	erc20 := model.ERC20Metadata{}
	erc20Bytes, err := stub.GetState(tokenName)
	if err != nil {
		return nil, model.NewCustomError(model.GetStateErrorType, "balance", err.Error())
	}
	err = json.Unmarshal(erc20Bytes, &erc20)
	if err != nil {
		return nil, model.NewCustomError(model.UnMarshalErrorType, "erc20Metadata", err.Error())
	}
	return erc20.GetTotalSupply(), nil
}

func SaveBalance(stub shim.ChaincodeStubInterface, owner, balance string) error {
	err := stub.PutState(owner, []byte(balance))
	if err != nil {
		return model.NewCustomError(model.PutStateErrorType, "balance", err.Error())
	}

	return nil
}

func GetBalanceBytes(stub shim.ChaincodeStubInterface, owner string, isZeror bool) ([]byte, error) {
	amountBytes, err := stub.GetState(owner)
	if err != nil {
		return nil, model.NewCustomError(model.GetStateErrorType, owner, err.Error())
	}
	if amountBytes == nil {
		amountBytes = []byte("0")
	}
	return amountBytes, nil
}

func GetBalance(stub shim.ChaincodeStubInterface, owner string, isZero bool) (*int, error) {
	amountBytes, err := stub.GetState(owner)
	if err != nil {
		return nil, model.NewCustomError(model.GetStateErrorType, "balance", err.Error())
	}

	if isZero && amountBytes == nil {
		amountBytes = []byte("0")
	}

	amount, err := strconv.Atoi(string(amountBytes))
	if err != nil {
		return nil, model.NewCustomError(model.ConvertErrorType, "amount", err.Error())
	}

	return &amount, nil
}
