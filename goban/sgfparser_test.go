package goban

import (
	"testing"
	"github.com/magiconair/properties/assert"
)

func TestParseSgf(t *testing.T)  {
	//正常解析测试
	sgf:="(;SZ[19]KM[6.5]HA[1]AB[ab]AW[bb]C[123132]GM[1];B[ba])"
	kifu:=ParseSgf(sgf)
	kifu.Last()
 	assert.Equal(t,kifu.Size,int32(19))
	assert.Equal(t,kifu.Komi,float32(6.5))
	assert.Equal(t,kifu.Handicap,1)
	assert.Equal(t,len(kifu.Root.Steup),2)
	assert.Equal(t,kifu.NodeCount,1)
	assert.Equal(t,kifu.Play(0,0,-1),false)
	assert.Equal(t,kifu.ToSgf(),"(;SZ[19]KM[6.5]HA[1]AB[ab]AW[bb]C[123132]GM[1];B[ba])")

	sgf="(;SZ[19]AP[WGo.js:2]FF[4]GM[1]CA[UTF-8];B[ab];W[cc];B[ba])"
	kifu=ParseSgf(sgf)
	kifu.Last()
	assert.Equal(t,kifu.CurPos.ShowBoard(),`  .  X  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  X  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  O  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .
  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .  .`)

	sgf="(;SZ[19]AB[ii]AW[jj];B[aa];W[ab];B[ac];W[ba];B[bb];W[ca];B[aa])"
	kifu=ParseSgf(sgf)
	kifu.Last()
	assert.Equal(t,kifu.NodeCount,7)
	assert.Equal(t,kifu.Play(0,1,-1),false)
	assert.Equal(t,kifu.Play(0,0,-1),false)
	assert.Equal(t,kifu.Play(-1,-1,-1),true)
	assert.Equal(t,kifu.Play(10,10,1),true)
	assert.Equal(t,kifu.Play(0,1,-1),true)
	assert.Equal(t,kifu.NodeCount,10)
	assert.Equal(t,kifu.ToSgf(),"(;SZ[19]KM[7.5]HA[0]AB[ii]AW[jj];B[aa];W[ab];B[ac];W[ba];B[bb];W[ca];B[aa];W[tt];B[kk];W[ab])")
	kifu.GoTo(0)

	assert.Equal(t,kifu.Play(18,18,1),true)
	assert.Equal(t,kifu.Play(15,15,-1),true)
	assert.Equal(t,kifu.NodeCount,12)
	assert.Equal(t,kifu.ToSgf(),"(;SZ[19]KM[7.5]HA[0]AB[ii]AW[jj](;B[aa];W[ab];B[ac];W[ba];B[bb];W[ca];B[aa];W[tt];B[kk];W[ab])(;B[ss];W[pp]))")
	assert.Equal(t,kifu.ToCurSgf(),"(;SZ[19]KM[7.5]HA[0]AB[ii]AW[jj];B[ss];W[pp])")
	kifu.GoTo(0)
	kifu.CurNode.LastSelect=2
	kifu.GoTo(1)
	assert.Equal(t,kifu.ToCurSgf(),"(;SZ[19]KM[7.5]HA[0]AB[ii]AW[jj];B[aa])")

	sgf="(;SZ[a9]KM[7aa]HA[aa](;B[pd];W[dd];B[qp];W[dq];B[oq];W[qf];B[nc];W[rd];B[fc];W[cf];B[qh];W[qc];B[qe];W[re];B[pf];W[pg];B[qg];W[rf];B[og];W[hc];B[pj];W[jq];B[cn];W[cp];B[dj];W[en];B[cl];W[ci];B[cj];W[mq];B[db];W[cc];B[ed];W[kc];B[de];W[ce];B[dc];W[cd];B[he];W[ke];B[ib];W[ic];B[hb];W[jb];B[gc];W[ia];B[ga];W[ie];B[no];W[qo];B[rp];W[ol];B[ml];W[ql];B[qk];W[oj];B[pl];W[pm];B[pk];W[ok];B[om];W[pi];B[qi];W[rk];B[qj];W[rm];B[pn];W[qm];B[ro];W[rj];B[ri];W[sj];B[sl];W[si];B[sh];W[sg];B[ph];W[rh];B[hq];W[io];B[ep];W[eq];B[fq];W[go];B[fo];W[gp];B[fp];W[fn];B[do];W[fr];B[gr];W[er];B[gq];W[ip];B[bi];W[cb];B[bg];W[hf];B[ge];W[mj];B[nm];W[nh];B[mf];W[lg];B[mg];W[mh];B[lf];W[kg];B[pc];W[qb];B[bo];W[bp];B[kr];W[kq];B[lr];W[mr];B[jr];W[lq];B[gn];W[gm];B[hn];W[il];B[fm];W[hm];B[em];W[in];B[dn];W[ho];B[lk];W[kj];B[hh];W[gf];B[ff];W[fg];B[ef];W[gh];B[gi];W[fi];B[gj];W[fj];B[gk];W[fk];B[jk];W[kk];B[jl];W[kl];B[jm];W[km];B[jj];W[ig];B[ki];W[lj];B[jn];W[kn];B[jo];W[ko];B[eh];W[fh];B[ae];W[ad];B[af];W[be];B[ab];W[bf];B[ag];W[cg];B[ch];W[dh];B[di];W[eg];B[dg];W[df];B[dh];W[bc];B[ac];W[bb];B[bd];W[ee];B[fe];W[ad];B[kf];W[jf];B[bd];W[eb];B[ec];W[ad];B[jh];W[jg];B[bd];W[md];B[nd];W[ad];B[hg];W[gg];B[bd];W[oe];B[ne];W[ad];B[fl];W[gl];B[bd];W[ng];B[nf];W[ad];B[ap];W[aq];B[bd];W[oh];B[pb];W[ad];B[ao];W[br];B[bd];W[mb];B[nb];W[ad];B[bq];W[cq];B[bd];W[rg];B[oi];W[ad];B[ni];W[mi];B[bd];W[of];B[pg];W[ad];B[ar];W[as];B[bd];W[qd];B[pe];W[ad];B[iq];W[jp];B[bd];W[oa];B[pa];W[ad];B[qa];W[aa];B[rb];W[sm];B[nr];W[ms];B[nq];W[qn];B[po];W[ek];B[el];W[hd];B[gd];W[da];B[ea];W[fa];B[fb];W[mc];B[ma];W[la];B[na];W[dk];B[ck];W[dl];B[dm];W[ls];B[ir];W[mn];B[mo];W[mm];B[nl];W[mk];B[ll];W[nn];B[on];W[mp];B[np];W[lo];B[lm];W[ln];B[le];W[ld];B[rc];W[sc];B[ns];W[sb];B[ra];W[co];B[bn];W[ks];B[js];W[gs];B[hs];W[fs];B[sn];W[rl];B[rn];W[ca];B[ea];W[ja];B[tt])(;B[aa]))"
	kifu=ParseSgf(sgf)
	kifu.Last()
	assert.Equal(t,kifu.ToSgf(),"(;SZ[19]KM[7.5]HA[0](;B[pd];W[dd];B[qp];W[dq];B[oq];W[qf];B[nc];W[rd];B[fc];W[cf];B[qh];W[qc];B[qe];W[re];B[pf];W[pg];B[qg];W[rf];B[og];W[hc];B[pj];W[jq];B[cn];W[cp];B[dj];W[en];B[cl];W[ci];B[cj];W[mq];B[db];W[cc];B[ed];W[kc];B[de];W[ce];B[dc];W[cd];B[he];W[ke];B[ib];W[ic];B[hb];W[jb];B[gc];W[ia];B[ga];W[ie];B[no];W[qo];B[rp];W[ol];B[ml];W[ql];B[qk];W[oj];B[pl];W[pm];B[pk];W[ok];B[om];W[pi];B[qi];W[rk];B[qj];W[rm];B[pn];W[qm];B[ro];W[rj];B[ri];W[sj];B[sl];W[si];B[sh];W[sg];B[ph];W[rh];B[hq];W[io];B[ep];W[eq];B[fq];W[go];B[fo];W[gp];B[fp];W[fn];B[do];W[fr];B[gr];W[er];B[gq];W[ip];B[bi];W[cb];B[bg];W[hf];B[ge];W[mj];B[nm];W[nh];B[mf];W[lg];B[mg];W[mh];B[lf];W[kg];B[pc];W[qb];B[bo];W[bp];B[kr];W[kq];B[lr];W[mr];B[jr];W[lq];B[gn];W[gm];B[hn];W[il];B[fm];W[hm];B[em];W[in];B[dn];W[ho];B[lk];W[kj];B[hh];W[gf];B[ff];W[fg];B[ef];W[gh];B[gi];W[fi];B[gj];W[fj];B[gk];W[fk];B[jk];W[kk];B[jl];W[kl];B[jm];W[km];B[jj];W[ig];B[ki];W[lj];B[jn];W[kn];B[jo];W[ko];B[eh];W[fh];B[ae];W[ad];B[af];W[be];B[ab];W[bf];B[ag];W[cg];B[ch];W[dh];B[di];W[eg];B[dg];W[df];B[dh];W[bc];B[ac];W[bb];B[bd];W[ee];B[fe];W[ad];B[kf];W[jf];B[bd];W[eb];B[ec];W[ad];B[jh];W[jg];B[bd];W[md];B[nd];W[ad];B[hg];W[gg];B[bd];W[oe];B[ne];W[ad];B[fl];W[gl];B[bd];W[ng];B[nf];W[ad];B[ap];W[aq];B[bd];W[oh];B[pb];W[ad];B[ao];W[br];B[bd];W[mb];B[nb];W[ad];B[bq];W[cq];B[bd];W[rg];B[oi];W[ad];B[ni];W[mi];B[bd];W[of];B[pg];W[ad];B[ar];W[as];B[bd];W[qd];B[pe];W[ad];B[iq];W[jp];B[bd];W[oa];B[pa];W[ad];B[qa];W[aa];B[rb];W[sm];B[nr];W[ms];B[nq];W[qn];B[po];W[ek];B[el];W[hd];B[gd];W[da];B[ea];W[fa];B[fb];W[mc];B[ma];W[la];B[na];W[dk];B[ck];W[dl];B[dm];W[ls];B[ir];W[mn];B[mo];W[mm];B[nl];W[mk];B[ll];W[nn];B[on];W[mp];B[np];W[lo];B[lm];W[ln];B[le];W[ld];B[rc];W[sc];B[ns];W[sb];B[ra];W[co];B[bn];W[ks];B[js];W[gs];B[hs];W[fs];B[sn];W[rl];B[rn];W[ca];B[ea];W[ja];B[tt])(;B[aa]))")
	assert.Equal(t,kifu.ToCurSgf(),"(;SZ[19]KM[7.5]HA[0];B[pd];W[dd];B[qp];W[dq];B[oq];W[qf];B[nc];W[rd];B[fc];W[cf];B[qh];W[qc];B[qe];W[re];B[pf];W[pg];B[qg];W[rf];B[og];W[hc];B[pj];W[jq];B[cn];W[cp];B[dj];W[en];B[cl];W[ci];B[cj];W[mq];B[db];W[cc];B[ed];W[kc];B[de];W[ce];B[dc];W[cd];B[he];W[ke];B[ib];W[ic];B[hb];W[jb];B[gc];W[ia];B[ga];W[ie];B[no];W[qo];B[rp];W[ol];B[ml];W[ql];B[qk];W[oj];B[pl];W[pm];B[pk];W[ok];B[om];W[pi];B[qi];W[rk];B[qj];W[rm];B[pn];W[qm];B[ro];W[rj];B[ri];W[sj];B[sl];W[si];B[sh];W[sg];B[ph];W[rh];B[hq];W[io];B[ep];W[eq];B[fq];W[go];B[fo];W[gp];B[fp];W[fn];B[do];W[fr];B[gr];W[er];B[gq];W[ip];B[bi];W[cb];B[bg];W[hf];B[ge];W[mj];B[nm];W[nh];B[mf];W[lg];B[mg];W[mh];B[lf];W[kg];B[pc];W[qb];B[bo];W[bp];B[kr];W[kq];B[lr];W[mr];B[jr];W[lq];B[gn];W[gm];B[hn];W[il];B[fm];W[hm];B[em];W[in];B[dn];W[ho];B[lk];W[kj];B[hh];W[gf];B[ff];W[fg];B[ef];W[gh];B[gi];W[fi];B[gj];W[fj];B[gk];W[fk];B[jk];W[kk];B[jl];W[kl];B[jm];W[km];B[jj];W[ig];B[ki];W[lj];B[jn];W[kn];B[jo];W[ko];B[eh];W[fh];B[ae];W[ad];B[af];W[be];B[ab];W[bf];B[ag];W[cg];B[ch];W[dh];B[di];W[eg];B[dg];W[df];B[dh];W[bc];B[ac];W[bb];B[bd];W[ee];B[fe];W[ad];B[kf];W[jf];B[bd];W[eb];B[ec];W[ad];B[jh];W[jg];B[bd];W[md];B[nd];W[ad];B[hg];W[gg];B[bd];W[oe];B[ne];W[ad];B[fl];W[gl];B[bd];W[ng];B[nf];W[ad];B[ap];W[aq];B[bd];W[oh];B[pb];W[ad];B[ao];W[br];B[bd];W[mb];B[nb];W[ad];B[bq];W[cq];B[bd];W[rg];B[oi];W[ad];B[ni];W[mi];B[bd];W[of];B[pg];W[ad];B[ar];W[as];B[bd];W[qd];B[pe];W[ad];B[iq];W[jp];B[bd];W[oa];B[pa];W[ad];B[qa];W[aa];B[rb];W[sm];B[nr];W[ms];B[nq];W[qn];B[po];W[ek];B[el];W[hd];B[gd];W[da];B[ea];W[fa];B[fb];W[mc];B[ma];W[la];B[na];W[dk];B[ck];W[dl];B[dm];W[ls];B[ir];W[mn];B[mo];W[mm];B[nl];W[mk];B[ll];W[nn];B[on];W[mp];B[np];W[lo];B[lm];W[ln];B[le];W[ld];B[rc];W[sc];B[ns];W[sb];B[ra];W[co];B[bn];W[ks];B[js];W[gs];B[hs];W[fs];B[sn];W[rl];B[rn];W[ca];B[ea];W[ja];B[tt])")
	assert.Equal(t,CoorToBoardNode(8,0,19),"J19")
	log:=`Thinking at most 36.0 seconds...
NN eval=0.468024

 Q17 ->       2 (V: 50.69%) (N:  9.74%) PV: Q17 R4
 C16 ->       0 (V:  0.00%) (N:  9.66%) PV: C16 
2.0 average depth, 3 max depth
2 non leaf nodes, 1.00 average children
3 visits, 1083 nodes, 2 playouts, 13 n/s
`
	result:=ParseBranch(log)
	assert.Equal(t,len(result),2)
	assert.Equal(t,result[0]["times"],"2")
	log=`  0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0 999   0   0   0   0
 0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0   0
pass: 0
winrate: 1.000000
`
	heatMap,wineRate:=ParseHeatMap(log)
	assert.Equal(t,len(heatMap),362)
	assert.Equal(t,heatMap[283],999.0)
	assert.Equal(t,wineRate,1.0)
}