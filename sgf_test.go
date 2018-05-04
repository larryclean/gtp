package gtp

import (
	"testing"
	"github.com/larry-dev/gtp/sgf"
)

func TestSgf(t *testing.T) {
	kifu := sgf.ParseSgf("(;KM[7.5]HA[0]AP[弈客围棋]FF[4]GM[1]CA[UTF-8]SZ[19]AW[��];B[pj];B[rj];B[ba];B[da];B[bd];W[bf];B[bh];W[bj];B[bk];W[cm];B[el];W[fk];B[fj];W[fi];B[ko];W[hp];B[gp];W[jq];B[no];W[nm];B[qq];W[nq];B[mo];W[mp];B[lp];W[lq];B[kq];W[ne];B[oe];W[qs];B[rs];W[lb];B[rf];W[nk]n];W[ro];B[ll];W[ml];B[nl];W[ol];B[om];W[pl];B[pm];W[pn];B[on];W[kn];B[qc];W[oo];B[nn];W[mn];B[np];W[op];B[lo];W[jp];B[jo];W[jn];B[ln];W[lm];B[io];W[in];B[ho];W[go];B[gn];W[fo];B[hn];W[hm];B[gm];W[gl];B[fm];W[em];B[fn];W[en];B[fl];W[ek];B[dl];W[dm];B[dj];W[dk];B[bm];W[cl];B[kp];W[kr];B[iq];W[ip];B[ji];W[ki];B[li];W[mi];B[ni];W[oi];B[oj];W[ih];B[ig];W[if];B[qf];W[qg];B[pf];W[pg];B[fh];W[ge];B[ic];W[md];B[kd];W[je];B[ie];W[nc];B[nb];W[ob];B[oc];W[pd];B[qd];W[nf];B[eq];W[fq];B[gq];W[hj];B[ij];W[kj];B[mj];W[gi];B[ql];W[ng])")
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
