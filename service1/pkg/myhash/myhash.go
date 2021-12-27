package myhash

import (
	"fmt"

	"golang.org/x/crypto/sha3"
)

func Hash(s string) string {
	sha3 := sha3.Sum256([]byte(s))
	return fmt.Sprintf("%x", sha3)
}
