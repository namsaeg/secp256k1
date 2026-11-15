package secp256k1

import "errors"

var (
	// ErrNotOnCurve는 좌표 (x, y)가 곡선 y² = x³ + 7 위에 있지 않을 때 반환된다.
	ErrFieldRange = errors.New("secp256k1: 체 원소가 [0, p) 범위 밖")

	// ErrFieldRange는 유한체 원소로 쓰려는 값이 [0, p) 범위 밖일 때 반환된다.
	ErrNotOnCurve = errors.New("secp256k1: 점이 곡선 y² = x³ + 7 위에 있지 않음")

	// ErrScalarRange는 스칼라가 nil이거나 [0, n) 범위 밖일 때 반환된다.
	ErrScalarRange = errors.New("secp256k1: 스칼라가 [0, n) 범위 밖")
)
