
# include "rasl.h"
# include "decl.h"
# include "ddecl.h"
# include "macros.h"
# include "dmacro.h"
# include "ifunc.h"
# include "tfunc.h"


# define MDEBUG 0

/* REFAL Interpreter (Tracer version): May 9, 1988. D.T. */

/* For copying of trace informations */
extern FILE * fp_debugInfo;

int ri_inter (void) {
	/* This is the shortest version of the interpreter, containing only those RASL operators which
	 * are produced by the REFAL compiler.
	 */

	/* July, 27, 1985. D.T. */
	/* Some macros have been expanded because the PC compiler can't handle too many macros. (its stack overflows.)
	 * DT July 1 1986.
	 */

	/* Some other macros have been replaced by functions to reduce  the size of object module. March 7 1987. DT. */

	register short n;
	int error;
	char c, ins = 0;
	long mdig, bifnum;
	char *arg;
/*	short m;*/


	error = 0;
restart:
	while (error == 0) {
		ins = *p;

# if MDEBUG
		if (dump_toggle) printf ("%lx: %d\n", p, ins);
# endif

		switch (ins) {
		case ACT1:
			ASGN_CHARP (++p, arg);
			p += sizeof (char *);
			act1 (arg);  
			curk ++;
			break;  

		case BL:
			p++; bl;
			break;

		case BLR:
			p++; 
			ri_blr ();
			break;

		case BR:
			p++; br; 
			break;

		case CL:
			p++; cl;
			break;

		case SYM:
			c = (unsigned char) * ++p;
			p++;
			sym (c);
			break;

		case SYMR:
			c = (unsigned char) * ++p;
			p++;
			symr (c);
			break;

		case EMP:
			p++; emp; break;

		case EST:
			curk --;
			est;
			p = break0;
			break;

		case MULE:
			/*
			n = * ++p;
			++p;
			mule ((int) n);
			*/
			ASGN_LONG (++p, mdig);
			p += sizeof (long);
			mule (mdig);
			break;

		case MULS:
			/*n = * ++p; muls (n); p++; break;*/
			ASGN_LONG (++p, mdig);
			p += sizeof (long);
			muls (mdig);
			break;

		case PLEN:
			p++; plen; p++; break;

		case PLENS:
			p++; plens; break;

		case PLENP:
			p++; plenp; break;

		case PS:
			++p; ps; break;

		case PSR:
			++p;
			psr;
			break;

		case OEXP:
			n = (unsigned char) * ++p;
			++p;
			oexp (n);
			break;

		case OEXPR:
			n = (unsigned char) * ++p;
			++p;
			oexpr (n);
			break;

		case OVSYM:
			n = (unsigned char) * ++p;
			ovsym (n);
			++p;
			break;

		case OVSYMR:
			n = (unsigned char) * ++p;
			ovsymr (n);
			p++;
			break;

		case TERM:
			p++;
			term;
			break;

		case TERMR:
			p++;
			termr;
			break;

		case RDY:
			n = (unsigned char) * ++p;
			++p;
			rdy (n);
			break;

		case SETB:
/*
			n = (unsigned char) * ++p;
			m = (unsigned char) * ++p;
			++p; setb (n,m); break;
*/
			{
				long l_n, l_m;

				ASGN_LONG (++p, l_n);
				p += sizeof (long);
				ASGN_LONG (p, l_m);
				p += sizeof (long); 
				setb (l_n, l_m);
			}
			break;

		case LEN:
			p++;
			len;
			break;

		case LENS:
			c = (unsigned char) *++p;
			p++;
			lens (c);
			break;

		case LENP:
			++p;
			lenp;
			break;

		case SYMS:
			n = (unsigned char) * ++p;
			p++;
			syms (n);
			break;

		case SYMSR:
			n = (unsigned char) * ++p;
			p++;
			symsr (n)
			break;

		case TEXT:
			n = (unsigned char) * ++p;
			p++;
			text (n);
			break;

		case NS:
			c = (unsigned char) * ++p;
			++p;
			ns (c);
			break;

		case TPLE:
			/*
			n = * ++p;
			p++;
			tple (n);  
			break;
			*/
			ASGN_LONG (++p, mdig);
			p += sizeof (long);
			tple (mdig);
			break;

		case TPLS:
			/*
			n = * ++p;
			p++;
			tpls (n);
			break;
			*/
			ASGN_LONG (++p, mdig);
			p += sizeof (long);
			tpls (mdig);
			break;

		case TRAN:
			ASGN_CHARP (++p, arg);
			p += sizeof (char *);
			tran (arg);
			break;

		case VSYM:
			p++;
			vsym;
			break;

		case VSYMR:
			++p;
			vsymr;
			break;

		case OUTEST:
			curk --;
			out (2); 
			est; 
			p = break0;
			break;

		case ECOND:
			if (tel - te + nel + 100 >= size_table_element) {
        			if (fp_debugInfo != NULL) ri_print_error_code(fp_debugInfo,13);
				ri_error (13);
			}
			ASGN_CHARP (++p, arg);
			b = st[sp].b1;
			b = st[sp].b1;
			act1 (arg);
			tel += (teoff = nel);
			est;
			p = break0;
			break;

		case POPVF:
			++p;
			tel -= teoff;
			nel = teoff + 3;
			sp = stoff-1;
			teoff = st[sp].nel;
			stoff = (long) st[sp].b2;
			break;

		case PUSHVF:
			if (sp + 20 >= size_local_stack) {
        			if (fp_debugInfo != NULL) ri_print_error_code(fp_debugInfo,14);
				ri_error (14);
			}
			++p;
			b = tbel (2) -> prec;
			blr;
			pushst (b->prec,b,NULL,NULL);
			sp++;
			pushst (b,stoff,teoff, IMP_);
			b = b -> prec;
			stoff = sp + 1;
			break;

		case STLEN:
			++p;
			sp = stoff;
			break;

		case CSYM:
			ASGN_CHARP (++p, arg);
			p += sizeof (char *);
			csym (arg);
			break;

		case CSYMR:
			ASGN_CHARP (++p, arg);
			p += sizeof (char *);
			csymr (arg);
			break;

		case NSYM:
			ASGN_LONG (++p, mdig);
			p += sizeof (long);
			nsym (mdig);
			break;

		case NSYMR:
			ASGN_LONG (++p, mdig);
			p += sizeof (long);
			nsymr (mdig);
			break;

		case NCS:
			ASGN_CHARP (++p, arg);
			p += sizeof (char *);
			ncs (arg);
			break;

		case NNS:
			ASGN_LONG (++p, mdig);
			p += sizeof (long);
			nns (mdig);
			break;

		/* builtin functions: R.N. - 20 Jul 85 */
		case BUILT_IN: /* a call to a built in function no arguments. */
			curk --;
			ASGN_LONG (p+1, bifnum);
			error = ri_bif (bifnum,NULL);
			p = break0;
			break;

		/* builtin functions with one argument: D.T. - July 27, 1985. */
		case BUILT_IN1:
			/* a call to a function with one argument. */
			/* Arguments are stored before function address. */
			curk --;
			ASGN_CHARP(++p, arg);
			ASGN_LONG (p + (sizeof (char *)), bifnum);
			error = ri_bif (bifnum, arg);
			p = break0;
			break;

		default:
			ri_default (ins, &error);
			break;
		}
	}

	if (error != 0) {
		fprintf (stderr,"RASL instruction:  %4d  at address:  %lx\n", *p, (unsigned long) p);
		if (fp_debugInfo != NULL)
			fprintf (fp_debugInfo,"RASL instruction:  %4d  at address:  %lx\n", *p, (unsigned long) p);
		ri_error(4);
	}

	return 0;
}


/* ri_default () deals with Special RASL instruction for the Tracer only. May 9 1988. D.T. */
int ri_default (int ins, int * error) {
	short n;
	long jk;
	char *arg;

	switch (ins) {
	case EQS:
		ASGN_LONG (++p, jk);
		p += sizeof (long);
		eqs (jk);
		break;

	case CHACT:
		ASGN_CHARP(++p, arg);
		p += sizeof (char *);
		chact (arg);
		break;

	case GT:
		ASGN_LONG (++p, jk);
		p += sizeof (long);
		gt (jk);
		break;

	case LT:
		ASGN_LONG (++p, jk);
		p += sizeof (long);
		lt (jk);
		break;

	case EBR:
		n = (short) (* ++p);
		ebr (n);
		break;

	case CURK:
		ASGN_LONG (++p, jk);
		p += sizeof (long);
		if (jk < curk) {rimp;}
		else rd_endcom ();
		break;

	case NO_BREAK: /* nobreak: debugging. */
		nobreak;
		break;

	default:
		*error = 2;
		break;
	}

restart:
	return 0;
}

int ri_frz (int code) {

	ri_frz1 (code);
	if (break_at_freeze) {
		ebr (-1);
	} else p = break0;
	return 0;
}

extern FILE* rdout;
FILE* ref_err_file (void) {
	return rdout;
}
