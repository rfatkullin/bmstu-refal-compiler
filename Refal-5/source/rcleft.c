

# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"


	/* Refal compiler function associated with translation of left sides.
		Separated into a separate file for the Tracer. */

char * mystrrev(s)
char * s;
   { char ch;
     char * sr = s;
     char * s0 = s;

     while(*sr != '\0') sr++;
     sr--;
     while( s != sr )
	  { ch = * s;
	    *s++  = *sr;
	    *sr   = ch;
	    if( sr-- == s ) break;
	  }
     return s0;
   }

int transl_left (e, number)
	struct element *e;
	int *number;
	{
	/*register*/ int i;
	int lend, rend, lnum, rnum;
	struct HOLES *holes, *hole, *h;
	int result;
	union param p1;
	int len;

#if MDEBUG_HOLES
	debug_print_expr(e);
#endif

	p1.n = 0L;

		/* search to the end of the expression. */
	/* 46 line. Was ...; ... != NULL; ... Shura. 29.01.98 */
	for (i = 0; e[i].type != 0; i ++) ;

		/* initialize the holes. */

	holes = (struct HOLES *) malloc (sizeof (struct HOLES));
	if (holes == NULL) 
		{
		fprintf (stderr, "ran out of memory allocating holes.\n");
		exit (1);
		}
	holes -> next = NULL;
	holes -> left = 0;	/* index of the left end. */
	holes -> right = i - 1;	/* index of the right end. */
	holes -> lte = lnum = *number; /* te number of the left end. */
	holes -> rte = rnum = *number+1; /* te number of the right end. */
	*number += 2;

	while (holes != NULL)
		{
#if MDEBUG_HOLES
	debug_print_holes(holes);
#endif
		hole = select_hole (holes, e, &len, &h);
		lend = hole -> left;
		rend = hole -> right;
		if (lnum != hole -> lte || rnum != hole -> rte)
			{
			lnum = hole -> lte;
			rnum = hole -> rte;
			p1.d.i1 = lnum;
			p1.d.i2 = rnum;
			rc_out (SETB, p1);
#if MDEBUG_HOLES
			printf("SETB %d %d\n", lnum, rnum);
#endif
			}
		holes = delete_hole (holes, hole);

		if (len)
			{
			rc_out (LEN, p1);

		/* add it to the table of variables. */
			table [table_len].index = e [lend].body.i;
			table [table_len].te_offset = (*number) + 1;
			table_len ++;
			lend ++;
			lnum = (*number) + 1;
			(*number) += 2;
			}
			/* otherwise if the expression in empty */
		else if (lend > rend)
			{
			p1.n = 0L;
			rc_out (EMP, p1);
			}
			
			/* go from left to right */
		while (lend <= rend)
			{
			if (e [lend].type == LPAR)
				{
					/* match parentheses */
				rc_out (PS, p1);

					/* going from left to right when we see structure 
						brackets, we open them up and push the rest of 
						the expression on the holes list. */
				holes = add_hole (holes, h, e[lend].body.i + 1, 
					rend, (*number)+1, rnum);

					/* start matching the inside of the parentheses. */
				rend = e[lend].body.i - 1;
				lend ++;
				lnum = (*number);
				rnum = (*number) + 1;
				(*number) += 2;

			/* check if the inside of parentheses is empty */
				if (lend > rend)
					{
					p1.n = 0L;
					rc_out (EMP, p1);
					}
				}
			else
				{
				result = match (e, lend, rend, &lnum, &rnum, 0, number);
				if (result == 0) break;	/* failed to match. */
				else lend ++;
				}
			}
			/* go from right to left */
		while (lend <= rend)
			{
			if (e [rend].type == RPAR)
				{
					/* match parentheses */
				rc_out (PSR, p1);

					/* going from right to left when we see structure brackets,
					we continue with the expression and throw the insides of the 
					brackets on the holes list. */
				holes = add_hole (holes, h, e[rend].body.i + 1, 
					rend - 1, *number, (*number) + 1);

					/* continue matching the outside preceding the parentheses. */
				rend = e[rend].body.i - 1; /* element preceding the pair */
				rnum = *number;	/* left parenthesis. */
				(*number) += 2;

			/* don't have to check if the expression before parentheses 
				is empty, because it should not happen: if it were so, then
				we would've been going from the left. */

				}
			else
				{
				result = match (e, lend, rend, &lnum, &rnum, 1, number);
				if (result == 0) break;	/* failed to match. */
				else rend --;
				}
			}
			/* add holes */
		if (lend < rend)
			{
				/* add one more hole to the list of holes: we hit
					an apparently open variable. */
			holes = add_hole (holes, h, lend, rend, lnum, rnum);
			}
		}
	free_holes (holes);
	return 0;
	}

int match (e, l, r, ln, rn, dir, n)

	struct element *e;	/* the whole expression. */
	int l;		/* index of the left end. */
	int r;		/* index of the right end. */
	int *ln;		/* te number of the  right end. */
	int *rn;		/* te number of the left end. */
	int dir;	/* 0 from left, 1 from right */
	int *n;	/* te number */

	{
	int x;
	int type;
	int occur;
	union param p1;

	if (dir) x = r;
	else x = l;
	type = e[x].type;

	switch (type)
		{
		case E_VAR:
			/* Check if this is an old variable. */
			occur = check_var (e[x].body.i);
			if (occur > -1)
				{
				p1.i = occur;
				rc_out (OEXP+dir, p1);
					/* add old var to the table of var. */
				table [table_len].index = e [x].body.i;
				table [table_len].te_offset = *n + 1;
				table_len ++;
				if (dir) *rn = *n;
				else *ln = *n + 1;
				*n += 2;
				if (l == r)
					{
					rc_out (EMP, p1);
					}
				return 1;
				}
			/* Check if this is a closed variable. */
			else if (l == r)
				{
				p1.n = 0L;
				rc_out (CL, p1);
				if (dir) *rn = *n;
				else *ln = *n + 1;
					/* add closed var to the table of var. */
				table [table_len].index = e [x].body.i;
				table [table_len].te_offset = *n + 1;
				table_len ++;
				*n += 2;
				return 1;
				}
				/* otherwise fail. */
			else return 0;

		case S_VAR:
			/* Check if this is an old variable. */
			occur = check_var (e[x].body.i);
			if (occur > -1)
				{
				p1.i = occur;
				rc_out (OVSYM+dir, p1);
				}
			else 
				{
				p1.n = 0L;
				rc_out (VSYM+dir, p1);
				}
				/* enter this variable into the list of vars. (even if old) */
			table [table_len].index = e[x].body.i;
			table [table_len].te_offset = *n;
			table_len ++;
			if (dir) *rn = *n;
			else *ln = *n;
			(*n) ++;
			if (l == r)
				{
				rc_out (EMP, p1);
				}
			return 1;

		case T_VAR:
			/* Check if this is an old variable. */
			occur = check_var (e[x].body.i);
			if (occur > -1)
				{
				p1.i = occur;
				rc_out (OEXP+dir, p1);
				}
			else 
				{
				p1.n = 0L;
				rc_out (TERM+dir, p1);
				}
				/* add to the table of vars. */
			table [table_len].index = e[x].body.i;
			table [table_len].te_offset = *n+1;
			table_len ++;
			if (dir) *rn = *n;
			else *ln = *n + 1;
			*n += 2;
			if (l == r)
				{
				rc_out (EMP, p1);
				}
			return 1;

		case CHAR:
			p1.c = e[x].body.c;
			rc_out (SYM+dir, p1);
/* A.P. Nemytykh after A.P. Konyshev. 28 January, 2002 */
                           if (dir) *rn = *n;
                           else *ln = *n;
			(*n) ++;
			if (l == r)
				{
				rc_out (EMP, p1);
				}
			return 1;

		case STRING:
			occur = strlen (e[x].body.f);
			if (occur > 1)
				{
				p1.f = e[x].body.f;
				rc_out (SYMS+dir, p1);
				if (dir) *rn = *n + occur - 1;
				else *ln = *n + occur - 1;
				*n += occur;
				}
			else 
				{
				p1.c = * (e[x].body.f);
				rc_out (SYM+dir, p1);
				if (dir) *rn = *n;
				else *ln = *n;
				(*n) ++;
				}
			if (l == r)
				{
				rc_out (EMP, p1);
				}
			return 1;

		case ATOM:
			p1.f = e[x].body.f;
			rc_out (CSYM+dir, p1);
			if (dir) *rn = *n;
			else *ln = *n;
			(*n) ++;
			if (l == r)
				{
				rc_out (EMP, p1);
				}
			return 1;

		case DIGIT:
			p1.n = e[x].body.n;
			rc_out (NSYM+dir, p1);
			if (dir) *rn = *n;
			else *ln = *n;
			(*n) ++;
			if (l == r)
				{
				rc_out (EMP, p1);
				}
			return 1;

		case LPAR:
		case RPAR:
			return 0;

		default:
			fprintf (stderr, "Internal Errors in MATCH: %d -- aborted: Line -- %d\n", type, line_no);
			exit (1);
		}
	return 0;
	}


		/* Holes manipulation routines. */
struct HOLES *delete_hole (holes, hole)

	struct HOLES *holes;
	struct HOLES *hole;

	{
	struct HOLES *h;

	if (holes == hole)
		{
		holes = hole -> next;
		free ((void *) hole);
		return holes;
		}
	h = holes;
	while (h != NULL)
		{
		if (h -> next == hole)
			{
			h -> next = hole -> next;
			free ((void *) hole);
			return holes;
			}
		else h = h -> next;
		}
	return holes;
	}

struct HOLES *select_hole (holes, e, len, prev_hole)

	struct HOLES *holes;
	struct element *e;
	int *len;
	struct HOLES **prev_hole;

	{
	struct HOLES *h;

	*len = 0;
	*prev_hole = NULL;
	h = holes;
	while (h != NULL)
		{
			/* see if there is only one closed variable (or empty). */ 
		if (h -> left >= h -> right) return h;
			/* select such a hole first that no lengthening is necessary. */
		if (no_lengthening (e + (h -> left)) || no_lengthening (e + (h -> right)))
			return h;
		*prev_hole = h;
		h = h -> next;
		}

		/* lenghtening is necessary */

	*len = 1;
		/* pick the first. */
	*prev_hole = NULL;
	return holes;
	}

		/* add one more hole to the list of holes. */
struct HOLES *add_hole (holes, h, lend, rend, lnum, rnum)

	struct HOLES *holes, *h;
	int lend, rend;
	int lnum, rnum;

	{
	struct HOLES *next_hole;

	next_hole = (struct HOLES *) malloc (sizeof (struct HOLES));
	if (next_hole == NULL)
		{
		fprintf (stderr, "no more memory for holes -- aborted\n");
		exit (1);
		}

	next_hole -> left = lend;
	next_hole -> right = rend;
	next_hole -> lte = lnum;
	next_hole -> rte = rnum;

#if MDEBUG_HOLES
	printf("adding hole: left= %2d  right= %2d   lte= %2d  rte= %2d\n",
		next_hole->left, next_hole->right, next_hole->lte, next_hole->rte);
#endif
	if (h == NULL)
		{
		next_hole -> next = holes;
		holes = next_hole;
		}
	else
		{
		next_hole -> next = h -> next;
		h -> next = next_hole;
		}

	return holes;
	}

int free_holes (holes)
	struct HOLES *holes;
	{
	struct HOLES *h;

	while (holes != NULL)
		{
		h = holes -> next;
		free ((void *) holes);
		holes = h;
		}
	return 0;
	}

int check_var (v)
	int v;
	{
		/* checks if variable v was defined before. if was then returns
			its offset in the table of elements. */
	/*register*/ int i;

	for (i = 0; i < table_len; i ++)
		{
		if (table [i].index == v) return table[i].te_offset;
		}
		/* not found. */
	return -1;
	}

int no_lengthening (e)
	struct element *e;
	{
	if (e -> type != E_VAR) return 1;
	else if (check_var (e -> body.i) != -1) return 1;
	else return 0;
	}

int rc_out (rasl_ins, pars)
	int rasl_ins;
	union param pars;
	{
	struct rasl_instruction *r;

		r = (struct rasl_instruction *)
			malloc (sizeof (struct rasl_instruction));
		/* 508 line. Not check result of the malloc. Shura. 29.01.98 */
		if (NULL == r) {
		  fprintf (stderr, "No memory for rasl instruction\n");
		  exit (1);
		}
		ftransl -> code = rasl_ins;
		ftransl -> next = r;
		/* 510 line. Was ... = NULL. Shura. 29.01.98 */
		r -> code = 0;
		r -> next = NULL;

		switch (rasl_ins)
			{
			case L: case E: case LABEL: case LBL: case ECOND: case SETB:
			case OEXP: case OEXPR: case OVSYM: case OVSYMR: case TRAN:
			case TPLE: case TPLS: case MULE: case MULS: case RDY:
			case NSYM: case NSYMR: case NNS:
			case CSYM: case CSYMR: case NCS: case ACT1:
			case SYMS: case TEXT: case SYM: case SYMR: case NS:

				ftransl -> p = pars;
				break;

			case SYMSR:
					/* reverse the string. */
				mystrrev (pars.f);
				ftransl -> p = pars;
				break;

					/* these instructions take no arguments. */
			case CL: case EMP: case VSYM: case VSYMR: case BL: case BLR:
			case BR: case LEN: case TERM: case TERMR: case PS: case PSR:
			case OUTEST: case PUSHVF: case POPVF: case STLEN: case B:

				break;

			default:
				fprintf (stderr, "Unknown RASL operator: %d -- ignored\n", rasl_ins);
				break;
			}
		ftransl = r;
		return 0;
	}

#if MDEBUG_HOLES
int debug_print_holes(holes)
	struct HOLES *holes;
	{
	struct HOLES *h;
	printf("HOLES:\n");
	for (h = holes; h != NULL; h = h->next)
		printf("left= %2d  right= %2d   lte= %2d  rte= %2d\n",
			h->left, h->right, h->lte, h->rte);
	return 0;
	}


int debug_print_expr(e)
	struct element *e;
	{
	int i;
	if (e == NULL || e[0].type == NULL) printf("Expression *empty*\n");
	else printf("Expression: \n");

	for (i = 0; e[i].type != NULL; i++) {
		switch(e[i].type) {
			case LPAR:
				printf("%2d: proj= %2d  element= (\n", i, e[i].number);
				break;
			case RPAR:
				printf("%2d: proj= %2d  element= ) pair= %d\n", i, e[i].number, e[i].body.i);
				break;
			case ACT_LEFT:
				printf("%2d: proj= %2d  element= <%s\n", i, e[i].number, e[i].body.f);
				break;
			case ACT_RIGHT:
				printf("%2d: proj= %2d  element= > pair= %d\n", i, e[i].number, e[i].body.i);
				break;
			case S_VAR:
				printf("%2d: proj= %2d  element= s.%d\n", i, e[i].number, e[i].body.i);
				break;
			case E_VAR:
				printf("%2d: proj= %2d  element= e.%d\n", i, e[i].number, e[i].body.i);
				break;
			case T_VAR:
				printf("%2d: proj= %2d  element= t.%d\n", i, e[i].number, e[i].body.i);
				break;
			case DIGIT:
				printf("%2d: proj= %2d  element= NUMBER: %lu\n", i, e[i].number, e[i].body.n);
				break;
			case ATOM:
				printf("%2d: proj= %2d  element= ATOM: %s\n", i, e[i].number, e[i].body.f);
				break;
			case CHAR:
				printf("%2d: proj= %2d  element= CHAR: %c (%d)\n", i, e[i].number, e[i].body.c, e[i].body.c);
				break;
			case STRING:
				printf("%2d: proj= %2d  element= STRING: %s\n", i, e[i].number, e[i].body.f);
				break;
			default: 
				printf("%2d: type= %d\n", i, e[i].type);
				break;
			}
		}
	return 0;
	}
#endif

