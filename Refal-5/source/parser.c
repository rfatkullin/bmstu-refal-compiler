
# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"


# define next_token (token = rc_gettoken ())


int rc_skip (t)
	int t;
	{
	while ((token != t) && (token != EOF)) next_token;
	return 0;
	}

char *rc_fname ()
	{
	char *fn;

	if (token != ID) rc_serror (104, NULL);
	while (token != ID)
		{
		rc_skip ('}');
		next_token;
		if (token == ';') next_token;
		if (token == EOF) return NULL;
		};
	fn = rc_deffn ();
	if (fn != NULL) strcpy (last_fn, fn); /*strncpy (last_fn, fn, MAXWS);*/
	else strcpy (last_fn, "***");
	next_token;
	return fn;
	}


struct node *rc_r_first ()
	{
	branch_t p1, p2, p3;

	p1.chunk = rc_r_side ();
	if (token != ':')
		{
		rc_serror (107, NULL);
		rc_skip (';');
		return Error;
		}
	else next_token;
	if (token == '{')
		{
		next_token;
		p2.tree = rc_l_list ();
		if (token != '}')
			{
			rc_serror (102, NULL);
			rc_skip (';');
			return Error;
			}
		else next_token;
		return rc_mknode (RSF1B, p1, p2, zero);
		};
	pushstk;
	p2.chunk = rc_l_side ();
	p3.tree = rc_r_tail ();
	popstk;
	return rc_mknode (RSF1, p1, p2, p3);
	}

struct node *rc_r_list ()
	{
	branch_t p1, p2;

	p1.tree = rc_r_first ();
	if (token == ';')
		{
		next_token;
		if (token == '}') p2.tree = NULL;
		else p2.tree = rc_r_list ();
		}
	else p2.tree = NULL;
	return rc_mknode (RSF, p1, p2, zero);
	}

struct node *rc_l_first ()
	{
	branch_t p1, p2;

	pushstk;
	p1.chunk = rc_l_side ();
	p2.tree = rc_r_tail ();
	if (p2.tree == Error)
		{
		popstk;
		rc_skip (';');
		return Error;
		};
	popstk;
	return rc_mknode (LSF1, p1, p2, zero);
	}

struct node *rc_l_list ()
	{
	branch_t p1, p2;

	p1.tree = rc_l_first ();
	if (token == ';')
		{
		next_token;
		if (token == '}') p2.tree = NULL;
		else p2.tree = rc_l_list ();
		}
	else p2.tree = NULL;
	return rc_mknode (LSF, p1, p2, zero);
	}


struct node *rc_r_tail ()
	{
	branch_t p1, p2, p3;

	if (token == '=')
		{
		next_token;
		p1.chunk = rc_r_side ();
		return rc_mknode (RCS1, p1, zero, zero);
		}
	else if (token != ',')
		{
		rc_serror (106, NULL);
		rc_skip (';');
		return (Error);
		}
	next_token;
	if (token == '{')
		{
		next_token;
		p1.tree = rc_r_list ();
		if (token != '}')
			{
			free_tree (p1.tree);
			rc_serror (102, NULL);
			rc_skip ('}');
			return (Error);
			}
		else next_token;
		return rc_mknode (RCS4, p1, zero, zero);
		};
	p1.chunk = rc_r_side ();
	if (token != ':')
		{
		rc_serror (107, NULL);
		rc_skip (';');
		free_tree (p1.tree);
		return Error;
		}
	else next_token;
	if (token == '{')
		{
		next_token;
		p2.tree = rc_l_list ();
		if (token != '}')
			{
			rc_serror (102, NULL);
			rc_skip ('}');
			free_tree (p1.tree);
			free_tree (p2.tree);
			return Error;
			}
		else next_token;
		return rc_mknode (RCS3, p1, p2, zero);
		};
	p2.chunk = rc_l_side ();
	p3.tree = rc_r_tail ();
	return rc_mknode (RCS2, p1, p2, p3);

	}


struct node *
rc_sentence (void) {
  branch_t ls, rs;

  clrstk;
  ls.chunk = rc_l_side ();
  rs.tree = rc_r_tail ();
  if (rs.tree == Error) {
    rc_skip (';');
    return Error;
  }
  return rc_mknode (ST, ls, rs, zero);
}


struct node *
rc_sents (void)	{
  branch_t st, sts, stsl, j;

  st.tree = rc_sentence ();
  stsl.tree = sts.tree = rc_mknode (SENTS, st, zero, zero);
  while (token == ';') {
    next_token;
    if (rc_1ofsent (token)) {
      st.tree = rc_sentence ();
      if (st.tree == Error) st.tree = NULL;
      j.tree = rc_mknode (SENTS, st, zero, zero);
      stsl.tree -> a3 = j;
      stsl.tree = j.tree;
    } else break;
  }
  return sts.tree;
}

int rc_1ofsent (t)
	int t;

	/* Determines if the token is in the FIRST of sentence. */

	{
	return (t != ';') && (t != '}');
	}

char *rc_idlist ()
	{

	if (token != ID)
		{
		rc_serror (104, NULL);
		rc_skip (';');
		}
	else
		{
		rc_mkextrn (str);
		next_token;};
	while (token == ',')
		{
		next_token;
		if (token != ID)
			{
			rc_serror (104, NULL);
			rc_skip (';');
			}
		else
			{
			rc_mkextrn (str);
			next_token;};
		};
	return NULL;
	}

struct node *
rc_fdef (void) {
  branch_t psents, fn;
  int nodetype;

  if (token == EOF) return NULL;
  if (token == EXTRN) {
    next_token;
    fn.pchar = rc_idlist ();
    if (token == EOF) return NULL;
    if (token != ';') {
      rc_serror (103, NULL);
      rc_skip (';');
    } else while (token == ';') next_token;
    return rc_fdef ();
  }
  else if (token == ENTRY) {
    nodetype = ENTRY;
    next_token;
  } else nodetype = FDEF;
  fn.pchar = rc_fname ();
  if (fn.pchar == NULL) return Error;
  else if (token != '{') {
    rc_serror (101, NULL);
    rc_skip ('{');
  } else {
    next_token;
    psents.tree = rc_sents ();
    if (token != '}') {
      rc_serror (102, NULL);
      rc_skip ('}');
    }
  }
  if (nerrors) return Error;
  else return rc_mknode (nodetype, fn, psents, zero);
}

struct node *
rc_parser (void) {
  struct node *res;

  next_token; /* token = rc_gettoken () */
  while (token == ';') next_token;
  if (token == EOF) res = NULL;
  else res = rc_fdef ();
  return res;
}

struct element *rc_r_side ()
	{
	return refal_expression (1);
	}

struct element *
rc_l_side (void) {
  return refal_expression (0);
}

