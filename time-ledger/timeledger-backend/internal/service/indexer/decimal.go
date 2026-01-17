package indexer

import (
	"fmt"
	"math/big"
	"strings"
)

// deltaWeiStr 可能带负号，例如 "-1000000000000000000"
func addDecimalIntString(a, b string) (string, error) {
	ai, ok := new(big.Int).SetString(strings.TrimSpace(a), 10)
	if !ok {
		return "", fmt.Errorf("invalid int string a=%s", a)
	}
	bi, ok := new(big.Int).SetString(strings.TrimSpace(b), 10)
	if !ok {
		return "", fmt.Errorf("invalid int string b=%s", b)
	}
	ai.Add(ai, bi)
	if ai.Sign() < 0 {
		// 余额不应该为负（除非你允许透支）。这里直接保护。
		return "", fmt.Errorf("balance became negative: %s + %s = %s", a, b, ai.String())
	}
	return ai.String(), nil
}

func negate(pos string) string {
	if strings.HasPrefix(pos, "-") {
		return strings.TrimPrefix(pos, "-")
	}
	return "-" + pos
}
