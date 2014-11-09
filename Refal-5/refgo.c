
# define DEFINE_EXTERNALS 1

#ifdef FOR_OS_WINDOWSNT
#	include <sys/timeb.h>
#	include <sys/types.h>
#endif

# include "rasl.h"
# include "decl.h"
# include "macros.h"
# include "fileio.h"
# include "ifunc.h"
# include <time.h>
# include <stdlib.h>

/* For translator pgraph */
unsigned long ul_local_calls;
int flag_local_calculation = 0;

FILE * fp_debugInfo = NULL;

/* initialized  to zero by default. */
static char def_imp [6 + 1 + sizeof (long)];/*[MAXWS + 1 + sizeof(long)];*/
static char def_stop [7 + 1 + sizeof (long)];/*[MAXWS + 1 + sizeof(long)];*/

int main (int argc, char * argv []) {      /* main entry to interpreter */

        ri_options(argc,argv, 0);

	/* define labels IMP and STOP. */
        ri_init_stop();
	def_imp [0] = 0;
	strcpy (def_imp + 1, "IMP$");
	IMP_ = def_imp + 1 + 4 + 1; /*MAXWS;*/
	*IMP_ = BUILT_IN;
	wrlong_to_mem (42L, (IMP_ + 1)); /* change this 42 if you change bif.c */
	def_stop [0] = 0;
	strcpy (def_stop + 1, "STOP$");
	STOP_ = def_stop + 1 + 5 + 1; /*MAXWS;*/
	*STOP_ = BUILT_IN;
	wrlong_to_mem (43L, (STOP_ + 1)); /* change this 43 if you change bif.c */

/*	memory_limit = 600;  / * about 4 meg. */
/*	memory_limit = 1200;  / * about 8 meg. */
	
	ri_memory();
	if (-1 == ri_common_stack ()) return -1;
	ri_init ();

#ifdef FOR_OS_WINDOWSNT
	_ftime (& whens);
	tm_localtime = whens;
        srand( (unsigned int)(tm_localtime.time) );
#else
	whens = time(NULL);
        srand((unsigned int)whens);
#endif	
	
	ri_inter ();
	exit (0);
	return 0;
   }
