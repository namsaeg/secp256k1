package secp256k1

import (
	"fmt"
	"math/big"
)

var (
	// curveN은 생성점 G로 만들어지는 군의 위수 n이다.
	curveN = mustHex("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")

	// curveB는 곡선 방정식 y² = x³ + 7의 상수 7이다.
	curveB = fieldElem{num: big.NewInt(7)}

	// feThree는 값이 3인 유한체 원소이다.
	feThree = fieldElem{num: big.NewInt(3)}

	// infinity는 무한원점이다.
	infinity = Point{}

	// generator는 secp256k1의 생성점 G다.
	generator = Point{
		x:        fieldElem{num: mustHex("79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798")},
		y:        fieldElem{num: mustHex("483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")},
		isFinite: true,
	}
)

// Point는 secp256k1 위의 점이다.
type Point struct {
	x, y     fieldElem
	isFinite bool
}

// NewPoint는 좌표 (x, y)로 점을 만든다.
// x, y가 nil이거나 [0, p) 범위 밖이면 ErrFieldRange를,
// 곡선 위에 있지 않으면 ErrNotOnCurve를 감싼 에러를 반환한다.
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

// LiftX는 x좌표가 x이고 y가 짝수인 점을 반환한다.
// x가 [0, p) 범위 밖이면 ErrFieldRange를,
// 그런 점이 없으면 ErrNotOnCurve를 감싼 에러를 반환한다.
func LiftXEvenY(x *big.Int) (Point, error) {
	fx, err := newFieldElem(x)
	if err != nil {
		return Point{}, fmt.Errorf("x 좌표: %w", err)
	}
	fy, ok := fx.Cube().Add(curveB).Sqrt()
	if !ok {
		return Point{}, fmt.Errorf("%w: x = %x", ErrNotOnCurve, x)
	}
	if fy.IsOdd() {
		fy = fy.Neg()
	}
	return Point{
		x:        fx,
		y:        fy,
		isFinite: true,
	}, nil
}

// ScalarBaseMul은 생성점에 대한 스칼라 곱셉의 결과를 반환한다.
// scalar가 [0, n) 범위 밖이면 ErrScalarRange를 감싼 에러를 반환한다.
func ScalarBaseMul(scalar *big.Int) (Point, error) {
	return generator.ScalarMul(scalar)
}

// X는 x좌표의 복사본을 반환한다. 무한원점이면 ok는 false다.
func (p Point) X() (x *big.Int, ok bool) {
	if !p.isFinite {
		return nil, false
	}
	return new(big.Int).Set(p.x.n()), true
}

// Y는 y좌표의 복사본을 반환한다. 무한원점이면 ok는 false다.
func (p Point) Y() (y *big.Int, ok bool) {
	if !p.isFinite {
		return nil, false
	}
	return new(big.Int).Set(p.y.n()), true
}

// IsFinite는 p가 무한원점이 아닌지 반환한다.
func (p Point) IsFinite() bool {
	return p.isFinite
}

// HasEvenY는 p가 유한한 점이고 y가 짝수인지 반환한다.
// 무한원점이거나 y가 홀수면 false다.
func (p Point) HasEvenY() bool {
	return p.isFinite && !p.y.IsOdd()
}

// Add는 p + q를 반환한다
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

// Double은 2p를 반환한다.
func (p Point) Double() Point {
	if !p.isFinite || p.y.IsZero() {
		return infinity
	}
	twoYInv, _ := p.y.Add(p.y).Inv()
	m := feThree.Mul(p.x.Square()).Mul(twoYInv)
	return p.thirdPoint(p.x, m)
}

// ScalarMul은 p에 대한 스칼라 곱셈의 결과를 반환한다.
func (p Point) ScalarMul(scalar *big.Int) (Point, error) {
	if scalar == nil || scalar.Sign() < 0 || scalar.Cmp(curveN) >= 0 {
		return Point{}, ErrScalarRange
	}
	reduced := new(big.Int).Mod(scalar, curveN)
	r := infinity
	for i := reduced.BitLen() - 1; i >= 0; i-- {
		r = r.Double()
		if reduced.Bit(i) == 1 {
			r = r.Add(p)
		}
	}
	return r, nil
}

// Neg는 −p를 반환한다.
func (p Point) Neg() Point {
	if !p.isFinite {
		return infinity
	}
	return Point{x: p.x, y: p.y.Neg(), isFinite: true}
}

// Equal은 p와 q가 같은 점인지 반환한다. 제로값과 무한원점은 같은 값이다.
func (p Point) Equal(q Point) bool {
	if !p.isFinite || !q.isFinite {
		return p.isFinite == q.isFinite
	}
	return p.x.Equal(q.x) && p.y.Equal(q.y)
}

// String은 디버깅용으로 좌표를 보여 준다. 무한원점이면 "Point(infinity)"다.
func (p Point) String() string {
	if !p.isFinite {
		return "Point(infinity)"
	}
	return fmt.Sprintf("Point(0x%x, 0x%x)", p.x.n(), p.y.n())
}

// thirdPoint는 기울기 m인 직선이 곡선과 만나는 세 번째 점을 x축에 대칭시킨 점,
// 즉 p와 (x2, …)를 지나는 직선으로 얻는 덧셈 결과를 반환한다.
func (p Point) thirdPoint(x2, m fieldElem) Point {
	x3 := m.Square().Sub(p.x).Sub(x2)
	y3 := m.Mul(p.x.Sub(x3)).Sub(p.y)
	return Point{x: x3, y: y3, isFinite: true}
}

// isOnCurve는 (x, y)가 y² = x³ + 7을 만족하는지 반환한다.
func isOnCurve(x, y fieldElem) bool {
	return y.Square().Equal(x.Cube().Add(curveB))
}

// N은 군의 위수 n의 복사본을 반환한다.
func N() *big.Int {
	return new(big.Int).Set(curveN)
}

// mustHex는 16진수 문자열 s를 정수로 바꾼다.
// 패키지 초기화용 상수에만 쓰며, s가 올바른 16진수가 아니면 패닉한다.
func mustHex(s string) *big.Int {
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		panic("secp256k1: 잘못된 16진수 상수: " + s)
	}
	return n
}
