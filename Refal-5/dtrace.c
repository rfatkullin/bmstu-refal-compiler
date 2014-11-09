 
# include "rasl.h"
# include "decl.h"
# include "cdecl.h"
# include "ddecl.h"
# include "macros.h"
# include "ifunc.h"
# include "tfunc.h"
 
 
 
/* suppress debugging */
# define MDEBUG 0
 
extern FILE * fp_debugInfo;

/* Interactive Refal Tracer. */
 
int rd_trace (void) {
	FILE *ftmp;
 
	fprintf (rdout, " Refal TRACER: Step #%-5ld  ", nst);
	if (curr_point != -1) {
		fprintf (rdout, " Break point %4d.\n", curr_point);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, " Break point %4d.\n", curr_point);
	} else {
		fprintf (rdout, " Freeze occured.\n");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, " Freeze occured.\n");
	}

	if ((curr_point == 0) && (pr_res_flag == 1)) {
		rd_print_res ();
		pr_res_flag = 0;
	}
 
/* Enter interactive loop.	*/
	while (1) {
 
	/* 1. Get an instruction.	*/
		rd_getins ();
 
	/* 2. Get the command.		*/
		rd_getcom (1);
 
	/* 3. Test command.		*/
		if (combf [0] == '\0') ;
		else if ((strcmp (combf, "SET_BREAK") == 0) || (strcmp (combf, "BREAK") == 0) || (strcmp (combf, "SET") == 0)) rd_set_break ();
		else if ((strcmp (combf, "PRINT") == 0) || (strcmp (combf, "PR") == 0) || (strcmp (combf, "P") == 0)) rd_display ();
		else if ((strcmp (combf, "EXIT") == 0) || (strcmp (combf, "EX") == 0)) {
				fprintf(rdout, "Exit: Step =%9ld.\n", nst);
				if (fp_debugInfo != NULL) fprintf(fp_debugInfo, "Exit: Step =%9ld.\n", nst);
				ri_error(0);
		} else if ((strcmp (combf, "QUIT") == 0) || (strcmp (combf, "Q") == 0)) exit (0);
		else if ((strcmp (combf, "GO") == 0) || (strcmp (combf, "G") == 0)) {
			res.active = 0; 
			break;
		} else if ((strcmp (combf, "DELETE") == 0) || (strcmp (combf, "DEL") == 0) || (strcmp (combf, "D") == 0)) rd_del_brk ();
		else if ((strcmp (combf, "COMPUTE") == 0) || (strcmp (combf, "COMP") == 0) || (strcmp (combf, "COM") == 0)) {
			rd_compute ();
			break;
		} else if ((strcmp (combf, "STEP") == 0) || (strcmp (combf, "S") == 0)) {
			rd_step (); 
			break;
		} else if ((strcmp (combf, "SHOW") == 0) || (strcmp (combf, "SHO") == 0) || (strcmp (combf, "SH") == 0)) rd_show ();
		else if ((strcmp (combf, "HELP") == 0) || (strcmp (combf, "H") == 0)) rd_help ();
		else if (strcmp (combf, "FREEZE") == 0) break_at_freeze = 1;
		else if (strcmp (combf, "NOFREEZE") == 0) break_at_freeze = 0;

# if MDEBUG
		else if (strcmp (combf, "Z") == 0) rd_toggle ();
		else if (strcmp (combf, "XV") == 0) {
			printf ("view field in detail.\n");
			prexp (vf, vfend);
		} else if (strcmp (combf, "XA") == 0) {
			printf ("active expression in detail.\n");
			prexp (tbel (1), tbel (2));
		}
# endif
		else if (combf [0] == '>') {
			if (combf [1] == '\0') ftmp = stderr;
			else ftmp = fopen (combf+1, "at");
			if (ftmp == NULL) {
				fprintf (rdout, "Can\'t open %s\n", combf+1);
				if (fp_debugInfo != NULL) fprintf (rdout, "Can\'t open %s\n", combf+1);
			} else {
				if (rdout != stderr) fclose (rdout);
				rdout = ftmp;
			}
		} else { 
			fprintf (rdout, "Illegal command: %s\n", combf);
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Illegal command: %s\n", combf);
		}
	}
	return 0;
}
 
int rd_toggle (void) {
	dump_toggle = !dump_toggle;
	fprintf (rdout, "Dump is %d\n", dump_toggle);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Dump is %d\n", dump_toggle);
	return 0;
}
 
/* get an instruction from the input.	*/
int rd_getins (void) {
	register int i;
	int c;
	char temp_ibf [RD_INBUFSIZ];
	static int repeat = 0;
 
	/* examine the rest of the input buffer. */
	while ((c = rd_sep (ibf [ibp])) > -1) {
		if (ibp++ > RD_INBUFSIZ-1) break;
		if (c == 2) return 0;
	}
	ibp = 0;

	/* check if repeat count is on. */
	if (repeat > 0) {
		repeat --;
		fprintf (rdout, "TRACE> %s\n", ibf);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "TRACE> %s\n", ibf);
		return 0;
	}

 
	readin:
 
	fflush (stdout);
	fflush (rdout);
	fprintf (rdout, "TRACE> ");
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "TRACE> ");

# ifdef IBM370
	fprintf (rdout, "\n");
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "\n");
# endif

 	if (fgets (temp_ibf, RD_INBUFSIZ-1, rdin) == NULL) {
		fprintf (rdout, "EOF encountered.\n");
		fprintf (rdout, "Exit: Step = %9ld.\n", nst);
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "EOF encountered.\n");
			fprintf (fp_debugInfo, "Exit: Step = %9ld.\n", nst);
		}
		ri_error(0);
	} else if (temp_ibf [0] == '.') { /* check if this is a 'repeat previous command' command. */
		i = atoi (temp_ibf+1);
		if (i <= 1) repeat = 0;
		else repeat = i-1;
		fprintf (rdout, "%s\n", ibf);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "%s\n", ibf);

		ibp = 0;
		return 0;
	}

	/*  Strip the trailing blanks from the input. */
	/*  See if it is to be continued. */
	i = strlen (temp_ibf);
	if (i > RD_INBUFSIZ) i = RD_INBUFSIZ;
	i--;
	while ((i >= 0) && (rd_sep (temp_ibf [i]) == 1)) i--;
	if (i < 0) goto readin;
	else if (temp_ibf [i] == '\\') {
		if (ibp+i >= RD_INBUFSIZ) goto buffer_overflow;
		temp_ibf [i++] = '\n';
		temp_ibf [i] = '\0';
		strcpy (ibf + ibp, temp_ibf);
		ibp += i;
		goto readin;
	} else if (temp_ibf [i] == '\n') temp_ibf [i] = '\0';
	if (ibp + i >= RD_INBUFSIZ) goto buffer_overflow;
	temp_ibf [++i] = '\0';
	strcpy (ibf+ibp, temp_ibf);
	ibp = 0;
	return 0;

	buffer_overflow:

	fprintf (rdout, "Input too large -- retype\n");
	fprintf (rdout, "TRACE> ");
	if (fp_debugInfo != NULL) {
		fprintf (fp_debugInfo, "Input too large -- retype\n");
		fprintf (fp_debugInfo, "TRACE> ");
	}
# ifdef IBM370
	fprintf (rdout, "\n");
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "\n");
# endif
	ibp = 0;
	goto readin;
}
 
int rd_sep (c)
	char c;
	/* RD_SEP returns:
			-1 if c is '\0', 
			1 if c is a parameter separator, 
			2 if c is a command separator, 
			0 otherwise.	*/
	{
	switch (c)
		{
		case '\0': return -1;
		case ' ': case '\f': case '\r': case '\v':
		case '\n': case '\t': case ',': return 1;
		case ';': return 2;
		default : return 0;
		};
	}

/** if 1 then convert to upper case. **/
int rd_getcom (int f) {
	int i = 0;
 
	/* skip blanks.	*/
	while (rd_sep (ibf [ibp]) == 1) ibp++;
	/* copy the command converting it to upper case. */
	while (rd_sep (ibf [ibp]) == 0) {
		if (f) combf [i++] = toupper (ibf [ibp]); /* ??? */
		else combf [i++] = ibf [ibp];
		ibp ++;
	}
	combf [i] = '\0';
	return 0;
}
 
int rd_get_num_break ()
	{
	int i;
 
	for (i=0; i < BR_TAB_SIZ; i++)
		if (break_table [i].active == 0) break;
	return i;
	}
 
	/*		SET BREAK POINT.		*/
 
int rd_set_break (void) {
	int err;	/* error code. */
	char func_name [MAXWS+2];
	int fnum, i, j;
	char * curr_break;
	char * pc;
	int num_break;	/* number of the new break.	*/

	/* list of functions with the given name. */
	struct listfun *fs, *fstmp;
 
	/* 1. Get the first available break number.	*/
	num_break = rd_get_num_break ();
	if (num_break >= BR_TAB_SIZ) {
		fprintf (rdout, "No more breaks. ");
		fprintf (rdout, "Delete some of the old breaks.\n");
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "No more breaks. ");
			fprintf (fp_debugInfo, "Delete some of the old breaks.\n");
		}
		return 1;
	}
 
	/* 2. Call function rd_parse () to parse the input string.	*/
	curr_break = NULL;
	ASGN_CHARP(last_break + 1, pc);
	err = rd_parse (func_name, num_break, pc, &curr_break);
	if (err != 0) {
		fprintf (rdout, "Errors found: No break point set.\n");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Errors found: No break point set.\n");
		return 0;
	}
 
	/* 4. Get the function name.	*/
	fs = rd_implode (func_name, &fnum);
	if (fnum == 0) {
		fprintf (rdout, "Error: function %s does not exist.\n", func_name);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Error: function %s does not exist.\n", func_name);

		return 1;
	} else if (fnum == 1) fstmp = fs;
	else {
		/* function occurs in more than one module. */
		fprintf (rdout, "Function %s occurs in %d modules.\n", func_name, fnum);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Function %s occurs in %d modules.\n", func_name, fnum);

		getmodule:
		fprintf (rdout, "Type a number to choose the module.\n");
		fprintf (rdout, "\t 0\tto abort.\n");
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "Type a number to choose the module.\n");
			fprintf (fp_debugInfo, "\t 0\tto abort.\n");
		}
		fstmp = fs;
		for (i = 1; i <= fnum; i++) {
			fprintf (rdout, "\t%2d\t%-16s\taddress %10lx.\n", i, fstmp -> module, (unsigned long)(fstmp -> function));
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "\t%2d\t%-16s\taddress %10lx.\n", i, fstmp -> module, (unsigned long)(fstmp -> function));
			fstmp = fstmp -> next;
		}
		fprintf (rdout, "NUMBER> ");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "NUMBER> ");

		j = atoi (fgets (combf, 16, rdin));
		if ((j < 0) || (j > fnum)) {
			fprintf (rdout, "Illegal number.\n");
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Illegal number.\n");
			goto getmodule;
		} else if (j == 0) {
			fprintf (rdout, "ABORT: No break is set.\n");
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "ABORT: No break is set.\n");
 
			/* 4b. free space. */
			last_br_num = num_break;
			while (fs != NULL) {
				fstmp = fs -> next;
				free ((void *) fs);
				fs = fstmp;
			}
			return 1;
		}
		for (fstmp = fs, i=1; i<j; i++) fstmp = fstmp -> next;
	}
	fprintf (rdout, "Function = %s, Module = %s\n", func_name, fstmp -> module);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Function = %s, Module = %s\n", func_name, fstmp -> module);

	/* write the function address into the code. */
	wrcharp_to_mem (fstmp -> function, curr_break + 2 * sizeof (char) + sizeof (char *));
 
	/* free memory. */
	while (fs != NULL) {
		fstmp = fs -> next;
		free ((void *) fs);
		fs = fstmp;
	}
 
	/* 5. Activate the break point.	*/
	break_table [num_break].code = curr_break;
	break_table [num_break].active = 1;
	strcpy (break_table [num_break].rexp, ibf);
 
	/* 6. Link it with the previous break points. */
	break_table [last_br_num].nxt_brk = num_break;
	break_table [num_break].prc_brk = last_br_num;
	break_table [num_break].nxt_brk = -1;
	wrcharp_to_mem (curr_break, last_break + 1);
	last_break = curr_break;
	last_br_num = num_break;

	/* 7. Echo message and return. */
	fprintf (rdout, "Break point #%d is set.\n", num_break);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Break point #%d is set.\n", num_break);
	return 0;
}
 
int rd_display (void) {
	char vtype;
	int i, j;
 
	j = ibp;		/** save the buffer pointer. **/
	rd_getcom (1);	/* copy the argument into the combf. */
	if ((strcmp (combf, "VIEW") == 0) || (strcmp (combf, "V") == 0)) {
		fprintf (rdout, "The View-Field:\n");
		ri_putmb (vf, vfend, rdout);
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "The View-Field:\n");
			ri_putmb (vf, vfend, fp_debugInfo);
		}
		return 0;
	} else if ((strcmp (combf, "EXP") == 0) || (strcmp (combf, "ACT") == 0)) {
		fprintf (rdout, "Active Expression:\n");
		ri_putmb (tbel (1), tbel (2), rdout);
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "Active Expression:\n");
			ri_putmb (tbel (1), tbel (2), fp_debugInfo);
		}
		return 0;
	} else if ((strcmp (combf, "RES") == 0) || (strcmp (combf, "RESULT") == 0)) {
		rd_print_res ();
		return 0;
	} else if (strcmp (combf, "CALL") == 0) {
		fprintf (rdout, "The value of the call:\n");
		ri_putmb (res.l_exp, res.r_exp, rdout);
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "The value of the call:\n");
			ri_putmb (res.l_exp, res.r_exp, fp_debugInfo);
		}
		return 0;
	} else if (strcmp (combf, "STOCK") == 0) {
		fprintf (rdout, "The contents of the stock:\n");
		ri_putmb (stock, stock -> pair.b, rdout);
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "The contents of the stock:\n");
			ri_putmb (stock, stock -> pair.b, fp_debugInfo);
		}
		return 0;
	}
 
	/*** Otherwise the argument is a variable. ***/
	if (curr_point == -1) {
		fprintf (rdout, "Illegal argument to PRINT\n");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Illegal argument to PRINT\n");
		return 0;
	}
	ibp = j;		/** restore the buffer pointer. ***/
	rd_getcom (0);	/* copy the variable into the combf. */
	vtype = combf [0];
	j = (combf [1] == '.' ? 2 : 1);
	i = rd_look_for_var (vtype, combf+j, curr_point);
	if (i == break_table [curr_point].num_var) {
		fprintf (rdout, "Variable %s not found.\n", combf);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Variable %s not found.\n", combf);
		return 0;
	}
	fprintf (rdout, "Variable %s:\n", combf);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Variable %s:\n", combf);
	if ((vtype == 'e') || (vtype == 't') || (vtype == 'w')) {
		ri_putmb (tbel (break_table [curr_point].lv_tab [i].end-1), tbel (break_table [curr_point].lv_tab [i].end), rdout);
		if (fp_debugInfo != NULL)
			ri_putmb (tbel (break_table [curr_point].lv_tab [i].end-1), tbel (break_table [curr_point].lv_tab [i].end), fp_debugInfo);
	} else {
		ri_putmb (tbel (break_table [curr_point].lv_tab [i].end), NULL, rdout);
		if (fp_debugInfo != NULL) ri_putmb (tbel (break_table [curr_point].lv_tab [i].end), NULL, fp_debugInfo);
	}
	return 0;
}
 
int rd_look_for_var (v, indx, n) 
	char v, indx [];
	int n;
 
	{
	int i;
 
	for (i = 0; i < break_table [n].num_var; i++)
		{
		if ((v == break_table [n].lv_tab [i].typ) &&
			(strcmp (indx, break_table [n].lv_tab [i].index) == 0))
				break;
		};
	return i;
	}
 
int rd_del_brk (void) {
	int i;
 
	rd_getcom (1);
	i = atoi (combf);
	rd_delete_break (i);
	return 0;
}

/* delete break # n	*/
int rd_delete_break (int n) {
	int prec, next;
 
	if ((n == 0) || (n >= BR_TAB_SIZ)) {
		fprintf (rdout, "can\'t delete break #%3d.\n", n);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "can\'t delete break #%3d.\n", n);
		return 1;
	}
	if (break_table [n].active == 0) {
		fprintf (rdout, "Break point #%3d is not active.\n", n);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Break point #%3d is not active.\n", n);
		return 1;
	}
	/* 1. Delete the break point from the chain. */
	prec = break_table [n].prc_brk;
	next = break_table [n].nxt_brk;
	break_table [prec].nxt_brk = next;
	/* deleting break in the middle */
	if (next != -1) {
		break_table [next].prc_brk = prec;
	} else { /* deleting the last break */
		last_br_num = prec;
		last_break = break_table [last_br_num].code;
	}

	/* 2. Deactivate it.	*/
	break_table [n].active = 0;

	/* 3. Change the transfer address of the preceding break point.	*/
	memcpy ((char *) (break_table [prec].code + 1), (char *) (break_table [n].code + 1), sizeof (char *));

	/* 4. Release the memory.	*/
	free (break_table [n].code);
	fprintf (rdout, "Break Point #%d deleted.\n", n);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Break Point #%d deleted.\n", n);
	return 0;
}
 
int rd_show (void) {
	int i;
 
	rd_getcom (1);
	if (strcmp (combf, "ALL") == 0) rd_shwall ();
	else if ((strcmp (combf, "MODULES") == 0) || (strcmp (combf, "MODULE") == 0) || (strcmp (combf, "MOD") == 0)) {
		rd_getcom (1);
		if (combf [0] == '\0') rd_list_modules ();
		else if (combf [0] == '*') rd_shmodall ();
		else rd_shwmod (combf);
	} else if ((strcmp (combf, "FUNCTION") == 0) || (strcmp (combf, "FUNC") == 0) || (strcmp (combf, "FUN") == 0)) {
		rd_getcom (1);
		if (combf [0] == '\0') {
			fprintf (rdout, "Illegal format for SHOW command.\n");
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Illegal format for SHOW command.\n");
		} else {
			rd_shwfunc (combf);
		}
	} else if (strcmp (combf, "POINT") == 0) {
		rd_shwbreak (curr_point);
	} else if (strcmp (combf, "STEP") == 0) {
		fprintf (rdout, "Step = %7ld.\n", nst);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Step = %7ld.\n", nst);
	} else if (strcmp (combf, "CURK") == 0) {
		fprintf (rdout, "Current number of functions = %7ld.\n", curk);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Current number of functions = %7ld.\n", curk);
	} else if (strcmp (combf, "STACK") == 0) {
		fprintf (rdout, "Current number of functions = %7ld.\n", curk);
		rd_shwstack (rdout);
		if (fp_debugInfo != NULL) {
			fprintf (fp_debugInfo, "Current number of functions = %7ld.\n", curk);
			rd_shwstack (fp_debugInfo);
		}
	} else if (strcmp (combf, "BREAK") == 0) {
		rd_getcom (1);
		i = atoi (combf);
		if (i >= BR_TAB_SIZ) {
			fprintf (rdout, "No such break.\n");
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "No such break.\n");
			return 1;
		}
		fprintf (rdout, "Break #%d.\n", i);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Break #%d.\n", i);
		rd_shwbreak (i);
	} else {
		if (isdigit (combf [0]) && (i=atoi (combf)) >= 0) {
			if (i >= BR_TAB_SIZ) {
				fprintf (rdout, "No such break.\n");
				if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "No such break.\n");
				return 1;
			}
			fprintf (rdout, "Break #%d.\n", i);
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Break #%d.\n", i);
			rd_shwbreak (i);
		} else {
			fprintf (rdout, "Illegal format for SHOW command.\n");
			if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Illegal format for SHOW command.\n");
		}
	}
	return 0;
}
 
int rd_shwstack (FILE *fp) {
	LINK *q;
 
	fprintf (fp, "      ");

	{
		char * cp = actfun - 1;

		for (cp --; * cp != '\0'; cp --);
		cp ++;
		ri_actput (cp, fp);
	}
	/*ri_actput (actfun, fp);*/

	fprintf (fp, "\n");
	if (actfun == STOP_) return 0;
	q = nextp;
	while (q != vfend) {
		fprintf (fp, "      ");
		{
			char * cp = q -> pair.b -> pair.f - 1;

			for (cp --; * cp != '\0'; cp --);
			cp ++;
			ri_actput (cp, fp);
		}
		/*ri_actput (q -> pair.b -> pair.f, fp);*/

		fprintf (fp, "\n");
		q = q -> foll -> prec;
	}
	return 0;
}
		
/* print the break point # n */
int rd_shwbreak (int n) {
	if (n < 0) {
		fprintf (rdout, "No such break point\n");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "No such break point\n");
		return 0;
	}
	if (break_table [n].active == 0) {
		fprintf (rdout, "Break point #%d is not active.\n", n);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Break point #%d is not active.\n", n);
		return 1;
	}
	fprintf (rdout, "Break Point #%d: %s\n", n, break_table [n].rexp);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Break Point #%d: %s\n", n, break_table [n].rexp);
	return 0;
}
 
int rd_shwall (void) {
	int i, t;
 
	t = 0;
	fprintf (rdout, "All Active Break Points:\n");
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "All Active Break Points:\n");
	for (i = 1; i < BR_TAB_SIZ; i++) {
		if (break_table [i].active == 1) {
			rd_shwbreak (i);
			putc ('\n', rdout);
			if (fp_debugInfo != NULL) putc ('\n', fp_debugInfo);
			t++;
		}
	}
	fprintf (rdout, "Total: %2d Active Break Points.\n", t);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Total: %2d Active Break Points.\n", t);
	return 0;
}

/* Print all available commands. */
int rd_help (void) {
	FILE *f;
	int c, d, d1;
	int lines = 0;
	char *helpfile;

# ifdef IBM370
	char hf_name [20];
# else
	static char hf_name [] = "reftr.hlp";
# endif
 
# ifdef IBM370
	strcpy (hf_name, "REFTRACE HELP");
	helpfile = hf_name;
# else
	helpfile = getenv ("REFAL_HELP");
	if (helpfile == NULL) helpfile = hf_name;
# endif
 
	f = fopen (helpfile, "rt");
	if (f == NULL) {
		fprintf (rdout, "Cannot open the help file.\n");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Cannot open the help file.\n");
# ifndef IBM370
		fprintf (rdout, "The environment variable REFAL_HELP must contain the file name of the help file.\n");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "The environment variable REFAL_HELP must contain the file name of the help file.\n");
# else
		fprintf (rdout, "Help File %s is not found\n", helpfile);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Help File %s is not found\n", helpfile);
# endif
	} else {
		putc ('\n', rdout);
		if (fp_debugInfo != NULL) putc ('\n', fp_debugInfo);
		while ((c = getc (f)) != EOF) {
			putc (c, rdout);
			if (fp_debugInfo != NULL) putc (c, fp_debugInfo);
			if (c == '\n') lines++;
			if (lines >= 20) { 
				fprintf (rdout, "\nHit \'q\' and RETURN to abort; RETURN to continue: ");
				if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "\nHit \'q\' and RETURN to abort; RETURN to continue: ");
				d1 = d = getc (rdin);
				while ((d1 != '\n')&& (d1 != EOF)) d1=getc (rdin);
				if (d == 'q' || d == 'Q') return 0;
				lines = 0;
				putc ('\n', rdout);
				if (fp_debugInfo != NULL) putc ('\n', fp_debugInfo);
			}
		}
		fclose (f);
	}
	return 0;
}
 
int rd_cplink (LINK * r) {
	switch (r -> ptype) {
 	case LINK_TYPE_LSTRUCTB: case LINK_TYPE_LACT:
		bl; break;
 
	case LINK_TYPE_RSTRUCTB:
		br; break;
 
	case LINK_TYPE_CHAR:
		ns (r -> pair.c); break;
 
	case LINK_TYPE_COMPSYM:
		ncs (r -> pair.f); break;

	case LINK_TYPE_NUMBER:
		nns (r -> pair.n); break;
 
	case LINK_TYPE_RACT:
		br;
		b -> pair.b -> pair.f = r -> pair.b -> pair.f;
		b -> pair.b -> ptype = 5;
		b -> ptype = 6;
		break;
 
	case LINK_TYPE_SVAR:
	case LINK_TYPE_EVAR:
	case LINK_TYPE_TVAR:
		nns (r -> pair.n); b -> ptype = r -> ptype; break;
 
	default:
		fprintf (rdout, "rd_cplink: strange link %4d.\n", r->ptype);
		fprintf (fp_debugInfo, "rd_cplink: strange link %4d.\n", r->ptype);
		break;
	}
	return 0;
}
 
LINK *rd_cp_refx (le, re)
	LINK *le, *re;
 
	{
	LINK *q, *r;
 
	if (le == NULL) return NULL;
	if (le == re)
		{
		q = lfm;
		all;
		q -> ptype = le -> ptype;
		q -> pair = le -> pair;
		q -> prec = q;
		q -> foll = q;
		return q;
		};
	q = b = lfm;
	all;
	r = le;
	while (r != re)
		{
		rd_cplink (r);
		r = r -> foll;
		};
	rd_cplink (re);
	r = q -> foll;
	q -> foll = lfm;
	lfm = q;
	r -> prec = b;
	b -> foll = r;
	return r;
	}
 
int rd_step (void) {
	long k;
	char * cp_w;
		
	/* 1. Get the parameters. */
	rd_getcom (1);
	if ((strcmp (combf, "RES") == 0) || (strcmp (combf, "RESULT") == 0)) {
		pr_res_flag = 1;
		rd_getcom (1);
	}
	k = atol (combf);
	rd_getcom (1);
	if ((strcmp (combf, "RES") == 0) || (strcmp (combf, "RESULT") == 0)) pr_res_flag = 1;
	if (k == 0L) k = 1L;
 
	/* 2. Set Break point 0. */
	wrlong_to_mem (k+nst, break0 + sizeof (char) + sizeof (char *) + sizeof (char));
 
	/* 3. Now set the global variable RES. */
	/* test if the active function is an auxiliary (i.e. $-function) */
	for (cp_w = actfun - 2; * cp_w != 0; cp_w --);
	cp_w ++;
	/*if (strnchr (actfun-MAXWS, '$', MAXWS) != NULL) {*/
	if (strchr (cp_w, '$') != NULL) {
		res.leftend = NULL;
		res.ritend = NULL;
		res.l_exp = NULL;
		res.r_exp = NULL;
		res.active = 0;
		res.nsteps = 0;
	} else {
		res.leftend = tbel (0);
		res.ritend = tbel (2) -> foll;
		if (res.l_exp != NULL) {
			res.r_exp -> foll = lfm;
			lfm = res.l_exp;
		}
		res.l_exp = rd_cp_refx (tbel (1), tbel (2));
		res.r_exp = res.l_exp -> prec;
		res.nsteps = nst;
		res.active = 1;
	}
	return 0;
}
 
int rd_compute (void) {
	char * cp_w;

	/* 1. See if there is a 'res' modifier. */
	rd_getcom (1);
	if ((strcmp (combf, "RES") == 0) || (strcmp (combf, "RESULT") == 0)) pr_res_flag = 1;
 
	/* 2. Set Break Point #0 in such a way that the Refal */
	/*    Interpreter will continue until the current active */
	/*    expression is rd_computed. */
 
	mwp = break0+1;
 
	/* Save the address of the next break point. */
	ASGN_CHARP (mwp, res.ra);
	wrmemi (NOBREAK_);	/* Nobreak pointer.	*/
	wrmemb (CURK);
	wrmemi (curk-1);	
	/* The rest is the same: CL; EBR, 0	*/
	/* 3. Now set the global variable RES.	*/
	/* test if the active function is a $-function. */
	for (cp_w = actfun - 2; * cp_w != 0; cp_w --);
	cp_w ++;
	/*if (strnchr (actfun - MAXWS, '$', MAXWS) != NULL)*/
	if (strchr (cp_w, '$') != NULL) {
		res.leftend = NULL;
		res.ritend = NULL;
		res.l_exp = NULL;
		res.r_exp = NULL;
		res.active = 0;
		res.nsteps = 0;
	} else {
		res.leftend = tbel (0);
		res.ritend = tbel (2) -> foll;
		if (res.l_exp != NULL) {
			res.r_exp -> foll = lfm;
			lfm = res.l_exp;
		}
		res.l_exp = rd_cp_refx (tbel (1), tbel (2));
		res.r_exp = res.l_exp -> prec;
		res.nsteps = nst;
		res.active = 1;
	}
	return 0;
}
 
LINK *rd_adjust_re (l, r)
 
	/* Adjust right end of a result expression. Aug 15 1986. DT */
	LINK *l, *r;
 
	{
	LINK *x;
 
	x = l;
	l = l -> foll;
	while (l != r)
		{
		x = l;
		if (LINK_LSTRUCTB (l)) l = l -> pair.b;
		else l = l -> foll;
		};
	return x;
	}
 
/* Print the result of an expression. */
int rd_print_res (void) {
	LINK *z;
 
	if (res.leftend == NULL) return 1;
	fprintf (rdout, "The Result is:\n");
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "The Result is:\n");
	if (res.leftend -> foll == res.ritend) {		
		fprintf (rdout, "***EMPTY***\n");
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "***EMPTY***\n");
	} else {
		z = rd_adjust_re (res.leftend, res.ritend);
		ri_putmb (res.leftend -> foll, z, rdout);
		if (fp_debugInfo != NULL) ri_putmb (res.leftend -> foll, z, fp_debugInfo);
	}
	return 0;
}

/* End of rd_compute cycle: restore Break Point 0. */
int rd_endcom (void) {
	mwp = break0;
	mwp++;

	/* restore the address of the next break point.	*/
	wrmemi (res.ra);
	wrmemb (EQS); /* EQS = 80 */
	wrmemi (1);	/* stop on the first rd_step */

	/* The rest is the same: CL; EBR, 0	*/
	return 0;
}

/* Initialize the tracer. */
int rd_init (void) {
# ifdef VMS
	char *str;
# endif
	int k;
 
	/* 1. Open the interactive I/O files.	*/
# ifdef VMS
	str = ctermid (0);
	rdin = fopen (str, "r");
# elif defined(PCAT)
	rdin = stdin;
# endif

	rdout = stderr;
	if ((rdin == NULL) || (rdout == NULL)) {
		fprintf (stderr, "cannot open terminal.\n");
		exit (3);
	}
 
	/* 2. Set up Break Point #0. */
	mwp = break0;
	wrmemb (TRAN);
	wrmemi (NOBREAK_); /* Nobreak pointer. */
	wrmemb (EQS);
	wrmemi (1);	/* stop on the first step */
	wrmemb (CL);
	wrmemb (EBR);
	wrmemb (0);	/* break = 0 */
 
	break_table [0].active = 1;
	strcpy (break_table [0].rexp, "set_break  at STEP = 1");
	break_table [0].code = break0;
	break_table [0].lv_tab [0].typ = 'e';
	break_table [0].lv_tab [0].end = 4;
	strcpy (break_table [0].lv_tab [0].index, ".1");
	break_table [0].num_var = 1;
 
	/* 3. Set some other variables.	*/
	module_list = NULL;
	curr_point = 0;
	last_br_num = 0;
	last_break = break0;
	pr_res_flag = 0; /* print result flag. */

	/* input buffer for the tracer.	*/
	ibf [RD_INBUFSIZ-1] = ibf [ibp = 0] = '\0';
 
	/* variable RES	*/
	res.leftend = NULL;
	res.ritend = NULL;
	res.l_exp = NULL;
	res.r_exp = NULL;
	res.active = 0;
	res.nsteps = 0;
 
	/* 4. Set all other breaks to inactive.	*/
	for (k=1; k < BR_TAB_SIZ; k++) break_table [k].active = 0;
 
	/* 5. Set various flags. */
	break_at_freeze = 0;
	dump_toggle = 0;
 
	return 0;
}

/* n must be not more than 100. */
char * strnchr (char * s, char c, int n) {
	char buf [100];

	sprintf (buf, "%-*.*s", n, n, s);
	buf [n] = '\0';
	return strchr (s, c);
}
