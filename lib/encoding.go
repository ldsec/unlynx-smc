package libunlynxsmc

import (
	"math/big"
	//"github.com/henrycg/prio/config"
	//"github.com/henrycg/prio/utils"
)

//Encode contains the encodings of different operations
func Encode(x *big.Int, /*input_parameters []*big.Int, */operation string) []*big.Int {
	//JS: input_parameters has the value x, the value y (label for linear regression), the global min and the global max
	/*x := input_parameters[0]
	y := input_parameters[1]
	global_min := input_parameters[2]
	global_max := input_parameters[3]*/

	//JS: to be seen later, what should be the number of bits of the x and y values?
	var LinRegBits []int
	//y value
	LinRegBits = append(LinRegBits, 64)
	//x value
	LinRegBits = append(LinRegBits, 64)


	//JS: use these default values for now, "countMinBuckets": 32, "countMinHashes": 8
	nHashes := 8
	nBuckets := 32

	result := make([]*big.Int, 1)
	result[0] = x

	lenR := len(result)
	// if we use a 63-bit modulus
	if lenR > 0 && lenR <= 64 {
		lenR = 63
		// if we use a 87-bit modulus
	} else if lenR > 64 && lenR <= 88 {
		lenR = 87
		// if we use a 102-bit modulus
	} else if lenR > 88 && lenR <= 104 {
		lenR = 103
		// if we use a 265-bit modulus
	} else if lenR > 104 && lenR <= 266 {
		lenR = 255
	}

	switch operation {
		case "sum", "mean":
			result = append(result, IntNew(lenR, x) ...)
			break

		case "variance":
			/*//JS: first include x and its bit representation
			result = append(result, IntNew(lenR, x) ...)
			//JS: then include x^2 and its bit representation
			square := make([]*big.Int, 1)
			square[0] = new(big.Int).Exp(x, big.NewInt(2), nil)
			result = append(result, square ...)
			result = append(result, computePows(2, x) ...)*/
			result = append(result, intPowNew(lenR, 2, x) ...)
			break

		case "bool_AND", "bool_OR":
			//JS: Should this be done this way? (lamda zeros in prio!)
			if (x == big.NewInt(1)) {
				result = append(result, boolNew(true) ...)
			} else {result = append(result, boolNew(false) ...)}
			break

		case "min":
			total := nHashes * nBuckets
			values := make([]bool, total)
			//JS: set values[i] to true for i <= x and to false otherwise (as in the prio design)
			for i := 0; i < total; i++ {
				if big.NewInt(int64(i)).Cmp(x) < 1 {values[i] = true} else {values[i] = false}
			}
			result = append(result, countMinNew(nHashes, nBuckets, values)...)
			break

		case "lin_reg":
			result = append(result, linRegNew_updated(result, LinRegBits) ...)
			break
		}

	return result
}

func Decode(output []*big.Int, operation string) *big.Int {
	result := big.NewInt(int64(0))
	switch operation {
		case "sum", "mean":
			result = output[0]
			break
		case "bool_AND", "bool_OR":
			result = output[1]
			break
		case "min":
			for i := 1; i < len(output); i++ {
				if big.NewInt(int64(0)).Cmp(output[i]) == 0 {
					result = big.NewInt(int64(i-1))
					break
				}
			}
			break
		case "lin_reg":
			break
		}

	return result
}