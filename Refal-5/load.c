
# include "version.h"
# include "rasl.h"
# include "decl.h"
# include "ldecl.h"
# include "macros.h"
# include "ifunc.h"

#ifdef FOR_OS_LINUX
#	define MAX_PATH_LEN 1024
#	include <sys/stat.h>
#	include <unistd.h>
#else
#	ifdef FOR_OS_DOS
#		define _MAX_PATH 1024
#	endif
#	include <sys/types.h>
#	include <sys/stat.h>
#endif

/* REFAL-5 SYSTEM: Feb. 24, 1987.  D.T. */

extern unsigned long ul_limitSizeCode;

# define MDEBUG 0

# if MDEBUG
	FILE *fDebug;
# endif

static FILE * search_file_path (char *, char *);

struct local_path {
	char * cp_path;
	struct local_path * next;
	struct local_path * back;
};
#define NEW_LOCAL_PATH(lcp,size)                                                       \
{                                                                                      \
	if (NULL == ((lcp) = (struct local_path *) malloc (sizeof (struct local_path)))) { \
		fprintf (stderr, "Cannot allocate memory for temporary path\n");               \
		exit (1);                                                                     \
	}                                                                                  \
	(lcp) -> next = (lcp) -> back = (lcp);                                             \
	if (NULL == ((lcp) -> cp_path = (char *) malloc ((size)))) {                       \
		fprintf (stderr, "Cannot allocate memory for local path\n");                   \
		exit (1);                                                                     \
	}                                                                                  \
}
#define DELETE_LOCAL_PATH(lcp_del)  \
{                                   \
	free ((lcp_del) -> cp_path);    \
	free ((lcp_del));               \
}
#define FREE_LOCAL_PATHS(lcp_beg)       \
{                                       \
	struct local_path * lcp;            \
                                        \
	lcp_beg -> back -> next = NULL;     \
	while (NULL != (lcp = (lcp_beg))) { \
                                        \
		lcp_beg = lcp -> next;          \
		DELETE_LOCAL_PATH (lcp);        \
	}                                   \
}
#define INSERT_LOCAL_PATH(lcp_beg,lcp_ins)  \
{                                           \
	(lcp_ins) -> next = (lcp_beg);          \
	(lcp_ins) -> back = (lcp_beg) -> back;  \
	(lcp_beg) -> back -> next = (lcp_ins);  \
	(lcp_beg) -> back = (lcp_ins);          \
}
#define DROP_LOCAL_PATH(lcp_drop)                    \
{                                                    \
	(lcp_drop) -> back -> next = (lcp_drop) -> next; \
	(lcp_drop) -> next -> back = (lcp_drop) -> back; \
	(lcp_drop) -> next = (lcp_drop) -> back = NULL;  \
}

#define MAX_ENVS 256
static char * cap_ref5env [MAX_ENVS];
static void read_config (void);

/* Only for ri_checkentry */
#include "bif_lex.h"
extern struct bitab bi [];
extern long nbi;
/* End. */

static char *
ri_checkentry (MODULE * mlist, MODULE * m) {
	MODULE * m_w;
	int i;

	if (mlist == NULL) return NULL;
	for (i = 0; i < m -> np_size; i ++) {
		int i_w;

		/* Skip standart functions */
		for (i_w = 0; i_w < nbi; i_w ++) {
			if (0 == strcmp (m -> entry_pts [i].name, bi [i_w].fname)) {
				break;
			}
		}
		if (i_w < nbi) continue;

		for (m_w = mlist; m_w != NULL; m_w = m_w -> next) {
			for (i_w = 0; i_w < m_w -> np_size; i_w ++) {
				if (0 == strcmp (m_w -> entry_pts [i_w].name, m -> entry_pts [i].name)) {
					return m -> entry_pts [i].name;
				}
			}
		}
	}
	return NULL;
}

char *ri_load (void) {
	int i, done;
	char lname [FILENAME_MAX], *s, *ch;
	FILE *fdtmpr;
	MODULE *module, *modlist, *next_module;
	LOAD_TABLE *table;
	int funnum, modnum, ntfnum;
	int unres;
	struct local_path * lcp_beg = NULL;

	long l_modlen;

	/*read_config ();*/

	/* garg [0] contains the name of the interpreted file. */
	if (garg [0] == NULL) {
		if (NULL == (garg [0] = (char *) malloc (FILENAME_MAX * FOPEN_MAX + FOPEN_MAX))) {
			fprintf (stderr, "No memory for RASL file name.\n");
			exit (1);
		}
		ri_inquire ("RASL file name (s)? ", garg [0], FILENAME_MAX * FOPEN_MAX + FOPEN_MAX);
	}
	s = garg [0];
	modlist = NULL;

	/* 2. Read in the headers of files. */
# if MDEBUG
	fprintf (fDebug, "Reading headers: %s\n", filenames);
# endif

	done = 0;
	while(!done) {

		/* 2a. Get the next file name into lname. */ 
		if (0 == * s) break;
		if (-1 == (i = strcspn (s, "+()"))) {
			/* Error */
			fprintf (stderr, "ERROR in argument string\n");
			exit (1);
		}
		if (i == 0) {
			if (')' == * s) {
				if (lcp_beg -> back == lcp_beg) {
					DELETE_LOCAL_PATH (lcp_beg);
					lcp_beg = NULL;
				} else {
					struct local_path * lcp_w = lcp_beg -> back;

					DROP_LOCAL_PATH (lcp_w);
					DELETE_LOCAL_PATH (lcp_w);
				}
			}
			s ++;
			continue;
		} else if (0 == s[i]) {
			/* End */
			done  = 1;
			strcpy (lname, s);
		} else if ('(' == s [i]) {
			struct local_path * lcp_w;

			NEW_LOCAL_PATH (lcp_w, i + 1);
			if (NULL == lcp_beg) {
				lcp_beg = lcp_w;
			} else {
				INSERT_LOCAL_PATH (lcp_beg, lcp_w);
			}
			strncpy (lcp_w -> cp_path, s, i);
			lcp_w -> cp_path [i] = 0;
			s += i + 1;
			continue;
		} else if (')' == s [i] || '+' == s [i]) {
			int i_len;
			lname [0] = 0;

			if (lcp_beg != NULL) {
				struct local_path * lcp_w;

				for (lcp_w = lcp_beg; lcp_w -> next != lcp_beg; lcp_w = lcp_w -> next) {

#if defined (FOR_OS_LINUX) || defined (FOR_OS_WINDOWSNT)
					sprintf (lname + strlen (lname), "%s/", lcp_w -> cp_path/*ca_dir*/);
#elif defined (FOR_OS_DOS)
					sprintf (lname + strlen (lname), "%s\\", lcp_w -> cp_path/*ca_dir*/);
#endif
				}

#if defined (FOR_OS_LINUX) || defined (FOR_OS_WINDOWSNT)
				sprintf (lname + strlen (lname), "%s/", lcp_w -> cp_path/*ca_dir*/);
#elif defined (FOR_OS_DOS)
				sprintf (lname + strlen (lname), "%s\\", lcp_w -> cp_path/*ca_dir*/);
#endif
			}
			i_len = strlen (lname);
			strncat (lname, s, i);
			lname [i_len + i] = 0;
			if (')' == s [i]) {
				if (lcp_beg -> back == lcp_beg) {
					DELETE_LOCAL_PATH (lcp_beg);
					lcp_beg = NULL;
				} else {
					struct local_path * lcp_w = lcp_beg -> back;

					DROP_LOCAL_PATH (lcp_w);
					DELETE_LOCAL_PATH (lcp_w);
				}
			}
			s += i + 1;
		}

		/* 2b. See if it has an extention. */
		ch = strrchr (lname, '.');
		if (NULL != ch) {
			if (! ((ch [1] == 'r' || ch [1] == 'R') && (ch [2] == 's' || ch [2] == 'S') && (ch [3] == 'l' || ch [3] == 'L') && ch [4] == '\0')) {
				strcat (lname, ".rsl");
			}
			/*
			if (0 != strcmp (ch, ".rsl")) {
				strcat (lname, ".rsl");
			}
			*/
		} else {
			strcat (lname, ".rsl");
		}

# ifdef IBM370
		if (ch != NULL) *ch = '\0';
		strcat (lname, " rsl (bin lrecl 1");
# else 
		/*if (ch == NULL) strcat (lname, ".rsl");*/
# endif

# if MDEBUG
		fprintf (fDebug, "File: %s\n", lname);
# endif

		fdtmpr = fopen (lname, "rb");
		if (fdtmpr == NULL) {
			/*
			if (NULL != strrchr (lname, '\\') || NULL != strrchr (lname, '/')) {
				fprintf (stderr, "Unable to open file %s\n", lname);
				exit (1);
			}
			*/
			if (NULL == (fdtmpr = search_file_path (lname, "rb"))) {
				fprintf (stderr, "Unable to open file %s\n", lname);
				exit (1);
			}
		}

		/* 2c. Allocate memory for module descriptor. */
		module = (MODULE *) malloc (sizeof (MODULE));
		if (module == NULL) {
			fprintf (stderr, "No more memory.\n");
			exit (1);
		}

		/* 2d. Read the header of that file. */
		ri_readhdr (fdtmpr, module);
		
		{
			char * cp_w;

			if (NULL != (cp_w = ri_checkentry (modlist, module))) {
				fprintf (stderr, "Run time ERROR: Same function \'%s\' are in difference modules.\n", cp_w);
				exit (1);
			}
		}
		module -> next = modlist;
		modlist = module;
	}

# if MDEBUG
	fprintf (fDebug, "Headers are read.\n");
# endif
	
	if (NULL != lcp_beg) {
		fprintf (stderr, "WARNING: Balance of brackets in command line is wrong!!!\n");
		FREE_LOCAL_PATHS (lcp_beg);
		lcp_beg = NULL;
	}

	/* 3. Reverse the list of modules, counting the number of 
			modules, total number of functions and number of
			entry functions. */

		if (modlist == NULL) ri_lerror (10);
		next_module = modlist -> next;
		modlist -> next = NULL;
		modnum = 1;
		funnum = modlist -> lf_size;
		ntfnum = modlist -> np_size;
		l_modlen = strlen (modlist -> module_name) + 1;
		while (next_module)
			{
			module = next_module -> next;
			next_module -> next = modlist;
			modlist = next_module;
			next_module = module;
			modnum ++;
			funnum += modlist -> lf_size;
			ntfnum += modlist -> np_size;
			l_modlen += strlen (modlist -> module_name) + 1;
			};

# if MDEBUG
		fprintf (fDebug, "Number of modules %d, functions %d\n", modnum, funnum);
# endif

	/* 4. Resolve external references and get the entry address. */
		/* As far as resolving external references is concerned, nothing
			should be done now: all addresses are in the entry tables. */
		/* The global ENTRY point should be label GO in the first module
			pointed by modlist. */

		ch = NULL;
		for (module = modlist; module != NULL; module = module -> next) {
			table = module -> entry_pts;
			for (i = 0; i < module -> np_size; i++) {
				if (table [i].name [0] == 'G' && (table [i].name [1] == 'O' || table [i].name [1] == 'o') && 
					table [i].name [2] == 0) {
					ch = table [i].addr;
					break;
				}
			}
			if (ch != NULL) break;
		}
		if (ch == NULL) ri_lerror (2);

# if MDEBUG
			fprintf (fDebug, "Entry point = %ld\n", ch);
# endif

	/* 5. Allocate storage for the list of local functions and create it.
			Also create the list of module names. */

		/* 5a. Allocate storage for the list of modules. */
			/*module_table = malloc (modnum * MAXWS);*/
			/*module_table = malloc (modnum * FILENAME_MAX);*/
			module_table = malloc (l_modlen);
			num_modules = modnum;
			if (module_table == NULL) ri_lerror (9);

		/* 5b. Allocate storage for the local function list. */
			local_functions = (char **) malloc (sizeof (char *) * 
				 (1 + 2*modnum + funnum));
			num_local_functions = funnum;
			if (local_functions == NULL) ri_lerror (9);

		/* 5c. Create the function and module tables. */
			ri_mkmodlist (modlist);

		/* 5d. Allocate storage for the list of entry functions. */
			num_entry = ntfnum;
			entry_functions = (char **) malloc (sizeof (char *) * 
				 (1 + ntfnum));
			if (entry_functions == NULL) ri_lerror (9);

		/* 5e. Create the list of entry functions. */
			ri_mkentlist (modlist);
			 

	/* 6. Load the codes for all modules. */
# if MDEBUG
	fprintf (fDebug, "loading ...\n");
# endif
	unres = 0;
	for (module = modlist; module != NULL; module = module -> next) 
		unres += ri_loadcode (module, modlist);

# if MDEBUG
	fprintf (fDebug, "... loading done\n");
# endif
	if (unres != 0)
		{
		fprintf (stderr, "%d unresolved references were found.", unres);

			/* if flags don't contain '-i' option exit. */
		if (strchr (flags, 'i')) fprintf (stderr, "\n");
		else 
			{
			fprintf (stderr, " Aborted\n");
			exit (1);
			}
		}

# if MDEBUG 
	{ char **ffff; fprintf (fDebug, "The table of local functions. %ld\n", local_functions);
		for (ffff = local_functions, 
			fprintf (fDebug, "Module: %ld: *%s*\n", *ffff, *ffff), ffff++; 
			*ffff != NULL; ffff ++)
			fprintf (fDebug, "Symbol: %ld ***%s***\n", *ffff, (*ffff)-MAXWS); };
# endif

	/* 7. Free the allocated memory and return. */
		/* ri_freemod frees the memory and returns the pointer to
			the next module descriptor. */
	while (modlist != NULL) modlist = ri_freemod (modlist);

	return ch;
	}

char ** ri_mkfunlist (MODULE *m, char *ma, char **fa) {
	int i, lf;
	char *base;
	FILE *fdtmpr;
	long k;

	lf = (int) (m -> lf_size);
	fdtmpr = m -> from;
	base = m -> base;

	/* 1. Copy the module name into its place. */
	strcpy (ma, m -> module_name);

	/* 2. Write the module address first. */
	*fa++ = ma;

	/* 3. For each function save its address in the list. */
	for (i = 0; i < lf; i ++) {
		read_long_to_mem (k);
		*fa++ = base + k;
	}

	/* 4. Write the terminating NULL and return. */
	*fa++ = NULL;
	return fa;
}

int ri_mkmodlist (MODULE *m) {
	char *maddr, **faddr;

	for (maddr = module_table, faddr = local_functions; m != NULL; m = m -> next, maddr += strlen (maddr) + 1) {
		m -> local_table = faddr;
		faddr = ri_mkfunlist (m, maddr, faddr);
	}
	*faddr = NULL;
	return 0;
}


int ri_mkentlist (m)
	MODULE *m;
	{
	LOAD_TABLE *table;
	char **faddr;
	int i;

	for (faddr = entry_functions; m != NULL; m = m -> next)
		{
		table = m -> entry_pts;
		for (i = 0; i < m -> np_size; i++) *faddr ++ = table [i].addr;
		};
	*faddr = NULL;
	return 0;
	}

MODULE *ri_freemod (m)
	MODULE *m;
	{
	MODULE *next;

	if (m == NULL) return NULL;
	fclose (m -> from);
	if (m -> cs_table) free ((char *) (m -> cs_table));
	if (m -> entry_pts) {
		int i;

		for (i = 0; i < m -> np_size; i ++) {
			free (m -> entry_pts [i].name);
		}
		free ((char *) (m -> entry_pts));
	}
	if (m -> ext_table) {
		int i;

		for (i = 0; i < m -> xt_size; i ++) {
			free (m -> ext_table [i].name);
		}
		free ((char *) (m -> ext_table));
	}
	next = m -> next;
	free ((char *) m);
	return next;
	}



MODULE *ri_readhdr (fdtmpr, m)
	/* This function reads in the header of the given file and stores 
		the information in the structure pointed by m. 3-8-1987. */

	FILE *fdtmpr;
	MODULE *m;

	{
	long z, /* Size of the code */
		np_size,  /* size of entry function table. */
		xt_size,  /* size of external function table. */
		cs_size,  /* size of compound symbol table. */
		lf_size,  /* size of compound symbol table. */
		knt,     /* counter */
		offset;		/* entry offset */
	unsigned int zz = 0;
	int i;
	LOAD_TABLE *table;
	COMPSYM_TABLE *comp_symbols;
	char *base, lname [MAXWS];



	/* 1. Read in the title and sizes of the tables and code. */

	/*for (i = 0; i < MAXWS; i ++) read_byte_to_mem (m -> module_name [i]);*/
	for (i = 0; i < FILENAME_MAX; i ++) {
		read_byte_to_mem (m -> module_name [i]);
		if (m -> module_name [i] == '\0') break;
	}
	read_long_to_mem (z);
	read_long_to_mem (np_size);
	read_long_to_mem (xt_size);
	read_long_to_mem (cs_size);
	read_long_to_mem (lf_size);
# if MDEBUG
	fprintf (fDebug, "Title: %s\n", m -> module_name);
	fprintf (fDebug, "CODE = %ld, ENTRY = %ld, EXTRN = %ld, CS = %ld\n", 
		z, np_size, xt_size, cs_size);
# endif



	/* 2. Check that these numbers are reasonable. */
	if (z <= 0L) ri_lerror (1);
	else m -> size = z;
	if (np_size <= 0L) ri_lerror (2);
	else m -> np_size = np_size;
	if (lf_size <= 0L) ri_lerror (2);
	else m -> lf_size = lf_size;
	if (feof (fdtmpr)) ri_lerror (5);



	/* 3. Allocate memory for code. */
	if (z > ul_limitSizeCode) ri_lerror (7);
	else zz = z;
	base = malloc (zz);
	if (base == NULL) ri_lerror (7);
	m -> base = base;
# if MDEBUG
	fprintf (fDebug, "Loader: BASE = %ld LENGTH = %d\n", base, zz);
# endif



	/* 4. Create the Entry Points table. */
	table = (LOAD_TABLE *) malloc ((unsigned int)(np_size * sizeof (LOAD_TABLE)));
	if (table == NULL) ri_lerror (9);
	for (knt = 0L; knt < np_size; knt ++)
		{
		char ca_buf [MAXWS];
		for (i = 0; i < MAXWS; i++) {
			/*
			read_byte_to_mem (table [knt].name [i]);
			if (table [knt].name [i] == 0) {
				break;
			}
			*/
			read_byte_to_mem (ca_buf [i]);
			if (ca_buf [i] == 0) {
				break;
			}
		}
		if (NULL == (table [knt].name = (char *) malloc (i + 1))) {
			ri_lerror (9);
		}
		strcpy (table [knt].name, ca_buf);

		read_long_to_mem (offset);
		table [knt].addr = base+offset;
		};
	
	m -> entry_pts = table;
# if MDEBUG
	fprintf (fDebug, "Entry table created.\n");
	for (knt = 0L; knt < np_size; knt ++) 
		fprintf (fDebug, "%ld: %s %ld\n", knt, table [knt].name, 
				table [knt].addr);
	fprintf (fDebug, "Total %ld entry points.\n", np_size);
# endif



	/* 5. Get the external functions. */
	if (xt_size < 0L) ri_lerror (1);
	else m -> xt_size = xt_size;
		/* We don't really need the external function table, 
			but we may need it in the future, so we simply skip it.  */
	for (zz=0,knt = 0L; knt < xt_size; knt ++) {
		for (i = 0; i < MAXWS; i++) {
			read_byte_to_mem (zz);
			if (zz == 0) {
				break;
			}
		}
	}
	m -> ext_table = NULL;
# if MDEBUG
	fprintf (fDebug, "External table skipped.\n");
# endif


	/* 6. Create dynamic compound symbol table. */
	m -> cs_size = cs_size;
	if (cs_size > 0L) {
		comp_symbols = (COMPSYM_TABLE *) malloc ((unsigned int)
			(sizeof (COMPSYM_TABLE) * cs_size));
		if (comp_symbols == NULL) ri_lerror (9);
		/* Read in the compound symbols. (they are in reverse order) */
		for (knt = 1L; knt <= cs_size; knt ++) {
			for (i = 0; i < MAXWS; i++) {
				read_byte_to_mem (lname [i]);
				if (lname [i] == 0) {
					break;
				}
			}
			comp_symbols [cs_size - knt].name = ri_cs_impl (lname);
			comp_symbols [cs_size - knt].ident = cs_size - knt;
		}
		m -> cs_table = comp_symbols;
	} else {
		m -> cs_table = NULL;
	}

# if MDEBUG
	fprintf (fDebug, "Print the compound symbol table.\n");
	for (knt = 0L; knt < cs_size; knt ++) 
		fprintf (fDebug, "%ld: %s %ld\n", knt, comp_symbols [knt].name, 
				comp_symbols [knt].ident);
	fprintf (fDebug, "Total %ld compound symbols.\n", cs_size);
# endif

	/* 7. Miscellaneous stuff, and return. */
	m -> from = fdtmpr;
	m -> local_table = NULL;
	return m;
	}


int ri_lerror (code)
	int code;
	{
	switch (code)
		{
		case 1:
			fprintf (stderr, "Illegal format in file.\n");
			break;
		case 2:
			fprintf (stderr, "No entry point.\n");
			break;
		case 3:
			fprintf (stderr, "Too many entry points. Not yet implemented.\n");
			break;
		case 4:
			fprintf (stderr, "EOF not reached.\n");
			break;
		case 5:
			fprintf (stderr, "Unexpected EOF encountered.\n");
			break;
		case 6:
			fprintf (stderr, "Illegal Compound symbol.\n");
			break;
		case 7:
			fprintf (stderr, "Too large block of code.\n");
			break;
		case 8:
			fprintf (stderr, "General error loading code.\n");
			break;
		case 9:
			fprintf (stderr, "Not enough memory.\n");
			break;
		case 10: 
			fprintf (stderr, "No files are found.\n");
			break;
		default:
			fprintf (stderr, "Error %d\n", code);
			break;
		};
	fprintf (stderr, "Loading aborted.\n");
	exit (1);
	return code;
	}

int ri_loadcode (m, mlist)
	MODULE *m, *mlist;
	{
	long k, size, z;
	LOAD_TABLE *table;
	FILE *fdtmpr;
	unsigned char opcode, c, d;
	char lname [MAXWS], *base, *ldptr, *addr;
	COMPSYM_TABLE *cstable;
	MODULE *mtmp;
	int i, unresolved;

# if MDEBUG
	fprintf (fDebug, "Module %s\n", m -> module_name);
# endif

	z = 0L;
	unresolved = 0;
	ldptr = base = m -> base;
	fdtmpr = m -> from;
	cstable = m -> cs_table;
	size = m -> size;
	while (read_byte_to_mem (opcode) == 1) {
		if (z >= size) ri_lerror (4);
		switch (opcode) {

		case ACT_EXTRN:
			/* This RASL instruction takes an address of a function as an argument. */
			for (i = 0; i < MAXWS; i++) {
				read_byte_to_mem (lname [i]);
				if (lname [i] == 0) {
					break;
				}
			}
			z += sizeof (char) + sizeof (char *);
			/* Search for this function in the list of entry tables. */
			addr = NULL;
			for (mtmp = mlist; (addr == NULL) && mtmp; mtmp = mtmp -> next) {
				table = mtmp -> entry_pts;
				for (i = 0; i < mtmp -> np_size; i++) {
					if (strcmp (table [i].name, lname) == 0) {
						addr = table [i].addr;
						break;
					}
				}
			}
			if (addr == NULL) {
				fprintf (stderr, "Unresolved external reference %s\n", lname);
				unresolved++;
				addr = IMP_;
			}
			write_byte_to_mem (ACT1);
			write_long_to_mem (addr);
			break;

		case ACT1: case TRAN: case ECOND:
			/* These RASL operators take as argument an address. */
			write_byte_to_mem (opcode);
			read_long_to_mem (k);
			addr = base+k;

# if MDEBUG
			fprintf (fDebug, "Load (%d) base = %ld k= %ld addr= %ld\n", opcode, base, k, addr);
# endif
			write_long_to_mem (addr);
			z += sizeof (char) + sizeof (char *);
			break;

		case CSYM: case CSYMR: case NCS:
			/* These RASL operators take a compound symbol as an argument. */
			write_byte_to_mem (opcode);
			read_long_to_mem (k);
			z += sizeof (char) + sizeof (char *);
			if (k >= m -> cs_size) ri_lerror (1);
			else write_long_to_mem (cstable [k].name);
			break;

		case NSYM: case NSYMR: case NNS: case BUILT_IN:
			/* These RASL operators require a (long) number as a parameter */
			write_byte_to_mem (opcode);
			read_long_to_mem (k);
			write_long_to_mem (k);
			z += sizeof (char) + sizeof (long);
			break;

		case BUILT_IN1:	/* Builtin function call with an argument. */
			write_byte_to_mem (BUILT_IN1);
			read_long_to_mem (k);	/*  read the address. (It is zero) */
			read_long_to_mem (k);	/*  Read the number. */
			write_long_to_mem (m -> local_table);	/* address of the local table. */
			write_long_to_mem (k);
			z +=  sizeof (char) + sizeof (char *) + sizeof (long);
			break;

			/* No arguments for these RASL operators. */
		case BL: case BLR: case BR: case CL: case EMP: case EST: case PLEN:
		case PLENS: case PLENP: case PS: case PSR: case TERM: case TERMR:
		case LEN: case LENP: case VSYM: case VSYMR: case OUTEST: case POPVF:
		case PUSHVF: case STLEN:
			write_byte_to_mem (opcode);
			z ++;
			break;

			/* These RASL operators require a byte (character) as parameter. */
		case SYM: case SYMR: case LENS: case NS:
			write_byte_to_mem (opcode);
			read_byte_to_mem (c);
			write_byte_to_mem (c);
			z += sizeof (char) + sizeof (char);
			break;

			/* These RASL operators take one operand of size 4 bytes. */
		case MULE: case MULS: case TPLE: case TPLS:
			write_byte_to_mem (opcode);
			read_long_to_mem (k);
			write_long_to_mem (k);
			z += sizeof (char) + sizeof (long);
			break;

			/* These RASL operators take one operand of size 1 byte. */
		/*case MULE: case MULS:*/ case OEXP: case OEXPR: case OVSYM: case OVSYMR:
		case RDY: case LENOS:/* case TPLE: case TPLS:*/
			write_byte_to_mem (opcode);
			read_byte_to_mem (c);
			write_byte_to_mem (c);
			z += sizeof (char) + sizeof (char);
			break;

		case SETB:
			/* This RASL operator require two operands of size 1 byte. */
			write_byte_to_mem (opcode);
			read_long_to_mem (k);
			write_long_to_mem (k);
			read_long_to_mem (k);
			write_long_to_mem (k);
			z += sizeof (char) + 2 * sizeof (long);
			break;

		/* These RASL operators take 1 byte and a variable number of bytes as parameters. */
		case SYMS: case SYMSR: case TEXT:
			write_byte_to_mem (opcode);
			read_byte_to_mem (d);
			write_byte_to_mem (d);
			for (i = 0; i < d; i++) {
				read_byte_to_mem (c);
				write_byte_to_mem (c);
			}
			z += sizeof (char) + sizeof (char) + d*sizeof (char);
			break;


		/* These operators define labels. */
		case LBL: case LABEL: case L: case E:
			read_byte_to_mem (c);
			write_byte_to_mem (c);
			for (i = 0; i < MAXWS; i++) {
				read_byte_to_mem (c);
				write_byte_to_mem (c);
				if (c == 0) break;
			}
			z += /*MAXWS*/(i+2) * sizeof (char);
			break;
  
		default:
			fprintf (stderr, "%d : Strange Opcode\n", opcode);
			break;
		}
	}

# if MDEBUG
	fprintf (fDebug, "z = %ld size = %ld.\n", z, size);
	fprintf (fDebug, "Base and Load pointer = %ld %ld\n", base, ldptr);
	/****
	fprintf (stderr, "Wanna dump? ");
	while ((i=getchar ()) == '\n');
	if (i == 'y') dumpcode (base, size);
	****/
# endif
	return unresolved;
}


int ri_init ()
	{
	LINK *nqup;
	char *go;

 /* Call Refal loader to load the interpretation file (s).
 	Function ri_load returns the address of the Entry point in memory. */

# if MDEBUG
	fDebug = fopen ("refload.dbg", "wt");
	if (fDebug == NULL)
		{
		fprintf (stderr, "can't open refload.dbg\n");
		exit (1);
		};
# endif

	go = ri_load ();
	if (go == NULL) exit (1);

	lfm = ri_fmout();
	if (lfm == NULL)
		{
		fprintf (stderr, "Can\'t get memory. exit\n");
		exit (1);
		};

 /*  ri_init initializes the view field.   */

 /*  initialize the transition stack.  */
	st [0].b1  =  b1;    /* left boundary */
	st [0].b2  =  b2;     /* right boundary */
	st [0].nel =  3L;    /* new element  */
	st [0].ra  =  IMP_;	/* return address   */

 /* create links quasipoint (quap) */

	quap = lfm;  all; nqup = lfm; all;
	quap->foll = nqup;
	precp = quap;
	nextp = quap;

 /*  initialize some other variables.   */
	nst = -1;   /*  number of steps.  */
	curk = 1;   /* current number of k-signs. */
	stoff = 1; /* transition stack offset.	 */
	/* 704 line. Was ... = NULL. Shura. 29.01.98 */
	teoff = 0; /* table of elements offset.	 */
	tel = & (te [teoff]); /* pointer to an element 
				in the table of elements */
	pa = NULL;

 /*	initialize the stock pointer to a set of empty brackets.	*/
	stock = lfm;
	all;
	stock -> ptype = 0;
	lfm -> ptype = 1;
	weld (stock, lfm);
	con (stock, lfm);
	all;
	stock -> prec = NULL;
	stock -> foll -> foll = NULL;

 /*  the view field will have the form: </STOP$/ </GO/ >>  */
	b = lfm;
	all;
	bl;
	vf = b;
	blr;
	act1 (go);
	br;
	vfend = b;
	act1 (STOP_);
	est;

# if MDEBUG
	fprintf (fDebug, "\nBy the way: STOP$ = %ld, IMP$$ = %ld, GO = %ld\n\n", 
		STOP_, IMP_, go);
	fclose (fDebug);
# endif

	return 0;
	}

	/* ask the user something. */
int ri_inquire (char * m, char * r, int i_len) {

/*  printf (VERSION, "System", _refal_build_version);*/
  *r = '\0';
  while (*r == '\0') {
	  int i;

    fprintf (stderr, "%s", m);
    if (NULL == fgets (r, i_len, stdin)) {
      fprintf (stderr, "Cannot read filename from standart input\n");
      return -1;
    }
    for (i = strlen (r) - 1; i >= 0 && r [i] == '\n'; i --);
    r [i + 1] = 0;
  }
  return 0;
}

static char * get_local_env (char * cp_var) {
	int i;

	for (i = 0; cap_ref5env [i] != NULL; i ++) {
		if (strcmp (cap_ref5env [i], cp_var) == 0) {
			return (cap_ref5env [i] + strlen (cap_ref5env [i]) + 1);
		}
	}
	return NULL;
}

static FILE * open_file (char * cp_path, char * cp_name, char * cp_mode) {
	FILE * fp;
	char * cp;
	char ca_buf [FILENAME_MAX];

#if defined (FOR_OS_LINUX) || defined (FOR_OS_OS2)
	for (; NULL != (cp = strchr (cp_path, ':')); cp_path = cp + 1) {
#elif defined (FOR_OS_DOS) || defined (FOR_OS_WINDOWSNT) 
	for (; NULL != (cp = strchr (cp_path, ';')); cp_path = cp + 1) {
#endif
		strncpy (ca_buf, cp_path, cp - cp_path);
		ca_buf [cp - cp_path] = 0;
		if (ca_buf [cp - cp_path - 1] != '\\' && ca_buf [cp - cp_path - 1] != '/') {
#if defined (FOR_OS_LINUX) || defined (FOR_OS_WINDOWSNT) || defined (FOR_OS_OS2)
			strcat (ca_buf, "/");
#elif defined (FOR_OS_DOS)
			strcat (ca_buf, "\\");
#endif
		}
		strcat (ca_buf, cp_name);
		if (NULL != (fp = fopen (ca_buf, cp_mode))) {
			return fp;
		}
	}
	if (cp_path [0] != 0) {
		strcpy (ca_buf, cp_path);
		if (cp_path [strlen (cp_path) - 1] != '\\' && cp_path [strlen (cp_path) - 1] != '/') {
#if defined (FOR_OS_LINUX) || defined (FOR_OS_WINDOWSNT)
			strcat (ca_buf, "/");
#elif defined (FOR_OS_DOS)
			strcat (ca_buf, "\\");
#endif
		}
		strcat (ca_buf, cp_name);
		return fopen (ca_buf, cp_mode);
	}
	return NULL;
}

static FILE * search_file_path (char * cp_name, char * cp_mode) {
	char * cp_ref5rsl;
	FILE * fp;

	if (NULL != (cp_ref5rsl = get_local_env ("REF5RSL"))) {
		if (NULL != (fp = open_file (cp_ref5rsl, cp_name, cp_mode))) {
			return fp;
		}
	}
	if (NULL == (cp_ref5rsl = getenv ("REF5RSL"))) {
		return fopen (cp_name, cp_mode);
	}
	return open_file (cp_ref5rsl, cp_name, cp_mode);
}

/* Aruments:
 *   1) Buffer for saving of a variable name.
 *   2) Pointer to FILE (config file)
 *
 * Result:
 *   '#', '\n', EOF, '='
 */
static int get_varname (char ca_nm [], FILE * fp) {
	int i = 0, i_c = getc (fp);

	ca_nm [0] = 0;
	if (i_c == '#' || i_c == ':') {
		/* comments */
		return '#';
	} else if (i_c == '\n' || i_c == EOF) {
		return i_c;
	}
	do {
		switch (i_c) {
		case '#': case '\n': case EOF:
			if (i != 0) {
				ca_nm [i] = 0;
				fprintf (stderr, "WARNING: Not value for variable \'%s\'\n", ca_nm);
			}
			return (i_c);

		case ' ':
		case '\t':
			break;

		case '=':
			ca_nm [i] = 0;
			return '=';

		default:
			if (isalnum (i_c)) {
				ca_nm [i ++] = i_c;
			} else {
				fprintf (stderr, "Unexpected character is used in variable name\n");
				ca_nm [0] = 0;
				return '#';
			}
		}
		i_c = getc (fp);
	} while (i_c != EOF && i_c != '\n');
	ca_nm [i] = 0;
	return i_c;
}

/* Aruments:
 *   1) Buffer for saving of a value of a variable.
 *   2) Pointer to FILE (config file)
 *
 * Result:
 *   '#', '\n', EOF
 */
static int get_values (char ca_val [], FILE * fp) {
	int i, i_c;

	ca_val [0] = 0;
	for (i = 0, i_c =  getc (fp); EOF != i_c && '\n' != i_c; i_c = getc (fp)) {
		switch (i_c) {
		case '#':
		case '\n':
		case EOF:
			if (i == 0) {
				fprintf (stderr, "WARNING: Value of the variable is empty\n");
			}
			ca_val [i] = 0;
			return i_c;

			/* white space */
		case ' ': case '\t':
			break;

			/* correct special symbols */
		case '.': case ';': case ':': case '-': case '_': case ',': case '\\': case '/':
			ca_val [i ++] = i_c;
			break;

		default:
			if (isalnum (i_c)) {
				ca_val [i ++] = i_c;
			} else {
				fprintf (stderr, "WARNING: Uncorrect special symbol \'%c\'(hex %x) will be missed\n", i_c, i_c);
			}
		}
	}
	ca_val [i] = 0;
	return i_c;
}

static int read_comments (FILE * fp) {
	int i;

	while (EOF != (i = getc (fp)) && '\n' != i);
	return i;
}

static FILE * open_config (void) {
	FILE * fp;
	char * cp_path, * cp_w;

#ifdef FOR_OS_WINDOWSNT
	struct _stat st_buf;
#else
	struct stat st_buf;
#endif
		
#ifdef FOR_OS_LINUX
	char ca_dir [MAX_PATH_LEN];
#else
	char ca_dir [_MAX_PATH];
#endif

	if (NULL != (fp = fopen ("refal5.cfg", "rt"))) {
		return fp;
	}
	/* No in a current directory */
	if (NULL == (cp_path = getenv ("PATH"))) {
		fprintf (stderr, "No global variable PATH, but there isn't config file in the current directory\n");
		return NULL;
	}

	for (; 
#ifdef FOR_OS_LINUX
		NULL != (cp_w = strchr (cp_path, ':'));
#else
		NULL != (cp_w = strchr (cp_path, ';'));
#endif
		cp_path = cp_w + 1) {

		strncpy (ca_dir, cp_path, cp_w - cp_path);
		ca_dir [cp_w - cp_path] = 0;
#ifdef FOR_OS_DOS
		strcat (ca_dir, "\\refal5.cfg");
#else
		strcat (ca_dir, "/refal5.cfg");
#endif

#ifdef FOR_OS_WINDOWSNT
		if (0 == _stat (ca_dir, & st_buf)) {
#else
		if (0 == stat (ca_dir, & st_buf)) {
#endif

			/* Found the config file */
			if (NULL == (fp = fopen (ca_dir, "rt"))) {
				fprintf (stderr, "Found config file, but it cannot open it!?\n");
				/*return NULL;*/
			}
			return fp;
		}
	}
	if (0 != cp_path [0]) {
		/* Last directory is rest */
		sprintf (ca_dir, "%s%s", cp_path,
#ifdef FOR_OS_DOS
								"\\refal5.cfg"
#else
								"/refal5.cfg"
#endif
				);

#ifdef FOR_OS_WINDOWSNT
		if (0 == _stat (ca_dir, & st_buf)) {
#else
		if (0 == stat (ca_dir, & st_buf)) {
#endif

			if (NULL == (fp = fopen (ca_dir, "rt"))) {
				fprintf (stderr, "Found config file, but it cannot open it.");
			}
			return fp;
		}
	}
	return NULL;
}

#ifndef MAXSTR
#	define MAXSTR 200
#endif

static void read_config (void) {
	FILE * fp;
	int i_c, i;
	char ca_name [MAXSTR + 1], ca_values [MAXSTR + 1];

	/*
	if (NULL == (fp = fopen ("refal5.cfg", "rt"))) {
 		if (NULL == (fp = search_file_path ("refal5.cfg", "rt"))) {
			fprintf (stderr, "WARNING: No environ file of Refal 5\n");
			return;
		}
	}
	*/
	if (NULL == (fp = open_config ())) {
		fprintf (stderr, "WARNING: No environ file of Refal 5\n");
		return;
	}
	for (i = 0;i_c = get_varname (ca_name, fp);) {
		switch (i_c) {
		case '#':
			if (EOF  == read_comments (fp)) {
				cap_ref5env [i] = NULL;
				fclose (fp);
				return;
			}
		case '\n':
			break;

		case EOF:
			cap_ref5env [i] = NULL;
			fclose (fp);
			return;

		case '=':
			i_c = get_values (ca_values, fp);
			switch (i_c) {
			default:
				fprintf (stderr, "WARNING: Unknown result was returned by GET_VALUES ()\n");
				break;

			case EOF:
			case '\n':
			case '#':
				if (0 != strlen (ca_values)) {
					cap_ref5env [i] = (char *) malloc (strlen (ca_name) + strlen (ca_values) + 2);
					if (NULL == cap_ref5env [i]) {
						fprintf (stderr, "ERROR: Cannot allocate memory for the local refal environ\n");
						fclose (fp);
						cap_ref5env [i] = NULL;
						exit (1);
					}
					strcpy (cap_ref5env [i], ca_name);
					strcpy (cap_ref5env [i] + strlen (ca_name) + 1, ca_values);
					i ++;
				}
				if (EOF == i_c) {
					cap_ref5env [i] = NULL;
					fclose (fp);
					return;
				}
			}
			break;

		default:
			fprintf (stderr, "WARNING: Unknown result of function GET_VARNAME ()\n");
		}
		if (EOF == i_c) {
			cap_ref5env [i] = NULL;
			fclose (fp);
			return;
		}
	} /* end of for */
}
#undef MAXSTR

# if MDEBUG
int dumpcode (start, size)
	char *start;
	long size;
	{
	unsigned char *p;
	int i;

	for (i = 0, p = start; i < size; i++, p++)
		{
		fprintf (fDebug, "%ld %u", p, *p);
		if (isprint (*p)) fprintf (fDebug, "   %c\n", *p);
		else putchar ('\n');
		};
	return 0;
	}
# endif

