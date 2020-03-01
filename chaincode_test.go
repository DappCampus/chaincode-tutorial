/*
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/erc20/model"
	"github.com/erc20/repository"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var function = []byte("mint")

const txMint = "txMint"
const address = "dappcampus"
const initAmount = 100000

const tokenName = "dappToken"

func Test_Init_success(t *testing.T) {
	cc := NewChaincode()
	stub := shim.NewMockStub("erc20", cc)
	res := stub.MockInit("1", [][]byte{[]byte("init"), []byte(tokenName), []byte("dt"), []byte(address), []byte(strconv.Itoa(initAmount))})
	if res.Status != shim.OK {
		t.FailNow()
	}

	// check totalSupply
	erc20 := model.ERC20Metadata{}
	erc20Bytes, _ := stub.GetState(tokenName)
	json.Unmarshal(erc20Bytes, &erc20)
	if *erc20.GetTotalSupply() != initAmount {
		t.FailNow()
	}

	// check dappcampus balance
	balance, _ := repository.GetBalance(stub, address, true)
	if *balance != initAmount {
		t.FailNow()
	}
}

func initERC20(t *testing.T) *shim.MockStub {
	cc := NewChaincode()
	stub := shim.NewMockStub("erc20", cc)
	res := stub.MockInit("1", [][]byte{[]byte("init"), []byte(tokenName), []byte("dt"), []byte(address), []byte(strconv.Itoa(initAmount))})
	if res.Status != shim.OK {
		t.FailNow()
	}
	return stub
}

func Test_Mint_lengthIsInvalid_failure(t *testing.T) {
	stub := initERC20(t)
	arguments := [][]byte{function, []byte(tokenName), []byte(address)}
	res := stub.MockInvoke(txMint, arguments)
	if res.Status != shim.ERROR {
		t.FailNow()
	}
}

func Test_Mint_amountIsNotPositive_failure(t *testing.T) {
	stub := initERC20(t)
	arguments := [][]byte{function, []byte(tokenName), []byte(address), []byte("-100")}
	arguments2 := [][]byte{function, []byte(tokenName), []byte(address), []byte("abcde")}
	res := stub.MockInvoke(txMint, arguments)
	res2 := stub.MockInvoke(txMint, arguments2)

	if res.Status != shim.ERROR && res2.Status != shim.ERROR {
		t.FailNow()
	}
}

func Test_Mint_success(t *testing.T) {
	stub := initERC20(t)
	const increaseAmount = 10000
	arguments := [][]byte{function, []byte(tokenName), []byte(address), []byte(strconv.Itoa(increaseAmount))}
	res := stub.MockInvoke(txMint, arguments)
	if res.Status != shim.OK {
		t.FailNow()
	}

	// increase TotalSupply
	totalSupply, _ := repository.GetERC20TotalSupply(stub, tokenName)
	if *totalSupply != initAmount+increaseAmount {
		t.FailNow()
	}

	// increase owner balance
	balance, _ := repository.GetBalance(stub, address, true)
	if *balance != initAmount+increaseAmount {
		t.FailNow()
	}

	// emit trasnfer event
	data := <-stub.ChaincodeEventsChannel
	if data.GetEventName() != repository.TransferEventKey {
		t.FailNow()
	}
	event := model.NewTransferEvent("admin", address, increaseAmount)
	eventBytes, _ := json.Marshal(event)
	if string(data.GetPayload()) != string(eventBytes) {
		t.FailNow()
	}
}
