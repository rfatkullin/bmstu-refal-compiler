

	/*** External definitions. ***/
# ifndef EXTERNAL
# ifdef DEFINE_EXTERNALS
# define EXTERNAL 
# else
# define EXTERNAL extern
# endif
# endif

# ifndef FOR_OS_SUNLINUX
# 	define wrmemi(i) { *((long *) mwp) = (long)(i); mwp += sizeof(long);}
# else
/*** MAC, Sun (Linux) Versions ***/
# 	define wrmemi(i) { long TMP; TMP = (long)(i); memcpy (mwp, (char *) &TMP, sizeof (long)); mwp += sizeof(long);}
/****/
# endif
# define wrmemb(b)  { *mwp++ = b;} 


			/* size of the table of break points.	*/
# define BR_TAB_SIZ 10
			/* maximum number of variables used in a break */
# define MAX_VAR_PER_BREAK 10
			/* number of character in the input line. */
# define RD_INBUFSIZ 128

		/* IRT - Interactive Refal Tracer.	*/

	struct var_tab	/* table of variables.		*/
		{
			char typ;	/* variable type. */
			char index [MAXWS];	/* index. */
			int end;		/* the entry of table of elements. */
		};

	struct breakpt 	/* break point data structure.	*/
		{
			int active;	/* 0 if not, 1 otherwise. */
			char rexp [RD_INBUFSIZ];	/* Refal expression. */
			struct var_tab lv_tab [MAX_VAR_PER_BREAK]; /* list of local variables. */
			int num_var;		/* number of variables used. */
			int nxt_brk,prc_brk; /* pointers to the next and preceding breaks. */
			char *code;	/* pointer to the code.	*/
		};

	EXTERNAL struct breakpt break_table[BR_TAB_SIZ];	
				/* table of break points.	*/

	struct respt	/* result point.	*/
		{
			LINK *leftend,*ritend;	/* pointers to
				left and right ends of the result 
				expression	*/
			LINK *l_exp,*r_exp;	/* pointers to
				left and right ends of the copy of the 
				original  expression	*/
			char *ra;	/* address of the next break 
				point.		*/
			long nsteps;	/* number of steps executed
				prior to the original expression.	*/
			char active;	/* flag if active.	*/
		};
	EXTERNAL struct respt res;

	/* listfun and modlist structures are needed to resolve possible
		clashes in specification of function names for breaks.
			August 9 1985.  DiTu.	*/
	struct listfun 	/* list of function and names of modules.	*/
		{
			char *function; /* addr of the function. */
			char *module; /* name of the module. */
			struct listfun *next;
		};

	struct modlist	/* list of all modules.	*/
		{
			char *module; /* name of the module. */
			char **funcs; /* pointer to the list of functions. */
			struct modlist *next;
		};
	EXTERNAL struct modlist *module_list;	/* list of modules  */


	EXTERNAL int last_br_num;		/* number of the last break.	*/
	EXTERNAL int curr_point;			/* the current break point.	*/
	EXTERNAL int pr_res_flag;		/* print result flag.		*/

	EXTERNAL char ibf [RD_INBUFSIZ];	/* input buffer.		*/
	EXTERNAL int ibp;	/* its pointer.			*/
	EXTERNAL char combf [16];	/* command buffer.		*/
	EXTERNAL char *mwp;	/* memory write pointer.	*/
	EXTERNAL char *last_break; /* pointer to the last break.	*/
	EXTERNAL char break0 [13]; /* pointer to break 0.		*/

	EXTERNAL char *NOBREAK_;	/* Nobreak.	*/

		/* Tracer Input and Output file pointers. */
	EXTERNAL FILE *rdin, *rdout;

		/* Debugging. */
	EXTERNAL int break_at_freeze;
	EXTERNAL int dump_toggle;



