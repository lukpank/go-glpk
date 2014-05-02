# glpk

Package glpk contains Go bindings for GLPK (GNU Linear Programming Kit).

The binding is not complete but enough for my purposes. Fill free to
contact me (email at the end) if there is some part of GLPK that you
would like to use and it is not yet covered by the glpk package.

Package glpk is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or (at
your option) any later version.

## Instalation

First install [GLPK](http://www.gnu.org/software/glpk/) (GNU Linear
Programming Kit). Package glpk is known to work with GLPK v4.54
(available in Debian Sid) and v4.45 (available in Debian Wheezy).  On
Debian GLPK can be installed by installing the package `libglpk-dev`.

To install glpk package run

    go get github.com/lukpank/go-glpk/glpk

## Documentation

Documentation for glpk package is
[available on godoc.org](http://godoc.org/github.com/lukpank/go-glpk/glpk).

## Example

This example is a Go rewrite of the PyGLPK example from
http://tfinley.net/software/pyglpk/discussion.html (Which is a Python
reimplementation of a C program from GLPK documentation).

```go
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
```

The output of this example is:

    GLPK Simplex Optimizer, v4.54
    3 rows, 3 columns, 9 non-zeros
    *     0: obj =   0.000000000e+00  infeas =  0.000e+00 (0)
    *     2: obj =   7.333333333e+02  infeas =  0.000e+00 (0)
    OPTIMAL LP SOLUTION FOUND
    Z = 733.3333333333333; x0 = 33.333333333333336; x1 = 66.66666666666666; x2 = 0

## About

Package glpk was written by Łukasz Pankowski (username at o2 dot pl;
where username is lukpank).
