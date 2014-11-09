
# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"
# include "memory.h"

/* initialize Refal parser.	*/
int rc_initrp (void) {
  nlv = 0;

  /* allocate discardable memory for the compiler */
  /* 13 line. Not check result of malloc. Shura. 29.01.98 */
  if ((block = (char *) malloc (MEM_BLK_SIZE)) == NULL) {
    fprintf (stderr, "No memory for compile\n");
    exit (1);
  }
  /* put NULL into the first location of the discardable memory */
  * ((char **) block) = NULL;

  /* initialize block pointer. */
  blptr = sizeof (char *);

  /* initialize other variables. */
  if (cbuf [sc] == 0) {
	  cbuf[sc = 0] = '\n';
  }
  vsp = 0;
  nerrors = 0;
  globsav = 0;
  token = 0;
  v = 0L;
  return 0;
}

/* memory that will not be freed back. */
char *rc_allmem (int l) {
	int i; char *pt;

	if ((l + (i = nrptr)) < MEM_BLK_SIZE) {
		nrptr += l;
		return nonret+i;
	} else {
		pt = (char *) malloc (MEM_BLK_SIZE);
		if (pt == NULL) {
			fprintf (stderr,"Out of memory.\n");
			exit (1);
		}
		nonret = pt;
		nrptr = l;
		return nonret;
	}
}

	/* memory allocator. (allocates memory from the block) 
		All this memory will be released after the function is compiled. */
char *rc_memral (l)
	int l;
	{
	int i; 
	char *pt;

	if ((l+ (i = blptr)) < MEM_BLK_SIZE)
		{
		blptr += l;
		return block+i;
		}
	else
		{
		pt = (char *) malloc (MEM_BLK_SIZE);
		if (pt == NULL)
			{
			fprintf (stderr,"Out of memory.\n");
			exit (1);
			};
		wrcharp_to_mem (block, pt);
		block = pt;
		blptr = sizeof (char *) + l;
		return block + sizeof (char *);
		}
	}


	/* Refal compiler post optimization. (separated for the Tracer). */

int rc_post_opt (translation, t)
	struct rasl_instruction *translation [];
	int t;
	{
	int i;

	for (i = 1; i < t; i ++)
		{
		rc_post (translation [i]);
		}
	return 0;
	}


int rc_post (q)
	struct rasl_instruction *q;
	{
	struct rasl_instruction *z, *q1;
	union param p1;

	while ((q1 = q -> next) != NULL)
		{
		switch (q1 -> code)
			{
				/* if it is a special mark: B, cut it out. */
			case B:
				z = q1 -> next;
				free ((void *) q1);
				q -> next = z;
				break;

				/* merge together BL, BR to BLR */
			case BL:
				z = q1 -> next;
				if (z -> code == BR)
					{
					q1 -> code = BLR;
					q1 -> next = z -> next;
					free ((void *) z);
					}
				q = q1;
				break;

			case LEN:
				z = q1 -> next;
				if (z -> code == SYM)
					{
					z -> code = LENS;
					q1 -> code = PLENS;
					}
				else if (z -> code == PS)
					{
					z -> code = LENP;
					q1 -> code = PLENP;
					}
				else
					{
					p1.n = 0L;
					insert_instruction (q, PLEN, p1);
					}
				q = q -> next;
				q = q -> next;
				break;

			case NS: case TEXT:

				merge_string_instr (NS, TEXT, q1);
				q = q -> next;
				break;

			case SYM: case SYMS:

				merge_string_instr (SYM, SYMS, q1);
				q = q -> next;
				break;

			case SYMR: case SYMSR:

				merge_string_instr (SYMR, SYMSR, q1);
				q = q -> next;
				break;

			default: 
				q = q1;
				break;
			}
		}
	return 0;
	}

int merge_string_instr (ins1, ins2, q1)
	int ins1, ins2;
	struct rasl_instruction *q1;
	{
	int code, l;
	char *s, temp_str [2];
	struct rasl_instruction *z;
		
		/* compute the total length of the string. */
	l = 0;
	z = q1;
	while (z != NULL)
		{
		if (z -> code == ins1) l ++;
		else if (z -> code == ins2) l += strlen (z -> p.f);
		else break;
		z = z -> next;
		}
	z = q1 -> next;
	code = z -> code;
		/* if there is no more ins1 or ins2 instructions, leave. */
	if (code != ins1 && code != ins2) return 0;

		/* allocate memory. */
	s = rc_memral (l+1);
	if (q1 -> code == ins2)
		{
		strcpy (s, q1 -> p.f);
		}
	else 
		{
			/* change the first ins1 to ins2. */
		s [0] = q1 -> p.c;
		s [1] = '\0';
		q1 -> code = ins2;
		}
	q1 -> p.f = s;

		/* merge all ins1 and ins2 instructions into one. */
	while (z != NULL)
		{
		if (code == ins1)
			{
			temp_str [0] = z -> p.c;
			temp_str [1] = '\0';
			strcat (s, temp_str);
			}
		else if (code == ins2)
			{
			strcat (s, z -> p.f);
			}
		else break;

			/* free the instruction. */
		q1 -> next = z -> next;
		free ((void *) z);
		z = q1 -> next;
		code = z -> code;
		}
	return 0;
	}


		/* expression translation. */
# define next_token (token = rc_gettoken ())

/* this file deals with parsing expressions. */
/* is x = 0 then pattern expression otherwise object expression. */
struct element *
refal_expression (int x) {

  /* maximum expression size */
# define MAX_EXPRESSION_SIZE 8192
  /* current expression. */
  static struct element curr_exp [MAX_EXPRESSION_SIZE];

  struct element *ret;
  int count, flag;
  char *ch;
  int last_par, prev_par;
  unsigned int size;

  count = -1;
  flag = 1;
  last_par = -1;

  while (flag) {
    if (count < MAX_EXPRESSION_SIZE) count ++;
    switch (token) {
    case '(':
      next_token;
      curr_exp [count].type = LPAR;
      curr_exp [count].body.i = last_par;
      curr_exp [count].number = -1;
      last_par = count;
      break;

    case ')':
      next_token;
      curr_exp [count].type = RPAR;
      curr_exp [count].body.i = last_par;
      if (last_par == -1) rc_serror (100, NULL);
      else {
	if (curr_exp [last_par].type != LPAR)
	  rc_serror (100, NULL);
	curr_exp [count].number = -1;
	/* This condition is added by Shura, but no body. 24.05.99 */
	if (curr_exp [last_par].type == LPAR ||
	    curr_exp [last_par].type == RPAR) {
	  prev_par = curr_exp [last_par].body.i;
	  curr_exp [last_par].body.i = count;
	  last_par = prev_par;
	}
      }
      break;
      
    case LCBRAK:
      if (x) {
	ch = rc_getact ();
	curr_exp [count].type = ACT_LEFT;
	curr_exp [count].body.f = ch;
	curr_exp [count].number = last_par;
	last_par = count;
      } else rc_serror (110, str);
      next_token;
      break;

    case '>':
      if (x) {
	curr_exp [count].type = ACT_RIGHT;
	curr_exp [count].body.i = last_par;
	if (last_par == -1) rc_serror (111, NULL);
	else {
	  if (curr_exp [last_par].type != ACT_LEFT)
	    rc_serror (111, NULL);
	  prev_par = curr_exp [last_par].number;
	  curr_exp [last_par].number = count;
	  curr_exp [count].body.f = curr_exp [last_par].body.f;
	  last_par = prev_par;
	}
	next_token;
      } else {
	flag = 0;
      }

      break;

    case VAR:
      rc_getvar (x, vtype, &(curr_exp [count]));
      curr_exp [count].number = -1;
      next_token;
      break;

	case STRING:
		for (ch = strings; length > 254; length -= 254, ch += 254) {
			char * p;

			p = rc_memral (255);
			strncpy (p, ch, 254);
			p [254] = 0;
			curr_exp [count].body.f = p;
			curr_exp [count].type = STRING;
			curr_exp [count].number = -1;
			if (count < MAX_EXPRESSION_SIZE) count ++;
		}
		if (length == 0) {
			count --;
		} else if (length == 1) {
			curr_exp [count].type = CHAR;
			curr_exp [count].body.c = * ch;
			curr_exp [count].number = -1;
		} else {
			char * p = rc_memral (length + 1);

			strcpy (p, ch);
			curr_exp [count].body.f = p;
			curr_exp [count].type = STRING;
			curr_exp [count].number = -1;
		}
		next_token;
		break;

    case MDIGIT:
    case COMPSYM:
    case ID:
    case NUMBER:
    case ASCIIV:
    /*case STRING:*/
	  getconst (token, & (curr_exp [count]));
      curr_exp [count].number = -1;
      next_token;
      break;

    default: /* empty expression */
      flag = 0;
      break;
    }	/* switch */
  }	/* while */

  /* check that all brackets are closed. */
  if (last_par != -1) rc_serror (12, NULL);

  /* set the last element to zero. */
  if (count >= MAX_EXPRESSION_SIZE) rc_serror (109, NULL);
  /* 359 line. Was ... = NULL. Shura. 29.01.98 */
  curr_exp [count].type = 0;
  curr_exp [count].body.n = -1;
  curr_exp [count].number = -1;
  count ++;

  /* allocate memory and copy there the expression. */
  size = count * sizeof (struct element);
  ret = (struct element *) malloc (size);
  /* 366 line. Not check result of malloc. Shura. 29.01.98 */
  if (NULL == ret) {
    fprintf (stderr, "No memory for refal expression");
    exit (1);
  }
  memcpy ((void *) ret, (void *) curr_exp, size);
  
  /* ... and return the pointer. */
  return ret;
}

char *rc_getact (void) { /* returns a pointer to function name */
	struct functab *p;

	p = searchf (str,fc);   /* search the function calls table.*/
	if (p == NULL) {
		/* insert into the table. */
		p = (struct functab *) rc_allmem (sizeof (struct functab));
		if (p == NULL) {
			fprintf (stderr,"Ran out of memory.\n");
			exit (1);
		}
		p -> next = fc;
		fc = p;

		if (NULL == (p -> name = (char *) malloc (strlen (str) + 1))) {
			fprintf (stderr, "No memory for function name\n");
			exit (1);
		}
		/*copyst (p -> name);*/
		strcpy (p -> name, str);
	}
	return (p -> name);
}

char *getcoms (void) { /* returns a pointer to compound symbol name */
	struct functab *p;

	p = searchf (str,cs);   /* search the compound symbol table.*/
	if (p == 0) {
	   /* insert into the table.  */
		p = (struct functab *) rc_allmem (sizeof (struct functab));
		if (p == NULL) {
			fprintf (stderr,"Ran out of memory.\n");
			exit (1);
		}
		p -> next = cs;
		cs = p;

		if (NULL == (p -> name = (char *) malloc (strlen (str) + 1))) {
			fprintf (stderr, "No memory for function name\n");
			exit (1);
		}
		/*copyst (p -> name);*/
		strcpy (p -> name, str);
		p -> offset = cscount ++;
	}
	return (p -> name);
}

int rc_getvar (mode, vtype, elem)
	int mode;
	char vtype;
	struct element *elem;

	{

	int k;
	char errstr [MAXWS+3];
	int var_type = 0;

	k = searchv (vtype,str);
	if (k == nlv)
		{
		if (mode)
				/* right size: undefined variable. */
			{
			errstr[0] = vtype;
			errstr[1] = '.';
			strcpy (errstr+2,str);
			rc_serror (10, errstr);
			k = -1;
			}
		else
				/* left side: define variable. */
			{
			copyst (lv[k].vindex);
			lv[nlv++].vt = vtype;
			}
		};
	if (vtype == 'E' || vtype == 'e') var_type = E_VAR;
	else if (vtype == 'S' || vtype == 's') var_type = S_VAR;
	else if (vtype == 'W' || vtype == 'w' || vtype == 'T' || vtype == 't')
		var_type = T_VAR;
	elem -> type = var_type;
	elem -> body.i = k;
	return 0;
	}

int getconst (int ct, struct element * elem) {
	unsigned long k;

	switch (ct) {
	case MDIGIT:
		k = atol (str);
		elem -> body.n = k;
		elem -> type = DIGIT;
		break;

	/* Compound symbols are inserted into the table cs.
	 * The pointer to the beginning of the string is returned.
	 */
	case ID:
		elem -> body.f = getcoms ();
		elem -> type = ATOM;
		break;

	case COMPSYM:
		{
			struct functab *p;

			/* p = searchf (str,cs);   / * search the compound symbol table.*/
			p = searchf (strings,cs);
			if (p == 0) {
				/* insert into the table.  */
				p = (struct functab *) rc_allmem (sizeof (struct functab));
				if (p == NULL) {
					fprintf (stderr,"Ran out of memory.\n");
					exit (1);
				}
				p -> next = cs;
				cs = p;

				if (NULL == (p -> name = (char *) malloc (strlen (strings) + 1))) {
					fprintf (stderr, "No memory for function name\n");
					exit (1);
				}
				strcpy (p -> name, strings);
				p -> offset = cscount ++;
			}
			elem -> body.f = p -> name;
			elem -> type = ATOM;
		}
		break;

	case ASCIIV:
		elem -> type = CHAR;
		elem -> body.c = (char) v;
		break;

	case NUMBER:
		elem -> type = DIGIT;
		elem -> body.n = v;
		break;

	default:
		break;
	}
	return 0;
}

int searchv (vt,ind)
	char vt, *ind;
	{
	int i;

	for (i=0; i<nlv; i++)
		if (strcmp (ind,lv[i].vindex) == 0)
			{
			if (vt == lv[i].vt) break;
			else rc_swarn (6);
			};
	return i;
	}

struct functab *searchf (char * fn, struct functab * table) {
	struct functab *fp;

	for (fp = table; fp != NULL; fp = fp -> next) {
		if (strcmp (fn,fp -> name) == 0) return fp;
	}
	/*
	fp = table;
	while (fp != NULL) {
		if (strcmp (fn,fp -> name) == 0) return fp;
		else fp = fp -> next;
	}
	*/
	return NULL;
}

	/* inserts after instrcution. */
int insert_instruction (z, opcode, params)
	struct rasl_instruction *z;
	int opcode;
	union param params;
	{
	struct rasl_instruction *r;

	r = (struct rasl_instruction *) 
		malloc (sizeof (struct rasl_instruction));
	/* 552 line. Not check result of malloc. Shura. 29.01.98 */
	if (NULL == r) {
	  fprintf (stderr, "No memory for rasl-instruction\n");
	  exit (1);
	}
	r -> code = opcode;
	r -> p = params;
	r -> next = z -> next;
	z -> next = r;

	return 0;
	}

int rc_help(void) {
 printf("Using: refc [options] MODULE1 MODULE2 ... \n");
 
 printf("\nOptions recognized by the Refal compiler:");
 printf("\n -l : produce the listing files.");
 printf("\n -h : print this help message.\n\n");
 exit(0);
}

int rc_options(int argc, char * argv []) {
	int i;

	c_flags [0] = '\0';

	for (i = 1; i < argc; i++) {
		if (*argv[i] == '-') {/* c_flags */
			 if ( argv [i][1] == 'l') {
			} else if ( argv [i][1] == 'h' ) { 
                                rc_help(); 
			} else {
				strcat (c_flags,argv[i]);
			}
		}
	}
        return 0;
}


