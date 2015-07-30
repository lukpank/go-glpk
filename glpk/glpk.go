// This code is part of glpk package (Go bindings for the GNU Linear Programming Kit).
//
// Copyright (C) 2014 Łukasz Pankowski <lukpank@o2.pl>
//
// Some comments/strings are taken or adapted from GLPK and thus are
// subject to the following copyright:
//
// Copyright (C) 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008,
// 2009, 2010, 2011, 2013, 2014 Andrew Makhorin, Department for Applied
// Informatics, Moscow Aviation Institute, Moscow, Russia. All rights
// reserved. E-mail: <mao@gnu.org>.
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

// Go bindings for GLPK (GNU Linear Programming Kit).
//
// For usage examples see https://github.com/lukpank/go-glpk#examples.
//
// The binding is not complete but enough for my purposes. Fill free
// to contact me if there is some part of GLPK that you would like to
// use and it is not yet covered by the glpk package.
//
// Package glpk is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
package glpk

import (
	"reflect"
	"runtime"
	"unsafe"
)

// #cgo LDFLAGS: -lglpk
// #include <glpk.h>
// #include <stdlib.h>
import "C"

// ObjDir is used to specify objective function direction
// (maximization or minimization).
type ObjDir int

const (
	MAX = ObjDir(C.GLP_MAX) // maximization
	MIN = ObjDir(C.GLP_MIN) // minimization
)

// BndsType is used to specify bounds type of a variable.
type BndsType int

const (
	FR = BndsType(C.GLP_FR) // a free (unbounded) variable
	LO = BndsType(C.GLP_LO) // a lower-bounded variable
	UP = BndsType(C.GLP_UP) // an upper-bounded variable
	DB = BndsType(C.GLP_DB) // a double-bounded variable
	FX = BndsType(C.GLP_FX) // a fixed variable
)

// SolStat specifies solution status.
type SolStat int

const (
	UNDEF  = SolStat(C.GLP_UNDEF)  // solution is undefined
	FEAS   = SolStat(C.GLP_FEAS)   // solution is feasible
	INFEAS = SolStat(C.GLP_INFEAS) // solution is infeasible
	NOFEAS = SolStat(C.GLP_NOFEAS) // there is no feasible solution
	OPT    = SolStat(C.GLP_OPT)    // solution is optimal
	UNBND  = SolStat(C.GLP_UNBND)  // problem has unbounded solution
)

// VarType is used to specify variable type (kind).
type VarType int

const (
	CV = VarType(C.GLP_CV) // Contineous Variable
	IV = VarType(C.GLP_IV) // Integer Variable
	BV = VarType(C.GLP_BV) // Binary Variable. Equivalent to IV with 0<=iv<=1
)

type prob struct {
	p *C.glp_prob
}

// Prob represens optimization problem. Use glpk.New() to create a new problem.
type Prob struct {
	p *prob
}

func finalizeProb(p *prob) {
	if p.p != nil {
		C.glp_delete_prob(p.p)
		p.p = nil
	}
}

// New creates a new optimization problem.
func New() *Prob {
	p := &prob{C.glp_create_prob()}
	runtime.SetFinalizer(p, finalizeProb)
	return &Prob{p}
}

// Delete deletes a problem.  Calling Delete on a deleted problem will
// have no effect (It is save to do so). But calling any other method
// on a deleted problem will panic. The problem will be deleted on
// garbage collection but you can do this as soon as you no longer
// need the optimization problem.
func (p *Prob) Delete() {
	if p.p.p != nil {
		C.glp_delete_prob(p.p.p)
		p.p.p = nil
	}
}

// Erase erases the problem. After erasing the problem is empty as if
// it were created with glpk.New().
func (p *Prob) Erase() {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_erase_prob(p.p.p)
}

// SetProbName sets (changes) the problem name.
func (p *Prob) SetProbName(name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_prob_name(p.p.p, s)
}

// SetObjName sets (changes) objective function name.
func (p *Prob) SetObjName(name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_obj_name(p.p.p, s)
}

// SetObjDir sets optimization direction (either glpk.MAX for
// maximization or glpk.MIN for minimization)
func (p *Prob) SetObjDir(dir ObjDir) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_obj_dir(p.p.p, C.int(dir))
}

// AddRows adds rows (constraints). Returns (1-based) index of the
// first of the added rows.
func (p *Prob) AddRows(nrs int) int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_add_rows(p.p.p, C.int(nrs)))
}

// AddCols adds columns (variables). Returns (1-based) index of the
// first of the added columns.
func (p *Prob) AddCols(nrs int) int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_add_cols(p.p.p, C.int(nrs)))
}

// SetRowName sets i-th row (constraint) name.
func (p *Prob) SetRowName(i int, name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_row_name(p.p.p, C.int(i), s)
}

// SetColName sets j-th column (variable) name.
func (p *Prob) SetColName(j int, name string) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))
	C.glp_set_col_name(p.p.p, C.int(j), s)
}

// SetColKind sets the kind of j-th column
// as specified by the VarType parameter kind.
func (p *Prob) SetColKind(j int, kind VarType) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_col_kind(p.p.p, C.int(j), C.int(kind))
}

// SetRowBnds sets row bounds
func (p *Prob) SetRowBnds(i int, type_ BndsType, lb float64, ub float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_row_bnds(p.p.p, C.int(i), C.int(type_), C.double(lb), C.double(ub))
}

// SetColBnds sets column bounds
func (p *Prob) SetColBnds(j int, type_ BndsType, lb float64, ub float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_col_bnds(p.p.p, C.int(j), C.int(type_), C.double(lb), C.double(ub))
}

// SetObjCoef sets objective function coefficient of j-th column.
func (p *Prob) SetObjCoef(j int, coef float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_obj_coef(p.p.p, C.int(j), C.double(coef))
}

// SetMatRow sets (replaces) i-th row. It sets
//
//     matrix[i, ind[j]] = val[j]
//
// for j=1..len(ind). ind[0] and val[0] are ignored. Requires
// len(ind) = len(val).
func (p *Prob) SetMatRow(i int, ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_set_mat_row(p.p.p, C.int(i), C.int(len(ind)-1), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
}

// SetMatCol sets (replaces) j-th column. It sets
//
//     matrix[ind[i], j] = val[i]
//
// for i=1..len(ind). ind[0] and val[0] are ignored. Requires
// len(ind) = len(val).
func (p *Prob) SetMatCol(j int, ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_set_mat_col(p.p.p, C.int(j), C.int(len(ind)-1), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
}

// LoadMatrix replaces all of the constraint matrix. It sets
//
//     matrix[ia[i], ja[i]] = ar[i]
//
// for i = 1..len(ia). ia[0], ja[0], and ar[0] are ignored. It
// requiers len(ia)=len(ja)=len(ar).
func (p *Prob) LoadMatrix(ia, ja []int32, ar []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ia) != len(ja) || len(ia) != len(ar) {
		panic("len(ia) and len(ja) and len(ar) should be equal")
	}
	ia_ := (*reflect.SliceHeader)(unsafe.Pointer(&ia))
	ja_ := (*reflect.SliceHeader)(unsafe.Pointer(&ja))
	ar_ := (*reflect.SliceHeader)(unsafe.Pointer(&ar))
	C.glp_load_matrix(p.p.p, C.int(len(ia)-1), (*C.int)(unsafe.Pointer(ia_.Data)), (*C.int)(unsafe.Pointer(ja_.Data)), (*C.double)(unsafe.Pointer(ar_.Data)))
}

// TODO:
// glp_check_dup
// glp_del_rows

// Copy returns a copy of the given optimization problem. If name is
// true also symbolic names are copies otherwise their not copied
func (p *Prob) Copy(names bool) *Prob {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	q := &Prob{&prob{C.glp_create_prob()}}
	var names_ C.int
	if names {
		names_ = C.GLP_ON
	} else {
		names_ = C.GLP_OFF
	}
	C.glp_copy_prob(q.p.p, p.p.p, names_)
	return q
}

// ProbName returns problem name.
func (p *Prob) ProbName() string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_prob_name(p.p.p))
}

// ObjName returns objective name.
func (p *Prob) ObjName() string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_obj_name(p.p.p))
}

// ObjDir returns optimization direction (either glpk.MAX or glpk.MIN).
func (p *Prob) ObjDir() ObjDir {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return ObjDir(C.glp_get_obj_dir(p.p.p))
}

// NumRows returns number of rows.
func (p *Prob) NumRows() int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_get_num_rows(p.p.p))
}

// NumCols returns number of columns.
func (p *Prob) NumCols() int {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return int(C.glp_get_num_cols(p.p.p))
}

// RowName returns row (constraint) name of i-th row.
func (p *Prob) RowName(i int) string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_row_name(p.p.p, C.int(i)))
}

// ColName returns column (variable) name of j-th column.
func (p *Prob) ColName(j int) string {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return C.GoString(C.glp_get_col_name(p.p.p, C.int(j)))
}

// ColKind returns the kind of j-th column
func (p *Prob) ColKind(j int) VarType {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return VarType(C.glp_get_col_kind(p.p.p, C.int(j)))
}

// RowType returns the type of i-th row, i.e. the type of the
// corresponding auxiliary variable.
func (p *Prob) RowType(i int) BndsType {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return BndsType(C.glp_get_row_type(p.p.p, C.int(i)))
}

// RowLB returns the lower bound of i-th row, i.e. the lower bound of
// the corresponding auxiliary variable. If the row has no lower bound
// -math.MaxFloat64 is returned.
func (p *Prob) RowLB(i int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_row_lb(p.p.p, C.int(i)))
}

// RowUB returns the upper bound of i-th row, i.e. the upper bound of
// the corresponding auxiliary variable. If the row has no upper bound
// +math.MaxFloat64 is returned.
func (p *Prob) RowUB(i int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_row_ub(p.p.p, C.int(i)))
}

// ColType returns the type of j-th column, i.e. the type of the
// corresponding structural variable.
func (p *Prob) ColType(j int) BndsType {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return BndsType(C.glp_get_col_type(p.p.p, C.int(j)))
}

// ColLB returns the lower bound of j-th column, i.e. the lower bound
// of the corresponding structural variable. I the column has no lower
// bound -math.MaxFloat64 is returned.
func (p *Prob) ColLB(j int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_col_lb(p.p.p, C.int(j)))
}

// ColUB returns the upper bound of j-th column, i.e. the upper bound
// of the corresponding structural variable. I the column has no upper
// bound +math.MaxFloat64 is returned.
func (p *Prob) ColUB(j int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_col_ub(p.p.p, C.int(j)))
}

// ObjCoef returns objective function coefficient of j-th column.
func (p *Prob) ObjCoef(j int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_obj_coef(p.p.p, C.int(j)))
}

// TODO:
// glp_get_num_nz

// MatRow returns nonzero elements of i-th row. ind[1]..ind[n] are
// column numbers of the nonzero elements of the row, val[1]..val[n]
// are their values, and n is the number of nonzero elements in the
// row.
func (p *Prob) MatRow(i int) (ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	length := C.glp_get_mat_row(p.p.p, C.int(i), nil, nil)
	ind = make([]int32, length+1)
	val = make([]float64, length+1)
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_get_mat_row(p.p.p, C.int(i), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
	return
}

// MatCol returns nonzero elements of j-th column. ind[1]..ind[n] are
// row numbers of the nonzero elements of the column, val[1]..val[n]
// are their values, and n is the number of nonzero elements in the
// column.
func (p *Prob) MatCol(j int) (ind []int32, val []float64) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	if len(ind) != len(val) {
		panic("len(ind) and len(val) should be equal")
	}
	length := C.glp_get_mat_col(p.p.p, C.int(j), nil, nil)
	ind = make([]int32, length+1)
	val = make([]float64, length+1)
	ind_ := (*reflect.SliceHeader)(unsafe.Pointer(&ind))
	val_ := (*reflect.SliceHeader)(unsafe.Pointer(&val))
	C.glp_get_mat_col(p.p.p, C.int(j), (*C.int)(unsafe.Pointer(ind_.Data)), (*C.double)(unsafe.Pointer(val_.Data)))
	return
}

// TODO:
// glp_create_index
// glp_find_row
// glp_find_col
// glp_delete_index
// glp_set_rii
// glp_set_sjj
// glp_get_rii
// glp_get_sjj
// glp_scale_prob
// glp_unscale_prob

// VarStat represents status of auxiliary/structural variable.
type VarStat int

const (
	BS = VarStat(C.GLP_BS) // basic variable
	NL = VarStat(C.GLP_NL) // non-basic variable on lower bound
	NU = VarStat(C.GLP_NU) // non-basic variable on upper bound
	NF = VarStat(C.GLP_NF) // non-basic free (unbounded) variable
	NS = VarStat(C.GLP_NS) // non-basic fixed variable
)

// SetRowStat sets the current status of i-th row (auxiliary variable)
// as specified by the stat argument.
func (p *Prob) SetRowStat(i int, stat VarStat) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_row_stat(p.p.p, C.int(i), C.int(stat))
}

// SetColStat sets the current status of j-th column (structural
// variable) as specified by the stat argument.
func (p *Prob) SetColStat(j int, stat VarStat) {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	C.glp_set_col_stat(p.p.p, C.int(j), C.int(stat))
}

// glp_std_basis
// glp_adv_basis
// glp_cpx_basis

// Optimization Error
type OptError int

const (
	EBADB   = OptError(C.GLP_EBADB)   // invalid basis
	ESING   = OptError(C.GLP_ESING)   // singular matrix
	ECOND   = OptError(C.GLP_ECOND)   // ill-conditioned matrix
	EBOUND  = OptError(C.GLP_EBOUND)  // invalid bounds
	EFAIL   = OptError(C.GLP_EFAIL)   // solver failed
	EOBJLL  = OptError(C.GLP_EOBJLL)  // objective lower limit reached
	EOBJUL  = OptError(C.GLP_EOBJUL)  // objective upper limit reached
	EITLIM  = OptError(C.GLP_EITLIM)  // iteration limit exceeded
	ETMLIM  = OptError(C.GLP_ETMLIM)  // time limit exceeded
	ENOPFS  = OptError(C.GLP_ENOPFS)  // no primal feasible solution
	ENODFS  = OptError(C.GLP_ENODFS)  // no dual feasible solution
	EROOT   = OptError(C.GLP_EROOT)   // root LP optimum not provided
	ESTOP   = OptError(C.GLP_ESTOP)   // search terminated by application
	EMIPGAP = OptError(C.GLP_EMIPGAP) // relative mip gap tolerance reached
	ENOFEAS = OptError(C.GLP_ENOFEAS) // no primal/dual feasible solution
	ENOCVG  = OptError(C.GLP_ENOCVG)  // no convergence
	EINSTAB = OptError(C.GLP_EINSTAB) // numerical instability
	EDATA   = OptError(C.GLP_EDATA)   // invalid data
	ERANGE  = OptError(C.GLP_ERANGE)  // result out of range
)

// Error implements the error interface.
func (r OptError) Error() string {
	switch r {
	case EBADB:
		return "invalid basis"
	case ESING:
		return "singular matrix"
	case ECOND:
		return "ill-conditioned matrix"
	case EBOUND:
		return "invalid bounds"
	case EFAIL:
		return "solver failed"
	case EOBJLL:
		return "objective lower limit reached"
	case EOBJUL:
		return "objective upper limit reached"
	case EITLIM:
		return "iteration limit exceeded"
	case ETMLIM:
		return "time limit exceeded"
	case ENOPFS:
		return "no primal feasible solution"
	case ENODFS:
		return "no dual feasible solution"
	case EROOT:
		return "root LP optimum not provided"
	case ESTOP:
		return "search terminated by application"
	case EMIPGAP:
		return "relative mip gap tolerance reached"
	case ENOFEAS:
		return "no primal/dual feasible solution"
	case ENOCVG:
		return "no convergence"
	case EINSTAB:
		return "numerical instability"
	case EDATA:
		return "invalid data"
	case ERANGE:
		return "result out of range"
	}
	return "unknown error"
}

// Simplex solves LP with Simplex method. The argument parm may by nil
// (means that default values will be used). See also NewSmcp().
// Returns nil if problem have been solved (not necessarly finding
// optimal solution) otherwise returns an error which is an instanse
// of OptError.
func (p *Prob) Simplex(parm *Smcp) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var err OptError
	if parm != nil {
		err = OptError(C.glp_simplex(p.p.p, &parm.smcp))
	} else {
		err = OptError(C.glp_simplex(p.p.p, nil))
	}
	if err == 0 {
		return nil
	}
	return err
}

// Exact solves LP with Simplex method using exact (rational)
// arithmetic. argument parm may by nil (means that default values
// will be used). See also NewSmcp().  Returns nil if problem have
// been solved (not necessarly finding optimal solution) otherwise
// returns an error which is an instanse of OptError.
func (p *Prob) Exact(parm *Smcp) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var err OptError
	if parm != nil {
		err = OptError(C.glp_exact(p.p.p, &parm.smcp))
	} else {
		err = OptError(C.glp_exact(p.p.p, nil))
	}
	if err == 0 {
		return nil
	}
	return err
}

// Smcp represents simplex solver control parameters, a set of
// parameters for Prob.Simplex() and Prob.Exact(). Please use
// NewSmcp() to create Smtp structure which is properly initialized.
type Smcp struct {
	smcp C.glp_smcp
}

// NewSmcp creates new Smcp struct (a set of simplex solver control
// parameters) to be given as argument of Prob.Simplex() or
// Prob.Exact().
func NewSmcp() *Smcp {
	s := new(Smcp)
	C.glp_init_smcp(&s.smcp)
	return s
}

// Message level
type MsgLev int

const (
	// Message levels (default: glpk.MSG_ALL). Usage example:
	//
	//     lp := glpk.New()
	//     defer lp.Delete()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetMsgLev(glpk.MSG_ERR)
	//     if err := lp.Simplex(smcp); err != nil {
	//             log.Fatal(err)
	//     }
	MSG_OFF = MsgLev(C.GLP_MSG_OFF) // no output
	MSG_ERR = MsgLev(C.GLP_MSG_ERR) // warning and error messages only
	MSG_ON  = MsgLev(C.GLP_MSG_ON)  // normal output
	MSG_ALL = MsgLev(C.GLP_MSG_ALL) // full output
	MSG_DBG = MsgLev(C.GLP_MSG_DBG) // debug output
)

// SetMsgLev sets message level displayed by the optimization function
// (default: glpk.MSG_ALL).
func (s *Smcp) SetMsgLev(lev MsgLev) {
	s.smcp.msg_lev = C.int(lev)
}

// Simplex method option
type Meth int

const (
	// Simplex method options (default: glpk.PRIMAL). Usage example:
	//
	//     lp := glpk.New()
	//     defer lp.Delete()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetMeth(glpk.DUALP)
	//     if err := lp.Simplex(smcp); err != nil {
	//             log.Fatal(err)
	//     }
	//
	PRIMAL = Meth(C.GLP_PRIMAL) // use primal simplex
	DUALP  = Meth(C.GLP_DUALP)  // use dual; if it fails, use primal
	DUAL   = Meth(C.GLP_DUAL)   // use dual simplex
)

// SetMeth sets simplex method option (default: glpk.PRIMAL).
func (s *Smcp) SetMeth(meth Meth) {
	s.smcp.meth = C.int(meth)
}

// Pricing technique
type Pricing int

const (
	// Pricing techniques (default: glpk.PT_PSE). Usage example:
	//
	//     lp := glpk.New()
	//     defer lp.Delete()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetPricing(glpk.PT_STD)
	//     if err := lp.Simplex(smcp); err != nil {
	//             log.Fatal(err)
	//     }
	//
	PT_STD = Pricing(C.GLP_PT_STD) // standard (Dantzig rule)
	PT_PSE = Pricing(C.GLP_PT_PSE) // projected steepest edge
)

// SetPricing sets pricing technique (default: glpk.PT_PSE).
func (s *Smcp) SetPricing(pricing Pricing) {
	s.smcp.pricing = C.int(pricing)
}

// Ratio test technique
type RTest int

const (
	// Ratio test techniques (default: glpk.RT_HAR). Usage example:
	//
	//     lp := glpk.New()
	//     defer lp.Delete()
	//     ...
	//     smcp := glpk.NewSmcp()
	//     smcp.SetRTest(glpk.RT_STD)
	//     if err := lp.Simplex(smcp); err != nil {
	//             log.Fatal(err)
	//     }
	//
	RT_STD = RTest(C.GLP_RT_STD) // standard (textbook)
	RT_HAR = RTest(C.GLP_RT_HAR) // two-pass Harris' ratio test
)

// SetRTest sets ratio test technique (default: glpk.RT_HAR)
func (s *Smcp) SetRTest(r_test RTest) {
	s.smcp.r_test = C.int(r_test)
}

// Status returns status of the basic solution.
func (p *Prob) Status() SolStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return SolStat(C.glp_get_status(p.p.p))
}

// PrimStat returns status of the primal basic solution.
func (p *Prob) PrimStat() SolStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return SolStat(C.glp_get_prim_stat(p.p.p))
}

// DualStat returns status of the dual basic solution.
func (p *Prob) DualStat() SolStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return SolStat(C.glp_get_dual_stat(p.p.p))
}

// ObjVal returns objective function value.
func (p *Prob) ObjVal() float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_obj_val(p.p.p))
}

// RowStat returns the current status of i-th row auxiliary variable.
func (p *Prob) RowStat(i int) VarStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return VarStat(C.glp_get_row_stat(p.p.p, C.int(i)))
}

// TODO:
// glp_get_row_prim
// glp_get_row_dual

// ColStat returns the current status of j-th column structural
// variable.
func (p *Prob) ColStat(j int) VarStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return VarStat(C.glp_get_col_stat(p.p.p, C.int(j)))
}

// ColPrim returns primal value of the variable associated with j-th
// column.
func (p *Prob) ColPrim(j int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return float64(C.glp_get_col_prim(p.p.p, C.int(j)))
}

// TODO:
// glp_get_col_dual
// ...

// Iocp represents MIP solver control parameters, a set of
// parameters for Prob.Intopt(). Please use
// NewIocp() to create Iocp structure which is properly initialized.
type Iocp struct {
	iocp C.glp_iocp
}

// Checks whether the optional MIP presolver is enabled.
func (p *Iocp) Presolve() bool {
	if p.iocp.presolve == C.GLP_ON {
		return true
	}
	return false
}

// Enables or disables the optional MIP presolver.
func (p *Iocp) SetPresolve(on bool) {
	if on {
		p.iocp.presolve = C.GLP_ON
	} else {
		p.iocp.presolve = C.GLP_OFF
	}
}

// Create and initialize a new Iocp struct, which is used
// by the branch-and-cut solver.
func NewIocp() *Iocp {
	p := new(Iocp)
	C.glp_init_iocp(&p.iocp)
	return p
}

// Solve MIP problem with the branch-and-cut method.
func (p *Prob) Intopt(params *Iocp) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	err := OptError(C.glp_intopt(p.p.p, &params.iocp))
	if err != 0 {
		return err
	}
	return nil
}

// MipStatus returns status of a MIP solution.
func (p *Prob) MipStatus() SolStat {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	return SolStat(C.glp_mip_status(p.p.p))
}

// Returns value of the j-th column for MIP solution.
func (p *Prob) MipColVal(i int) float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	val := C.glp_mip_col_val(p.p.p, C.int(i))
	return float64(val)
}

// Returns value of the objective function for MIP solution.
func (p *Prob) MipObjVal() float64 {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	val := C.glp_mip_obj_val(p.p.p)
	return float64(val)
}

// MPS file format: either fixed (ancient) or free (modern) format.
type MPSFormat int

const (
	// MPS file format type (fixed or free). To read an MPS
	// (fixed) file and switch to maximization (as MPS format does
	// not specify objective function direction and GLPK assumes
	// minimization) run
	//
	//     lp := glpk.New()
	//     defer lp.Delete()
	//     lp.ReadMPS(glpk.MPS_DECK, nil, "someMaximizationProblem.mps")
	//     lp.SetObjDir(glpk.MAX)
	//     if err := lp.Simplex(nil); err != nil {
	//             log.Fatal(err)
	//     }
	//
	MPS_DECK = MPSFormat(C.GLP_MPS_DECK) // fixed (ancient) MPS format
	MPS_FILE = MPSFormat(C.GLP_MPS_FILE) // free (modern) MPS format
)

// PathError is the error used by methods reading and writing MPS,
// CPLEX LP, and GPLK LP/MIP formats.
type PathError struct {
	Op      string // operation (either "read" or "write")
	Path    string // name of the file on which the operation was performed
	Message string // short description of the problem
}

// Error implements the error interface.
func (e *PathError) Error() string {
	return e.Op + " " + e.Path + ": " + e.Message
}

// MPSCP represent MPS format control parameters
type MPSCP struct {
	mpscp C.glp_mpscp
}

// NewMPSCP creates new initialized MPSCP struct (MPS format control
// parameters)
func NewMPSCP() *MPSCP {
	m := new(MPSCP)
	C.glp_init_mpscp(&m.mpscp)
	return m
}

// WriteMPS writes the problem instance into a file in MPS file
// format.  The format argument specifies either the fixed or free MPS
// format.  The params argument can be nil (could also be a value
// returned by NewMPSCP() but at this point GLPK package does not
// allow to specify any MPS parameters available in GLPK).
//
// Note that MPS format does not specify objective function direction
// (minimization or maximization).
func (p *Prob) WriteMPS(format MPSFormat, params *MPSCP, filename string) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var parm *C.glp_mpscp
	if params != nil {
		parm = &params.mpscp
	}
	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_write_mps(p.p.p, C.int(format), parm, fname) != 0 {
		return &PathError{"write", filename, "MPS writing error"}
	}
	return nil
}

// ReadMPS reads the problem instance from a file in MPS file format.
// The format argument specifies either the fixed or free MPS format.
// The params argument can be nil (could also be a value returned by
// NewMPSCP() but at this point GLPK package does not allow to specify
// any MPS parameters available in GLPK).
//
// Note that MPS format does not specify objective function direction
// (minimization or maximization). GLPK assumes minimization, use
// SetObjDir(glpk.MAX) to switch to maximization if needed.
func (p *Prob) ReadMPS(format MPSFormat, params *MPSCP, filename string) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var parm *C.glp_mpscp
	if params != nil {
		parm = &params.mpscp
	}
	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_read_mps(p.p.p, C.int(format), parm, fname) != 0 {
		return &PathError{"read", filename, "MPS reading error"}
	}
	return nil
}

// CPXCP represent CPLEX LP format control parameters
type CPXCP struct {
	cpxcp C.glp_cpxcp
}

// NewCPXCP creates new initialized CPXCP struct (CPLEX LP format
// control parameters)
func NewCPXCP() *CPXCP {
	m := new(CPXCP)
	C.glp_init_cpxcp(&m.cpxcp)
	return m
}

// WriteLP writes the problem instance into a file in CPLEX LP file
// format. The params argument can be nil (could also be a value
// returned by NewCPXCP() but it is reserved for future use and at
// this point GLPK does allow to specify any CPLEX LP parameters).
func (p *Prob) WriteLP(params *CPXCP, filename string) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var parm *C.glp_cpxcp
	if params != nil {
		parm = &params.cpxcp
	}
	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_write_lp(p.p.p, parm, fname) != 0 {
		return &PathError{"write", filename, "CPLEX LP writing error"}
	}
	return nil
}

// ReadLP reads the problem instance from a file in CPLEX LP file
// format. The params argument can be nil (could also be a value
// returned by NewCPXCP() but it is reserved for future use and at
// this point GLPK does allow to specify any CPLEX LP parameters).
func (p *Prob) ReadLP(params *CPXCP, filename string) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	var parm *C.glp_cpxcp
	if params != nil {
		parm = &params.cpxcp
	}
	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_read_lp(p.p.p, parm, fname) != 0 {
		return &PathError{"read", filename, "CPLEX LP reading error"}
	}
	return nil
}

// ProbRWFlags represents flags used for reading and writing of the
// problem instance in the GLPK LP/MIP format. Reserved for future use
// for now zero value should be used.
type ProbRWFlags int

// WriteProb writes the problem instance into a file in GLPK LP/MIP
// file format. The flags argument is reserved for future use, for now
// zero value should be used.
func (p *Prob) WriteProb(flags ProbRWFlags, filename string) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_write_prob(p.p.p, C.int(flags), fname) != 0 {
		return &PathError{"write", filename, "GLPK LP/MIP writing error"}
	}
	return nil
}

// ReadProb reads the problem instance from a file in GLPK LP/MIP file
// format. The flags argument is reserved for future use, for now zero
// value should be used.
func (p *Prob) ReadProb(flags ProbRWFlags, filename string) error {
	if p.p.p == nil {
		panic("Prob method called on a deleted problem")
	}
	fname := C.CString(filename)
	defer C.free(unsafe.Pointer(fname))
	if C.glp_read_prob(p.p.p, C.int(flags), fname) != 0 {
		return &PathError{"read", filename, "GLPK LP/MIP reading error"}
	}
	return nil
}
