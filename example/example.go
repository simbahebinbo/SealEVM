package main

import (
	"encoding/hex"
	"fmt"
	"github.com/SealSC/SealEVM"
	"github.com/SealSC/SealEVM/crypto/hashes"
	"github.com/SealSC/SealEVM/environment"
	"github.com/SealSC/SealEVM/evmInt256"
	"github.com/SealSC/SealEVM/storage"
	"os"
	"time"
)

func logPrinter(logCache *storage.LogCache) {
	for _, l := range *logCache {
		for _, t := range l.Topics {
			fmt.Println("topic:", t)
		}
		fmt.Println("data:", l.Data)
		fmt.Println("data as string:", string(l.Data))
	}
}

//store result to memStorage
func storeResult(result *SealEVM.ExecuteResult, storage *memStorage) {
	for addr, cache := range result.StorageCache.CachedData {
		for key, v := range cache {
			storage.storage[addr+key] = v.Bytes()
		}
	}
}

//create a new evm
func newEvm(code []byte, callData []byte, caller []byte, ms *memStorage) *SealEVM.EVM {
	hash := hashes.Keccak256(code)
	hashInt := evmInt256.New(0)
	hashInt.SetBytes(hash)

	//same contract code has same address in this example
	cNamespace := hashInt
	contract := environment.Contract{
		Namespace: cNamespace,
		Code:      code,
		Hash:      hashInt,
	}

	var callHash [32]byte
	copy(callHash[12:], caller)
	callerInt, _ := evmInt256.HashBytesToEVMInt(callHash)
	evm := SealEVM.New(SealEVM.EVMParam{
		MaxStackDepth:  0,
		ExternalStore:  ms,
		ResultCallback: nil,
		Context: &environment.Context{
			Block: environment.Block{
				ChainID:    evmInt256.New(0),
				Coinbase:   evmInt256.New(0),
				Timestamp:  evmInt256.New(int64(time.Now().Second())),
				Number:     evmInt256.New(0),
				Difficulty: evmInt256.New(0),
				GasLimit:   evmInt256.New(10000000),
				Hash:       evmInt256.New(0),
			},
			Contract: contract,
			Transaction: environment.Transaction{
				Origin:   evmInt256.New(2),
				GasPrice: evmInt256.New(1),
				GasLimit: evmInt256.New(10000000),
			},
			Message: environment.Message{
				Caller: callerInt,
				Value:  evmInt256.New(0),
				Data:   callData,
			},
		},
	})

	return evm
}

func main() {
	//load SealEVM module
	SealEVM.Load()

	//create memStorage
	ms := &memStorage{}
	ms.storage = make(map[string][]byte)
	ms.contracts = make(map[string][]byte)

	//deploy contract
	evm := newEvm(deployCode, nil, caller, ms)
	ret, err := evm.ExecuteContract(false)

	//check error
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	//result data of ret is the deployed code of example contract
	contractCode := ret.ResultData

	//call Counter() to get current counter's value
	evm = newEvm(contractCode, callCounter, caller, ms)
	ret, _ = evm.ExecuteContract(false)

	//result of Counter()
	fmt.Println("counter: ", hex.EncodeToString(ret.ResultData))

	//call increaseFor("example")
	evm = newEvm(contractCode, callIncreaseFor, caller, ms)
	ret, _ = evm.ExecuteContract(false)

	//store the result to ms
	storeResult(&ret, ms)

	//the event logs
	logPrinter(ret.StorageCache.Logs)

	//call Counter to get counter's value after increase
	evm = newEvm(contractCode, callCounter, caller, ms)
	ret, err = evm.ExecuteContract(false)

	//result of Counter()
	fmt.Println("counter: ", hex.EncodeToString(ret.ResultData))
}
