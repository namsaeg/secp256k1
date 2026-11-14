package secp256k1

import (
	"fmt"
	"math/big"
)

var (
	fieldP       = mustHex("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f")
	fieldPMinus1 = new(big.Int).Sub(fieldP, big.NewInt(1))
)

type fieldElem struct {
	num *big.Int
}

func newFieldElem(num *big.Int) (fieldElem, error) {
	if num.Sign() < 0 || num.Cmp(fieldP) >= 0 {
		return fieldElem{}, fmt.Errorf("%w: %x", ErrFieldRange, num)
	}
	return fieldElem{num: new(big.Int).Set(num)}, nil
}

func (a fieldElem) n() *big.Int {
	if a.num == nil {
		return big.NewInt(0)
	}
	return a.num
}

func (a fieldElem) Add(b fieldElem) fieldElem {
	num := new(big.Int).Add(a.n(), b.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

func (a fieldElem) Sub(b fieldElem) fieldElem {
	num := new(big.Int).Sub(a.n(), b.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

func (a fieldElem) Neg() fieldElem {
	num := new(big.Int).Neg(a.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

func (a fieldElem) Mul(b fieldElem) fieldElem {
	num := new(big.Int).Mul(a.n(), b.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

func (a fieldElem) Inv() (r fieldElem, ok bool) {
	if a.n().Sign() == 0 {
		return fieldElem{}, false
	}
	num := new(big.Int).ModInverse(a.n(), fieldP)
	return fieldElem{num: num}, true
}

func (a fieldElem) Pow(e *big.Int) fieldElem {
	if a.n().Sign() == 0 {
		if e.Sign() == 0 {
			return fieldElem{num: big.NewInt(1)}
		}
		return fieldElem{}
	}
	reduced := new(big.Int).Mod(e, fieldPMinus1)
	num := new(big.Int).Exp(a.n(), reduced, fieldP)
	return fieldElem{num: num}
}

func (a fieldElem) Square() fieldElem {
	num := new(big.Int).Mul(a.n(), a.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

func (a fieldElem) Cube() fieldElem {
	num := new(big.Int).Mul(a.n(), new(big.Int).Mul(a.n(), a.n()))
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

func (a fieldElem) Sqrt() (r fieldElem, ok bool) {
	num := new(big.Int).ModSqrt(a.n(), fieldP)
	if num == nil {
		return fieldElem{}, false
	}
	return fieldElem{num: num}, true
}

func (a fieldElem) Equal(b fieldElem) bool {
	return a.n().Cmp(b.n()) == 0
}

func (a fieldElem) IsZero() bool {
	return a.n().Sign() == 0
}

func (a fieldElem) IsOdd() bool {
	return a.n().Bit(0) != 0
}

func (a fieldElem) String() string {
	return fmt.Sprintf("fieldElem(0x%064x)", a.n())
}

func P() *big.Int {
	return new(big.Int).Set(fieldP)
}
