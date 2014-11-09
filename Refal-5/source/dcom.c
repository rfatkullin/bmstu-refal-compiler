
# include "rasl.h"
# include "cdecl.h"
# include "decl.h"
# include "ddecl.h"
# include "ifunc.h"
# include "tfunc.h"
# include "cfunc.h"


# ifndef min
# define min(a,b)  (((a) < (b)) ? (a) : (b))
# endif

# define MDEBUG 0

	/* this file contains the main parser routine for the Tracer. */

int rd_vyvod (struct rasl_instruction *, char *, int, char **);

/* pointer to the beginning of the results. */
int rd_parse (char func_name [], int break_no, char * last_break, char ** results) {
	struct element *e;
	int number;
	int err;
	struct rasl_instruction *translation, *q;
	char *blk_temp;

	/* don't need to allocate memory for tables etc. */

	/* set the output stream. */
	fdlis = rdout;

	/* initialize the parser. */
	err = 0;
	rc_initrp ();
	last_label = 1;

	/* get the first token and extract the function name. */
	token = rc_gettoken ();
	if (token != LCBRAK) {
		fprintf (rdout, "Illegal pattern.\n");
		return 1;
	}
	strcpy (func_name, str);
	nonret = (char *) malloc (MEM_BLK_SIZE);
	if (nonret == NULL) {
		fprintf (stderr,"Can\'t allocate memory.\n");
		return 1;
	}
	nrptr = 0;

	/* parse the pattern expression. */
	token = rc_gettoken ();
	e = refal_expression (0);
	if (token != '>') {
		fprintf (rdout, "Missing \'>\'\n");
		free ((void *) e);
		err = 1;
	}
	if (e == NULL) err = 2;
	else if (nerrors > 0) {
		free ((void *) e);
		err = 2;
	}

	/* set up some variables for the compiler. */
	if (!err) {
		number = 1;
		ftransl = (struct rasl_instruction *) malloc (sizeof (struct rasl_instruction));
		/* No checking of malloc result. Shura. 27.01.98 */
		if (NULL == ftransl) {
			exit (1);
		}
		/* Was ftransl -> code = NULL. Shura. 27.01.98 */
		ftransl -> code = 0;
		ftransl -> next = NULL;

		/* translation -> next is the first element of translation. */
		translation = (struct rasl_instruction *) malloc (sizeof (struct rasl_instruction));
		/* No checking of malloc result. Shura. 27.01.98 */
		if (NULL == translation) {
			exit (1);
		}
		/* Was translation -> code = NULL. Shura. 27.01.98 */
		translation -> code = 0;
		translation -> next = ftransl;

		table_len = 0;

		/* perform the actual translation. */
		transl_left (e, &number);

		/* copy variables to the break table with their offsets. */
		rd_cr_lvtab (break_no);

		/* free the memory for the expression. */
		free ((void *) e);

		/* perform the post optimizing ritual. */
		rc_post (translation);

		/* output the result. */
		rd_vyvod (translation, last_break, break_no, results);

		/* free the memory. */
		while (translation != NULL) {
			q = translation -> next;
			free ((void *) translation);
			translation = q;
		}
	}

	/* free the discardable memory: for strings etc. */
	while (block != NULL) {
		blk_temp = * ((char **) block);
		free ((void *) block);
		block = blk_temp;
	}

	/* free the compiler non-returnable memory: the Tracer doesn't need it. */
	if (nonret != NULL) free ((void *) nonret);

	return err;
}

	/* Alternative definition of rc_gchar () and rc_ungchar () */
int rc_gchar ()
	/* Get the next character from the terminal.  */
	{
	return (ibf[ibp] == '\0'? ';' : ibf[ibp++]);
	}

	/* Return one character to the input stream.  */
int rc_ungchar (c)
	char c;
	{ ibp --; return 0; }

		/* rd_vyvod () copies translation into memory. */
int rd_vyvod (struct rasl_instruction * translation, char * last_break,
	      int break_no, char ** results) {

	struct rasl_instruction *z;
	struct
		{
		int n;
		int o;
		} labels [30];	/* no more then 30 labels. */
	int l, num_labels, i, opcode, j;
	char *cs;

		/* estimate the size of translation and resolve all labels. */
	l = 0;
	num_labels = 0;
	z = translation;
	while (z != NULL)
		{
		switch (z -> code)
			{
				/* these are not instructions. */
			/*case NULL:*/case 0: case B:
				break;

				/* define labels: store the offset in the instruction. */
			case LABEL: case LBL:

				labels [num_labels].n = z -> p.i;
				labels [num_labels].o = l;
				num_labels ++;
				if (num_labels >= 30) 
					{
					fprintf (rdout, "Expression too  large.\n");
					return 1;
					}
				break;

				/* takes 2 integer arguments. */
			case SETB:
				l += sizeof (char) + 2 * sizeof (long);
				break;

			/* these instructions take 1 integer argument (it is stored in 1 byte). */
			case OEXP: case OEXPR: case OVSYM: case OVSYMR: 
				/* these instructions take 1 character argument (also 1 byte). */
			case SYM: case SYMR: case LENS:
				l += sizeof (char) + sizeof (char);
				break;

				/* these instructions take 1 long number argument. */
			case NSYM: case NSYMR:
				l += sizeof (char) + sizeof (long);
				break;

				/* these instructions take 1 character pointer argument. */
			case CSYM: case CSYMR: case TRAN:
				l += sizeof (char) + sizeof (char *);
				break;

				/* these instructions take 2 arguments: length and pointer. */
			case SYMS: case SYMSR:
				l += 2 + strlen (z -> p.f);
				break;

				/* these instructions take no arguments. */
			case CL: case EMP: case VSYM: case VSYMR: case BL: case BLR:
			case BR: case LEN: case TERM: case TERMR: case PS: case PSR:
			case OUTEST: case PUSHVF: case POPVF: case STLEN: case PLENS:
			case PLENP: case PLEN: case LENP:

				l ++;
				break;

			default:
				fprintf (stderr, "Invalid code in rd_vyvod (): %d\n", z -> code);
				exit (1);
				break;
			}
		z = z -> next;
		}

		/* Now we are ready to proceed. */
		/* allocate memory. */
	l += 30;		/* more or less arbitrary number. Needed for some overhead */
	*results = (char *) malloc (l);
	if (*results == NULL)
		{
		fprintf (rdout, "Not enough memory\n");
		return 1;
		}
	mwp = *results;

		/* write the necessary stuff. */
	wrmemb (TRAN);
	wrmemi (last_break);
	wrmemb (CHACT);
	wrmemi (NULL);	/* for now write Zero here -- it will be replaced
							by the function name later. */
	j = 2 * sizeof (char) + 2 * sizeof (char *);
	z = translation;

		/* make the second pass. */
	while (z != NULL)
		{
		opcode = z -> code;
		switch (opcode)
			{
				/* these are not instructions. */
			/*case NULL:*/case 0: case B:
				/* skip labels: */
			case LABEL: case LBL:
				break;

				/* takes 2 integer arguments. */
			case SETB:
				wrmemb (opcode);
				wrmemi (z -> p.d.i1);
				wrmemi (z -> p.d.i2);
				break;

			/* these instructions take 1 integer argument (it is stored in 1 byte). */
			case OEXP: case OEXPR: case OVSYM: case OVSYMR: 
				wrmemb (opcode);
				wrmemb (z -> p.i);
				break;

				/* these instructions take 1 character argument (also 1 byte). */
			case SYM: case SYMR: case LENS:
				wrmemb (opcode);
				wrmemb (z -> p.c);
				break;

				/* these instructions take 1 long number argument. */
			case NSYM: case NSYMR:
				wrmemb (opcode);
				wrmemi (z -> p.n);
				break;

				/* these instructions take 1 character pointer argument. */
			case CSYM: case CSYMR:

				wrmemb (opcode);
			/* implode the compound symbol */
				cs = ri_cs_impl (z -> p.f);
				wrmemi (cs);
				break;

				/* this instruction takes 1 character pointer argument. */
			case TRAN:
				wrmemb (opcode);
					/* find the offset of this label. */
				for (i = 0; i < num_labels; i ++)
					if (labels [i].n == z -> p.i) break;
				if (i == num_labels)
					{
					fprintf (rdout, "Label %d not found\nUnable to create a break point\n", z->p.i);
					free (*results);
					*results = NULL;
					return 1;
					}
					/* save the offset. */
				wrmemi ((*results) + j + labels [i].o);
				break;

				/* these instructions take 2 arguments: length and pointer. */
			case SYMS: case SYMSR:

				wrmemb (opcode);
				i = strlen (z -> p.f);
				wrmemb (i);
				strncpy (mwp, z -> p.f, i);
				mwp += i;
				break;

				/* these instructions take no arguments. */
			case CL: case EMP: case VSYM: case VSYMR:
			case LEN: case TERM: case TERMR: case PS: case PSR:
			case PLENS: case PLENP: case PLEN: case LENP:

				wrmemb (opcode);
				break;

			default:
				fprintf (stderr, "Invalid code in rd_vyvod (): %d\n", opcode);
				exit (1);
				break;
			}
		z = z -> next;
		}

	wrmemb (EBR);
	wrmemb (break_no);

# if MDEBUG
	if (*results + l < mwp)
		{
		fprintf (stderr, "memory overwrite occured.\n");
		exit (1);
		}
# endif

	return 0;
	}

	/* copy from the table of elements into the break point structure. */
int rd_cr_lvtab (b)
	int b;
	{
	int i, m, k;

	m = min (table_len, MAX_VAR_PER_BREAK);
	for (i = 0; i < m; i ++)
		{
		k = table [i].index;
			/* copy index of  the variable. */
		strcpy (break_table [b].lv_tab [i].index, lv [k].vindex);
			/* copy the type of the variable. */
		break_table [b].lv_tab [i].typ = lv [k].vt;
			/* copy the table of elements entry. */
		break_table [b].lv_tab [i].end = table [i].te_offset;
		}
		/* the total number of variables. */
	break_table [b].num_var = m;
	return 0;
	}

