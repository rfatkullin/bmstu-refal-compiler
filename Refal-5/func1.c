 
#ifdef FOR_OS_WINDOWSNT
#	include <sys/timeb.h>
#endif

# include "decl.h"
# include "macros.h"
# include "fileio.h"
# include "ifunc.h"

#include "version.h"

#ifdef FOR_OS_WINDOWSNT
/* For _stat () */
#	include <sys/types.h>
#	include <sys/stat.h>
/* For _getcwd () */
#	include <direct.h>
/* For getpid () */
#	include <process.h>
#else
/* For stat () */
#	include <sys/stat.h>

/* For stat (), getcwd () */
#	include <unistd.h>

/* For 	WIFEXITED, WEXITSTATUS */
#	include <sys/types.h>
#	include <sys/wait.h>

#       define  MAX_PATH_LEN 1024
#       ifdef FOR_OS_DOS
#               define _MAX_PATH 1024
#       endif
#endif

#	define MAX_SYSTEM_CALL 1024

#include <time.h>
#include <errno.h>

# define MDEBUG 0
# define MDEBUG1 0

  /* For translator pgraph */
extern unsigned long ul_local_calls;
extern int flag_local_calculation;


int rf_prout ()
	/* RF_PROUT prints its argument onto standart output file
		and returns empty. Format :
		<PROUT  e.arg>	*/

	{
	check_freeze;

	if (b1->foll != b2) ri_put (stdout);
	putc ('\n', stdout);
	rdy (0);
	out (2);
	est;
	return 0;
	}


int rf_put ()
	/* RF_PUT prints the second argument onto file REFALxx.DAT, 
		where xx is the first argument (must be a MACRODIGIT)
		and returns the second argument. Format :
		<PUT s.first e.second>	*/

	{
	int i;
	char s [11];
	FILE *f;

	check_freeze;

	b1 = b1->foll;
	if (!LINK_NUMBER (b1)) ri_error(8);
	i = (b1->pair.n) % FILE_LIMIT;
	if (i == 0) f = stderr;
	else
		{
		if (file_table [i] == NULL)
			{
			sprintf (s, "REFAL%d.DAT", i);
			ri_open (i, "w", s);
			};
		f = file_table [i];
		};
	tbel (3) = ((b1->foll == b2) ? NULL : b1->foll);
	tbel (4) = b2->prec;
	if (tbel (3) != NULL) ri_put (f);
	putc ('\n', f);
	rdy (0);
	tple (4);
	out (2);
	est;
	return 0;
	}


/* BR: transplant the argument to the stock.	*/
int rf_br (void) {
	LINK * lnkp_w;

	/* Check arguments */
	for (lnkp_w = b1 -> foll; lnkp_w != b2; lnkp_w = lnkp_w -> foll) {
		if (LINK_CHAR (lnkp_w) && lnkp_w -> pair.c == '=') {
			break;
		}
	}
	if (lnkp_w == b2) {
		ri_error (8);
		exit (1);
	}

	cl;
	b = stock;
	rend = b->foll;
	bl;
	tple (4);
	br;
	weld (b, rend);
	rdy (0);
	out (2);
	est;
	return 0;
}


int rf_putout (char outputNL)
	/* RF_PUTOUT prints the second argument onto file REFALxx.DAT, 
		where xx is the first argument (must be a MACRODIGIT)
		and returns the empty. Format :
		<PUTOUT s.first e.second>	*/

	{
	int i;
	char s [11];
	FILE *f;

	check_freeze;

	b1 = b1->foll;
	if (!LINK_NUMBER (b1)) ri_error(8);
	i = (b1->pair.n) % FILE_LIMIT;
	if (i == 0) f = stderr;
	else
		{
		if (file_table [i] == NULL)
			{
			sprintf (s, "REFAL%d.DAT", i);
			ri_open (i, "w", s);
 		};
		f = file_table [i];
		};
	tbel (3) = ((b1->foll == b2) ? NULL : b1->foll);
	tbel (4) = b2->prec;
	if (tbel (3) != NULL) ri_put (f);
	if (outputNL == YES) putc ('\n', f);
	rdy (0);
	out (2);
	est;
	return 0;
	}



int rf_card (void) { 
	check_freeze;
	return ri_get (stdin);
}

/* RF_PRINT prints its argument onto standart output file and returns the argument. Format : <PRINT  e.arg> */
int rf_print (void) {
	check_freeze;

	tbel (3) = ((b1->foll == b2) ? NULL : b1->foll);
	tbel (4) = b2->prec;
	if (tbel (3) != NULL) ri_put (stdout);
	putc ('\n', stdout);
	rdy (0);
	tple (4);
	out (2);
	est;
	return 0;
}

/*
static int
is_function (char * cp_id) {
	char ** cpp_w, * cp_func;

	for (cpp_w = local_functions; * cpp_w != NULL; cpp_w ++) {
		for (cpp_w ++; * cpp_w != NULL; cpp_w ++) {
			for (cp_func = (* cpp_w) - 2; * cp_func != 0; cp_func --);
			cp_func ++;
			if (0 == strcmp (cp_func, cp_id)) {
				return 1;
			}
		}
	}
	return 0;
}
*/

int rf_type (void) {
	char c [3] = {'\0', '0', '\0'}, x, i;

	cl;
	if (tbel (3) == NULL) c [0] = '*';
	else {
		x = tbel (3)->ptype;
		switch (x) {
		case LINK_TYPE_LSTRUCTB:
			c [0] = 'B';
			break;
		case LINK_TYPE_CHAR:
			i = tbel (3)->pair.c;
			if (isupper(i)) {
				c [0] = 'L'; c [1] = 'u';
			} else if (islower(i)) {
				c [0] = 'L'; c [1] = 'l';
			} else if (isdigit(i)) {
				c [0] = 'D';
			} else if (' ' <= i && i <= 127) {
				c [0] = 'P';
				c [1] = (isupper (i))? 'u': 'l';
			} else {
				c [0] = 'O'; 
				c [1] = (isupper (i))? 'u': 'l';
			}
			break;
		case LINK_TYPE_COMPSYM:
			c [0] = 'W';
/*  Nemytykh		if (isalpha (tbel (3) -> pair.f [0])) { */
			if (isIdentifier  (tbel (3) -> pair.f )) { 
				c [1] = 'i';
			} else {
				c [1] = 'q';
			}
			break;
		case LINK_TYPE_NUMBER:
			c [0] = 'N';
			break;
		case LINK_TYPE_SVAR:
			c [0] = 'S';
			break;
		case LINK_TYPE_EVAR:
			c [0] = 'E';
			break;
		case LINK_TYPE_TVAR:
			c [0] = 'T';
		default:
			ri_error(10);
			return 1;
		}
	}
	rdy (0);
	{
		int i_w;

		for (i_w = 0; c [i_w] != '\0'; i_w ++) ns (c [i_w]);
	}
	tple (4);
	out (2);
	est;
	return 0;
}

int rf_dg ()
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
		tple (11);
		out (2);
		b = tbel (5)->prec;
		out (6);
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


int rf_get ()
	{
	/*register*/ int i;
	char s [11];

	check_freeze;

	b1 = b1->foll;
	if (!LINK_NUMBER (b1)) ri_error (8);
	i = (b1->pair.n) % FILE_LIMIT;
	if (i == 0) ri_get (stdin);
	else
		{
		if (file_table [i] == NULL)
			{
			sprintf (s, "REFAL%d.DAT", i);
			ri_open (i, "r", s);
			};
		ri_get (file_table [i]);
		};
	return 0;
	}


int rf_arg ()
	{
	/*register*/ int i = 0;
	char *c;

	check_freeze;

	cl;
	if (tbel (3) == NULL) i = 1;
	else if (LINK_NUMBER (tbel (3))) i = tbel (3)->pair.n;
	else ri_error (8);
	rdy (0);
	if (i < gargc) for (c=garg [i]; *c != '\0'; c++) ns (*c);
	out (2);
	est;
	return 0;
	}

char *mu_find (char *str, char **pi) {
	char *pc;

	/*pi ++;*/ /* skip the title. */
	for (pi ++; *pi != NULL; pi ++) {
		pc = *pi;

		if (pc [-1] != '\0') {
			fprintf (stderr, "WARNING: Uncorrect format names(MU_FIND)\n");
		}
		for (pc -= 2; * pc != '\0'; pc --);
		pc ++;

		/*if (strncmp (str, pc - MAXWS, MAXWS) == 0) return pc;*/
		if (strcmp (str, pc) == 0) {
			for (; * pc != '\0'; pc ++);
			pc ++;
			return pc;
		}
		/*pi ++;*/
	}
	return NULL;
}
	
/* RF_MU accepts 2 formats:
 *	<MU (e.STRING) e.ARG>
 *	<MU s.SYMBOL e.ARG>
 * The result is <F e.ARG>, where F is the function whose name
 * is the same as the STRING in the first case or SYMBOL in the second.
 */
/* pointer to the beginning of the list of functions visible from this module. */
int rf_mu (char ** muaddr, int freeze_flag) {
	char str [MAXWS+1], c, *pc;
	int i;

	if (freeze_flag) check_freeze;
	i = 0;
	b1 = b1->foll;
	if (LINK_VAR (b1)) goto call_freeze;
	/* first format */
	if (LINK_LSTRUCTB (b1)) {
		/* copy the string into array str []; */
		for (b1 = b1 -> foll; LINK_CHAR (b1); b1 = b1 -> foll) {
			c = b1->pair.c;

			if (isalnum (c) || c == '_' || c == '$' || c == '-') str [i++] = c;
			else ri_imp ();

			if (i >= MAXWS) ri_error (7);
		}
		if (LINK_VAR (b1)) goto call_freeze;
		str [i] = '\0';
		/* check that the string is not empty. */
		if (i == 0) ri_imp ();
		/* check that b1 now points to a right parenthesis. */
		if (!LINK_RSTRUCTB (b1)) ri_imp ();
	} else if (LINK_COMPSYM (b1)) { /* second format */
		strcpy (str, b1->pair.f);
	} else if (LINK_SYMBOL (b1) && mu_special (str, b1->pair.c) == 0) {/* special symbol */
		;
	} else { /* error */ 
		ri_imp ();
	}
	/* Now b1 points just before the argument. e.ARG. */
	cl;
	/* Get the address of the function stored in str. */
	pc = mu_find (str, muaddr);
	if (pc == NULL) pc = mu_find (str, entry_functions);
	if (pc != NULL) {
		/* the function is found. */
		rdy (0)
		bl;
		tple (4);
		br;
		act1 (pc);
		out (2);
		est;
		return 0;
	} else {/* we got here if the function with the same name was not found. */
		ri_error (7);
		return 1;
	}

call_freeze:

	ri_frz (2);
	return 0;
}

	/* check if c_value is a special acceptable value for arithmetic
		function. Return 0 if success, 1 = failure
		Important: if FAIL, do NOT overwrite s_value */

int mu_special (char *s_value, char c_value) {
	switch (c_value) {
	case '+':
		strcpy (s_value, "ADD");
		break;

	case '-':
		strcpy (s_value, "SUB");
		break;

	case '*':
		strcpy (s_value, "MUL");
		break;

	case '/':
		strcpy (s_value, "DIV");
		break;

	case '%':
		strcpy (s_value, "MOD");
		break;

	case '?':
		strcpy(s_value, "RESIDUE");
		break;

	default:
		return 1;	/* failure */
	}
	return 0;	/* success */
}


int
rf_getenv (void) {
	char ca_buf [120];
	char * cp;

	check_freeze;

	for (cp = ca_buf, b1 = b1->foll; b1 != b2; b1 = b1 -> foll, cp ++) {
		if (! LINK_CHAR (b1)) ri_error (8);
		* cp = b1 -> pair.c;
	}
	* cp = 0;
	rdy (0);
	if (NULL != (cp = getenv (ca_buf))) {
		for (; * cp != 0; cp ++) ns (* cp);
	}
	out (2);
	est;
	return 0;
}

/* Get current process id. A.A. Vladimirov 10.10.2004 */
int rf_getpid (void) {

	check_freeze;

	b1 = b1->foll;
	if (b1 != b2) ri_error (8);
	rdy (0);

	nns ((unsigned long) getpid());
	out (2);
	est;
	return 0;
}

/* Get parent process id. A.A. Vladimirov 10.10.2004.
   This function is undefined under Windows operating sysytem.
   In fact, it returns current process id. 
*/
int rf_getppid (void) {
        unsigned long i;

	check_freeze;

	b1 = b1->foll;
	if (b1 != b2) ri_error (8);
	rdy (0);
# ifdef FOR_OS_WINDOWSNT
	i = (unsigned long) getpid();
# else
	i = (unsigned long) getppid();
# endif
        nns(i);
	out (2);
	est;
	return 0;
}

int
rf_system (void) {
	char ca_buf [MAX_SYSTEM_CALL+1];
	int i;

	check_freeze;

	for (b1 = b1 -> foll, i = 0; b1 != b2; b1 = b1 -> foll, i ++) {
		if (i == MAX_SYSTEM_CALL) ri_error (1);	/* A. Vladimirov , 10.10.2004 */
		if (! LINK_CHAR (b1)) ri_error (8);
		ca_buf [i] = b1 -> pair.c;
	}
	ca_buf [i] = 0;
	i = system (ca_buf);

#ifdef FOR_OS_LINUX
	/* madmax */
	if (WIFEXITED (i) != 0) {
		i = WEXITSTATUS (i);
	} else {
		i = 1;
	}
#endif

	rdy (0);

	if (i < 0) {
		ns ('-');
		i = abs(i);
	}
	nns (i);
	out (2);
	est;

	return 0;
}

extern FILE* ref_err_file();
#ifdef FOR_OS_WINDOWSNT
	extern void ri_informatiom (struct _timeb, int);
#else
	extern void ri_information (long, int);
#endif
int
rf_exit (void) {
	int i;
#ifdef FOR_OS_WINDOWSNT
	struct _timeb l_t;
	
	_ftime (& l_t);
#else
	long l_t = time (NULL);
#endif

	check_freeze;

	for (i = 0; i < FILE_LIMIT; i ++) {
		if (file_table [i] != NULL) {
			fclose (file_table [i]);
		}
	}

	b1 = b1 -> foll;
	if (LINK_CHAR (b1)) {
		if (b1 -> pair.c == '-') {
			i = -1;
		} else if (b1 -> pair.c == '+') {
			i = 1;
		}
		b1 = b1 -> foll;
	} else {
		i = 1;
	}
	if (! LINK_NUMBER (b1)) {
			ri_error (8);
	}
	i *= b1 -> pair.n;

	rdy (0);
	/* Print all informations which are wanted to see to user */
	ri_information (l_t, 0);
	exit (i);
	out (2);
	est;
}

int rf_close (void) {
	int i;

	check_freeze;

	b1 = b1 -> foll;
	if (! LINK_NUMBER (b1)) {
		ri_error (8);
	}
	i = b1 -> pair.n % FILE_LIMIT;
	if (NULL != file_table [i]) {
		fclose (file_table [i]);
		file_table [i] = NULL;
	}
	b1 = b2;
	rdy (0);
	out (2);
	est;
	return 0;
}

int rf_existfile (void) {
	char ca_file [FILENAME_MAX];
	char * cp;
	int i;

	check_freeze;

	b1 =  b1 -> foll;

	for (i = 0; i < FILENAME_MAX - 1 && b1 != b2; i ++, b1 = b1 -> foll) {
		if (! LINK_CHAR (b1)) {
			ri_error (8);
		}
		ca_file [i] = b1 -> pair.c;
	}
	ca_file [i] = 0;

	{
#ifdef FOR_OS_WINDOWSNT
		struct _stat st_buf;

		if (0 == _stat (ca_file, & st_buf)) {
			/* Exist */
			if (st_buf.st_mode & _S_IFREG) {
				cp = ri_cs_impl ("True");
			} else {
				cp = ri_cs_impl ("False");
			}
		} else {
			cp = ri_cs_impl ("False");
		}
	#else
		struct stat st_buf;

		if (0 == stat (ca_file, & st_buf)) {
			/* Exist */
			if (st_buf.st_mode & S_IFREG) {
				cp = ri_cs_impl ("True");
			} else {
				cp = ri_cs_impl ("False");
			}
		} else {
			cp = ri_cs_impl ("False");
		}
#endif

	}
	rdy (0);
	ncs (cp);
	out (2);
	est;
	return 0;
}

int rf_getcurrentdirectory (void) {
	int i;

#ifdef FOR_OS_WINDOWSNT
	char ca_curdir [_MAX_PATH];

	check_freeze;

	if (NULL == _getcwd (ca_curdir, _MAX_PATH)) {
		fprintf (stderr, "WARNING: Cannot get a current directory\n");
		ca_curdir [0] = 0;
	}
#else
	char ca_curdir [MAX_PATH_LEN];

	check_freeze;

	if (NULL == getcwd (ca_curdir, MAX_PATH_LEN)) {
		fprintf (stderr, "WARNING: Cannot get a current directory\n");
		ca_curdir [0] = 0;
	}
#endif

	rdy (0);
	for (i = 0; ca_curdir [i] != 0; i ++) {
		ns (ca_curdir [i]);
	}
	out (2);
	est;
	return 0;
}


int rf_removefile (void) {
	char ca_buf [FILENAME_MAX], * cp, * cp_err;
	int i;

	check_freeze;

	for (i = 0, b1 = b1 -> foll; i < FILENAME_MAX - 1 && b1 != b2; b1 = b1 -> foll, i++) {
		if (! LINK_CHAR (b1)) {
			ri_error (8);
		}
		ca_buf [i] = b1 -> pair.c;
	}
	ca_buf [i] = 0;
	if (b1 != b2) {
		fprintf (stderr, "WARNINIG: File name is too long.\n\tNow File name is \'%s\'\n", ca_buf);
		b1 = b2;
	}

	if (-1 == remove (ca_buf)) {
		/*fprintf (stderr, "Cannot delete file \'%s\'\n", ca_buf);*/
		cp = ri_cs_impl ("False");
		cp_err = strerror (errno);
	} else {
		cp = ri_cs_impl ("True");
		cp_err = NULL;
	}
	rdy (0);
	ncs (cp);
	bl;
	if (cp_err != NULL) {
		for (i = 0; cp_err [i] != 0; i ++) {
			ns (cp_err [i]);
		}
	}
	br;
	out (2);
	est;
	return 0;
}

int
rf_implodeExt (void) {
	char ca_str [MAXWS + 1], * cp_w;
	int i;

	for (i = 0, b1 = b1 -> foll; i < MAXWS && b1 != b2; i ++, b1 = b1 -> foll) {
		if (LINK_VAR (b1)) {
			ri_frz (2);
			return 0;
		} else if (! LINK_CHAR (b1)) break;
		ca_str [i] = b1 -> pair.c;
	}
	if (i == MAXWS && b1 != b2) {
		fprintf (stderr, "WARNING: the length of the string is too long.\n");
	}
	ca_str [i] = 0;
	cp_w = ri_cs_impl(ca_str);
	rdy (0);
	ncs (cp_w);
	out (2);
	est;
	return 0;
}

int
rf_explodeExt (void) {
	char * cp_w;

	if ((b1 = b1 -> foll) == b2) ri_imp ();
	else if (LINK_VAR (b1)) {
		ri_frz (2);
		return 0;
	} else if (! LINK_COMPSYM (b1) || b1->foll != b2) ri_imp ();
	rdy (0);
	for (cp_w = b1 -> pair.f; * cp_w != 0; cp_w ++) {
		ns (* cp_w);
	}
	out (2);
	est;
	return 0;
}

int
rf_timeExt (void) {
#ifdef FOR_OS_WINDOWSNT
	char /* * cp_time,*/ ca_refTime [50];
	int i;
	struct _timeb tm;
	unsigned long ul_rslt;
	
	_ftime (& tm);
	ul_rslt = (tm.millitm < tm_localtime.millitm)? tm.time - tm_localtime.time - 1: tm.time - tm_localtime.time;
	sprintf (ca_refTime, "%u.%03u", ul_rslt, 
		(tm.millitm < tm_localtime.millitm)? 1000 + tm.millitm - tm_localtime.millitm: tm.millitm - tm_localtime.millitm);
	
	/*cp_time = ctime (& (tm.time));*/
	/*sprintf (ca_refTime, "%.19s.%hu %s", cp_time, tm.millitm, cp_time + 20);*/
#endif

#ifdef FOR_OS_WINDOWSNT
	
	
	if ((b1 = b1 -> foll) != b2) {
		if (!LINK_NUMBER (b1)) {
			ri_error (8);
			return 1;
		}
		if (b1 -> pair.n == 0) {
			tm_localtime = tm;
		}
	}
#endif

	rdy (0);

#ifdef FOR_OS_WINDOWSNT
	for (i = 0; ca_refTime [i] != '\0'; i ++) {
		ns (ca_refTime [i]);
	}
#endif

	out (2);
	est;
	return 0;
}

int rf_XMLParse (void) {
  char ca_file [FILENAME_MAX];
  int i;
  FILE * fp_in, * fp_out;

  for (i = 0, b1 = b1 -> foll; b1 != b2; b1 = b1 -> foll, i ++) {
    if (! LINK_CHAR (b1)) {
      ri_error (8);
      return 1;
    }
    ca_file [i] = b1 -> pair.c;
  }
  ca_file [i] = 0;
  if (NULL == (fp_in = fopen (ca_file, "r"))) {
    fprintf (stderr, "Cannot open file \'%s\'.\n", ca_file);
    exit (1);
  }
  if (NULL == (fp_out = fopen ("REFAL5XML", "w"))) {
    fprintf (stderr, "Cannot open file \'REFAL5XML\'\n");
    exit (1);
  } 
  if (0 != ri_xml2ref (fp_in, fp_out)) {
    weld (b, rend);
    ri_error (8);
    fclose (fp_in);
    return 1;
  }
  fclose (fp_out);
  if (NULL == (fp_out = fopen ("REFAL5XML", "r"))) {
    fprintf (stderr, "Cannot open file \'REFAL5XML\'.\n");
    exit (1);
  }
  fclose (fp_in);
  return getXML (fp_out);
}

#define ESC 27


