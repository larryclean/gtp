package gtp

import (
	"testing"
	"github.com/larry-dev/gtp/sgf"
)

func TestSgf(t *testing.T) {
	kifu := sgf.ParseSgf("(;KM[7.5]HA[0]SZ[19]AP[WGo.js:2]FF[4]GM[1]CA[UTF-8];B[cp];W[pp];B[ce];W[pd];B[eq];W[ic];B[qf];W[nd];B[qi];W[qe];B[cg];W[rf];B[cm];W[qm];B[jp];W[oi];B[pk];W[ol];B[ok];W[qk];B[pi];W[nk];B[nj];W[oj];B[pj];W[nl];B[qg];W[oh];B[og];W[ng];B[nf];W[of];B[oe];W[pg];B[pf];W[pe];B[og];W[ph];B[jc];W[id];B[gb];W[hb];B[ha];W[ia];B[ga];W[jb];B[eb];W[dd];B[db];W[fc];B[bb];W[fb];B[fa];W[ea];B[da];W[gc];B[pl];W[pm];B[ql];W[rl];B[rk];W[rj];B[sk];W[qh];B[rm];W[rn];B[sl];W[ri];B[qn];W[pn];B[qo];W[ro];B[qp];W[qq];B[rp];W[rq];B[so];W[po];B[sn];W[lq];B[kq];W[lp];B[kr];W[hp];B[hr];W[cq];B[dq];W[bp];B[cr];W[co];B[bq];W[jo];B[nh];W[of];B[al];W[mh];B[ii];W[hn];B[di];W[ei];B[gi];W[dh];B[ch];W[dj];B[ci];W[hh];B[hi];W[jg];B[jh];W[kh];B[hm];W[in];B[gn];W[go];B[fo];W[fp];B[eo];W[ep];B[gq];W[gp];B[dp];W[gm];B[fn];W[fm];B[dn];W[hl];B[jn];W[ko];B[im];W[jm];B[il];W[ik];B[jl];W[jk];B[km];W[kl];B[lm];W[jm];B[io];W[kp];B[mq];W[lr];B[pq];W[oq];B[op];W[pr];B[mn];W[ll];B[lo];W[mo];B[mp];W[no];B[ip];W[iq];B[jq];W[ir];B[hq];W[ho];B[im];W[do];B[en];W[hs];B[fq];W[kn];B[fk];W[gs];B[fj];W[el];B[dk];W[ek];B[ej];W[dl];B[ck];W[qj];B[rl];W[ji];B[rh];W[rg];B[jj];W[ki];B[ij];W[kj];B[hj];W[ih];B[gj];W[fg];B[ff];W[ef];B[ie];W[de];B[kd];W[cd];B[bd];W[be];B[ec];W[ed];B[qc];W[pb];B[gg];W[gh];B[kg];W[mf];B[lg];W[mg];B[fh];W[eg];B[li];W[lh];B[mj];W[cl];B[fl];W[em];B[cf];W[bc];B[ee];W[fe];B[ld];W[gf];B[pc];W[ob];B[oc];W[nc];B[lc];W[qb];B[kc];W[ke];B[hc];W[ea];B[ac];W[ad];B[mc];W[le];B[nb];W[mb];B[gd];W[hd];B[he];W[bk];B[dg];W[eh];B[bh];W[cj];B[ai];W[bj];B[bi];W[aj];B[ni];W[cn];B[gk];W[dm];B[ah];W[bo];B[sh];W[sg];B[ig];W[hg];B[hf];W[sq];B[je];W[ds];B[rd];W[aq];B[me];W[br];B[kk];W[lk];B[hm];W[gn];B[cq];W[cs];B[er];W[es];B[gr];W[fs];B[fr];W[dr];B[hq];W[si];B[kf];W[lf];B[od];W[rb];B[md];W[na];B[jd];W[jf];B[ae];W[ab];B[ge];W[fd];B[lb];W[ma];B[qa];W[ra];B[dc];W[fa];B[oa];W[pa];B[kb];W[la];B[ka];W[ne];B[sf];W[se];B[bf];W[af];B[ag];W[ae];B[ak];W[il];B[ml];W[mk];B[nm];W[ln];B[om];W[on];B[fo];W[gl];B[fn];W[hk];B[eo];W[fi];B[bn];W[bm];B[an];W[ao])")
	kifu.Last()
	kifu.CurPos.PrintBoard()
	kifu.Play(sgf.Node{
		X:2,
		Y:12,
		C:1,
	})


	kifu.CurPos.PrintBoard()
	totals:=kifu.ToCurSgf()
	//fmt.Printf()
	t.Logf("%v",totals)

	s:=kifu.ToSgf()
	t.Logf(s)
}
