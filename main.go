package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/blockcypher/gox11hash"
	"golang.org/x/crypto/scrypt"
)

//DefaultMantissa2 - default mantissa for system
const DefaultMantissa2 = 0x0FFFF0

//MantissaMax - maximum mantissa value
const MantissaMax = 0xFFFFFF

//DefaultExponent - default exponent for system
const DefaultExponent = 29

var currentMantissa float64 = DefaultMantissa2
var currentExponent = DefaultExponent

//MaxTarget - max target for hash comparision
const MaxTarget string = "0ffff00000000000000000000000000000000000000000000000000000000000" //0x1dffff00
//Target = mantissa * 2^(8 * (exponent – 3)) (exponent - 1d)
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
	rand.Seed(time.Now().UnixNano())
	var hashStr string
	var transJSON string
	var avgDuration time.Duration

	maxUint64 := ^uint64(0)

	processingStart := time.Now()
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
			val := gox11hash.Sum(jsonStr)
			scryptHash, err := scrypt.Key(val, nil, 32768, 8, 1, 32)
			if err != nil {
				fmt.Println(err)
				return
			}

			hashStr = hex.EncodeToString(scryptHash)
			//fmt.Println(hashStr)

			//fmt.Println(countZeros(hashStr))
			//zeros := countZeros(hashStr)
			transJSON = string(jsonStr)

			if hashStr < MaxTarget {
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

	//currentMantissa = DefaultMantissa / newDifficulty
	SetExponent(DefaultMantissa2 / newDifficulty)
	fmt.Printf("new target: %x\n", uint64(currentMantissa))
	fmt.Printf("new exponent: %x\n", currentExponent)

	fmt.Printf("Target string: %s\n", TargetString())
}

func countZeros(hash string) int16 {
	var firstZeros int16
	for i := 0; i < len(hash); i++ {
		num, err := strconv.ParseInt(string(hash[i]), 10, 16)
		if err != nil {
			return 0
		}

		if num == 0 {
			firstZeros++
		} else {
			break
		}
	}
	return firstZeros
}

//SetExponent - s4et current exponent
func SetExponent(mantissa float64) {
	fmt.Printf("mantissa float hex: %x\n", math.Float64bits(mantissa))
	fmt.Printf("mantissa dec: %f\n", mantissa)
	fmt.Printf("mantissa dec exp: %e\n", mantissa)
	fmt.Printf("mantissa: %x\n", uint64(mantissa))
	//maximum mantissa minus 0xFF0000 to find first value after it
	for mantissa < (MantissaMax - 0xFF0000) {
		mantissa = mantissa * 2
		mantissa = mantissa * 2
		mantissa = mantissa * 2
		fmt.Printf("mantissa dec: %e\n", mantissa)
		fmt.Printf("mantissa: %x\n", uint64(mantissa))
		currentExponent--
	}
	currentMantissa = mantissa
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

	expPos := targetLen - currentExponent*2

	mIdx := 0
	for sIdx := expPos - len(mantissaStr); sIdx < expPos; sIdx++ {
		b[sIdx] = mantissaStr[mIdx]
		mIdx++
	}

	return string(b)
}
