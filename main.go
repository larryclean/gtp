package gtp

import (
	"github.com/larry-dev/gtp/sgf"
)

func GetKifu(s string)sgf.Kifu  {
	return sgf.ParseSgf("(;SZ[19])")
}