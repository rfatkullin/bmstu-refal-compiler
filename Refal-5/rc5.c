
# define DEFINE_EXTERNALS
/* FOR unlink */
#ifdef FOR_OS_LINUX
  #include <unistd.h>
#endif

# include "version.h"
# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"


	/* suppress debugging. */
# define MDEBUG 0

int
main (int argc, char * argv []) {
	struct node *q;
	char * blk_temp;
	int total_errors, i;
	char fname [FILENAME_MAX];
  
	/* print the version and copyright information. */
	printf (VERSION, "Compiler", _refal_build_version);
	printf ("Copyright: Refal Systems Inc.\n");

	rc_initcom (argc, argv); /* see pass2.c */
	i = rc_getbeginfile (argv);
	do { /* added by Shura for Korlukov */
		i ++; /* added by Shura for Korlukov */

		zero.number = 0;	/* create empty record. */
		ftransl = NULL;	/* initialize the translation pointer. */
		Error = rc_mknode (0, zero, zero, zero);

		/* initialize. */
		fc = fb = ft = fx = fe = cs = NULL;
		z = cscount = ntcount = xtcount = btcount = 0L;
		total_errors = 0;

		/* allocate memory for tables: this memory is needed throughout 
		 * the execution of the program -- therefore it is not explicitly released. 
		 * (it is released in the exit() function)
		 */
		if ((nonret = (char *) malloc (MEM_BLK_SIZE)) == NULL) {
			fprintf (stderr,"Can\'t allocate memory.\n");
		    exit (1);
		}
		nrptr = 0;

		/* loop parse, compile. some loop on rc_parse */
		while (1) {
			/* clear */
			rc_initrp ();
			q = rc_parser ();

			if (q == NULL) {
				if (nerrors == 0) break;
			}
			if (nerrors) {
				total_errors += nerrors;
				fprintf (stderr, "%d errors found in function %s\n", nerrors, last_fn);
				continue;
			}

			/* reset the label */
			last_label = 1;
			refcom (q);

			/* free allocated memory. */
			free_tree (q);

			/* free dicardable memory. */
			while (block != NULL) {
				blk_temp = * ((char **) block);
				free ((void *) block);
				block = blk_temp;
			}
		}
		/* if there were errors delete temp files. */
		if (total_errors) {
			fclose (fdref);
			fclose (fdlis);
			fclose (fdtmpw);
			sprintf (fname, "%s.tmp", title);
			unlink (fname);
			fprintf (stderr, "%d syntax errors found.\n", total_errors);
			exit (1);
		} else {
			/* otherwise do the 2nd pass. */
			rc_end ();
			if (nerrors) {
				fprintf (stderr, "%d errors in pass 2.\n", nerrors);
				sprintf (fname, "%s.rsl", title);
				unlink (fname);
				exit (1);
			} else if (!strchr (c_flags, 'l')) {
				/* unlink the unnecessary .lis file. */
				sprintf (fname, "%s.lis", title);
				unlink (fname);
			}
		}
	} while (-1 != (i = rc_getnextfile (i, argv)));/* added by Shura for Korlukov */
	exit (0);
}


int ri_inquire(char * m, char * r, int i_len) {
  *r = '\0';
  while (*r == '\0') {
	int i;

    fprintf(stderr,"%s",m);
    if (NULL == fgets(r, i_len - 1, stdout)) {
      return -1;
    }
    for (i = strlen (r) - 1; i >= 0 && r [i] == '\n'; i --);
    r [i + 1] = 0;
  }
  return 0;
}

/* define input routines. (they are redefined for the tracer). */
/* Puts character C back on the input buffer.  */
/*	July 27 1985. D.T.  */
int rc_ungchar(char c) {
  cbuf[--sc] = c;
  return 0;
}

/* get the next character from the input stream. */
int
rc_gchar (void) {
  /*register*/ int c = 0;

  if (cbuf[sc] == '\n') {
    c = getc(fdref);
    if (c == EOF) return EOF;
    sc = 1;
    fprintf(fdlis,"%5d:   ",++line_no);
    while ((c != '\n') && (c != EOF))  {
      if (sc < MAXSTR) cbuf[sc++] = c;
      putc(c, fdlis);
      c = getc (fdref);
    }
    putc ('\n', fdlis);
    if (sc >= MAXSTR) {
      rc_serror(11,NULL);
      sc = MAXSTR;
    }
    cbuf[sc++] = ' ';  /* pad with 1 blank. 8-4-85. DiTu. */
    cbuf[sc] = '\n';
    sc = 1;
    return '\n';
  } else return cbuf[sc++];
}


/* for debugging only. */
# if MDEBUG
  #ifndef FOR_OS_LINUX
    #include <dos.h> /* This include is strange for UNIX */
  #endif

int print_expr (struct element * e) {
  printf ("Expression: %04x:%04x\n", FP_SEG(e), FP_OFF(e));
  while (e -> type != NULL) {
    switch (e -> type) {
    case NULL:
      return 0;
    case E_VAR:
      printf (" e.%d", e -> body.i);
      break;
    case S_VAR:
      printf (" s.%d", e -> body.i);
      break;
    case T_VAR:
      printf (" t.%d", e -> body.i);
      break;
    case CHAR:
      printf (" t.%d", e -> body.i);
      break;

    case ATOM:
      printf (" %s", e -> body.f);
      break;

    case DIGIT:
      printf (" %ld", e -> body.n);
      break;

    case LPAR:
      printf (" (");
      break;

    case RPAR:
      printf (")");
      break;

    case ACT_LEFT:
      printf (" <%s", e -> body.f);
      break;
      
    case ACT_RIGHT:
      printf (" %s>", e -> body.f);
      break;
      
    case STRING:
      printf (" \"%s\"", e -> body.f);
      break;

    default:
      printf (" (TYPE=%d)", e -> type);
    } /* end of switch */
    e ++;
  } /* end of while */
  printf ("\n");
  return 0;
}

int print_tree (struct node * q) {
  static int id = 0;
  int save_id;
  int t;

  t = q -> nt;
  switch (t) {
  case ENTRY:
  case FDEF:
    save_id = id;
    printf ("node %d entry/fdef=%d func name %s\n", id ++, t, q -> a2.pchar);
    print_tree (q -> a3.tree);
    printf ("end node %d entry/fdef=%d func name %s\n",
	    save_id, t, q -> a2.pchar);
    break;

  case SENTS:
    save_id = id;
    printf ("node: %d sentences: first sentence:\n", id++);
    print_tree (q -> a2.tree);
    if (q -> a3.tree != NULL) {
      printf ("node: %d other sentences:\n", save_id);
      print_tree (q -> a3.tree);
    } else {
      printf ("node: %d No other sentences:\n", save_id);
    }
    break;

  case ST:
    save_id = id;
    printf ("node %d sentence: Left side:\n", id ++);
    print_expr (q -> a2.chunk);
    printf ("node %d right side tail:\n", save_id);
    print_tree (q -> a3.tree);
    break;

  case RCS1:
    save_id = id;
    printf ("node %d sentence: Right side:\n", id ++);
    print_expr (q -> a2.chunk);
    break;

  case RCS4:
  case RCS3:
  case RCS2:
  case RSF1B:
  case RSF1:
  case RSF: 
  case LSF1:
  case LSF: 
    printf ("Complex sentence: %d\n", t);
    break;

  default:
    printf ("Strange Node %d\n", t);
    break;
  } /* end of switch */
  return 0;
}

# define NUM_OF_RASL_INSTR 53

char *mnemonics [NUM_OF_RASL_INSTR] = {
  "ACT1", "BL", "BLR", "BR", "CL", "SYM", "SYMR", "EMP", "EST",
  "MULE", "MULS", "PLEN", "PLENS", "PLENP",	"PS", "PSR", "OEXP",
  "OEXPR", "OVSYM", "OVSYMR", "TERM", "TERMR", "RDY", "SETB",
  "LEN", "LENS", "LENP", "LENOS", "SYMS", "SYMSR", "TEXT", "NS",
  "TPLE", "TPLS", "TRAN", "VSYM", "VSYMR", "OUTEST", "ECOND",
  "POPVF", "PUSHVF", "STLEN", "CSYM", "CSYMR", "NSYM", "NSYMR",
  "NCS", "NNS", "LBL", "L", "E", "LABEL", "Special Mark B"
};

int rasl_numbers [NUM_OF_RASL_INSTR] = {
  ACT1, BL, BLR, BR, CL, SYM, SYMR, EMP, EST,
  MULE, MULS, PLEN, PLENS, PLENP,	PS, PSR, OEXP,
  OEXPR, OVSYM, OVSYMR, TERM, TERMR, RDY, SETB,
  LEN, LENS, LENP, LENOS, SYMS, SYMSR, TEXT, NS,
  TPLE, TPLS, TRAN, VSYM, VSYMR, OUTEST, ECOND,
  POPVF, PUSHVF, STLEN, CSYM, CSYMR, NSYM, NSYMR,
  NCS, NNS, LBL, L, E, LABEL, B
};

char unknown [10];

char *rasl_code (int id) {
  int i;

  for (i = 0; i < NUM_OF_RASL_INSTR; i ++)
    if (id == rasl_numbers[i]) return mnemonics[i];
  sprintf (unknown, "RASL-%d", id);
  return unknown;
}

int print_translation (struct rasl_instruction * r) {

  while (r != NULL) {
    print_rasl_inst (r);
    if (r -> code == NULL) break;
    r = r -> next;
  }
  return 0;
}

int print_rasl_inst (struct rasl_instruction *r) {
  int rasl_ins;

  rasl_ins = r -> code;
  switch (rasl_ins) {
  case NULL:
    printf ("end of translation\n");
    break;

  case B:
    printf ("Special Mark: B\n");
    break;

  case L:
  case E:
    printf ("Label %s\n", r -> p.f);
    break;

  case LABEL:
    printf ("Label: #%d\n", r -> p.i);
    break;

  case LBL:
    printf ("LBL: #%d\n", r -> p.i);
    break;

  case ECOND:
    printf ("   %s LABEL#%d\n", rasl_code (rasl_ins), r -> p.i);
    break;

    /* takes 2 integer arguments. */
  case SETB:
    printf ("   %s %d %d\n", rasl_code (rasl_ins), r -> p.d.i1, r -> p.d.i2);
    break;

    /* these instructions take 1 integer argument. */
  case OEXP: case OEXPR: case OVSYM: case OVSYMR: case TRAN:
  case TPLE: case TPLS: case MULE: case MULS: case RDY:
    printf ("   %s %d\n", rasl_code (rasl_ins), r -> p.i);
    break;

    /* these instructions take 1 long number argument. */
  case NSYM: case NSYMR: case NNS:
    printf ("   %s %lu\n", rasl_code (rasl_ins), r -> p.n);
    break;
    
    /* these instructions take 1 character pointer argument. */
  case CSYM: case CSYMR: case NCS: case ACT1:
    printf ("   %s %s\n", rasl_code (rasl_ins), r -> p.f);
    break;
    
    /* these instructions take 2 arguments: length and pointer. */
  case SYMS: case SYMSR: case TEXT:
    printf ("   %s %d %s\n", rasl_code (rasl_ins), strlen (r -> p.f),
	    r -> p.f);
    break;

    /* these instructions take 1 character argument. */
  case SYM: case SYMR: case NS: case LENS:
    printf ("   %s %c\n", rasl_code (rasl_ins), r -> p.c);
    break;
    
    /* these instructions take no arguments. */
  case CL: case EMP: case VSYM: case VSYMR: case BL: case BLR:
  case BR: case LEN: case TERM: case TERMR: case PS: case PSR:
  case OUTEST: case PUSHVF: case POPVF: case STLEN: case PLENS:
  case PLENP: case PLEN: case LENP:

    printf ("   %s\n", rasl_code (rasl_ins));
    break;

  default:
    printf ("UNKNOWN: %d\n", rasl_ins);
  } /* end of switch */
  return 0;
}

int print_var_table (unsigned char * bits) {
  /*register*/ int i;

  printf ("Table of variables: %d\n", table_len);
  for (i = 0; i < table_len; i ++) {
    printf ("%2d: index = %2d offset = %2d checked = %1d\n",
	    i, table [i].index, table [i].te_offset, is_bit_checked (bits, i));
  }
  printf ("\n");
  return 0;
}

int prftab (struct functab * table) {
  for (; table != NULL; table = table -> next) 
    printf ("%s %ld\n", table -> name, table -> offset);
  printf ("\n");
  return 0;
}

int print_holes (struct HOLES * holes, struct element * e) {
  struct HOLES *h;
  int i = 0;

  printf ("The list of holes\n");
  h = holes;
  while (h != NULL) {
    i ++;
    printf ("%d  [%2d, %2d] - %1d, %1d\n", i, h -> left, h -> right,
	    no_lengthening (e + (h -> left)),
	    no_lengthening (e + (h -> right)));
    h = h -> next;
  }
  printf ("\n");
  return 0;
}

# endif

