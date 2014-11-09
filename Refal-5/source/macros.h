
#define LINK_TYPE_LSTRUCTB   0
#define LINK_TYPE_RSTRUCTB   1
#define LINK_TYPE_CHAR       2
#define LINK_TYPE_COMPSYM    3
#define LINK_TYPE_NUMBER     4
#define LINK_TYPE_LACT       5
#define LINK_TYPE_RACT       6
#define LINK_TYPE_SVAR       7
#define LINK_TYPE_EVAR       8
#define LINK_TYPE_TVAR       9
#define LINK_TYPE_MLSTRUCTB   10
#define LINK_TYPE_MRSTRUCTB   11

#define META_SVAR            's'
#define META_EVAR            'e'
#define META_TVAR            't'

#define LINK_CONST_SYMBOL(q) ((((q)->ptype >= 2) && ((q)->ptype <= 4)))

# define LINK_SYMBOL(q) ((((q)->ptype >= 2) && ((q)->ptype <= 4)) || ((q)->ptype == 7))
# define LINK_STRUCTB(q) ((q)->ptype < 2)
# define LINK_LSTRUCTB(q) ((q)->ptype == LINK_TYPE_LSTRUCTB)
# define LINK_RSTRUCTB(q) ((q)->ptype == LINK_TYPE_RSTRUCTB)
# define LINK_MLSTRUCTB(q) ((q)->ptype == LINK_TYPE_MLSTRUCTB)
# define LINK_MRSTRUCTB(q) ((q)->ptype == LINK_TYPE_MRSTRUCTB)
# define LINK_CHAR(q) ((q)->ptype == LINK_TYPE_CHAR)
# define LINK_COMPSYM(q) ((q)->ptype == LINK_TYPE_COMPSYM)
# define LINK_NUMBER(q) ((q)->ptype == LINK_TYPE_NUMBER)
# define LINK_EQTYPES(q1,q2) ((q1)->ptype == (q2)->ptype)
# define LINK_LACT(q) ((q)->ptype == LINK_TYPE_LACT)
# define LINK_RACT(q) ((q)->ptype == LINK_TYPE_RACT)
# define LINK_SVAR(q) ((q)->ptype == LINK_TYPE_SVAR)
# define LINK_EVAR(q) ((q)->ptype == LINK_TYPE_EVAR)
# define LINK_VAR(q) ((q)->ptype >= 7 && 10 > (q) -> ptype)

#define NEXT(q)    ((q)->foll)
#define PREV(q)    ((q)->prec)
#define PAIR(q)    ((q)->pair.b)

static
#ifdef FOR_OS_LINUX
 inline
#endif
      int
EPAR (LINK * q) {
  return (LINK_MLSTRUCTB(q)
	  && (NEXT(q)->ptype == LINK_TYPE_CHAR)
	  && (META_EVAR == NEXT(q)->pair.c));
}

static
#ifdef FOR_OS_LINUX
 inline
#endif
      int
TPAR (LINK * q) {
  return (LINK_MLSTRUCTB(q) &&
	  (NEXT(q)->ptype == LINK_TYPE_CHAR) &&
	  (META_TVAR == NEXT(q)->pair.c));
}

static
#ifdef FOR_OS_LINUX
 inline
#endif
      int
SPAR (LINK * q) {
  return (LINK_MLSTRUCTB(q) &&
	  (NEXT(q)->ptype == LINK_TYPE_CHAR) &&
	  (META_SVAR == NEXT(q)->pair.c));
}

static
#ifdef FOR_OS_LINUX
 inline
#endif
      int
UNKNOWN_PAR (LINK * q) {
  return (LINK_MLSTRUCTB(q) &&
	  (NEXT(q)->ptype == LINK_TYPE_MLSTRUCTB));
}

	/*	code 2: ACT1.	*/
# define act1(x) ri_act1(x)

	/*	code 3: BL.	*/
# define bl { lfm->pair.b = pa; pa=lfm; ns1;}

	/*	code 4: BLR.	*/
# define blr ri_blr()

	/*	code 5: BR.	*/
# define br  {nextb = pa->pair.b; con (pa,lfm); pa=nextb; ns1;}

	/*	code 6: CL.	*/
# define cl ri_cl()

	/*	code 7: SYM.	*/
# define sym(arg)   { movb1;		\
		if (LINK_VAR(b1) || EPAR(b1) || TPAR (b1) || SPAR (b1)) freeze 	\
		if (!LINK_CHAR(b1) || (b1->pair.c != arg)) rimp; \
		tbel(nel++) = b1;}

	/*	code 8: SYMR.	*/
# define symr(arg)  { movb2;		\
		if (LINK_VAR(b2) || EPAR(b2) || TPAR (b2) || SPAR (b2)) freeze	\
		if(! LINK_CHAR(b2) || (b2->pair.c != arg)) rimp; \
		tbel(nel++) = b2;}

	/*	code 10: EMP.	*/
# define emp  if (b1->foll != b2) { \
		for (b1 = b1 -> foll; LINK_EVAR(b1) || EPAR (b1);b1 = b1->foll) {\
		  if (b1->foll == b2) freeze;\
                } \
		rimp;\
              }

	/*	code 11: EST.	*/
# define est  ri_est()

	/*	code 13: MULE.	*/
# define mule(arg) ri_mule(arg)

	/* code 14: MULS. */
# define muls(n) \
{lfm->ptype = tbel(n)->ptype;	lfm->pair = tbel(n)->pair; ns1;}

	/* code 16: PLEN. */
# define plen {pushst(b1,b2,nel,p); sp++; \
	tbel(nel) = NULL; tbel(++nel) = b1;	nel++; }

	/* code 17: PLENS. */
# define plens {pushst(b1,b2,nel,p); tbel(nel) = NULL; tbel(nel+2) = b1;}

	/* code 18: PLENP. */
# define plenp {pushst(b1,b2,nel,p); tbel(nel) = NULL; tbel(nel+3) = b1;}


	/*	code 19: PS.	*/
# define ps  { movb1; \
	if (LINK_EVAR (b1) || EPAR (b1)) freeze \
	if (!LINK_STRUCTB (b1)) rimp	\
	else {b2 = b1->pair.b;	\
		tbel(nel++) = b1;	\
		tbel(nel++) = b2;}} 

	/*	code 20: PSR.	*/
# define psr {	movb2;	\
	if (LINK_EVAR (b2) || EPAR (b2)) freeze \
	if (!LINK_STRUCTB(b2)) rimp \
	else {tbel(nel++) = b2->pair.b; \
	tbel(nel++) = b2; \
	b2 = b2->pair.b;}}


/*	code 23: OEXP.	*/
# define oexp(q) {if (ri_oexp(q) == 1) goto restart; }

/* code 24: OEXPR. */
# define oexpr(q) {if (ri_oexpr(q) == 1) goto restart; }

/* code 25: OVSYM. */
# define ovsym(n) { if (ri_ovs(n) == 1) goto restart; }

/* code 26: OVSYMR. */
# define ovsymr(n) { if (ri_ovsr(n) == 1) goto restart; }

/* code 27: TERM. */
# define term { \
	movb1; \
	if (LINK_EVAR(b1) || EPAR (b1)) freeze \
	else { \
		tbel(nel) = b1; \
		if(LINK_STRUCTB(b1)) b1 = b1->pair.b; \
		tbel(++nel) = b1; \
		nel++; \
	} \
}

/* code 28: TERMR. */
# define termr { \
	movb2; \
	if (LINK_EVAR(b2) || EPAR (b2)) freeze \
	else { \
		tbel(nel+1) = b2; \
		if(LINK_STRUCTB(b2)) b2 = b2->pair.b; \
		tbel(nel++) = b2; \
		nel++; \
	} \
}

/*	code 29: RDY.	*/
# define rdy(index)  {b=tbel(index);}

	/*	code 34: SETB.	*/
# define setb(n,m)  { b1 = tbel(n);		\
		b2 = tbel(m); if (tbel(m) == NULL) b2 = tbel(m+1)->foll;}

	/* code 35: LEN. */
# define len {if (tbel(nel) == NULL) tbel(nel) = b1->foll; \
		b1 = tbel(nel+1);	movb1; \
		if (LINK_EVAR (b1) || EPAR (b1)) freeze	\
		else {if (LINK_STRUCTB(b1)) b1 = b1->pair.b; \
		++sp; tbel(++nel) = b1; nel++; }}

	/* code 36: LENS. */
# define lens(c) { if (ri_lens(c) == 1) goto restart; }

	/* code 37: LENP. */
# define lenp { if (ri_lenp() == 1) goto restart; }

	/* code 38: LENOS. deleted (10/11/1987). */

	/* code 39: SYMS. */
# define syms(n) {while(n-- > 0 ) {movb1; \
	if (LINK_VAR(b1) || EPAR (b1) || TPAR (b1) || SPAR (b1)) freeze \
	else if (!LINK_CHAR(b1) || (b1->pair.c != *p++)) rimp \
	else tbel(nel++) = b1;}}

	/* code 40: SYMSR. */
# define symsr(n) {while(n-- > 0) {movb2; \
	if (LINK_VAR(b2) || EPAR (b2) || TPAR (b2) || SPAR (b2)) freeze \
	else if (!LINK_CHAR(b2) || (b2->pair.c != *p++)) rimp \
	else tbel(nel++) = b2;}}

	/* code 41: TEXT. */
# define text(n) {rend = b->foll; \
	while (n-- > 0) {lfm->ptype = LINK_TYPE_CHAR; lfm->pair.c = *p++; \
		weld (b,lfm); b = lfm; all; }; weld(b,rend); } 

	/*	code 43: NS.	*/
# define ns(x) {lfm->pair.c = x; lfm->ptype = LINK_TYPE_CHAR; ns1;}
# define ns1 ri_ns1()

	/*	code 45: TPLE.	*/
# define tple(arg) ri_tple(arg)

	/*	code 46: TPLS.	*/
# define tpls(arg) ri_tpls(arg)

	/* code 47: TRAN. */
# define tran(arg) {pushst(b1,b2,nel,arg); sp++;}

	/*	code 48: VSYM.	*/
/*# define vsym  {movb1; if (LINK_EVAR(b1)) freeze \ */
/*	else if (!LINK_SYMBOL(b1)) rimp else tbel(nel++) = b1;} */

	/*	code 48: VSYM.	*/
# define vsym  \
{\
  movb1; \
  if (EPAR(b1) || TPAR(b1) || UNKNOWN_PAR(b1)) { freeze;} \
  else if (LINK_CONST_SYMBOL(b1)) { tbel(nel++) = b1;} \
  else if (SPAR(b1)) { tbel(nel++) = b1; b1 = PAIR (b1); } \
  else rimp;\
}

	/*	code 49: VSYMR.	*/
# define vsymr {movb2; if (LINK_EVAR(b2) || EPAR (b2)) freeze \
	else if (!LINK_SYMBOL(b2)) rimp else tbel(nel++) = b2;}

	/*	code 50: OUTEST.	*/
# define out(n) ri_out(n)

	/*	code 55: CSYM.	*/
# define csym(arg)  { movb1; if (LINK_VAR(b1) || EPAR (b1) || TPAR (b1) || SPAR (b1)) freeze \
		else if (!LINK_COMPSYM(b1) || (b1->pair.f != arg)) rimp;\
		tbel(nel++) = b1;}

	/*	code 56: CSYMR.	*/
# define csymr(arg) { movb2; if (LINK_VAR(b2) || EPAR (b2) || TPAR (b2) || SPAR (b2)) freeze \
		else if (!LINK_COMPSYM(b2) || (b2->pair.f != arg)) rimp; \
		tbel(nel++) = b2;}

	/*	code 57: NSYM.	*/
# define nsym(arg) { \
	movb1; \
	if (LINK_VAR(b1) || EPAR (b1) || TPAR (b1) || SPAR (b1)) { \
		freeze; \
	} \
	if(!LINK_NUMBER(b1) || (b1->pair.n != (unsigned long) arg)) { \
		rimp; \
	} \
	tbel(nel++) = b1; \
}

	/*	code 58: NSYMR.	*/
# define nsymr(arg) { \
	movb2; \
	if (LINK_VAR(b2) || EPAR (b2) || TPAR (b2) || SPAR (b2)) { \
		freeze; \
	} \
	if(!LINK_NUMBER(b2) || (b2->pair.n != (unsigned long)arg)) { \
		rimp; \
	} \
	tbel(nel++) = b2; \
}

	/*	code 59: NCS.	*/
# define ncs(x) {lfm->pair.f = x; lfm->ptype=3; ns1;}

	/*	code 60: NNS.	*/
# define nns(x) {lfm->pair.n = x; lfm->ptype=4; ns1;}

  /*   -------------------------------------------------   */

# define tbel(x) (*(tel + (x)))

# define all {if ((lfm = lfm->foll) == NULL) lfm = ri_fmout();}

# define movb1 {if ((b1 = b1->foll) == b2) rimp;}

# define movb2 {if ((b2 = b2->prec) == b1) rimp;}

# define con(a1,a2) { \
		a1->ptype = LINK_TYPE_LSTRUCTB; \
		a2->ptype = LINK_TYPE_RSTRUCTB; \
		PAIR(a1) = a2;	PAIR(a2) = a1;}

# define weld(a1,a2) {NEXT(a1) = a2; PREV(a2) = a1; }

# define rimp  {ri_rimp(); goto restart;}

# define freeze  {ri_frz(1); goto restart;}

# define pushst(g1,g2,ne,ar)     	\
        { st[sp].b1 = (LINK *) g1;	\
          st[sp].b2 = (LINK *) g2;	\
          st[sp].nel = (long) ne;	\
          st[sp].ra = (char *) ar;}

# define popst(g1,g2,ne,ar)	\
    { g1 = st[sp].b1;			\
      g2 = st[sp].b2;			\
      ne = st[sp].nel;			\
      ar = st[sp].ra;}

# define check_freeze if (exists_freeze()) {ri_frz(1); return 0; }
# define check_frz_args(frz_code) if (contains_vars (tbel (1), tbel (2))) \
	{ri_frz(frz_code); return 0;}

# define YES 1
# define NO  0

  /*   -------------------------------------------------   */

