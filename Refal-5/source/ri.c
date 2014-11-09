
# include "rasl.h"
# include "decl.h"
# include "macros.h"
# include "ifunc.h"


# define MDEBUG 0

# if MDEBUG
	FILE *fp_dbg;
# endif

  /* For translator pgraph */
extern unsigned long ul_local_calls;
extern int flag_local_calculation;


	/*  REFAL Interpreter: June 6, 1986.   D.T.	*/

int ri_inter ()
	{

	/* This is the shortest version of the interpreter, 
		containing only those RASL operators which
		are produced by the REFAL compiler.		*/
	/*		July, 27, 1985.	D.T.		*/
	/* Some macros have been expanded because the PC compiler
		can't handle too many macros. (its stack overflows.)
		DT July 1 1986. */
	/* Some other macros have been replaced by functions to reduce 
		the size of object module. March 7 1987. DT. */

		/*register*/ short n;
		/*register*/ int error; 
		char c;
		long bifnum, mdig;
		char *arg;
/*		short m;*/


	error = 0;

# if MDEBUG
	fp_dbg = fopen ("refgo.dbg", "wt");
	fprintf (fp_dbg, "Initial view field:\n");
	deb_prexp (fp_dbg, vf, vfend);
	fprintf (fp_dbg, "\n");
# endif

restart:
	while (1)
		{

		/*printf ("Instruction = %d\n", * p);*/
# if MDEBUG
	printf ("DEBUG POINTER 1\n");
	fprintf (fp_dbg, "p = %ld, instruction = %d View Field:\n", p, *p);
	deb_prexp (fp_dbg, vf, vfend);
	fprintf (fp_dbg, "\n");
	printf ("DEBUG POINTER 2\n");
# endif
		switch (*p)
			{
			case ACT1:
				ASGN_CHARP (++p, arg);
				p += sizeof (char *);
				act1 (arg);

				/* For translator pgraph */
				ul_local_calls ++;

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
				est;
# if MDEBUG
				fprintf (fp_dbg, "end of step %ld\n", nst);
				{
					char * cp_w;
					
					for (cp_w = p - 2; * cp_w != 0; cp_w --);
					cp_w ++;
					fprintf (fp_dbg, "active function is %s\n", cp_w);
				}
				printf ("DEBUG POINTER 3\n");
				deb_prexp (fp_dbg, vf, vfend);
# endif
				/* FOR DEBUG * /
				{
					char * cp_w;

					for (cp_w = b1 -> pair.f - 2; * cp_w != 0; cp_w --);
					cp_w ++;
					fprintf (stderr, "Active func: %s\n", cp_w);
				}*/
				break;

			case MULE:
				/*n = * ++p;
				++p;
				mule ((int) n);*/
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
				/*n = * ++p;
				p++;
				tple (n);  
				break;*/
				ASGN_LONG (++p, mdig);
				p += sizeof (long);
				tple (mdig);
				break;

			case TPLS:
				/*n = * ++p;
				p++;
				tpls (n);
				break;*/
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
				out (2);
				est;

				/* For translator pgraph */
				ul_local_calls --;
				if (flag_local_calculation) {
					if (ul_local_calls == 0) {
						/* local function is calculated */
						return 0;
					}
					/* else the local function isn't calculated */
				}
# if MDEBUG
				fprintf (fp_dbg, "end of step %ld\n", nst);
				{
					char * cp_w;

					for (cp_w = p - 2; cp_w != 0; cp_w --);
					cp_w ++;
					fprintf (fp_dbg, "active function is %s\n", cp_w);
				}
				printf ("DEBUG POINTER 5\n");
				deb_prexp (fp_dbg, vf, vfend);
# endif
				/* FOR DEBUG * /
				{
					char * cp_w;

					for (cp_w = b1 -> pair.f - 2; * cp_w != 0; cp_w --);
					cp_w ++;
					fprintf (stderr, "Active func: %s\n", cp_w);
				}*/
				break;

			case ECOND:
				if (tel - te + nel + 100 >= size_table_element)	ri_error(13);
				ASGN_CHARP (++p, arg);
				b = st[sp].b1;
				act1 (arg);
				tel += (teoff = nel);
				est;
# if MDEBUG
				fprintf (fp_dbg, "end of step %ld\n", nst);
				{
					char * cp_w;

					for (cp_w = p - 2; * cp_w != 0; cp_w --);
					cp_w ++;
					fprintf (fp_dbg, "active function is %s\n", cp_w);
				}
				printf ("DEBUG POINTER 7\n");
				deb_prexp (fp_dbg, vf, vfend);
# endif
				/* FOR DEBUG * /
				{
					char * cp_w;

					for (cp_w = b1 -> pair.f - 2; * cp_w != 0; cp_w --);
					cp_w ++;
					fprintf (stderr, "Active func: %s\n", cp_w);
				}
				printf ("DEBUG POINTER 8\n"); */
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
				if (sp + 20 >= size_local_stack) ri_error (14);
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
				ASGN_LONG (p+1, bifnum);
				/* FOR DEBUG * /
				fprintf (stderr, "built function: %d\n", bifnum); */

				error = ri_bif (bifnum,NULL);
				if (error) goto error_handler;

				/* For translator pgraph */
				ul_local_calls --;
				if (flag_local_calculation) {
					if (ul_local_calls == 0) {
						/* local function is calculated */
						return 0;
					}
					/* else the local function isn't calculated */
				}

				break;

					/* builtin functions with one argument: D.T. - July 27, 1985. */
			case BUILT_IN1: /* a call to a function with one argument. */
				/* Arguments are stored before function address. */
				ASGN_CHARP(++p, arg);
				ASGN_LONG (p + (sizeof (char *)), bifnum);
				/* FOR DEBUG * /
				fprintf (stderr, "built function: %d\n", bifnum);*/

				error = ri_bif (bifnum, arg);
				if (error) goto error_handler;

				/* For translator pgraph */
				ul_local_calls --;
				if (flag_local_calculation) {
					if (ul_local_calls == 0) {
						/* local function is calculated */
						return 0;
					}
					/* else the local function isn't calculated */
				}

				break;

			default:
				error = 2;
				goto error_handler;
			}
		}

 error_handler:

# if MDEBUG
	fprintf (fp_dbg, "Error p = %ld *p = %d\n", p, *p);
	fclose (fp_dbg);
# endif

	if (error != 0)
		{
		fprintf (stderr,"RASL instruction:  %4d  at address:  %lx\n", *p,
			(unsigned long) p);
		ri_error(4);
		}

# if MDEBUG
	fclose (fp_dbg);
# endif

	return 0;
	}

int ri_frz (code)
	int code;
	{
	return ri_frz1 (code);
	}




# if MDEBUG

# ifdef IBM370
# include "freeze.h"
# else
# include "freeze.h"
# endif

	/*  REFAL SYSTEM: Auxiliary functions to
		print view field. June 12, 1986.   D.T.	*/

int deb_prexp (fp, q1, q2)
	FILE * fp;
	LINK *q1, *q2;
	{
		LINK *j;

		if (q1 == NULL) return 0;
		if (q2 == NULL) deb_prlnk(fp, q1);
		else for (j = q1; j != q2; j = j -> foll)
			if (j) deb_prlnk(fp, j);
			else
				{
				fprintf (fp, "NULL ptr. aborted.\n");
				break;
				};
		deb_prlnk(fp, q2);
		return 0;
	}

int deb_prlnk (fp, lk) /* print a link */
	FILE *fp;
	LINK *lk;
	{

		if (lk == NULL)
			{
			fprintf(fp, "Link: NULL.\n");
			return 0;
			}

		fprintf(fp, "link: %ld. Type=%d, Foll=%ld, Prec=%ld. ",
			lk, lk->ptype, lk->foll, lk->prec);

		switch (lk -> ptype)
			{
			case 0:
			case 1:
			case 6:

		 		fprintf(fp, "Pair=%ld\n", lk -> pair.b);
				break;

			case 2:

				fprintf(fp, "Symb=[%c], or %d\n",lk->pair.c,lk->pair.c);
				break;

			case 3:

				fprintf (fp, "Csym=%.32s\n", lk -> pair.f);
				break;

			case 4:

				fprintf(fp, "Numb=%lu\n",lk -> pair.n);
				break;

			case 5:

				fprintf(fp, "Func=");
				{
					char * cp_w = lk -> pair.f - 1;

					for (cp_w --; * cp_w != '\0'; cp_w --);
					cp_w ++;
					ri_actput (cp_w, fp);
				}
				/*ri_actput(lk->pair.f, fp);*/
				fprintf(fp, "\n");
				break;

			case 7:
			case 8:

				fprintf (fp, "lev=%2lu, ind=%lu\n",
					level_of (lk -> pair.n), index_of (lk -> pair.n));
				break;

			default:
				fprintf(fp, "STRANGE\n");
				break;
			};

		return 0;
	}
# endif

FILE *ref_err_file()
	{
	return stderr;
	}

