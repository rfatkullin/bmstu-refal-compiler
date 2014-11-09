
#ifdef FOR_OS_WINDOWSNT
#	include <sys/timeb.h>
#	include <sys/types.h>
#endif

# include "version.h"
# include "decl.h"
# include "macros.h"
# include "ifunc.h"
# include <time.h>


/* Refal interpreter auxiliary function. */

/* Memory allocation routines. */

# define NUM_OF_LINKS_AT_A_GULP 512
# define RI_MEMORY_LIMIT 16
/* Real space is RI_MEMORY_LIMIT Mbytes */
		
unsigned long ul_limitSizeCode = DEFAULT_CODE_LIMIT * 1024;

/* For copying of trace informations */
extern FILE * fp_debugInfo;

/* interactive stop on an internal error: YES/NO */
static char stop_on_error;	

LINK *ri_getmem (void) {
	register int i;
	register LINK *q, *qn;
	LINK *first;

	/* allocate space for 512 links */
	q = (LINK *) malloc(NUM_OF_LINKS_AT_A_GULP * sizeof(LINK));
	if (q == NULL) return NULL;
	first = qn = q;
	/*	structure it. */
	i = NUM_OF_LINKS_AT_A_GULP;
	while (-- i) q = q -> foll = ++ qn;
	q -> foll = NULL;
	return first;
}

static int memory_used = 0;
LINK *ri_fmout()
	{
	int error = 0;
	LINK *newp = NULL;

	if ((memory_limit > 0) && (memory_used >= memory_limit)) error = 11;
	else if ((newp = ri_getmem()) == NULL) error = 9;
	if (error)
		{
		if (lfm == NULL)
			{
			lfm = malloc(sizeof(LINK));
			if (NULL == lfm) {
			  fprintf (stderr, "No memory for module\n");
			  exit (1);
			}
			lfm->foll = NULL;
			}
			/* clear out the partially formed right side. */
		rdy(0);
		ri_out(2);
		ri_error(error);
		return NULL;
		}
	memory_used++; 
	return newp;
	}

/* Stop function. */
int ri_stop (void) {
/*	if (flags[0] != '\0')
		fprintf(ref_err_file(), "REFAL: NORMAL STOP.\n");*/
	ri_error(0);
	return 0;
}


/* IMP function. */
int ri_imp (void) {
	/* if there is a freezer, freeze and return. */
	if (exists_freeze()) {
		ri_frz (2); 
		return 0;
	}

	/* print message and call error (which calls exit()) */
	fprintf(ref_err_file(), "REFAL ERROR:  RECOGNITION IMPOSSIBLE\n");
	ri_error(12);
	return 0;
}

extern FILE* ref_err_file();

#ifdef FOR_OS_WINDOWSNT
void ri_information (struct _timeb tm, int code) {
#else
void ri_information (long l_time, int code) {
#endif
	FILE * ferr = ref_err_file();

#ifdef FOR_OS_WINDOWSNT
	unsigned long ul_rslt;
	
	ul_rslt = (tm.millitm < whens.millitm)? tm.time - whens.time - 1: tm.time - whens.time;
#else
	l_time -=  whens;
#endif

	if (code != 0 && flags[0] == '\0') {
		char *token, line[100];

                if ( stop_on_error == YES ) {
        		fprintf(ferr, "Do you want to print view-field? [Yn]");
	        	if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Do you want to print view-field? ...\n");

        		fgets(line, sizeof(line)-1, stdin);
	        	token = strtok(line, " \n\t\f");
        		if (token == NULL || *token == 'y' || *token == 'Y') strcpy(flags, "a");
                };
	}
	if (strchr(flags, 'a')) {
		fprintf(ferr, VERSION, "System", _refal_build_version);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, VERSION, "System", _refal_build_version);
	}
	if (code > 0) {
		char * cp_w = actfun - 1;

		fprintf(ferr, "*** Active function: ");
		for (cp_w --; '\0' != * cp_w; cp_w --);
		cp_w ++;
		ri_actput (cp_w, ferr);
		putc ('\n', ferr);
		fprintf(ferr, "*** Active expression:\n");
		ri_putmb (tbel (1), tbel (2), ferr);
		if (fp_debugInfo != NULL) {
			fprintf(fp_debugInfo, "*** Active function: ");
			ri_actput (cp_w, fp_debugInfo);
			fprintf(fp_debugInfo, "\n*** Active expression:\n");
			ri_putmb (tbel (1), tbel (2), fp_debugInfo);
		}
	}
	if (strchr (flags, 'a') || strchr (flags, 'v')) {
		fprintf(ferr, "*** The View Field:\n");
		ri_putmb (vf, vfend, ferr);
		if (fp_debugInfo != NULL) {
			fprintf(fp_debugInfo, "*** The View Field:\n");
			ri_putmb (vf, vfend, fp_debugInfo);
		}
	}
	if (strchr (flags, 'a') || strchr (flags, 'n')) {
		fprintf(ferr, "*** Number of Steps = %2ld\n", nst);
		if (fp_debugInfo != NULL) fprintf(fp_debugInfo, "*** Number of Steps = %2ld\n", nst);
	}
	if (strchr (flags, 't') || strchr (flags, 'a')) {

#ifdef FOR_OS_WINDOWSNT
		fprintf (ferr, "Elapsed system time: %u.%03u seconds\n", ul_rslt,
			(tm.millitm < whens.millitm)? 1000 + tm.millitm - whens.millitm: tm.millitm - whens.millitm);
		if (fp_debugInfo != NULL)
			fprintf (fp_debugInfo, "Elapsed system time: %u.%03u seconds\n", ul_rslt,
				(tm.millitm < whens.millitm)? 1000 + tm.millitm - whens.millitm: tm.millitm - whens.millitm);
#else
		fprintf (ferr, "Elapsed system time: %ld seconds.\n", l_time);
		if (fp_debugInfo != NULL) fprintf (fp_debugInfo, "Elapsed system time: %ld seconds.\n", l_time);
#endif
	}
	if (strchr (flags, 'a') || strchr (flags, 'k')) {
		fprintf(ferr, "*** Buried:\n");
		ri_putmb (stock, stock -> pair.b, ferr);
		if (fp_debugInfo != NULL) {
			fprintf(fp_debugInfo, "*** Buried:\n");
			ri_putmb (stock, stock -> pair.b, fp_debugInfo);
		}
	}
	if (strchr (flags, 'a') || strchr (flags, 's')) {
		unsigned long space = memory_used * sizeof (LINK) * NUM_OF_LINKS_AT_A_GULP;
		fprintf(ferr, "Memory allocated = %9ld Bytes\n", space);
		if (fp_debugInfo != NULL) fprintf(fp_debugInfo, "Memory allocated = %9ld Bytes\n", space);
	}
}

int ri_error (int code) {
	FILE *ferr = ref_err_file();

#ifdef FOR_OS_WINDOWSNT
	struct _timeb t;
#else
	long t;
#endif

#ifdef FOR_OS_WINDOWSNT
	_ftime (& t);
#else
	t = time (NULL);
#endif

	ri_print_error_code(ferr, code);
	if (fp_debugInfo != NULL) ri_print_error_code(fp_debugInfo, code);
	ri_information (t, code);

# ifdef VMS
	/* On VAX/VMS the return code must be odd, otherwise the
		operating system complains. July 1 1986. DT. */
	code = 2*code + 1;
	if (code == 1) exit (code);
#endif

# ifdef PCAT
	if (flags[0] != '\0')
# endif
		putc ('\n', ferr);
	exit (code);
	return 0;
}

/* exploding or imploding functions. */

/* Get the address of the compound symbol *s, assuming that *s consists of only upper-case letters, digits or underscore. */
char *ri_cs_impl (char * s) {
	int i;
	char *q;

	i = ri_hash (s);
	q = ri_cs_exist (i, s);
	if (q != NULL) return q;

	/* if we got here then the compound symbol was not found. */
	/* create it and insert it into the hast table. */

	q = malloc (strlen (s) + 1);
	/* 191 line. Not check result of malloc. Shura. 29.01.98 */
	if (NULL == q) {
	  fprintf (stderr, "No memory for string\n");
	  exit (1);
	}
	strcpy (q, s);
	ri_cs_ins (i, q);
	return q;
}


/* check to see if the compound symol is already in the table. */
char *ri_cs_exist (int i, char * s) {
	struct cshte *ptr;
	char *q;

	for (ptr = csht[i]; ptr != NULL; ptr = ptr -> next) {
		q = ptr -> cs;
		if (strcmp (s, q) == 0) return q;
	}
	return NULL;
}

/* insert the compound symbol q with hash value hv into the table. */
int ri_cs_ins (int hv, char * q) {
	struct cshte *ptr;

	ptr = (struct cshte *) malloc (sizeof (struct cshte));
	/* 221 line. Not check result of malloc. Shura. 29.01.98 */
	if (NULL == ptr) {
	  fprintf (stderr, "No memory for symbols\n");
	  exit (1);
	}
	ptr -> cs = q;
	ptr -> next = csht[hv];
	csht[hv] = ptr;
	return 0;
}

int ri_hash (char * s) {
	int base = 38, h = 0, i;

	while (*s != '\0') {
		switch (*s) {
		case '-':
		case '_':
			i = 36;
			break;

		case '$':
			i = 37;
			break;

		default:
			if (isupper (*s)) i = *s - 'A';
			else if (islower (*s)) i = *s - 'a';
			else if (isdigit (*s)) i = *s - '0' + 26;
			else return h;
		}
		s++;
		h = (h*base + i) % HASH_SIZE;
	}
	return h;
}

# define OPTION_MESS "The option given for definition of the size of the "
# define INCREASE_MESS " (maybe, by default). Try to increase the limit.\n\n"

int ri_print_error_code(FILE * fp, int err_code) {
        unsigned long m;

	if (err_code == 0) return 0;
	fprintf(fp, "Refal system Error: ");
	switch(err_code) {
	case 0:
		fprintf(fp, "Refal normal stop\n");
		break;
	case 1:
		fprintf(fp, "Memory allocation error\n");
		break;
	case 2:
		fprintf(fp, "File error.\n");
		break;
	case 3:
		fprintf(fp, "Error: Freeze occured and no Freezer found.\n");
		break;
	case 4:
		fprintf(fp, "Illegal instruction\n");
		break;
	case 5:
		fprintf(fp, "Error in FREEZER.\n");
		break;
	case 6:
		fprintf(fp, "Error in Freezer metacode.\n");
		break;
	case 7:
		fprintf(fp, "Error in function MU\n");
		break;
	case 8:
		fprintf(fp, "Format error in built-in function.\n");
		break;
	case 9:
		fprintf(fp, "No more free memory\n");
		break;
	case 10:
		fprintf(fp, "Unknown link type.\n");
		break;
	case 11:
                m = sizeof (LINK) * NUM_OF_LINKS_AT_A_GULP;
		fprintf(fp, "Memory limit reached\n");
		fprintf(fp, "Memory limit: %9d*%d\n",memory_limit,m);
		fprintf(fp, "Memory used: %9d*%d\n",memory_used,m);
		fprintf(fp, OPTION_MESS);
		fprintf(fp, "memory used for view field is \'-l%d\'",(m*memory_limit)/(1024*1024));
		fprintf(fp, INCREASE_MESS);
		break;
	case 12:
		fprintf(fp, "Recognition impossible.\n");
		break;
	case 13:
		fprintf(fp, "Stack of the variable tables overflow. ");
		fprintf(fp, OPTION_MESS);
		fprintf(fp, "stack is \'-V%d\'", size_table_element/DEFAULT_TABLE_SIZE);
		fprintf(fp, INCREASE_MESS);
		break;
	case 14:
		fprintf(fp, "Stack of the function calls from the left sides of sentences overflow. ");
		fprintf(fp, OPTION_MESS);
		fprintf(fp, "stack is \'-C%d\'", size_local_stack/DEFAULT_STACK_SIZE);
		fprintf(fp, INCREASE_MESS);
		break;
	default:
		fprintf(fp, "Error code %d\n", err_code);
		break;
	}
	return 0;
}

int ri_memory()
	{
	unsigned long limit = 0;
	char * fl = flags;
	
	while( ((* fl) != '\0') && !( ((* fl) == '-') && ((* ++fl) == 'l') )
	     ) fl++;
	
	/* -lnnn -- memory_limit */
	if ((* fl) == 'l') {
		for (fl ++; *fl != '\0' && isdigit(* fl); fl ++) {
			limit = (limit + ((* fl) - '0') ) * 10; 
		}
		limit /= 10;
	} else {
		limit = RI_MEMORY_LIMIT;
	}
	memory_limit = (limit * 1024 * 1024)/(sizeof (LINK) * NUM_OF_LINKS_AT_A_GULP) ;
	return 0;
}

int ri_common_stack (void) {
	if (size_local_stack == 0) {
		size_local_stack = DEFAULT_STACK_SIZE;
	}
	if (size_table_element == 0) {
		size_table_element = DEFAULT_TABLE_SIZE;
	}
	if (NULL == (te = (LINK **) malloc (sizeof (LINK *) * size_table_element))) {
		fprintf (stderr, "Cannot allcate memory for a variable stack\n");
		return -1;
	}
	if (NULL == (st = (STACK_STRUCTURE *) malloc (sizeof (STACK_STRUCTURE) * size_local_stack))) {
		fprintf (stderr, "Cannot allocate memory for a calling stack\n");
		return -1;
	}
	return 0;
}

int ri_init_stop() {
	char * fl = flags;
	
	while( ((* fl) != '\0') && !( ((* fl) == '-') && ((* ++fl) == 'e') ) ) fl++;
	
        /* interactive stop on an internal error: YES/NO */
        stop_on_error = ((* fl) == 'e') ? NO : YES;
        return 0;
}

int ri_help(char reftr) {
 if ( reftr ) printf("Using: reftr [options] MODULE1+MODULE2+... [arguments] \n");
 else         printf("Using: refgo [options] MODULE1+MODULE2+... [arguments] \n");
 
 printf("\nOptions recognized by the Refal system:");
 printf("\n -n : upon normal stop print the number of steps.");
 printf("\n -v :   \"   \"   \"   \"   \"   \"    view field.");
 printf("\n -k :   \"   \"   \"   \"   \"   \"    content of the stock.");
 printf("\n -t :   \"   \"   \"   \"   \"   \"    elapsed time.");
 printf("\n -s :   \"   \"   \"   \"   \"   \"    maximum allocated storage.");
 printf("\n -a : all of the above.");
 printf("\n -e : upon internal error stop do not use interactive mode -");
 printf("\n      \"Do you want to print view-field? [Yn]\", use [n] as default.");
 printf("\n -i : ignore unresolved external references.");

 printf("\n -C[ ]nn or --\"call_stack=nn\" :");
 printf("\n  parameter for the size of the stack of function calls from the left sides");
 printf("\n  of Refal sentences, where \'nn\' is a number ranging from 1 to 10 (recommended).");

 printf("\n -c[ ]nnn or --\"code_limit=nnn\" :");
 printf("\n  limit of the rsl-code size, where \'nnn\' is a number of Kbs.");
 printf("\n  By default it is 64Kb.");

 printf("\n -lnnn : , where \'nnn\' is a number - memory_limit for the view field,");
 printf("\n           memory = 512*sizeof(LINK)*memory_limit");

 printf("\n -V[ ]nn or --\"var_stack=nn\" :");
 printf("\n  parameter for the size of the stack of variable tables,");
 printf("\n  where \'nn\' is a number ranging from 1 to 10 (recommended).");

 if ( reftr ) printf("\n -f[ ]file_name: copy trace messages to the file name.");

 printf("\n -h : print this help message.\n\n");
}

int ri_options(int argc, char * argv [], char reftr) {
	int i;

	gargc = 0;
	flags [0] = '\0';
	garg = (char **) malloc ((argc + 1) * sizeof (char *));
	if (garg == NULL) {
		fprintf(stderr, "Error: Can\'t allocate memory.\n");
		exit (3);
	}
	for (i = 0; i < argc; i ++) garg [i] = NULL;

	for (i = 1; i < argc; i++) {
		if (*argv[i] == '-') {/* flags */
			if (argv [i][1] == '-') {
				if (strncmp (argv [i] + 2, "call_stack=", 11) == 0) {
					size_local_stack = DEFAULT_STACK_SIZE * atoi (argv [i] + 13);
				} else if (strncmp (argv [i] + 2, "var_stack=", 10) == 0) {
					size_table_element = DEFAULT_TABLE_SIZE * atoi (argv [i] + 12);
				} else if (strncmp (argv [i] + 2, "code_limit=", 11) == 0) {
					ul_limitSizeCode = atoi (argv [i] + 13) * 1024;
				} else {
					strcat (flags,argv[i]);
				}
			} else if (argv [i][1] == 'V') {
				char * cp = NULL;

				if (argv [i][2] == 0) {
					if ((cp = argv [i + 1]) == NULL) {
						fprintf (stderr, "WARNING: Uncorrect format for \'-V\'. It must be \'-V nn\' or \'-Vnn\'.");
						fprintf (stderr, " Ignoring the potion.\n");
						continue;
					}
					i ++;
				} else {
					cp = argv [i] + 2;
				}
				if (! isdigit (* cp)) {
					fprintf (stderr, "WARNING: Uncorrect format for \'-V\'. It must be \'-V nn\' or \'-Vnn\'.");
					fprintf (stderr, " Ignoring the potion.\n");
					continue;
				}
				size_table_element = DEFAULT_TABLE_SIZE * atoi (cp);
			} else if (argv [i][1] == 'C') {
				char * cp = NULL;

				if (argv [i][2] == 0) {
					if ((cp = argv [i + 1]) == NULL) {
						fprintf (stderr, "WARNING: Uncorrect format for \'-C\'. It must be \'-C nn\' or \'-Cnn\'.");
						fprintf (stderr, " Ignoring the potion.\n");
						continue;
					}
					i ++;
				} else {
					cp = argv [i] + 2;
				}
				if (! isdigit (* cp)) {
					fprintf (stderr, "WARNING: Uncorrect format for \'-C\'. It must be \'-C nn\' or \'-Cnn\'.");
					fprintf (stderr, " Ignoring the potion.\n");
					continue;
				}
				size_local_stack = DEFAULT_STACK_SIZE * atoi (cp);
			} else if (argv [i][1] == 'c') {
				char * cp = NULL;

				if (argv [i] [2] == 0) {
					if ((cp = argv [i + 1]) == NULL) {
						fprintf (stderr, "WARNING: Uncorrect format for \'-c\'.\n\t");
						fprintf (stderr, " It must be \'-cnnn\' or \'-c nnn\' or \'--code_limit=nnn\'\n\t Where nnn is number of Kb.");
						fprintf (stderr, " Ignoring the potion.\n");
						continue;
					}
					i ++;
				} else {
					cp = argv [i] + 2;
				}
				if (! isdigit (*cp)) {
					fprintf (stderr, "WARNING: Uncorrect format for \'-c\'.\n\t");
					fprintf (stderr, " It must be \'-cnnn\' or \'-c nnn\' or \'--code_limit=nnn\'\n\t Where nnn is number of Kb.");
					fprintf (stderr, " Ignoring the potion.\n");
					continue;
				}
				ul_limitSizeCode = atoi (cp) * 1024;
			} else if ( (argv [i][1] == 'f') && reftr ) {
				char * cp;

				if (argv [i][2] == 0) {
					if ((cp = argv [i + 1]) == NULL) {
						fprintf (stderr, "WARNING: Uncorrect foramt for \'-f\'.\n\t");
						fprintf (stderr, " It must be \'-f<file_name>\' or \'-f <file_name>\'");
						fprintf (stderr, " Ignoring the potion.\n");
						continue;
					}
					i ++;
				} else {
					cp = argv [i] + 2;
				}
				if (NULL == (fp_debugInfo = fopen (cp, "w"))) {
					fprintf (stderr, "Cannot open file \'%s\' for copying of trace informations\n", cp);
				}
			} else if ( argv [i][1] == 'h' ) { 
                                ri_help( reftr ); 
			} else {
				strcat (flags,argv[i]);
			}
		} else garg [gargc++] = argv [i];
	}
}

