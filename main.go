package gtp

import (
	"github.com/larry-dev/gtp/sgf"
)

func GetLastKifu(s string) sgf.Kifu {
	kifu := sgf.ParseSgf(s)
	kifu.Last()
	return kifu
}
