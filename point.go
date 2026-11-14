package secp256k1

import (
	"fmt"
	"math/big"
)

var (
	curveN    = mustHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")
	curveB    = fieldElem{num: big.NewInt(7)}
	feThree   = fieldElem{num: big.NewInt(3)}
	infinity  = Point{}
	generator = Point{
		x:        fieldElem{num: mustHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")},
		y:        fieldElem{num: mustHex("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")},
		isFinite: true,
	}
)

type Point struct {
	x, y     fieldElem
	isFinite bool
}

func NewPoint(x, y *big.Int) (Point, error) {
	fx, err := newFieldElem(x)
	if err != nil {
		return Point{}, fmt.Errorf("x 좌표: %w", err)
	}
	fy, err := newFieldElem(y)
	if err != nil {
		return Point{}, fmt.Errorf("y 좌표: %w", err)
	}
	if !isOnCurve(fx, fy) {
		return Point{}, fmt.Errorf("%w: (%x, %x)", ErrNotOnCurve, x, y)
	}
	return Point{x: fx, y: fy, isFinite: true}, nil
}

func (p Point) X() (x *big.Int, ok bool) {
	if !p.isFinite {
		return nil, false
	}
	return new(big.Int).Set(p.x.n()), true
}

func (p Point) Y() (y *big.Int, ok bool) {
	if !p.isFinite {
		return nil, false
	}
	return new(big.Int).Set(p.y.n()), true
}

func (p Point) IsFinite() bool {
	return p.isFinite
}

func (p Point) thirdPoint(x2, m fieldElem) Point {
	x3 := m.Square().Sub(p.x).Sub(x2)
	y3 := m.Mul(p.x.Sub(x3)).Sub(p.y)
	return Point{x: x3, y: y3, isFinite: true}
}

func (p Point) Add(q Point) Point {
	if !p.isFinite {
		return q
	}
	if !q.isFinite {
		return p
	}
	if p.x.Equal(q.x) {
		if !p.y.Equal(q.y) || p.y.IsZero() {
			return infinity
		}
		return p.Double()
	}
	dxInv, _ := q.x.Sub(p.x).Inv()
	m := q.y.Sub(p.y).Mul(dxInv)
	return p.thirdPoint(q.x, m)
}

func (p Point) Double() Point {
	if !p.isFinite || p.y.IsZero() {
		return infinity
	}
	twoYInv, _ := p.y.Add(p.y).Inv()
	m := feThree.Mul(p.x.Square()).Mul(twoYInv)
	return p.thirdPoint(p.x, m)
}

func (p Point) ScalarMul(scalar *big.Int) Point {
	reduced := new(big.Int).Mod(scalar, curveN)
	r := infinity
	for i := reduced.BitLen() - 1; i >= 0; i-- {
		r = r.Double()
		if reduced.Bit(i) == 1 {
			r = r.Add(p)
		}
	}
	return r
}

func (p Point) Neg() Point {
	if !p.isFinite {
		return infinity
	}
	return Point{x: p.x, y: p.y.Neg(), isFinite: true}
}

func (p Point) Equal(q Point) bool {
	if !p.isFinite || !q.isFinite {
		return p.isFinite == q.isFinite
	}
	return p.x.Equal(q.x) && p.y.Equal(q.y)
}

func (p Point) String() string {
	if !p.isFinite {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(0x%064x, 0x%064x)", p.x.n(), p.y.n())
}

func isOnCurve(x, y fieldElem) bool {
	return y.Square().Equal(x.Cube().Add(curveB))
}

func ScalarBaseMul(scalar *big.Int) Point {
	return generator.ScalarMul(scalar)
}

func mustHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("secp256k1: 잘못된 16진수 상수: " + s)
	}
	return n
}
