package secp256k1

import (
	"fmt"
	"math/big"
)

// fieldP는 유한체의 위수 p = 2²⁵⁶ − 2³² − 977이다.
var fieldP = mustHex("fffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f")

// fieldElem은 유한체 F_p의 원소다.
type fieldElem struct {
	num *big.Int
}

// newFieldElem은 num의 복사본으로 원소를 만든다.
// num이 nil이거나 [0, p) 범위 밖이면 ErrFieldRange를 반환한다.
func newFieldElem(num *big.Int) (fieldElem, error) {
	if num.Sign() < 0 || num.Cmp(fieldP) >= 0 {
		return fieldElem{}, fmt.Errorf("%w: %x", ErrFieldRange, num)
	}
	return fieldElem{num: new(big.Int).Set(num)}, nil
}

// n은 a의 값을 반환한다. 제로값이면 0을 반환한다.
func (a fieldElem) n() *big.Int {
	if a.num == nil {
		return big.NewInt(0)
	}
	return a.num
}

// Num은 a의 값의 복사본을 반환한다. 제로값이면 0을 반환한다.
func (a fieldElem) Num() *big.Int {
	return new(big.Int).Set(a.n())
}

// Add는 a + b (mod p)를 반환한다.
func (a fieldElem) Add(b fieldElem) fieldElem {
	num := new(big.Int).Add(a.n(), b.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

// Sub는 a − b (mod p)를 반환한다.
func (a fieldElem) Sub(b fieldElem) fieldElem {
	num := new(big.Int).Sub(a.n(), b.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

// Neg는 −a (mod p)를 반환한다.
func (a fieldElem) Neg() fieldElem {
	num := new(big.Int).Neg(a.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

// Mul은 a · b (mod p)를 반환한다.
func (a fieldElem) Mul(b fieldElem) fieldElem {
	num := new(big.Int).Mul(a.n(), b.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

// Inv는 a⁻¹ (mod p)를 반환한다. a가 0이면 역원이 없으므로 ok는 false다.
func (a fieldElem) Inv() (r fieldElem, ok bool) {
	if a.n().Sign() == 0 {
		return fieldElem{}, false
	}
	num := new(big.Int).ModInverse(a.n(), fieldP)
	return fieldElem{num: num}, true
}

// Square는 a² (mod p)를 반환한다.
func (a fieldElem) Square() fieldElem {
	num := new(big.Int).Mul(a.n(), a.n())
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

// Cube는 a³ (mod p)를 반환한다.
func (a fieldElem) Cube() fieldElem {
	num := new(big.Int).Mul(a.n(), new(big.Int).Mul(a.n(), a.n()))
	num.Mod(num, fieldP)
	return fieldElem{num: num}
}

// Sqrt는 r² = a (mod p)인 r을 반환한다. 제곱근이 없으면 ok는 false다.
// 두 제곱근 r, p − r 중 어느 쪽이 반환되는지는 정하지 않는다.
func (a fieldElem) Sqrt() (r fieldElem, ok bool) {
	num := new(big.Int).ModSqrt(a.n(), fieldP)
	if num == nil {
		return fieldElem{}, false
	}
	return fieldElem{num: num}, true
}

// Equal은 a와 b가 같은 원소인지 반환한다. 제로값과 0은 같다.
func (a fieldElem) Equal(b fieldElem) bool {
	return a.n().Cmp(b.n()) == 0
}

// IsZero는 a가 0인지 반환한다.
func (a fieldElem) IsZero() bool {
	return a.n().Sign() == 0
}

// IsOdd는 a가 홀수인지 반환한다.
func (a fieldElem) IsOdd() bool {
	return a.n().Bit(0) != 0
}

// String은 디버깅용으로 a의 값을 16진수로 보여 준다.
func (a fieldElem) String() string {
	return fmt.Sprintf("fieldElem(0x%x)", a.n())
}

// P는 유한체의 위수 p의 복사본을 반환한다.
func P() *big.Int {
	return new(big.Int).Set(fieldP)
}
