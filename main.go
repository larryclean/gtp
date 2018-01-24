package gtp

import (
	"github.com/larry-dev/gtp/sgf"
)

func GetKifu(s string) sgf.Kifu {
	kifu := sgf.ParseSgf("(;SZ[19])")
	kifu.Last()
	return kifu
}
