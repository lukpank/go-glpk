/*
 * This code is part of glpk package (Go bindings for the GNU Linear Programming Kit).
 *
 * Copyright (C) 2014 Łukasz Pankowski <lukpank@o2.pl>
 *
 * Pacakge glpk is free software: you can redistribute it and/or
 * modify it under the terms of the GNU General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Package glpk is distributed in the hope that it will be useful, but
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
 * General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with glpk package. If not, see <http://www.gnu.org/licenses/>.
 *
 * This example is a C rewrite of the PyGLPK example from
 * http://tfinley.net/software/pyglpk/discussion.html (Which is a
 * Python reimplementation of a C program from GLPK documentation)
 */

#include <stdio.h>
#include <glpk.h>

int main()
{
	glp_prob *lp = glp_create_prob();
	glp_set_prob_name(lp, "sample");
	glp_set_obj_name(lp, "Z");
	glp_set_obj_dir(lp, GLP_MAX);

	glp_add_rows(lp, 3);
	for (int i = 0; i < 3; i++) {
		char name[2] = { 'p' + i, 0 };
		glp_set_row_name(lp, i + 1, name);
	}
	glp_set_row_bnds(lp, 1, GLP_UP, 0, 100.0);
	glp_set_row_bnds(lp, 2, GLP_UP, 0, 600.0);
	glp_set_row_bnds(lp, 3, GLP_UP, 0, 300.0);
	
	glp_add_cols(lp, 3);
	for (int i = 0; i < 3; i++) {
		char name[3] = { 'x', '0' + i, 0 };
		glp_set_col_name(lp, i + 1, name);
		glp_set_col_bnds(lp, i + 1, GLP_LO, 0.0, 0.0);
	}

	int ind[] = {0, 1, 2, 3};

	glp_set_obj_coef(lp, 1, 10.0);
	glp_set_obj_coef(lp, 2, 6.0);
	glp_set_obj_coef(lp, 3, 4.0);

	double mat[3][4] = {
		{0, 1.0, 1.0, 1.0},
		{0,10.0, 4.0, 5.0},
		{0, 2.0, 2.0, 6.0}};
	for (int i = 0; i < 3; i++) {
		glp_set_mat_row(lp, i + 1, 3, ind, mat[i]);
	}

	glp_smcp parm;
	glp_init_smcp(&parm);
	glp_simplex(lp, &parm);

	printf("Z = %g", glp_get_obj_val(lp));
	for (int i = 0; i < 3; i++) {
		printf("; %s = %g", glp_get_col_name(lp, i + 1),
		       glp_get_col_prim(lp, i + 1));
	}
	putchar('\n');

	glp_delete_prob(lp);
}
