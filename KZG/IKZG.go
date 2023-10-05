package kzg

import "math/big"

type KZGZSender interface {
	SendZ() *big.Int
}
