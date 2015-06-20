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

## Examples

### Example with real-valued variables

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

### Reading the problem instance from a file

The above problem in the CPLEX LP format has the follow form

```
\* Problem: sample *\

Maximize
 Z: + 10 x0 + 6 x1 + 4 x2

Subject To
 p: + x2 + x1 + x0 <= 100
 q: + 5 x2 + 4 x1 + 10 x0 <= 600
 r: + 6 x2 + 2 x1 + 2 x0 <= 300

End
```

let us save it into `sample.lp` file. Then we can do the computation
analogous to the previous example with the following shorter program:

```go
package main

import (
	"fmt"
	"log"

	"github.com/lukpank/go-glpk/glpk"
)

func main() {
	lp := glpk.New()
	defer lp.Delete()
	lp.ReadLP(nil, "sample.lp")

	if err := lp.Simplex(nil); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s = %g", lp.ObjName(), lp.ObjVal())
	for i := 0; i < 3; i++ {
		fmt.Printf("; %s = %g", lp.ColName(i+1), lp.ColPrim(i+1))
	}
	fmt.Println()
}
```

Analogously you can read from a file in MPS or GPLK LP/MIP formats
using `ReadMPS` or `ReadProb` methods.  You can also write the problem
instance in MPS, CPLEX LP or GPLK LP/MIP formats by the corresponding
`WriteMPS`, `WriteLP` and `WriteProb` methods.

The output of this example is:

```
Reading problem data from 'sample.lp'...
3 rows, 3 columns, 9 non-zeros
11 lines were read
GLPK Simplex Optimizer, v4.55
3 rows, 3 columns, 9 non-zeros
*     0: obj =   0.000000000e+00  infeas =  0.000e+00 (0)
*     2: obj =   7.333333333e+02  infeas =  0.000e+00 (0)
OPTIMAL LP SOLUTION FOUND
Z = 733.3333333333333; x0 = 33.333333333333336; x1 = 66.66666666666666; x2 = 0
```

### MIP (Mixed Integer Programming) example

This example is a Go rewrite of the glpk MIP (Mixed Integer
Programming) example written by Masahiro Sakai. See
[glpk-mip-sample.c](https://gist.github.com/msakai/2450935).

```go
package main

import (
	"fmt"
	"log"

	"github.com/lukpank/go-glpk/glpk"
)

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
func main() {
	lp := glpk.New()
	lp.SetProbName("sample")
	lp.SetObjName("Z")
	lp.SetObjDir(glpk.MAX)

	lp.AddRows(3)
	lp.SetRowName(1, "c1")
	lp.SetRowBnds(1, glpk.UP, 0.0, 20.0)
	lp.SetRowName(2, "c2")
	lp.SetRowBnds(2, glpk.UP, 0.0, 30.0)
	lp.SetRowName(3, "c3")
	lp.SetRowBnds(3, glpk.FX, 0.0, 0)

	lp.AddCols(4)
	lp.SetColName(1, "x1")
	lp.SetColBnds(1, glpk.DB, 0.0, 40.0)
	lp.SetObjCoef(1, 1.0)
	lp.SetColName(2, "x2")
	lp.SetColBnds(2, glpk.LO, 0.0, 0.0)
	lp.SetObjCoef(2, 2.0)
	lp.SetColName(3, "x3")
	lp.SetColBnds(3, glpk.LO, 0.0, 0.0)
	lp.SetObjCoef(3, 3.0)
	lp.SetColName(4, "x4")
	lp.SetColBnds(4, glpk.DB, 2.0, 3.0)
	lp.SetObjCoef(4, 1.0)
	lp.SetColKind(4, glpk.IV)

	fmt.Printf("col1: %v\n", lp.ColKind(1) == glpk.CV)

	ind := []int32{0, 1, 2, 3, 4}
	mat := [][]float64{
		{0, -1, 1.0, 1.0, 10},
		{0, 1.0, -3.0, 1.0, 0.0},
		{0, 0.0, 1.0, 0.0, -3.5}}
	for i := 0; i < 3; i++ {
		lp.SetMatRow(i+1, ind, mat[i])
	}

	iocp := glpk.NewIocp()
	iocp.SetPresolve(true)

	if err := lp.Intopt(iocp); err != nil {
		log.Fatalf("Mip error: %v", err)
	}

	fmt.Printf("%s = %g", lp.ObjName(), lp.MipObjVal())
	for i := 0; i < 4; i++ {
		fmt.Printf("; %s = %g", lp.ColName(i+1), lp.MipColVal(i+1))
	}
	fmt.Println()

	lp.Delete()
}
```

The output of this example is:

```
GLPK Integer Optimizer, v4.55
3 rows, 4 columns, 9 non-zeros
1 integer variable, none of which are binary
Preprocessing...
3 rows, 4 columns, 9 non-zeros
1 integer variable, none of which are binary
Scaling...
 A: min|aij| =  1.000e+00  max|aij| =  1.000e+01  ratio =  1.000e+01
Problem data seem to be well scaled
Constructing initial basis...
Size of triangular part is 3
Solving LP relaxation...
GLPK Simplex Optimizer, v4.55
3 rows, 4 columns, 9 non-zeros
      0: obj =   2.300000000e+01  infeas =  1.400e+01 (0)
*     1: obj =   3.700000000e+01  infeas =  0.000e+00 (0)
*     5: obj =   1.252083333e+02  infeas =  0.000e+00 (0)
OPTIMAL LP SOLUTION FOUND
Integer optimization begins...
+     5: mip =     not found yet <=              +inf        (1; 0)
+     6: >>>>>   1.225000000e+02 <=   1.225000000e+02 < 0.1% (2; 0)
+     6: mip =   1.225000000e+02 <=     tree is empty   0.0% (0; 3)
INTEGER OPTIMAL SOLUTION FOUND
Z = 122.5; x1 = 40; x2 = 10.5; x3 = 19.5; x4 = 3
```

## About

Package glpk was written by Łukasz Pankowski (username at o2 dot pl;
where username is lukpank).
