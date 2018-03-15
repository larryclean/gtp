package gtp

import (
	"testing"
	"github.com/larry-dev/gtp/sgf"
)

func TestSgf(t *testing.T) {
	kifu := sgf.ParseSgf("(;SZ[19]AB[cc][dd][ee](;B[aa](;W[jj])(;W[ii]))(;B[bb]))")
	kifu.Last()
	kifu.CurPos.PrintBoard()
	//for _,v:=range result{
	//	println(v)
	//}
	totals:=kifu.GetCleanSgf()
	//fmt.Printf()
	t.Logf("%v",totals)

	s:=kifu.ToSgf()
	t.Logf(s)
}
