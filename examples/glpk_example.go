// This code is part of glpk package (Go bindings for the GNU Linear Programming Kit).
//
// Copyright (C) 2014 Łukasz Pankowski <lukpank@o2.pl>
//
// Pacakge glpk is free software: you can redistribute it and/or
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

// This example is a Go rewrite of the PyGLPK example from
// http://tfinley.net/software/pyglpk/discussion.html (Which is a
// Python reimplementation of a C program from GLPK documentation)

package main

import "fmt"
import "github.com/lukpank/go-glpk/glpk"

func main() {
	lp := glpk.New()
	lp.SetProbName("sample")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)

	lp.AddRows(3)
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("%c", 'p'+i)
		lp.SetRowName(i+1, name)
	}
	lp.SetRowBnds(1, glpk.UP, 0, 100.0)
	lp.SetRowBnds(2, glpk.UP, 0, 600.0)
	lp.SetRowBnds(3, glpk.UP, 0, 300.0)

	lp.AddCols(3)
	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("x%d", i)
		lp.SetColName(i+1, name)
		lp.SetColBnds(i+1, glpk.LO, 0.0, 0.0)
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

	lp.Simplex(nil)

	fmt.Printf("%s = %g", lp.ObjName(), lp.ObjVal())
	for i := 0; i < 3; i++ {
		fmt.Printf("; %s = %g", lp.ColName(i+1), lp.ColPrim(i+1))
	}
	fmt.Println()

	lp.Delete()
}
