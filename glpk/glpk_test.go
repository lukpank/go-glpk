// This code is part of glpk package (Go bindings for the GNU Linear Programming Kit).
//
// Copyright (C) 2014 Łukasz Pankowski <lukpank@o2.pl>
//
// Package glpk is free software: you can redistribute it and/or
// modify it under the terms of the GNU General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// Package glpk is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with glpk package. If not, see <http://www.gnu.org/licenses/>.

package glpk

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"testing"
)

func TestNewDelete(t *testing.T) {
	lp := New()
	lp.Erase()
	lp.Delete()
	lp.Delete() // second delete has no effect
}

func TestSetGetProbName(t *testing.T) {
	lp := New()
	name := "problem"
	lp.SetProbName(name)
	got := lp.ProbName()
	if got != name {
		t.Errorf("Got name %#v but %#v was set", got, name)
	}
	lp.Delete()
}

func TestSetGetObjName(t *testing.T) {
	lp := New()
	name := "objective"
	lp.SetObjName(name)
	got := lp.ObjName()
	if got != name {
		t.Errorf("Got name %#v but %#v was set", got, name)
	}
	lp.Delete()
}

func TestSetGetObjDir(t *testing.T) {
	lp := New()
	lp.SetObjDir(MAX)
	if lp.ObjDir() != MAX {
		t.Errorf("Got %d instead of %d (MAX)", lp.ObjDir(), MAX)
	}
	lp.SetObjDir(MIN)
	if lp.ObjDir() != MIN {
		t.Errorf("Got %d instead of %d (MIN)", lp.ObjDir(), MIN)
	}
	lp.Delete()
}

func TestSetGetNumRows(t *testing.T) {
	lp := New()
	lp.AddRows(11)
	if n := lp.NumRows(); n != 11 {
		t.Errorf("Got %d rows expected 11", n)
	}
	lp.Delete()
}

func TestSetGetNumCols(t *testing.T) {
	lp := New()
	lp.AddCols(11)
	if n := lp.NumCols(); n != 11 {
		t.Errorf("Got %d columns expected 11", n)
	}
	lp.Delete()
}

func TestSetGetRowName(t *testing.T) {
	lp := New()
	lp.AddRows(1)
	name := "a constraint"
	lp.SetRowName(1, name)
	got := lp.RowName(1)
	if got != name {
		t.Errorf("Got name %#v but %#v was set", got, name)
	}
	lp.Delete()
}

func TestSetGetColName(t *testing.T) {
	lp := New()
	lp.AddCols(1)
	name := "a variable"
	lp.SetColName(1, name)
	got := lp.ColName(1)
	if got != name {
		t.Errorf("Got name %#v but %#v was set", got, name)
	}
	lp.Delete()
}

var bndsExpected = []struct {
	typ    BndsType
	lb, ub float64
}{
	{FR, -math.MaxFloat64, math.MaxFloat64},
	{LO, 3.2, math.MaxFloat64},
	{UP, -math.MaxFloat64, 7.5},
	{DB, 3.2, 7.5},
	{FX, 3.2, 3.2},
}

func TestSetGetRowBnds(t *testing.T) {
	lp := New()
	lp.AddRows(1)
	for _, expected := range bndsExpected {
		lp.SetRowBnds(1, expected.typ, 3.2, 7.5)
		typ := lp.RowType(1)
		if typ != expected.typ {
			t.Errorf("Got type %d but %d was set", typ, expected.typ)
		}
		lb := lp.RowLB(1)
		if lb != expected.lb {
			t.Errorf("Got lower bound %g but %g was expected", lb, expected.lb)
		}
		ub := lp.RowUB(1)
		if ub != expected.ub {
			t.Errorf("Got upper bound %g but %g was expected", ub, expected.ub)
		}
	}
	lp.Delete()
}

func TestSetGetColBnds(t *testing.T) {
	lp := New()
	lp.AddCols(1)
	for _, expected := range bndsExpected {
		lp.SetColBnds(1, expected.typ, 3.2, 7.5)
		typ := lp.ColType(1)
		if typ != expected.typ {
			t.Errorf("Got type %s but %s was set", typ, expected.typ)
		}
		lb := lp.ColLB(1)
		if lb != expected.lb {
			t.Errorf("Got lower bound %g but %g was expected", lb, expected.lb)
		}
		ub := lp.ColUB(1)
		if ub != expected.ub {
			t.Errorf("Got upper bound %g but %g was expected", ub, expected.ub)
		}
	}
	lp.Delete()
}

func TestSetGetRowStat(t *testing.T) {
	lp := New()
	lp.AddRows(1)
	// such values was selected for which get returns what was set
	for _, stat := range []VarStat{BS, NF} {
		lp.SetRowStat(1, stat)
		got := lp.RowStat(1)
		if got != stat {
			t.Errorf("Got stat %d but %d was set", got, stat)
		}
	}
	lp.Delete()
}

func TestSetGetColStat(t *testing.T) {
	lp := New()
	lp.AddCols(1)
	// such values was selected for which get returns what was set
	for _, stat := range []VarStat{BS, NS} {
		lp.SetColStat(1, stat)
		got := lp.ColStat(1)
		if got != stat {
			t.Errorf("Got stat %d but %d was set", got, stat)
		}
	}
	lp.Delete()
}

func CmpIndicesData(ind []int32, data []float64, ind2 []int32, data2 []float64) bool {
	if len(ind) != len(data) || len(ind2) != len(data2) || len(ind) != len(ind2) {
		return false
	}
	m := make(map[int32]float64)
	for i := 1; i < len(ind); i++ {
		m[ind[i]] = data[i]
	}
	for i := 1; i < len(ind2); i++ {
		v, ok := m[ind2[i]]
		if !ok {
			return false
		}
		if v != data2[i] {
			return false
		}
		delete(m, ind2[i])
	}
	return true
}

func TestSetGetMatRow(t *testing.T) {
	lp := New()
	lp.AddRows(1)
	lp.AddCols(10)
	ind := []int32{0, 3, 7, 5, 2}
	row := []float64{9.0, 7.5, 11.0, 5.0, 12.0}
	lp.SetMatRow(1, ind, row)
	ind2, row2 := lp.MatRow(1)
	if !CmpIndicesData(ind, row, ind2, row2) {
		t.Errorf("Indices and values (%v, %v) does not match (%v, %v)", ind2, row2, ind, row)
	}
	lp.Delete()
}

func TestSetGetMatCol(t *testing.T) {
	lp := New()
	lp.AddRows(10)
	lp.AddCols(1)
	ind := []int32{0, 3, 7, 5, 2}
	col := []float64{9.0, 7.5, 11.0, 5.0, 12.0}
	lp.SetMatCol(1, ind, col)
	ind2, col2 := lp.MatCol(1)
	if !CmpIndicesData(ind, col, ind2, col2) {
		t.Errorf("Indices and values (%v, %v) does not match (%v, %v)", ind2, col2, ind, col)
	}
	lp.Delete()
}

func TestSetGetMatix(t *testing.T) {
	lp := New()
	lp.AddRows(2)
	lp.AddCols(20)
	ia := []int32{0, 1, 1, 1, 1, 2, 2, 2, 2}
	ja1 := []int32{0, 3, 7, 5, 2}
	ja2 := []int32{0, 11, 3, 7, 15}
	ja := append(ja1, ja2[1:]...)
	ar1 := []float64{9.0, 7.5, 11.0, 5.0, 12.0}
	ar2 := []float64{3.0, 5.5, 1.0, 4.0, 11.0}
	ar := append(ar1, ar2[1:]...)
	lp.LoadMatrix(ia, ja, ar)
	ind1, val1 := lp.MatRow(1)
	if !CmpIndicesData(ja1, ar1, ind1, val1) {
		t.Errorf("Indices and values (%v, %v) does not match (%v, %v)", ind1, val1, ja1, ar1)
	}
	ind2, val2 := lp.MatRow(2)
	if !CmpIndicesData(ja2, ar2, ind2, val2) {
		t.Errorf("Indices and values (%v, %v) does not match (%v, %v)", ind2, val2, ja2, ar2)
	}
	lp.Delete()
}

func TestCopy(t *testing.T) {
	lp := New()
	lp.AddRows(4)
	lp.AddCols(3)
	lp.SetProbName("problem")
	lp2 := lp.Copy(false)
	if n := lp2.NumRows(); n != 4 {
		t.Errorf("Got %d rows expected 4", n)
	}
	if n := lp2.NumCols(); n != 3 {
		t.Errorf("Got %d columns expected 3", n)
	}
	if s := lp2.ProbName(); s != "" {
		t.Errorf("names=false but got name %#v", s)
	}
	lp2.Delete()
	lp3 := lp.Copy(true)
	if n := lp3.NumRows(); n != 4 {
		t.Errorf("Got %d rows expected 4", n)
	}
	if n := lp3.NumCols(); n != 3 {
		t.Errorf("Got %d columns expected 3", n)
	}
	if s := lp3.ProbName(); s != "problem" {
		t.Errorf("names=true but got name %#v instead of \"problem\"", s)
	}
	lp3.Delete()
	lp.Delete()
}

func TestSetGetObjCoef(t *testing.T) {
	lp := New()
	lp.AddCols(1)
	coef := 3.5
	lp.SetObjCoef(1, coef)
	got := lp.ObjCoef(1)
	if got != coef {
		t.Errorf("Got coef %#v but %#v was set", got, coef)
	}
}

func CheckClose(t *testing.T, v1, v2 float64) {
	if math.Abs(v1-v2) > 1e-10 {
		t.Errorf("values %g and %g differ by %g", v1, v2, v1-v2)
	}
}

func CheckSolution(t *testing.T, lp *Prob) {
	if lp.Status() != OPT {
		t.Errorf("expected optimal solution, but got %d", lp.Status())
	}
	if lp.PrimStat() != FEAS {
		t.Errorf("expected optimal solution, but got %d", lp.PrimStat())
	}
	if lp.DualStat() != FEAS {
		t.Errorf("expected optimal solution, but got %d", lp.DualStat())
	}

	CheckClose(t, lp.ObjVal(), 733+1.0/3)
	CheckClose(t, lp.ColPrim(1), 33+1.0/3)
	CheckClose(t, lp.ColPrim(2), 66+2.0/3)
	CheckClose(t, lp.ColPrim(3), 0)
}

// PrepareTestExample is a Go rewrite of the PyGLPK example from
// http://tfinley.net/software/pyglpk/discussion.html (Which is a
// Python reimplementation of a C program from GLPK documentation)
func PrepareTestExample(t *testing.T) *Prob {
	lp := New()
	lp.SetProbName("sample")
	lp.SetObjName("Z")
	lp.SetObjDir(MAX)

	if n := lp.AddRows(3); n != 1 {
		t.Errorf("expected 0 but got %d", n)
	}
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("%c", 'p'+i)
		lp.SetRowName(i+1, name)
	}
	lp.SetRowBnds(1, UP, 0, 100.0)
	lp.SetRowBnds(2, UP, 0, 600.0)
	lp.SetRowBnds(3, UP, 0, 300.0)

	if n := lp.AddCols(3); n != 1 {
		t.Errorf("expected 0 but got %d", n)
	}
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("x%d", i)
		lp.SetColName(i+1, name)
		lp.SetColBnds(i+1, LO, 0.0, 0.0)
	}

	lp.SetObjCoef(1, 10.0)
	lp.SetObjCoef(2, 6.0)
	lp.SetObjCoef(3, 4.0)

	ind := []int32{0, 1, 2, 3}
	mat := [][]float64{
		{0, 1.0, 1.0, 1.0},
		{0, 10.0, 4.0, 5.0},
		{0, 2.0, 2.0, 6.0}}
	for i := 0; i < 3; i++ {
		lp.SetMatRow(i+1, ind, mat[i])
	}
	return lp
}

func CheckSimplexSolution(t *testing.T, lp *Prob) {
	smcp := NewSmcp()
	smcp.SetMsgLev(MSG_ERR)

	if err := lp.Simplex(smcp); err != nil {
		t.Errorf("Simplex error: %v", err)
	}
	CheckSolution(t, lp)
}

func TestExample(t *testing.T) {
	lp := PrepareTestExample(t)

	lp2 := lp.Copy(true)
	CheckSimplexSolution(t, lp)
	lp.Delete()

	smcp := NewSmcp()
	smcp.SetMsgLev(MSG_ERR)

	if err := lp2.Exact(smcp); err != nil {
		t.Errorf("Exact error: %v", err)
	}
	CheckSolution(t, lp2)
	lp2.Delete()
}

func TestReadWriteMPS(t *testing.T) {
	lp := PrepareTestExample(t)
	f1, err := ioutil.TempFile("", "glpk-test-")
	if err != nil {
		t.Fatal(err)
	}
	f1.Close()
	defer os.Remove(f1.Name())
	f2, err := ioutil.TempFile("", "glpk-test-")
	if err != nil {
		t.Fatal(err)
	}
	f2.Close()
	defer os.Remove(f2.Name())
	err = lp.WriteMPS(MPS_DECK, nil, f1.Name())
	if err != nil {
		t.Error(err)
	}
	err = lp.WriteMPS(MPS_FILE, NewMPSCP(), f2.Name())
	if err != nil {
		t.Error(err)
	}
	lp.Delete()

	lp1 := New()
	defer lp1.Delete()
	err = lp1.ReadMPS(MPS_DECK, NewMPSCP(), f1.Name())
	if err != nil {
		t.Error(err)
	} else {
		lp1.SetObjDir(MAX)
		CheckSimplexSolution(t, lp1)
	}

	lp2 := New()
	defer lp2.Delete()
	err = lp2.ReadMPS(MPS_FILE, nil, f2.Name())
	if err != nil {
		t.Error(err)
	} else {
		lp2.SetObjDir(MAX)
		CheckSimplexSolution(t, lp2)
	}
}

func TestReadWriteLP(t *testing.T) {
	CheckReadWriteLP(t, nil)
	CheckReadWriteLP(t, NewCPXCP())
}

func CheckReadWriteLP(t *testing.T, cpxcp *CPXCP) {
	lp := PrepareTestExample(t)
	f, err := ioutil.TempFile("", "glpk-test-")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())
	err = lp.WriteLP(cpxcp, f.Name())
	if err != nil {
		t.Error(err)
	}
	lp.Delete()

	lp1 := New()
	defer lp1.Delete()
	err = lp1.ReadLP(cpxcp, f.Name())
	if err != nil {
		t.Fatal(err)
	}
	CheckSimplexSolution(t, lp1)
}

func TestReadWriteProb(t *testing.T) {
	lp := PrepareTestExample(t)
	f, err := ioutil.TempFile("", "glpk-test-")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	defer os.Remove(f.Name())
	err = lp.WriteProb(0, f.Name())
	if err != nil {
		t.Error(err)
	}
	lp.Delete()

	lp1 := New()
	defer lp1.Delete()
	err = lp1.ReadProb(0, f.Name())
	if err != nil {
		t.Fatal(err)
	}
	CheckSimplexSolution(t, lp1)
}

func TestSetGetColKind(t *testing.T) {
	lp := New()
	lp.AddCols(3)
	for i, v := range []VarType{CV, IV, BV} {
		lp.SetColKind(i+1, v)
		if got := lp.ColKind(i + 1); got != v {
			t.Errorf("expected %v but got %v", v, got)
		}
	}
	lp.Delete()

}

func TestIocp(t *testing.T) {
	iocp := NewIocp()
	for _, v := range []bool{false, true} {
		iocp.SetPresolve(v)
		if got := iocp.Presolve(); got != v {
			t.Errorf("expected %v but got %v", v, got)
		}
	}
}

// TestExample is a Go rewrite of the glpk mip example written
// by Masahiro Sakai. https://gist.github.com/msakai/2450935
// (glpk-mip-sample.c).
func TestIntop(t *testing.T) {

	// Maximize
	//
	//      obj: x1 + 2 x2 + 3 x3 + x4
	//
	// Subject To
	//
	//      c1: 0 <= - x1 + x2 + x3 + 10 x4 <= 20
	//      c2: 0 <= x1 - 3 x2 + x3 <= 30
	//      c3: x2 - 3.5 x4 = 0
	//
	// Bounds
	//
	//      0 <= x1 <= 40
	//      x2 >= 0
	//      x3 >= 0
	//      2 <= x4 <= 3
	//
	// Type
	//
	//      x1, x2, x3 real
	//      x4 integer
	//
	// End

	lp := New()
	lp.SetProbName("sample")
	lp.SetObjName("Z")
	lp.SetObjDir(MAX)

	if n := lp.AddRows(3); n != 1 {
		t.Errorf("expected 0 but got %d", n)
	}
	lp.SetRowName(1, "c1")
	lp.SetRowBnds(1, DB, 0.0, 20.0)
	lp.SetRowName(2, "c2")
	lp.SetRowBnds(2, DB, 0.0, 30.0)
	lp.SetRowName(3, "c3")
	lp.SetRowBnds(3, FX, 0.0, 0)

	if n := lp.AddCols(4); n != 1 {
		t.Errorf("expected 0 but got %d", n)
	}

	lp.SetColName(1, "x1")
	lp.SetColBnds(1, DB, 0.0, 40.0)
	lp.SetObjCoef(1, 1.0)
	lp.SetColName(2, "x2")
	lp.SetColBnds(2, LO, 0.0, 0.0)
	lp.SetObjCoef(2, 2.0)
	lp.SetColName(3, "x3")
	lp.SetColBnds(3, LO, 0.0, 0.0)
	lp.SetObjCoef(3, 3.0)
	lp.SetColName(4, "x4")
	lp.SetColBnds(4, DB, 2.0, 3.0)
	lp.SetObjCoef(4, 1.0)
	lp.SetColKind(4, IV)

	ind := []int32{0, 1, 2, 3, 4}
	mat := [][]float64{
		{0, -1, 1.0, 1.0, 10},
		{0, 1.0, -3.0, 1.0, 0.0},
		{0, 0.0, 1.0, 0.0, -3.5}}
	for i := 0; i < 3; i++ {
		lp.SetMatRow(i+1, ind, mat[i])
	}

	iocp := NewIocp()
	iocp.SetPresolve(true)

	if err := lp.Intopt(iocp); err != nil {
		t.Errorf("Mip error: %v", err)
	}

	CheckMipSolution(t, lp)

	lp.Delete()
}

func CheckMipSolution(t *testing.T, lp *Prob) {
	state := lp.MipStatus()
	if state != OPT && state != FEAS {
		t.Errorf("expected optimal solution, but got %d", lp.MipStatus())
	}

	// z = 122.5; x1 = 40; x2 = 10.5; x3 = 19.5, x4 = 3
	CheckClose(t, lp.MipObjVal(), 122.5)
	CheckClose(t, lp.MipColVal(1), 40)
	CheckClose(t, lp.MipColVal(2), 10.5)
	CheckClose(t, lp.MipColVal(3), 19.5)
	CheckClose(t, lp.MipColVal(4), 3)
}

func TestGarbageCollection(t *testing.T) {
	// this loop should create enough objects to trigger garbage collection
	for i := 0; i < 2000; i++ {
		lp := New()
		_ = lp
		lp2 := New()
		lp2.Delete()
	}
}
