package types

import (
	"math/big"
)

// Constants

// MAX_TICK_INDEX is the maximum tick index supported by the clmmpool program.
const MAX_TICK_INDEX int = 443636

// MIN_TICK_INDEX is the minimum tick index supported by the clmmpool program.
const MIN_TICK_INDEX int = -443636

// MAX_SQRT_PRICE is the maximum sqrt-price supported by the clmmpool program.
var MAX_SQRT_PRICE *big.Int

// TICK_ARRAY_SIZE is the number of initialized ticks that a tick-array account can hold.
const TICK_ARRAY_SIZE int = 64

// MIN_SQRT_PRICE is the minimum sqrt-price supported by the clmmpool program.
var MIN_SQRT_PRICE *big.Int

// FEE_RATE_DENOMINATOR is the denominator which the fee rate is divided on.
var FEE_RATE_DENOMINATOR *big.Int = big.NewInt(1_000_000)

func init() {
	MAX_SQRT_PRICE, _ = big.NewInt(0).SetString("79226673515401279992447579055", 10)
	MIN_SQRT_PRICE, _ = big.NewInt(0).SetString("4295048016", 10)
}
