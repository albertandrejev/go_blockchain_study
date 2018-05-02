package main

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"bitbucket.org/albert_andrejev/free_info/merkle"
	"bitbucket.org/albert_andrejev/free_info/wrappers"

	"bitbucket.org/albert_andrejev/free_info/factory"
	"bitbucket.org/albert_andrejev/free_info/types"
	"bitbucket.org/albert_andrejev/free_info/utils"
)

//DefaultMantissa - default mantissa for system
const DefaultMantissa = 0xFFFF00

//DefaultExponent - default exponent for system
const DefaultExponent = 56

var currentMantissa int64 = DefaultMantissa
var currentExponent int64 = DefaultExponent

//DefaultDifficulty Starting calculation difficulty
const DefaultDifficulty float64 = 1

var currentDifficulty = DefaultDifficulty

//AvgBlocksDuration - duration in seconds for processing 10 blocks. Used for difficulty calculation
const AvgBlocksDuration float64 = 300

//AvgBlocksAmount amount of blocks used to calculate power of the network
const AvgBlocksAmount int = 10

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

var blockChain []*types.Block

//Transaction Simple transaction type

func main() {
	var bigI = new(big.Int)
	factory := factory.NewMainFactory()

	rand.Seed(time.Now().UnixNano())
	var hashStr string
	var avgDuration time.Duration

	maxUint64 := ^uint64(0)

	processingStart := time.Now()
	simpleHash := factory.GetSimpleHash()
	x12Hash := factory.GetX11Hash()

	prevBlockID := "0000000000000000000000000000000000000000000000000000000000000000"

	for i := 0; i < 120; i++ {
		start := time.Now()
		fmt.Printf("-------------New block. Height: %d-------------\n", len(blockChain))
		block := CreateBlock(prevBlockID)
		currentTarget := GetTarget(currentMantissa, currentExponent)
		fmt.Printf("Target: %x\n", currentTarget)

		jsonWrapper := new(wrappers.JSONWrapper)
		block.Transactions = append(block.Transactions, CreateTransaction(simpleHash, jsonWrapper))
		block.Transactions = append(block.Transactions, CreateTransaction(simpleHash, jsonWrapper))
		block.Transactions = append(block.Transactions, CreateTransaction(simpleHash, jsonWrapper))
		block.Transactions = append(block.Transactions, CreateTransaction(simpleHash, jsonWrapper))
		block.Transactions = append(block.Transactions, CreateTransaction(simpleHash, jsonWrapper))
		block.Transactions = append(block.Transactions, CreateTransaction(simpleHash, jsonWrapper))

		merkle := merkle.NewTree(factory, jsonWrapper)
		allSums := merkle.Init(block.Transactions)
		block.Data.MerkleRoot = hex.EncodeToString(merkle.CalcRoot(allSums))
		nonce := uint64(0)
		extraNonce := uint64(0)
		for {
			block.Data.Nonce = nonce
			block.Data.ExtraNonce = string(extraNonce)
			blockDataJSON, err := jsonWrapper.Encode(block.Data)
			if err != nil {
				fmt.Println(err)
				continue
			}
			dataHash, err := x12Hash.Sum256(blockDataJSON)
			if err != nil {
				fmt.Println(err)
				continue
			}
			hashStr = hex.EncodeToString(dataHash)
			//fmt.Println(hashStr)

			hashI := bigI.SetBytes(dataHash)

			if hashI.Cmp(currentTarget) == -1 {
				block.BlockID = hashStr
				prevBlockID = hashStr
				if CheckBlock(block, factory) {
					blockChain = append(blockChain, block)
					fmt.Println("Block was added to blockchain")
				} else {
					fmt.Println("Block was NOT added to blockchain")
				}

				break
			}

			if nonce == maxUint64 {
				nonce = 0
				extraNonce++
			} else {
				nonce++
			}
		}
		stop := time.Now()
		fmt.Printf("Elapsed: %v\n", stop.Sub(start))

		/*blockJSON, err := jsonWrapper.Encode(block)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(blockJSON))*/
		fmt.Printf("Block ID: %v\n", hashStr)
		avgDuration = (avgDuration*time.Duration(i) + stop.Sub(start)) / time.Duration(i+1)
		fmt.Printf("Average Duration: %v\n", avgDuration)
	}
	processDuration := time.Now().Sub(processingStart)

	fmt.Printf("Average Elapsed: %v\n", avgDuration)
	fmt.Printf("Total duration: %v\n", processDuration)
}

//CheckBlock for validity
func CheckBlock(block *types.Block, factory factory.IMainFactory) bool {
	//1. check block id hash
	//1.2. check hash difficulty
	//2. check transactions

	//1. check block id hash
	if CheckBlockID(block, factory) == false {
		return false
	}

	return true
}

//GetCurrentNetworkDifficulty network difficulty calculation
func GetCurrentNetworkDifficulty(blockTimestamp int64) *big.Float {
	totalBlocks := len(blockChain)
	if totalBlocks > AvgBlocksAmount {
		startBlockIdx := totalBlocks - AvgBlocksAmount
		startBlock := blockChain[startBlockIdx]
		lastBlock := blockChain[totalBlocks-1]

		startTime := time.Unix(startBlock.Data.Timestamp, 0)
		endTime := time.Unix(blockTimestamp, 0)

		processDuration := endTime.Sub(startTime)
		prevMantissa, prevExponent := SeparateTarget(lastBlock.Data.Target)
		prevTarget := GetTarget(prevMantissa, prevExponent)
		fmt.Printf("prev block target: %x\n", prevTarget)
		defaultTarget := GetTarget(DefaultMantissa, DefaultExponent)

		lastBlockDifficulty := new(big.Int)
		lastBlockDifficulty.Div(defaultTarget, prevTarget)
		fmt.Printf("last difficulty: %v\n", lastBlockDifficulty)

		fmt.Printf("Actual %d blocks duration: %v\n", AvgBlocksAmount, processDuration)
		difficultyChange := big.NewFloat(AvgBlocksDuration / processDuration.Seconds())
		difficultyChangeF, accuracyD := difficultyChange.Float64()
		fmt.Printf("difficulty change: %f, accuracy: %x\n", difficultyChangeF, accuracyD)

		newDifficulty := new(big.Float)
		lastBlockDifficultyFloat := new(big.Float).SetInt(lastBlockDifficulty)
		newDifficulty.Mul(lastBlockDifficultyFloat, difficultyChange)
		newDifficultyF, accuracy := newDifficulty.Float64()
		fmt.Printf("new difficulty: %f, accuracy: %x\n", newDifficultyF, accuracy)
		//newDifficulty := big.NewFloat(lastBlockDifficulty.Bytes()) * difficultyChange
		return newDifficulty
	}

	return big.NewFloat(DefaultDifficulty)
}

/*func CheckBlockDifficulty(block *types.Block) bool {
	totalBlocks := len(blockChain)
	if totalBlocks > AvgBlocksAmount {
		startBlockIdx := totalBlocks - AvgBlocksAmount
		startBlock := blockChain[startBlockIdx]
		lastBlock := blockChain[totalBlocks-1]

		startTime := time.Unix(startBlock.Data.Timestamp, 0)
		endTime := time.Unix(lastBlock.Data.Timestamp, 0)

		processDuration := endTime.Sub(startTime)
		prevMantissa, prevExponent := SeparateTarget(lastBlock.Data.Target)
		prevTarget := GetTarget(prevMantissa, prevExponent)
		fmt.Printf("prev block target: %x\n", prevTarget)
		defaultTarget := GetTarget(DefaultMantissa, DefaultExponent)

		lastBlockDifficulty := defaultTarget.Div(defaultTarget, prevTarget)
		fmt.Printf("last difficulty: %x\n", lastBlockDifficulty.Int64())

		newDifficulty := float64(lastBlockDifficulty.Int64()) * (AvgBlocksDuration / processDuration.Seconds())
		fmt.Printf("new difficulty: %f\n", newDifficulty)

		SetTarget(newDifficulty)
	}
}*/

//CheckBlockID - current block id
func CheckBlockID(block *types.Block, factory factory.IMainFactory) bool {
	var bigI = new(big.Int)
	jsonWrapper := new(wrappers.JSONWrapper)
	x12Hash := factory.GetX11Hash()

	blockDataJSON, err := jsonWrapper.Encode(block.Data)
	if err != nil {
		fmt.Println(err)
		return false
	}
	dataHash, err := x12Hash.Sum256(blockDataJSON)
	if err != nil {
		fmt.Println(err)
		return false
	}

	hashI := bigI.SetBytes(dataHash)

	blockMantissa, blockExponent := SeparateTarget(block.Data.Target)
	blockTarget := GetTarget(blockMantissa, blockExponent)

	if hashI.Cmp(blockTarget) == -1 {
		return true
	}

	return false
}

//CreateBlock create block of transactions
func CreateBlock(prevBlockID string) *types.Block {
	/*totalBlocks := len(blockChain)
	if totalBlocks > AvgBlocksAmount {
		startBlockIdx := totalBlocks - AvgBlocksAmount
		startBlock := blockChain[startBlockIdx]
		lastBlock := blockChain[totalBlocks-1]

		startTime := time.Unix(startBlock.Data.Timestamp, 0)
		endTime := time.Unix(lastBlock.Data.Timestamp, 0)

		processDuration := endTime.Sub(startTime)
		fmt.Printf("10 blocks duration: %v\n", processDuration)
		fmt.Printf("Average for one block: %v\n", processDuration.Seconds()/float64(AvgBlocksAmount))

		prevMantissa, prevExponent := SeparateTarget(lastBlock.Data.Target)
		prevTarget := GetTarget(prevMantissa, prevExponent)
		fmt.Printf("prev block target: %x\n", prevTarget)
		defaultTarget := GetTarget(DefaultMantissa, DefaultExponent)

		lastBlockDifficulty := defaultTarget.Div(defaultTarget, prevTarget)
		fmt.Printf("last difficulty: %x\n", lastBlockDifficulty.Int64())

		newDifficulty := float64(lastBlockDifficulty.Int64()) * (AvgBlocksDuration / processDuration.Seconds())
		fmt.Printf("new difficulty: %f\n", newDifficulty)

		SetTarget(newDifficulty)
	}*/

	blockTimestamp := time.Now().Unix()

	newDifficulty := GetCurrentNetworkDifficulty(blockTimestamp)
	if newDifficulty.Cmp(big.NewFloat(DefaultDifficulty)) != 0 {
		SetTarget(newDifficulty)
	}

	target := uint32(currentMantissa<<8 + currentExponent)

	return &types.Block{
		Data: types.BlockData{
			PrevBlockID: prevBlockID,
			Target:      target,
			Timestamp:   blockTimestamp,
		},
	}
}

//CreateTransaction return transaction
func CreateTransaction(simpleHash utils.ISimpleHash, json wrappers.IJSONWrapper) *types.Transaction {

	trans := &types.Transaction{
		Data: types.TransactionData{
			PubKey:    RandStringRunes(16),
			Timestamp: time.Now().Unix(),
		},
	}
	transDataJSON, err := json.Encode(trans.Data)
	if err != nil {
		fmt.Println(err)
	}
	txID := simpleHash.Sum256(transDataJSON)
	trans.TxID = hex.EncodeToString(txID)

	return trans
}

//GetTarget - calculate target number based on difficulty bits
func GetTarget(mantissa int64, exponent int64) *big.Int {
	var target = new(big.Int)
	target.Exp(big.NewInt(16), big.NewInt(exponent), nil)
	target.Mul(target, big.NewInt(mantissa))

	return target
}

//SeparateTarget return mantissa and exponent from target
func SeparateTarget(target uint32) (int64, int64) {
	exponent := target & 0x000000ff
	mantissa := target >> 8

	fmt.Printf("mantissa: %x, exponent: %x\n", mantissa, exponent)

	return int64(mantissa), int64(exponent)
}

//SetTarget - set current exponent
func SetTarget(difficulty *big.Float) {
	if difficulty.Cmp(big.NewFloat(DefaultDifficulty)) <= 0 {
		fmt.Println("Set default target")
		currentMantissa = DefaultMantissa
		currentExponent = DefaultExponent
		return
	}
	var val []byte
	//var trimmedVal []byte
	defaultTargetInt := GetTarget(DefaultMantissa, DefaultExponent)
	defaultTarget := new(big.Float).SetInt(defaultTargetInt)

	/*if int64(difficulty+0.5) <= 0 {
		return
	}*/

	newTarget := new(big.Float)
	newTarget.Quo(defaultTarget, difficulty)
	newTargetInt := new(big.Int)
	newTarget.Int(newTargetInt)

	fmt.Printf("old target: %x\n", defaultTargetInt)
	fmt.Printf("new target: %x\n", newTargetInt)

	targetBytes := newTargetInt.Bytes()
	/*first := true
	for idx := 0; idx < len(targetBytes); idx++ {
		if targetBytes[idx] == 0 && first {
			continue
		} else {
			first = false
			trimmedVal = append(trimmedVal, targetBytes[idx])
		}

	}

	fmt.Printf("target bytes: %s\n", hex.EncodeToString(trimmedVal))*/
	for idx := 0; (idx < 3) && (idx < (len(targetBytes) - 1)); idx++ {
		val = append(val, targetBytes[idx])
	}

	tmpMantissa, err := strconv.ParseInt(hex.EncodeToString(val), 16, 64)
	if err == nil {
		currentMantissa = tmpMantissa
	} else {
		fmt.Println(err)
	}
	//currentMantissa = newTarget.MantExp
	fmt.Printf("current mantissa: %x\n", currentMantissa)

	currentExponent = (int64(len(targetBytes) - len(val))) * 2

	currentTarget := GetTarget(currentMantissa, currentExponent)
	fmt.Printf("new current target: %x\n", currentTarget)

	currentNormalDifficulty := new(big.Int)
	currentNormalDifficulty.Div(defaultTargetInt, currentTarget)
	fmt.Printf("actual current difficulty: %d\n", currentNormalDifficulty.Int64())
}

/*
RandStringRunes Random string function
*/
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

//TargetString return target string representation
func TargetString() string {
	targetLen := 64
	b := make([]byte, targetLen)

	mantissaStr := fmt.Sprintf("%x", uint64(currentMantissa))
	fmt.Printf("mantissa str: %s\n", mantissaStr)

	for i := range b {
		b[i] = '0'
	}

	return string(b)
}
