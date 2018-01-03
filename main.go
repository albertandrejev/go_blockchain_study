package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"time"

	"bitbucket.org/albert_andrejev/free_info/utils"
)

//DefaultMantissa - default mantissa for system
const DefaultMantissa = 0xFFFF00

//MantissaMax - maximum mantissa value
const MantissaMax = 0xFFFFFF00

//DefaultExponent - default exponent for system
const DefaultExponent = 57

var currentMantissa int64 = DefaultMantissa
var currentExponent int64 = DefaultExponent

//0ffff00000000000000000000000000000000000000000000000000000000000

//DefaultDifficulty Starting calculation difficulty
const DefaultDifficulty float64 = 1

//Avg10BlocksDuration - duration in seconds for processing 10 blocks. Used for difficulty calculation
const Avg10BlocksDuration float64 = 3000

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

//Transaction Simple transaction type
type Transaction struct {
	PrevTxID   string
	TxID       string
	Nonce      uint64
	ExtraNonce uint64
	Sign       string
	PubKey     string
}

func main() {
	var bigI = new(big.Int)

	rand.Seed(time.Now().UnixNano())
	var hashStr string
	var transJSON string
	var avgDuration time.Duration

	maxUint64 := ^uint64(0)

	processingStart := time.Now()
	currentTarget := GetTarget(currentMantissa, currentExponent)
	fmt.Printf("current target: %x\n", currentTarget)

	for i := 0; i < 10; i++ {
		start := time.Now()
		trans := &Transaction{Sign: RandStringRunes(16)}
		nonce := uint64(0)
		extraNonce := uint64(0)
		for {
			trans.Nonce = nonce
			trans.ExtraNonce = extraNonce
			jsonStr, err := json.Marshal(trans)
			if err != nil {
				fmt.Println(err)
				return
			}
			dataHash, err := utils.X12Hash(jsonStr)
			if err != nil {
				fmt.Println(err)
				return
			}

			hashStr = hex.EncodeToString(dataHash)
			//fmt.Println(hashStr)

			transJSON = string(jsonStr)
			hashI, success := bigI.SetString(hashStr, 16)
			if !success {
				fmt.Println("Not successed")
			}

			if hashI.Cmp(currentTarget) == -1 {
				fmt.Println("by max target")
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
		fmt.Println(string(transJSON))
		fmt.Println(hashStr)
		avgDuration = (avgDuration*time.Duration(i) + stop.Sub(start)) / time.Duration(i+1)
	}
	processDuration := time.Now().Sub(processingStart)

	newDifficulty := DefaultDifficulty * (Avg10BlocksDuration / processDuration.Seconds())

	fmt.Printf("Average Elapsed: %v\n", avgDuration)
	fmt.Printf("Total duration: %v\n", processDuration)
	fmt.Printf("new difficulty: %f\n", newDifficulty)

	SetTarget(newDifficulty)
	fmt.Printf("new target: %x\n", currentMantissa)
	fmt.Printf("new exponent: %x\n", currentExponent)
}

//GetTarget - calculate target number based on difficulty bits
func GetTarget(mantissa int64, exponent int64) *big.Int {
	var target = new(big.Int)
	//target := big.NewInt(DefaultMantissa)
	target.Exp(big.NewInt(16), big.NewInt(exponent), nil)
	target.Mul(target, big.NewInt(mantissa))

	return target
}

//SetTarget - set current exponent
func SetTarget(difficulty float64) {
	var val []byte
	var trimmedVal []byte
	oldTarget := GetTarget(currentMantissa, currentExponent)
	newTarget := new(big.Int)
	newTarget.Div(oldTarget, big.NewInt(int64(difficulty)))
	fmt.Printf("old target: %x\n", oldTarget)
	fmt.Printf("new target: %x\n", newTarget)

	targetBytes := newTarget.Bytes()
	first := true
	for idx := 0; idx < len(targetBytes); idx++ {
		if targetBytes[idx] == 0 && first {
			continue
		} else {
			first = false
			trimmedVal = append(trimmedVal, targetBytes[idx])
		}

	}

	fmt.Printf("target bytes: %s\n", hex.EncodeToString(trimmedVal))
	for idx := 0; (idx < 3) && (idx < (len(targetBytes) - 1)); idx++ {
		val = append(val, trimmedVal[idx])
	}

	tmpMantissa, err := strconv.ParseInt(hex.EncodeToString(val), 16, 64)
	if err == nil {
		currentMantissa = tmpMantissa
	} else {
		fmt.Println(err)
	}
	fmt.Printf("current mantissa: %x\n", currentMantissa)

	currentExponent = (int64(len(trimmedVal) - len(val))) * 2

	currentTarget := GetTarget(currentMantissa, currentExponent)
	fmt.Printf("new current target: %x\n", currentTarget)
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

	//expPos := targetLen - currentExponent*2

	//mIdx := 0
	/*for sIdx := expPos - len(mantissaStr); sIdx < expPos; sIdx++ {
		b[sIdx] = mantissaStr[mIdx]
		mIdx++
	}*/

	return string(b)
}
