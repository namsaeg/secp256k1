package secp256k1

import "errors"

var (
	ErrFieldRange  = errors.New("secp256k1: 체 원소가 [0, p) 범위 밖")
	ErrNotOnCurve  = errors.New("secp256k1: 점이 곡선 y² = x³ + 7 위에 있지 않음")
	ErrScalarRange = errors.New("secp256k1: 스칼라가 [0, n) 범위 밖")
)
