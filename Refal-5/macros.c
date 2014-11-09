
# include "decl.h"
# include "macros.h"
# include "ifunc.h"


		/* This file contains operators of RASL which are defined
			as functions. Feb 28 1987. D.T. */

  /* For translator pgraph */
/*extern unsigned long ul_local_calls;*/

void ri_act1 (x)
     char *x;
{
  (b -> pair.b) -> pair.f = x;
  (b -> pair.b) -> ptype = 5;
  b -> ptype  =  6;
  (precp -> foll) -> prec = b;
  precp = b;
}

int ri_blr ()
	{
	rend = b -> foll;
	weld (b, lfm);
	b = lfm;
	all;
	con (b, lfm);	
	weld (b, lfm);
	b = lfm;
	all;
	weld (b, rend);
	return 0;
	}

void ri_cl ()
{
  if (b1 -> foll == b2) {
    tel [nel ++] = NULL;
    tel [nel ++]  = b1;
  } else {
    tel [nel ++] = b1 -> foll;
    tel [nel ++] = b2 -> prec;
  }
}

int ri_est (void) {
	(precp -> foll) ->  prec = nextp;
	b2 = (quap -> foll) -> prec;
	b1 = b2 -> pair.b;
	p = b1 -> pair.f;
	nextp = (b2 -> foll) -> prec;
	precp = quap;
	(b2 -> foll) -> prec = b2;
	sp = stoff;
	tbel (0) = b1  -> prec;
	tbel (1) = b1;
	tbel (2) = b2;
	nel = 3L;
	actfun=p;
	nst++;

	return 0;
}


int ri_mule (arg)
     int arg;
{
  LINK *j, *oexa;

  if (tbel (arg-1) != NULL) {
    oexa = tbel (arg-1);
    rend = b -> foll;
    j = tbel (arg);
    while (oexa != j) {
      ri_cpelt (oexa);
      oexa = oexa -> foll;
    }
    ri_cpelt (oexa);
    weld (b, rend);
  }
  return 0;
}

/* Copy some non-empty expression from lp_beg to lp_end */
int ri_mulpar (LINK * lp_beg, LINK * lp_end)
{
  LINK *j, *oexa;

  oexa = lp_beg;
  rend = b -> foll;
  j = lp_end;
  while (oexa != j) {
    ri_cpelt (oexa); /* Copy an element */
    oexa = oexa -> foll; /* move to a free element */
  }
  ri_cpelt (oexa);
  weld (b, rend);
  return 0;
}

/* Shura. Add META_BRACKETS */
int ri_cpelt (oexa)
	LINK *oexa;
	{
	if (!LINK_STRUCTB (oexa))
		{
		lfm -> pair = oexa -> pair;
		lfm -> ptype = oexa -> ptype;
		weld (b, lfm);
		b = lfm;
		all;
		}
	else if (LINK_LSTRUCTB (oexa))
		{
		lfm -> pair.b = pa;
		pa = lfm;
		weld (b, lfm);
		b = lfm;
		all;
		}
	else if (LINK_RSTRUCTB (oexa))
		{
		nextb = pa -> pair.b;
		con (pa, lfm);
		pa = nextb;
		weld (b, lfm);
		b = lfm;
		all;	
		};
	return 0;
	}


int ri_ns1 ()
{
  rend = b -> foll; /* b=(<) -> rend=(....>) */
  weld (b, lfm);    /* b=(<) <-> [list of free memory] */
  b=lfm; /* (<) <-> b=[list of free memory] */
  all;
  /* all is:
   * if (NULL == (lfm = lfm -> foll)) lfm = ri_fmount ();
   */
  weld (b, rend); 
  /* (<) <-> b=(*) <-> rend=(..>) ... [list of free memory] */
  return 0;
}

void ri_tple (arg)
     int arg;
{
  if (tel [arg - 1] != NULL) {
    NEXT (PREV (tel [arg - 1])) = NEXT (tel [arg]);
    PREV (NEXT (tel [arg])) = PREV (tel [arg - 1]);
    /* weld (tel [arg - 1] -> prec, tel [arg] -> foll); */
    NEXT (tel [arg]) = NEXT (b);
    PREV (NEXT (b)) = tel [arg];
    /* weld (tel [arg], NEXT (b)); */
    NEXT (b) = tel [arg - 1];
    PREV (tel [arg - 1]) = b;
    /* weld (b, tel [arg - 1]); */
    b = tel [arg];
  }
}

void ri_out (n)
     int n;
{
  LINK *j;

  j = tel [n];
  if (j != b) {
    rend = j -> foll;
  	if (rend != b) {
		weld (j, lfm);
		lfm = b -> foll;
		weld (b, rend);
	}
  }
}
/*
int ri_tpls (arg)
int arg;
{
weld (tbel (arg) -> prec, tbel (arg) -> foll);
weld (tbel (arg), b -> foll);
weld (b, tbel (arg));
b = tbel (arg);
return 0;
}
*/

int ri_tpls (arg)
     int arg;
{
  /* b is the current end of the built right part */
  if (SPAR (tel [arg])) {
    LINK * lp_pair = PAIR (tel [arg]);

    /* NEXT (PREV (tel [arg])) = NEXT (PAIR (tel [arg]));
     * PREV (NEXT (PAIR (tel [arg]))) = PREV (tel [arg]);
     */
    weld (PREV (tel [arg]), NEXT (lp_pair)); /* CUT s-par. */
    /* NEXT (PAIR (tel [arg])) = NEXT (b);
     * PREV (NEXT (b)) = PAIR (tel [arg]);
     */
    weld (lp_pair, NEXT (b)); /* Glue the tail of s-par. */
    /* NEXT (b) = tel [arg];
     * PREV (tel [arg]) = b;
     */
    weld (b, tel [arg]); /* Glue the head of s-par. */
    b = lp_pair;
  } else {
    weld (tbel (arg) -> prec, tbel (arg) -> foll);
    weld (tbel (arg), b -> foll);
    weld (b, tbel (arg));
    b = tbel (arg);
  }
  return 0;
}

int ri_rimp ()
	{
	/*printf ("DEBUG. Stack depth: %d\n", sp);*/
	sp--;
	/*printf ("DEBUG. Stack depth: %d\n", sp);*/
	popst (b1, b2, nel, p);
	/* FOR DEBUG */
	/*if (p != NULL) printf ("Instruction (stack): %d\n", * p);*/
	while (p == NULL)
		{
		rend = b2 -> foll;
		weld (b2, lfm);
		lfm = b1;
		weld (b1 -> prec, rend);

		/*printf ("DEBUG. Stack depth: %d\n", sp);*/
		sp--;
		/*printf ("DEBUG. Stack depth: %d\n", sp);*/
		popst (b1, b2, nel, p);
		/* FOR DEBUG */
		/*if (p != NULL) printf ("Instruction (stack): %d\n", * p);*/
		};
	/* FOR DEBUG */
	/*printf ("Restore instruction from stack\n");*/
	return 0;
	}


int ri_oexp (q)
	int q;
	{
	LINK *j, *oexa;

	if (tbel (q-1) == NULL) tbel (nel)=NULL;
	else
		{
		tbel (nel) = b1  -> foll;
		oexa = tbel (q-1) -> prec;
		j = tbel (q);
		while (oexa != j)
			{
			oexa = oexa -> foll;
			movb1;
			if (LINK_VAR (b1) || LINK_VAR (oexa))
				{
				if (!LINK_EQTYPES (b1, oexa)) freeze
				else if (b1 -> pair.n != oexa -> pair.n) freeze
				}
			else if (!LINK_EQTYPES (b1, oexa)) rimp
			else if (LINK_STRUCTB (b1));
			else if (LINK_CHAR (b1))
				{
				if (b1 -> pair.c!=oexa -> pair.c)rimp;
				}
			else if (b1 -> pair.n!=oexa -> pair.n) rimp;
			}
		};
	tbel (++nel) = b1;
	nel++;
	return 0;

	restart:
		return 1;
	}


int ri_oexpr (n)
	int n;
	{
	LINK *j, *oexa;

	tbel (nel+1) = b2 -> prec;
	if (tbel (n-1) == NULL)
		{
		tbel (nel)=NULL;
		nel += 2;
		}
	else
		{
		oexa = tbel (n) -> foll;
		j = tbel (n-1);
		while (oexa != j)
			{
			oexa = oexa -> prec;
			movb2;
			if (LINK_VAR (b2) || LINK_VAR (oexa))
				{
				if (!LINK_EQTYPES (b2, oexa)) freeze
				else if (b2 -> pair.n != oexa -> pair.n) freeze
				}
			else if (!LINK_EQTYPES (b2, oexa)) rimp
			else if (LINK_STRUCTB (b2)) ;
			else if (LINK_CHAR (b2))
				{
				if (b2 -> pair.c != oexa -> pair.c) rimp
				}
			else if (b2 -> pair.n != oexa -> pair.n) rimp;
			};
		tbel (nel) = b2;
		nel += 2;
		}
	return 0;

	restart:
		return 1;
	}

int ri_ovs (n)
	int n;
	{
		movb1;
		if (LINK_VAR (b1) || LINK_VAR (tbel (n)))
			{
			if (LINK_SVAR (b1) && LINK_SVAR (tbel (n)) &&
					b1 -> pair.n == tbel (n) -> pair.n) ;
			else freeze
			}
		else if (!LINK_EQTYPES (b1, tbel (n))) rimp
		else if (LINK_CHAR (b1))
			{
			if (b1 -> pair.c != tbel (n) -> pair.c) rimp
			}
		else if (b1 -> pair.n != tbel (n) -> pair.n) rimp
		tbel (nel) = b1;
		nel++;
		return 0;

	restart:
		return 1;
	}

int ri_ovsr (n)
	int n;
	{
		movb2;
		if (LINK_VAR (b2) || LINK_VAR (tbel (n)))
			{
			if (LINK_SVAR (b2) && LINK_SVAR (tbel (n)) &&
					b2 -> pair.n == tbel (n) -> pair.n) ;
			else freeze
			}
		else if (!LINK_EQTYPES (b2, tbel (n))) rimp
		else if (LINK_CHAR (b2))
			{
			if (b2 -> pair.c != tbel (n) -> pair.c) rimp
			}
		else if (b2 -> pair.n != tbel (n) -> pair.n) rimp
		tbel (nel) = b2;
		nel++;
		return 0;

	restart:
		return 1;
	}

int ri_lens (c)
	char c;
	{
		LINK *j;

		if (tbel (nel) == NULL) tbel (nel) = b1 -> foll;
		b1 = tbel (nel+2);
		do
			{
			movb1;
			while (LINK_STRUCTB (b1))
				{
				b1 = b1 -> pair.b;
				movb1;
				}
			if (LINK_VAR (b1)) freeze
			}
			while (!LINK_CHAR (b1) || b1 -> pair.c != c);
		j = tbel (nel);
		if (j == b1) tbel (nel) = NULL;
		sp++;
		tbel (++nel) = b1 -> prec;
		tbel (++nel) = b1;
		nel++;
		return 0;

	restart:
		return 1;
	}


int ri_lenp ()
	{
		LINK *j;

		if (tbel (nel) == NULL) tbel (nel) = b1 -> foll;
		b1 = tbel (nel+3);
		movb1;
		while (!LINK_LSTRUCTB (b1))
			{
			if (LINK_EVAR (b1)) freeze;
			movb1;
			}
		b2 = b1 -> pair.b;
		j = tbel (nel);
		if (j == b1) tbel (nel) = NULL;
		++sp;
		tbel (nel+1) = b1 -> prec;
		tbel (nel+2) = b1;
		tbel (nel+3) = b2;
		nel += 4;
		return 0;

	restart:
		return 1;
	}

