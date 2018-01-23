package gtp

import (
	"testing"
	"github.com/larry-dev/gtp/sgf"
)

func TestSgf(t *testing.T) {
	kifu := sgf.ParseSgf("(;SZ[19])")
	t.Logf(kifu.SgfStr)
}
