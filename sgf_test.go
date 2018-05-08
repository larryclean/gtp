package gtp

import (
	"testing"
	"github.com/larry-dev/gtp/sgf"
)

func TestSgf(t *testing.T) {
	kifu := sgf.ParseSgf("(;SZ[7])")
	kifu.Last()
	kifu.CurPos.PrintBoard()
	kifu.Play(sgf.Node{
		X:3,
		Y:3,
		C:1,
	})
	kifu.Play(sgf.Node{
		X:4,
		Y:2,
		C:-1,
	})
	kifu.Play(sgf.Node{
		X:3,
		Y:2,
		C:1,
	})
	kifu.Play(sgf.Node{
		X:4,
		Y:3,
		C:-1,
	})
	kifu.Play(sgf.Node{
		X:3,
		Y:4,
		C:1,
	})
	kifu.Play(sgf.Node{
		X:4,
		Y:4,
		C:-1,
	})
	kifu.Play(sgf.Node{
		X:4,
		Y:5,
		C:1,
	})
	kifu.Play(sgf.Node{
		X:5,
		Y:5,
		C:-1,
	})

	kifu.CurPos.PrintBoard()
	totals:=kifu.ToCurSgf()
	//fmt.Printf()
	t.Logf("%v",totals)

	s:=kifu.ToSgf()
	t.Logf(s)
}
