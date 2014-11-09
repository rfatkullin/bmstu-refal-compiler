
/* File TRACE.C
	Main function for the Refal Tracer.
	March 8 1987.
*/

# define DEFINE_EXTERNALS 1

#ifdef FOR_OS_WINDOWSNT
#	include <sys/timeb.h>
#	include <sys/types.h>
#endif

# include "version.h"
# include "rasl.h"
# include "decl.h"
# include "cdecl.h"
# include "ddecl.h"
# include "macros.h"
# include "dmacro.h"
# include "fileio.h"
# include "ifunc.h"
# include "tfunc.h"
# include <time.h>
# include <stdlib.h>


 /* For translator pgraph */
unsigned long ul_local_calls;
int flag_local_calculation = 0;

/* For copying of trace informations */
FILE * fp_debugInfo = NULL;

/* initialized  to zero by default. */
/*
static char def_imp [MAXWS + 1 + sizeof(long)];
static char def_stop [MAXWS + 1 + sizeof(long)];
static char def_nobreak [MAXWS + 1];
*/
static char def_imp [6/*MAXWS*/ + 1 + sizeof(long)];
static char def_stop [7/*MAXWS*/ + 1 + sizeof(long)];
static char def_nobreak [10/*MAXWS*/ + 1];

/* main entry to Refal tracer. */
int main (int argc, char * argv []) {
	int i;

	/* print the version and copyright information. */
	printf (VERSION, "Tracer", _refal_build_version);
	printf ("Copyright: Refal Systems Inc.\n");
        
        ri_options(argc,argv,1);

	/* define labels IMP and STOP. */
	def_imp [0] = 0;
	strcpy (def_imp + 1, "IMP$");
	IMP_ = 1 + def_imp + 4 + 1;/*MAXWS;*/
	*IMP_ = BUILT_IN;
	wrlong_to_mem (42L, (IMP_ + 1)); /* change this 42 if you change bif.c */
	def_stop [0] = 0;
	strcpy (def_stop + 1, "STOP$");
	STOP_ = 1 + def_stop + 5 + 1;/*MAXWS;*/
	*STOP_ = BUILT_IN;
	wrlong_to_mem (43L, (STOP_ + 1)); /* change this 43 if you change bif.c */
	/* define NOBREAK_ */
	def_nobreak [0] = 0;
	strcpy (def_nobreak + 1, "NOBREAK$");
	NOBREAK_ = 1 + def_nobreak + 8 + 1;/*MAXWS;*/
	*NOBREAK_ = NO_BREAK;

/*	memory_limit = 600;  / * about 4 meg. */
/*	memory_limit = 1200;  / * about 8 meg. */
/*	memory_limit = 1500;  / * about 10 meg. */
/*	memory_limit = 1800;  / * about 12 meg. */
/*	memory_limit = 2400;  / * about 16 meg. */
	memory_limit = 2800;  /* about 20 meg. */

    ri_memory();
	if (-1 == ri_common_stack ()) exit (1);
	ri_init ();
	rd_init ();

#ifdef FOR_OS_WINDOWSNT
	_ftime (& whens);
	tm_localtime = whens;
        srand( (unsigned int)(tm_localtime.time) );
#else
	whens = time (NULL);
        srand((unsigned int)whens);
#endif

 	ri_inter ();

	return 0; /*exit (0);*/
}
