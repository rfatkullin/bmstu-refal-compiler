
# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"

extern struct bitab bi [];
extern long nbi;

struct node *
rc_mknode (int t, branch_t arg2, branch_t arg3, branch_t arg4) {
  struct node *p;

  /* Was. Shura 20.05.99
   *  if (nerrors) return NULL;
   */
  if (nerrors) return Error;

  /* 14 line. Not check result of malloc. Shura. 29.01.98 */
  if ((p = (struct node *) malloc (sizeof (struct node))) == NULL) {
    fprintf (stderr, "No memory for node\n");
    exit (1);
  }
  p -> nt = t;
  p -> a2 = arg2;
  p -> a3 = arg3;
  p -> a4 = arg4;
  return p;
}

/* returns a pointer to function name */
char *rc_deffn (void) {
	struct functab *p;
	/* char * cp_w; */
	int flag_del = 0;

	/* search the function table.  */
	p = searchf (str, ft);     

	/* cheking in standart functions */
	if (p == NULL) {
		int i;

		for (i = 1; i < nbi; i ++) {
			if (strcmp (str, bi [i].fname) == 0) {
				if (NULL == (p = (struct functab *) malloc (sizeof (struct functab)))) {
					fprintf (stderr, "No enough memory\n");
					exit (1);
				}
				if (NULL == (p -> name = (char *) malloc (strlen (bi [i].fname) + 1))) {
					fprintf (stderr, "No enough memory\n");
					exit (1);
				}
				strcpy (p -> name, bi [i].fname);
				flag_del = 1;
				break;
			}
		}
	}

	/* insert into the ft and fb tables.    */
	if (p == NULL) {
		p = (struct functab *) rc_allmem (sizeof (struct functab));
		if (p == NULL) {
			fprintf (stderr,"Ran out of memory.\n");
			exit (1);
		}
		p -> next = fb;
		fb = p;

		if (NULL == (p -> name = (char *) malloc (strlen (str) + 1))) {
			fprintf (stderr, "No memory for function name\n");
			exit (1);
		}
		/*copyst (p -> name);*/
		strcpy (p -> name, str);
		btcount ++;
		p = (struct functab *) rc_allmem (sizeof (struct functab));
		if (p == NULL) {
			fprintf (stderr,"Ran out of memory.\n");
			exit (1);
		}
		p -> next = ft;
		ft = p;

		if (NULL == (p -> name = (char *) malloc (strlen (str) + 1))) {
			fprintf (stderr, "No memory for function name\n");
			exit (1);
		}
		/*copyst (p -> name);*/
		strcpy (p -> name, str);
		return (p -> name);
	} else {
		char * cp_w;

		rc_serror (201, (cp_w = p -> name));
		if (flag_del) {
			/*free (p -> name);*/
			free (p);
		}
		return (cp_w);
		/*return NULL;*/
	}
}

/* puts the function name on external list */
/* DT June 10 1986. */
/* s is the name of the function. */
char *rc_mkextrn (char *s) {
	struct functab *p;

	p = searchf (s,fx);   /* search the external table.   */

	if (p == NULL) { /* insert into table. */
		p = (struct functab *) rc_allmem (sizeof (struct functab));
		if (p == NULL) {
			fprintf (stderr,"Ran out of memory.\n");
			exit (1);
		}
		p -> next = fx;
		fx = p;
		xtcount ++;

		if (NULL == (p -> name = (char *) malloc (strlen (s) + 1))) {
			fprintf (stderr, "No memory for function name\n");
			exit (1);
		}
		strcpy (p -> name,s);
		return (p -> name);
	} else {
		rc_serror (201,p -> name);
		/*
		free (p -> name);
		free (p);
		*/
		return p -> name;
	}
}

/* puts the function name on entry list */
char *rc_mkentry (char *ptr) {
	struct functab *p;

	p = searchf (ptr,fe);   /* search the entry table.   */
	if (p == NULL) {   /* insert into table.           */
		p = (struct functab *) rc_allmem (sizeof (struct functab));
		if (p == NULL) {
			fprintf (stderr,"Ran out of memory.\n");
			exit (1);
		}
		p -> next = fe;
		fe = p;
		ntcount ++;

		if (NULL == (p -> name = (char *) malloc (strlen (ptr) + 1))) {
			fprintf (stderr, "No memory for function name\n");
			exit (1);
		}
		strcpy (p -> name,ptr);
		/*return (p -> name);*/
	}
	return p -> name;
}

int free_tree (q)
	struct node *q;
	{
		/* frees the nodes of the parse tree pointed by *q */
	struct node *z;
	struct element *e;

	while (q != NULL)
		switch (q -> nt)
			{
			case 0:
				q = NULL;
				break;

			case FDEF: case ENTRY:

				z = q -> a3.tree;
				free ((void *) q);
				q = z;
				break;

			case RSF1B: case LSF1: case RCS3: case ST:

				e = q -> a2.chunk;
				free ((void *) e);
				z = q -> a3.tree;
				free ((void *) q);
				q = z;
				break;

			case RCS2: case RSF1:

				e = q -> a2.chunk;
				free ((void *) e);
				e = q -> a3.chunk;
				free ((void *) e);
				z = q -> a4.tree;
				free ((void *) q);
				q = z;
				break;

			case LSF: case RSF:

				z = q -> a2.tree;
				free_tree (z);
				z = q -> a3.tree;
				free ((void *) q);
				q = z;
				break;

			case RCS1:

				e = q -> a2.chunk;
				free ((void *) e);
				free ((void *) q);
				q = NULL;
				break;
				
			case SENTS: case RCS4:

				z = q -> a2.tree;
				free ((void *) q);
				q = z;
				break;

			default:
				fprintf (stderr, "Error Releasing parse tree: node type = %d\n",
					q -> nt);
				q = NULL;
				return 0;
			}
	return 0;
	}
