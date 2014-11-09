# include "decl.h"
# include "macros.h"
# include "fileio.h"
# include "ifunc.h"
# include <time.h>

int rf_cp ()
	{
	LINK *end;

	check_frz_args (1);

	cl;
	b1 = stock->foll;
	end = stock->pair.b;
	while (b1 != end)
		{
		tbel (5) = b1;
		tbel (6) = b2 = b1->pair.b;
		nel=7;
		oexp (4);
		sym ('=');
		cl;
		rdy (0);
		mule (11);
		out (2);
		est;
		return 0;
	restart:
		sp++;
		b1 = tbel (6)->foll;
		};
	rdy (0);
	out (2);
	est;
	return 0;
	}

int rf_chr ()
	{

	check_frz_args (1);

	cl;
	b1 = b1->foll;
	while (b1 != b2)
		{
		if (LINK_NUMBER (b1))
			{
			b1->ptype = LINK_TYPE_CHAR;
			b1->pair.c = (char)b1->pair.n;
			}
		b1 = b1->foll;
		};
	rdy (0);
	tple (4);
	out (2);
	est;
	return 0;
	}

int rf_ord ()
	{

	check_frz_args (1);

	cl;
	b1 = b1->foll;
	while (b1 != b2)
		{
		if (LINK_CHAR (b1))
			{
			b1->ptype = LINK_TYPE_NUMBER;
			b1->pair.n = b1->pair.c;
			};
		b1 = b1->foll;
		};
	rdy (0);
	tple (4);
	out (2);
	est;
	return 0;
	}

int rf_last ()
	{
	long k;

	movb1;
	if (LINK_VAR (b1)) goto call_freeze;
	else if (!LINK_NUMBER (b1)) goto restart;
	k = b1->pair.n;
        if( b1 != b2 ) b1 = b1->foll;

	while ( (k-- > 0L) )
		{
		if (LINK_EVAR (b2)) goto call_freeze;
		if (LINK_RSTRUCTB (b2)) b2 = b2->pair.b;
                if ( b1 == b2 ) break;
		b2 = b2->prec;
		};

	if (LINK_RSTRUCTB (b2)) b2 = b2->pair.b;
	tbel (4) = (b1 == b2 ? NULL : b1);
	tbel (5) = b2->prec;

	if (b2 == tbel (2))
		{
		tbel (6) = NULL;
		tbel (7) = b2;
		}
	else
		{
		tbel (6) = b2; 
		tbel (7) = tbel (2)->prec;
		}

	rdy (0);
	bl;
	tple (5);
	br;
	tple (7);
	out (2);
	est;
	return 0;

	restart:
		ri_imp ();
		return 1;

	call_freeze:
		ri_frz (2);
		return 0;
	}

int rf_first ()
	{
	long k;

	vsym;
	if (LINK_VAR (b1)) goto call_freeze;
	else if (!LINK_NUMBER (b1)) goto restart;
	k = b1->pair.n;
	tbel (4) = (b1 = b1->foll);

	if (b1 == b2) /*** empty exspression. **/
		{
		tbel (4) = NULL;
		tbel (5) = b2->prec;
		tbel (6) = NULL;
		tbel (7) = b2->prec;
		}
	else
		{
                if ( k == 0L ) tbel (4) = NULL;
		while ((b1 != b2) && (k-- > 0L))
			{
			if (LINK_EVAR (b1)) goto call_freeze;
			else if (LINK_LSTRUCTB (b1)) b1 = b1->pair.b;
			b1 = b1->foll;
			};
		tbel (5) = b1->prec;
		tbel (6) = (b1 == b2 ? NULL : b1);
		tbel (7) = b2->prec;
		};

	rdy (0);
	bl;
	tple (5);
	br;
	tple (7);
	out (2);
	est;
	return 0;

	restart:
		ri_imp ();
		return 1;

	call_freeze:
		ri_frz (2);
		return 0;
	}


int rf_implode (void) {
	char str [MAXWS+1];
	int i;
	char c, *pc;

	cl;
	if ((b1 = tbel (3)) == NULL) goto zero_end;
	else if (LINK_VAR (b1)) goto call_freeze;
	else if (!LINK_CHAR (b1) || !isalpha (b1->pair.c)) goto zero_end;
	i = 0;
	while (b1 != b2) {
		if (LINK_VAR (b1)) goto call_freeze;
		else if (!LINK_CHAR (b1)) break;
		c = b1->pair.c;
		/* convert c to upper case checking that it is a letter, digit or $,-,_	*/
		
		/* I'm not sure in needing of the condition. */
		if (i == 0) {
			if (isalpha (c)) str [i ++] = c; /*toupper (c);*/
			else break;
		} else {
			if (isalnum (c) || c == '_' || c == '$' || c == '-') str [i++] = c;
			else break;
		}

		b1 = b1->foll;
		if (i == MAXWS) break;
	}
	str [i] = '\0';
	if (b1 == b2)  tbel (3) = NULL;
	else tbel (3) = b1;
	pc = ri_cs_impl (str);
	rdy (0);
	ncs (pc);

end:

	tple (4);
	out (2);
	est;
	return 0;

zero_end:

	rdy (0);
	nns (0L);
	goto end;

call_freeze:

	ri_frz (2);
	return 0;
}

int rf_explode ()
	{
	char *pc;

	cl;
	if ((b1 = tbel (3)) == NULL) ri_imp ();
	else if (LINK_VAR (b1))
		{
		ri_frz (2);
		return 0;
		}
	else if (!LINK_COMPSYM (b1)) ri_imp ();
	else if (b1->foll != b2) ri_imp ();
	pc = b1->pair.f;
	rdy (0);
	while (*pc != 0) /* 252 line. Was *pc != NULL. Shura. 27.01.98 */
		{
		ns (*pc);
		pc ++;
		};
	out (2);
	est;
	return 0;
	}


int rf_lower ()
	{

	check_frz_args (1);

	cl;
	b1 = b1->foll;
	while (b1 != b2)
		{
		if (LINK_CHAR (b1)) b1->pair.c = tolower (b1->pair.c);
		b1 = b1->foll;
		};
	rdy (0);
	tple (4);
	out (2);
	est;
	return 0;
	}

int rf_upper (void) {
	check_frz_args (1);
	cl;
	b1 = b1->foll;
	while (b1 != b2) {
		if (LINK_CHAR (b1)) b1->pair.c = toupper (b1->pair.c);
		b1 = b1->foll;
	}
	rdy (0);
	tple (4);
	out (2);
	est;
	return 0;
}

int rf_step ()
	{
	cl;
	rdy (0);
	nns (nst);
	out (2);
	est;
	return 0;
	}

int rf_time ()
	{
	char *s;
	long k;

	cl;
	rdy (0);
# ifndef IBM370
	k = time (NULL);
	s = ctime (&k);
# else
	s = "n/a\n";
# endif
	while (*s != '\n')
		{
		ns (*s);
		s++;
		};
	out (2);
	est;
	return 0;
	}


int rf_lenw ()
	/*	LENW computes the number of terms in the Refal expression. */
	/*	Format:
			<LENW e.x> ==>  s.n e.x
	*/

	{
	long k;

	cl;
	k = 0L;
	b1 = b1->foll;
	while (b1 != b2)
		{
		k++;
		if (LINK_LSTRUCTB (b1)) b1 = b1->pair.b;
		else if (LINK_EVAR (b1))
			{
			ri_frz (2);
			return 0;
			}
		b1 = b1->foll;
		};
	rdy (0);
	nns (k);
	tple (4);
	out (2);
	est;
	return 0;
	}


int rf_dgall ()
	{

	if (b1->foll != b2) ri_imp ();
	b1 = stock;
	b2 = stock->pair.b;
	cl;
	rdy (0);
	tple (4);
	out (2);
	est;
	return 0;
	}


int rf_rp ()
	/* RP: replace the argument in the stock.	*/
	{
	LINK *end;

	check_frz_args (1);

	while (!LINK_CHAR (b1) || (b1->pair.c != '='))
		{
		b1 = b1->foll;
		if (b1 == b2) ri_imp ();
		};

	/* b1 now points to the first '=' in the argument. */

	tbel (4) = b1->prec;
	if (tbel (4) ==  tbel (1)) tbel (3) = NULL;
	else tbel (3) = tbel (1)->foll;

		/* tbel (3) - (4) contain the expression before '='. */

	nel = 15;
	cl;
		/* tbel (15) - (16) contain the expression after '='. */

	b1 = stock->foll;
	end = stock->pair.b;
	while (b1 != end)
		{
		tbel (5) = b1;
		tbel (6) = b2 = b1->pair.b;
		nel=7;
		oexp (4);
		sym ('=');
		cl;

			/* set b to '=', transplant the new expression and
				delete the old expression */
		b = tbel (9);
		tple (16);
                if (tbel(10) != NULL) /* Nemytykh A.P., 15.05.2004 */
		        out (11);
		rdy (0);
		out (2);
		est;
		return 0;

 restart:
		sp++;
		b1 = tbel (6)->foll;
		};
	b = stock;
	rend = b->foll;
	bl;
	tple (4);
	ns ('=');
	tple (16);
	br;
	weld (b,rend);
	rdy (0);
	out (2);
	est;
	return 0;
	}


int rf_open ()
	{
	char s [FILENAME_MAX];
	char m[2];
        char * mode = m;
	int i, lognum = 0;

	check_freeze;

	b1 = b1->foll;	/* Get the mode.*/
	if (LINK_CHAR (b1)) 
           { m [0] = b1->pair.c; 
             m [1] = '\0';
           }
	else if (LINK_COMPSYM (b1)) mode = b1->pair.f;
	else ri_error(8);

	b1 = b1->foll;	/* Get the file number. */
	if (!LINK_NUMBER (b1)) ri_error (8);
	else lognum = b1->pair.n;
	lognum %= FILE_LIMIT;
	for (b1 = b1->foll, i=0; (b1 != b2) && (i < FILENAME_MAX); b1 = b1->foll, i++)
		if (!LINK_CHAR (b1)) ri_error(8);
		else s[i] = b1->pair.c;
	if (i) s[i] = '\0';
	else sprintf (s,"REFAL%d.DAT",lognum);
	ri_open (lognum,mode,s);
	rdy (0);
	out (2);
	est;
	return 0;
	}


