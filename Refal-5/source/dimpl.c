
# include "rasl.h"
# include "decl.h"
# include "ddecl.h"
# include "macros.h"
# include "tfunc.h"
# include "ifunc.h"
# include <ctype.h>

/* hash table entry. */
struct hte {
	char *value;     /* pointer value. */
	char *module;    /* pointer to the module name. */
	struct hte *nxt; /* next hte. */
};

struct hte ht [HASH_SIZE]; /* hash table. */
short ht_flag = 0;         /* show if ht was created. */

/* For copying of trace informations */
extern FILE * fp_debugInfo;

/* returns pointer to the list of functions with name *str and its length */
struct listfun * rd_implode (char *str, int *function_number) {
	struct  listfun *k, *ks, *kr; /* kr points to the first in the list, ks to the last. */
	int hv; /* hash value of the string	*/
	struct hte *r;
	char *q;

	if (ht_flag == 0) rd_create_ht ();
	*function_number = 0;
	kr = ks = NULL;
	hv = ri_hash (str);
	r = &(ht [hv]);
	while (r != NULL) {
		q = r -> value;
		if (q == NULL) break;
		for (q -= 2; * q != 0; q --);
		q ++;
		if (strcmp (str, q) == 0) {
			k = (struct listfun *) malloc (sizeof (struct listfun));
			/* No checking for NULL. Shura. 27.01.98 */
			if (NULL == k) {
			  exit (1);
			}
			k -> module = r -> module;
			k -> function = q + strlen (q) + 1;
			k -> next = NULL;
			if (ks != NULL) ks -> next = k;
			else kr = k;
			ks = k;
			(*function_number) ++;
		}
		r = r -> nxt;
	}
	return kr;
}

/* creates hash table.	*/
int rd_create_ht (void) {
	char *mod, **q;

	ht_flag = 1;
	for (q = local_functions; *q != NULL; q++) {
		mod = *q++;
		rd_ins_mod (mod, q);
		while (*q != NULL) {
			rd_ins_ht (*q, mod);
			++q;
		}
	}
	return 0;
}

int rd_ins_mod (char *modn, char **funcp) {
	struct modlist *q;

	q = (struct modlist *) malloc (sizeof (struct modlist));
	/* No checking for NULL. Shura. 27.01.98 */
	if (NULL == q) {
	  exit (1);
	}
	q -> module = modn;
	q -> funcs = funcp;
	q -> next = module_list;
	module_list = q;
	return 0;
}

/* insert a symbol into the hash table.	*/
int rd_ins_ht (char *s, char * mod) {
	short hv;
	struct hte *r;
	char * cp_w;

	for (cp_w = s - 2; * cp_w != 0; cp_w --);
	cp_w ++;
	hv = ri_hash (cp_w);
	if (ht [hv].value == NULL) {
		ht [hv].value = s;
		ht [hv].module = mod;
		ht [hv].nxt = NULL;
	} else {
		r = (struct hte *) malloc (sizeof (struct hte));
		/* No checking for NULL. Shura. 27.01.98 */
		if (NULL == r) {
		  exit (1);
		}
		r -> nxt = ht [hv].nxt;
		ht [hv].nxt = r;
		r -> value = s;
		r -> module = mod;
	}
	return 0;
}

int rd_pr_module (m)
		/* print the name of the module m and the 
			names of all functions in it.  */
	struct modlist *m;
	{
	char **q;
	char *s;

	fprintf (rdout, " Module %s:\n", m -> module);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, " Module %s:\n", m -> module);

	q = m -> funcs;
	while ((s = *q) != NULL)
		{
		/*fprintf (rdout, "\t%16s\n", (s-MAXWS));*/
		for (s -= 2; * s != 0; s --);
		s ++;
		fprintf (rdout, "\t%s\n", s);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "\t%s\n", s);

		q++;
		};
	putc ('\n', rdout);
	if (fp_debugInfo != NULL) putc ('\n', fp_debugInfo);

	return 0;
	}

/* print information about module *s. */
int rd_shwmod (char * s) {
	struct modlist *m;

	if (ht_flag == 0) rd_create_ht ();
	for (m = module_list; m != NULL; m = m -> next) {
		int i;

		if (strcmp (s, m -> module) == 0) break;
		for (i = 0; s [i] != 0 && m -> module [i] != 0; i ++) {
			if (isalpha (s [i]) && isalpha (m -> module [i])) {
				if (toupper (s [i]) == toupper (m -> module [i])) {
					break;
				}
			} else {
				break;
			}
		}
		if (s [i] == 0 && m -> module [i] == 0) break;
	}
	if (m == NULL) {
		fprintf (rdout, "Module %s does not exist.\n", s);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Module %s does not exist.\n", s);

	} else rd_pr_module (m);
	return 0;
}

/* print information about all modules. */
int rd_shmodall (void) {
	struct modlist *m;

	if (ht_flag == 0) rd_create_ht ();
	m = module_list;
	while (m != NULL) {
		rd_pr_module (m);
		m = m -> next;
	}
	return 0;
}

/* print only the names of all modules. */
int rd_list_modules (void) {
	struct modlist *m;

	fprintf (rdout, " Modules are:\n");
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, " Modules are:\n");

	if (ht_flag == 0) rd_create_ht ();
	for (m = module_list; m != NULL; m = m -> next) {
		fprintf (rdout, "\t%s.\n", m -> module);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "\t%s.\n", m -> module);
	}
	fputc ('\n', rdout);
	if (fp_debugInfo != NULL) fputc ('\n', fp_debugInfo);

	return 0;
}

/* Print information about functions. */
int rd_shwfunc (char * s) {
	struct listfun *fl;
	int i, fn;

	fl = rd_implode (s, &fn);
	if (fn == 0) {
		fprintf (rdout, "Function %s is not defined.\n", s);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Function %s is not defined.\n", s);

		return 1;
	}
	fprintf (rdout, "Function %s occurs in %d modules.\n", s, fn);
	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Function %s occurs in %d modules.\n", s, fn);

	for (i = 1; i <= fn; i++) {
		/*
		fprintf (rdout, "\t%2d\t%-16s\taddress %8lx.\n", 
			i, fl -> module, (unsigned long)(fl -> function));
		*/
		fprintf (rdout, "\t%2d\t%s\taddress %8lx.\n", i, fl -> module, (unsigned long)(fl -> function));
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "\t%2d\t%s\taddress %8lx.\n", i, fl -> module, (unsigned long)(fl -> function));

		fl = fl -> next;
	}
	return 0;
}


