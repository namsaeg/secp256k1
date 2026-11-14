package secp256k1

import "errors"

var (
	ErrFieldRange = errors.New("secp256k1: 체 원소가 0 이상 p 미만이 아님")
	ErrNotOnCurve = errors.New("secp256k1: 점이 곡선 y² = x³ + 7 위에 있지 않음")
)
