
# include "decl.h"
# include "macros.h"
# include "freeze.h"
# include "ifunc.h"


# define MDEBUG  0
#if MDEBUG
int prexp(LINK *, LINK *);
int prlnk(LINK *);
#endif

# if MDEBUG
# define FP_SEG(fp) (*((unsigned short *)&(fp) + 1))
# define FP_OFF(fp) (*((unsigned short *)&(fp)))
# endif


static int freezer_number = 0;

static char freezer_function_def [MAXWS + 5] = "\0Freezer";


int is_freezer (LINK *q) {
	char *s;
	long *l;

	/* If q-ptype != LINK_TYPE_LACT. No function. */
	if (!LINK_LACT (q)) return 0;
	s = q->pair.f;
	if (*s != 100) return 0; /* unbuilt function, != BUILT_IN */
	l = (long *) (s+1);	/* now must point to 46 */
	if (*l != 46L) return 0; /* != TPLS */
	else return 1; /* 1: Function (LINK_TYPE_LACT), ptr to BUILT_IN, TPLS */
}

		/* test if there is a freezer. */
int exists_freeze (void) {
	return freezer_number > 0;
	
# if 0
		LINK *q;

		if ((q = nextp) == vfend) return 0;
		while (!is_freezer (q->pair.b))
			{
			q = q->foll->prec;
			if (q == vfend) return 0;
			}
		/* Freezer found. */
		return 1;
# endif
}


/* metacode_make_var() makes a variable out of index, level and elevation */
/*static unsigned long metacode_make_var(unsigned long, unsigned long, unsigned long);*/
static unsigned long metacode_make_var(unsigned long index, unsigned long level, unsigned long elev) {
	return (INDEX_MASK & index) | (ELEV_MASK & (elev << 16)) | (LEVEL_MASK & (level << 24));
}

/* check if an expression contains free variables. */
int contains_vars (LINK *left_bound, LINK *right_bound) {
	register LINK *q;

	q = left_bound; /* Only if left_bound < right_bound */
	do {
		if (LINK_VAR(q)) return 1; /* q is variable */
		q = q->foll;
	} while (q != right_bound);
	return 0;
}
		
int rf_frz (void) {
	/* 1. Check the insides. */
	emp; 

	/* 2. excise the <FREEZE> function call. */
	rdy(0);
	out(2);

	/* 3. call ri_frz */
	ri_frz (1);
	return 0;
	
restart:
	ri_error(5);
	return 0;
}

int ri_frz1 (int code) {
	LINK *q, *j;

#if MDEBUG
	printf("Freeze occurs: view field is:\n");
	prexp(vf, vfend);
#endif

	/* 1. find the next active <FREEZER >  */
	q = nextp;
	while (!is_freezer (q->pair.b)) {
		/* Seem that active point isn't ring.  */
		j = q->foll->prec; /* ..q --> (*) --> ||..(q) <--> (*) -->  */
	    q->foll->prec = q; /*   .. j <_|      ||.. q = j            */   
	    q = j;
	    curk --;
	    if (q == vfend) ri_error(3); /* Don't find freezer */
	}
	/* ^-- That it may loop itself, when j == q and q != end of view
	 * field and freezer. It is a potentional error. Shura. 10.02.98
	 */
	/* Freezer found. */
#if MDEBUG
printf("freezer = "); prlnk(q);
#endif

	/* 2. fix the next function call. */
	curk --;
	nextp = q->foll->prec;
	q->foll->prec = q;

	/* 3. excise the <FREEZER function call */
	j = q->pair.b;
	b = j->prec;
	rend = j->foll;

	j->foll = lfm;

	lfm = j;

	weld (b, rend);

	/* 4. Insert code macrodigit  */
	nns ((long) code);

	/* 5. excise the > of the <FREEZER */
	rend = q->foll;
	q->foll = lfm;
	lfm = q;
	q = q->prec;
	weld (q, rend);

	/* 6. Downgrade the expression. */
	if (q != b) {
		/* expression between b and q is not empty. */
		ri_dn (b->foll, q);
	}

	/* 7. End of step. */
	est;

	/* 8. decrease the count. */
	freezer_number --;		

	/* 9. return. */
	return 0;
}
			
int rf_dn (void) {
	cl;

	/* Clean up the function call. */
	rdy (4);
	out (2);
	rdy (0);
	out (1);	
		/* call ri_dn for the insides. */

	if (tbel (3) != NULL) ri_dn (tbel (3), tbel (4));
	/* there are no function calls inside, so we don't have to worry about fixing the stack. */
	est;
	return 0;
}


int ri_dn (LINK *q1, LINK *q2) {
	LINK *r, *rend, *j;
	char *s, c;
	unsigned long variable, index, level, elev;

#if MDEBUG
	printf("ri_dn() 1: "); prlnk(q1);
	printf("ri_dn() 2: "); prlnk(q2);
#endif

	r = q1;
	do {
    	q1 = r;
    	c = r->ptype;
    	/* metacode the expression. */
	    switch(c) {
		/* (e.x) ==> (* mu(e.x)) */
		case LINK_TYPE_LSTRUCTB: /* + */
			rend = NEXT(r);
			weld(r, lfm);
			r = lfm;
			all;
			weld(r, rend);
			r->ptype = LINK_TYPE_CHAR;
			r->pair.c = META_BRACKET;
			r = rend;
			break;

		/* <f e.x> ==> (! f mu(e.x)) */
		case LINK_TYPE_LACT: /* + */
			{
				char * cp_w = r -> pair.f;
				for (cp_w = cp_w - 2; * cp_w != '\0'; cp_w --);
				cp_w ++;
				/*s = ri_cs_impl(r->pair.f - MAXWS);*/
				s = ri_cs_impl(cp_w);
			}
			/* leave this link as is, and insert 2 symbols after it. */
			rend = NEXT(r);
			weld(r, lfm);
			r = lfm;
			all;
			r->ptype = LINK_TYPE_CHAR;
			r->pair.c = META_ACTIVE;
			weld(r, lfm);
			r = lfm;
			all;
			r->ptype = LINK_TYPE_COMPSYM;
			r->pair.f = s;
			weld(r, rend);
			r = rend;
			break;

		/* link this link with the pair link. */
		case LINK_TYPE_RACT: /* + */
      		j = PAIR(r);
      		con(j, r);
			/* fix the linked list structure. */
      		r -> foll -> prec = r;
			r = r -> foll;
			break;

		/* [ex *] ==> [dn (ex)] */
		/* [ex -] ==> (dn (ex)) */
		/* ['e'|'s'|'t' e.Exp] ==> ('e'|'s'|'t' dn (e.Exp)) */
		case LINK_TYPE_MLSTRUCTB:
			{
				LINK * lnk_nxt = r -> foll, * lnk_w = r -> pair.b -> prec;

				if (! LINK_CHAR (lnk_w)) {
					ri_error (6);
					return (1);
				}
				if (META_BRACKET != lnk_w -> pair.c && META_FIRST_BRACKET != lnk_w -> pair.c) {
					ri_error (6);
					return (1);
				}
				if (META_BRACKET != lnk_w -> pair.c) {
					/* ['e'|'s'|'t' e.Exp '-'] ==> ('e'|'s'|'t' e.Exp) */

					if (! LINK_CHAR (lnk_nxt)) {
						/* [e.Exp '-'] ==> I don't know that to do */
						/* return its unknowing about the situation */

					} else if (META_SVAR != lnk_nxt -> pair.c && META_EVAR != lnk_nxt -> pair.c && META_TVAR != lnk_nxt -> pair.c) {
						/* Uncorrect type of parametr. Error */
						ri_error (6);
						return (1);
					} else {
						/* ['e'|'s'|'t' e.Exp '-'] = ('e'|'s'|'t' e.Exp) */
						r -> ptype = LINK_TYPE_LSTRUCTB;
						r -> pair.b -> ptype = LINK_TYPE_RSTRUCTB;
					}
				} /* else : [e.Exp '*'] ==> [e.Exp] */

				/* Delete a symbol of '*' or '-' before closing meta-bracket ']' */
				lnk_w -> foll -> prec = lnk_w -> prec;
				lnk_w -> prec -> foll = lnk_w -> foll;
				weld (lnk_w, lfm);
				lfm = lnk_w;
				r = lnk_nxt;
			}
			break;

		/* variable ==> (s.type s.index s.elev) ... or lower the level. */
		/* [ex *] ==> [dn (ex)] */
		/* [ex -] ==> (dn (ex)) */
		/* ['e'|'s'|'t' e.Exp] ==> ('e'|'s'|'t' dn (e.Exp)) */
		/* [s.T ex  ] ==> ??? --- problem of Inreffs */
		case LINK_TYPE_SVAR: /* This is s-parametr, no s-variable */
		case LINK_TYPE_EVAR: /* This is e-parametr, no e-variable */
			variable = r->pair.n;
			index = index_of(variable);
			level = level_of(variable);
			elev = elevation_of(variable);
			rend = NEXT(r);
			/* see if we get away with simple reduction of level. */
			if (level > 1) {
				/* [ex *] ==> [dn (ex)] */
				level -= 1;
				r->pair.n = metacode_make_var(index, level, elev);
			} else {
				/* [ex -] ==> (dn (ex)) */
				/* ['e'|'s'|'t' e.Exp] ==> ('e'|'s'|'t' dn (e.Exp)) */
				/* [s.T ex  ] ==> ??? --- problem of Inreffs */
				/* otherwise insert a chain of symbols. */
				j = r;  /* save this parenthesis. */
				/* next symbol: type. */
				weld(r, lfm);
				r = lfm;
				all;
				r->ptype = LINK_TYPE_CHAR;
				if (c == LINK_TYPE_SVAR) r->pair.c = META_SVAR;
				else r->pair.c = META_EVAR;
				/* next symbol: index. */
				weld(r, lfm);
				r = lfm;
				all;
				r->ptype = LINK_TYPE_NUMBER;
				r->pair.n = index;
				/* next symbol: elevation. */
				if (elev != MAX_VAR_ELEV) {
					weld(r, lfm);
					r = lfm;
					all;
					r->ptype = LINK_TYPE_NUMBER;
					r->pair.n = elev;
				}
				/* next symbol: closing bracket. */
				weld(r, lfm);
				r = lfm;
				all;
				con(j, r);
				weld(r, rend);
			}
			r = rend;
			break;
 
		/* otherwise leave as is. */
		default: r = r -> foll; break;
		}
	} while (q1 != q2);
	return 0;
}

/* Ev-met = rf_setfrz */
int rf_setfrz (char ** muaddr) {
	LINK * lnkp_bWork;

	cl;
	/*rdy (0); Shura */
	/* Where b is boundary. tel is ptr to an entry in the table of elements */

	/* bl; Shura */

	/*if (tel [3] == NULL) tel [4] = b;
	else tple (4);
	out (2);  Shura */
	/*rdy (0);*/

	/* call ri_up for the insides. */
	if (tel [3] != NULL) {
		ri_up (tel [3], tel [4], muaddr);
		lnkp_bWork = b;
		rdy (0);
		bl;
		if (tel [3] == NULL) tel [4] = b;
		else tple (4);
	} else {
		bl;
	}
	out (2);

	/* there are no function calls inside, so we don't have to worry about fixing the stack. */

	/* write in definition of FREEZER function. */
	{
  		char * cp_w = freezer_function_def + 1;
  	
  		cp_w = cp_w + strlen (cp_w) + 1;
  		/*freezer_function_def [MAXWS] = 100;*/
  		/*wrlong_to_mem (46L, (freezer_function_def + MAXWS + 1));*/
  		cp_w [0] = 100;
  		wrlong_to_mem (46L, (cp_w + 1));

		/* finish the freezer function. */
		br;

  		act1 (cp_w);
	}
	freezer_number ++;
	est;
	return 0;
}
		

int rf_up (char ** muaddr) {

	cl;
	/* Clean up the function call. */
	rdy (4);
	out (2);
	rdy (0);
	out (1);	

	/* call ri_up for the insides. */
	if (tbel (3) != NULL) ri_up (tbel (3), tbel (4), muaddr);
	/* there are no function calls inside, so we don't have to worry about fixing the stack. */
	est;
	return 0;
}

	
/* ri_up () ups (increases the metacode level of) an expression.
 * Note: it set global variable b to the last link in the expression. d.t. 12-2-89
 */
/* pointer to the beginning of the list of functions visible from this module. */
int ri_up (LINK *q1, LINK *q2, char **muaddr) {
  LINK *r, *rp, *j, *j1, *rend, *x1, *x2;
  char c, *s;
  unsigned long variable, index, level, elev;
  char f_name [MAXWS+1];

	r = q1;
	do {
		q1 = r;
		/* if it is a variable, increase its level ... */
		if (LINK_VAR(r)) {
			/* [e.Exp] ==> [up (e.Exp) '*'] */
			variable = r->pair.n;
			index = index_of(variable);
			level = level_of(variable);
			elev = elevation_of(variable);
			level ++;
			r->pair.n = metacode_make_var(index, level, elev);
			r = NEXT(r);
		} else if (LINK_LSTRUCTB(r)) {
			/* if it is an expression in parentheses apply metacoder */
			rp = r -> pair.b;
			j = r -> foll;
			if (!LINK_CHAR(j)) {
      			/* May be we need to modify the block of the code. See below */
				ri_error(6);
				return 1;
			}
			c = j->pair.c;
      
			/* perform de-metacode translation */
			/* ('*' e.Exp) ==> (up (e.Exp))
			 * ('!' (FN s.F e.Arg)) ==> <s.F e.Arg>
			 * ('e'|'s'|'t' e.Exp) ==> ['e'|'s'|'t' up (e.Exp) '-']
			 * ('e'|'s'|'t' 5 1) ==> ['e'|'s'|'t' 5 1 -]
			 * (s.T ....) ==> ??? --- may be this is problem of Inreffs.
			 * [e.Exp] ==> [e.Exp '*']
			 */
			switch(c)	{
	
			/* activize the expression inside. */
			/* (! f e.x) ==>  <f de-mu(e.x)> */
			case META_ACTIVE:	
      			j1 = j -> foll;

				/* check that it is either a compound symbol or a special symbol (one of "+/-*%") */
				if (LINK_COMPSYM(j1)) {
					/*
					char * cp_w = j1 -> pair.f;
					for (cp_w = cp_w - 2; * cp_w != '\0'; cp_w --);
					cp_w ++;
					*/
					strncpy(f_name, j1->pair.f, MAXWS);
					/*strcpy(f_name, cp_w);*/
				} else if (LINK_SYMBOL(j1) && mu_special(f_name, j1->pair.c) == 0) {
					;
				} else {
					ri_error(6);
					return 1;
				}
				s = mu_find (f_name, muaddr);
				if (s == NULL) s = mu_find (f_name, entry_functions);
				if (s == NULL) {
					ri_error(5);
					return 1;
				}
				/* remove the list between j and j1. */
				j1 = j1 -> foll;
				j1 -> prec -> foll = lfm;
				weld(r, j1);
				/*j -> foll = lfm;*/
				lfm = j;
		
				/* delay activation of the brackets r, rp */
				/* activation will occur when r reaches rp, see below */
				r -> ptype = LINK_TYPE_LACT;
				r -> pair.f = s;
				rp -> ptype = LINK_TYPE_RACT;
				r = r -> foll;
				break;

			/* put parentheses around the expession. */
			/* (* e.x) ==> (de-mu(e.x)) */
			case META_BRACKET:
      			/* simply remove the current symbol and leave the old parentheses in place. */
      			j1 = j -> foll;
      			weld(r, j1);
      			r = j1;
				NEXT(j) = lfm;
				lfm = j;
				break;

			/* quote the expression inside. */
			/* (m e.x) ==> e.x */
			case META_QUOTE:
				x1 = NEXT(j);
				rend = NEXT(rp);
				/* check that expression inside is not empty */
				if (x1 == rp) {
					j = PREV(r);
					weld(j, rend);
					NEXT(rp) = lfm;
					lfm = r;
				} else {
					x2 = PREV(rp);
					/* link all garbage together. */
					NEXT(j) = rp;
					NEXT(rp) = lfm;
					lfm = r;
					/* transplant expr [x1,x2] between PREV(r) and rend */
					j = PREV(r);
					weld(x2, rend);
					weld(j, x1);
				}
				r = rend;
				break;

			/* create a variable */
			/* (s.type s.index s.elev) ==> variable */
			/* ('e'|'s'|'t' e.Exp) ==> ['e'|'s'|'t' e.Exp '-'] */
			case META_TVAR:
			case META_SVAR:
			case META_EVAR:
				/* This is needed for short time. */
				{
					LINK * lnk_w = lfm;

					all;
					/* Insert a META_FIRST_BRACKET = '-' */
					weld (rp -> prec, lnk_w);
					weld (lnk_w, rp);
					/* Set META_FIRST_BRACKET */
					lnk_w -> ptype = LINK_TYPE_CHAR;
					lnk_w -> pair.c = META_FIRST_BRACKET;
					r -> ptype = LINK_TYPE_MLSTRUCTB;
					rp -> ptype = LINK_TYPE_MRSTRUCTB;
					r = NEXT (j);
				}
				/* parse the expression inside. */
				/*
				 *j1 = NEXT(j);
				 *if (!LINK_NUMBER(j1)) {
				 * ri_error(6);
				 * return 1;
				 *}
				 *index = j1->pair.n;
				 *j1 = NEXT(j1);
				 *elev = -2;
				 *if (j1 == rp) elev = MAX_VAR_ELEV;
				 *else if (LINK_NUMBER(j1)) elev = j1->pair.n;
				 *else {
				 * ri_error(6);
				 * return 1;
				 *}
				 *level = 1;
				 */
				/* check index and elevation here !!! */
				/*variable = metacode_make_var(index, level, elev);*/
				/* remove all garbage and insert a variable here. */
				/*rend = NEXT(rp);
				 *NEXT(rp) = lfm;
				 *lfm = j;
				 *if (c == META_SVAR) r->ptype = LINK_TYPE_SVAR;
				 *else r->ptype = LINK_TYPE_EVAR;
				 *r->pair.n = variable;
				 *weld(r, rend);
				 *r = rend;
				 */
				break;

			/* otherwise error */
			default:
				ri_error(6);
				return 1;
			}
		} else if (LINK_RACT(r)) {
			/* activate the delayed function call see case META_ACTIVE above */
			precp -> foll -> prec = r;
			precp = r;
			curk++;
			r = NEXT(r);
		} else if (LINK_MLSTRUCTB (r)) {
			/* [e.Exp] ==> [e.Exp '*'] */
			LINK * lnk_rb = PAIR (r), * lnk_ins = lfm;

			all;
			/* Set META_BRACKET */
			lnk_ins -> ptype = LINK_TYPE_CHAR;
			lnk_ins -> pair.c = META_BRACKET;
			/* Insert META_BRACKET before a square bracket */
			weld (lnk_rb -> prec, lnk_ins);
			weld (lnk_ins, lnk_rb);
			r = NEXT (r);
		} else {
			/* otherwise leave the rest alone. */
			r = NEXT(r);
		}
	} while (q1 != q2);
	b = PREV(r);
	return 0;
}
					
int rf_frzr (void) {

	cl;
	/* downgrade the expression. */
	if (tbel (3) != NULL) {
		ri_dn (tbel(3), tbel(4));
		/* reset the bounds the execute the macro cl() again. */
		setb (1, 2);
		nel = 3;
		cl;
	}

	rdy(0);
	nns (0L);
	tple (4);
	out(2);
	est;
	freezer_number --;
	return 0;
}

# if MDEBUG

	/*  REFAL SYSTEM: Auxiliary functions to
		print view field. June 12, 1986.   D.T.	*/

int prexp(q1,q2)
	LINK *q1, *q2;
	{
		LINK *j;

		if (q1 == NULL) return 0;
		if (q2 == NULL) prlnk(q1);
		else for (j = q1; j != q2; j = j->foll)
			if (j) prlnk(j);
			else
				{
				printf("NULL ptr. aborted.\n");
				break;
				};
		prlnk(q2);
		return 0;
	}

int char_type (i)
     int i;
{
  static chars [12] = { '(',')','C','A','N','<','>','S','E','T','[',']'};

  if (i >= 0 && i < 12 ) return chars [i];
  else return '*';
}

/* print a link */
int prlnk(LINK *lk) {
	if (lk == NULL) {
		printf("Link: NULL.\n");
		return 0;
	}

	printf("link: %04x:%04x. Type=%c, NXT=%04x:%04x, PRV=%04x:%04x. ",
		FP_SEG(lk), FP_OFF(lk), char_type (lk->ptype),
		FP_SEG(lk->foll), FP_OFF(lk->foll),
		FP_SEG(lk->prec), FP_OFF(lk->prec));

	switch (lk->ptype) {
	case 0:
	case 1:
	case 6:
 		printf("PAIR=%04x:%04x\n", FP_SEG(lk->pair.b), FP_OFF (lk->pair.b));
		break;

	case 2:
		printf("SYMB=[%c], or %d\n",lk->pair.c,lk->pair.c);
		break;

	case 3:
		printf("ATOM=%.32s\n", lk->pair.f);
		break;

	case 4:
		printf("NUMB=%lu\n",lk->pair.n);
		break;

	case 5:
		printf("FUNC=(%04x:%04x) ", FP_SEG(lk->pair.b), FP_OFF (lk->pair.b));
		{
			char * cp = lk -> pair.f - 1;

			for (cp --; * cp != '\0'; cp --);
			cp ++;
			ri_actput (cp, stdout);
		}
		/*ri_actput(lk->pair.f,stdout);*/
		printf("\n");
		break;

	case 7:
	case 8:
		printf ("lev=%2lu, ind=%lu elev=%lu\n", 
				level_of (lk->pair.n), index_of (lk->pair.n),
				elevation_of(lk->pair.n));
		break;

	default:
		printf("STRANGE\n");
		break;
	}

	return 0;
}
# endif

