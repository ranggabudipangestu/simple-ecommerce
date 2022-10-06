package helper

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateTransactionNumber() string {
	source := rand.NewSource(time.Now().UnixNano())
	randomNumber := rand.New(source)

	return `TRX-` + strconv.Itoa(randomNumber.Int())
}
