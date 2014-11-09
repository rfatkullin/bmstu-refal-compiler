#ifdef FOR_OS_DOS
#	ifndef FILENAME_MAX
#		define FILENAME_MAX 256
#	endif
#	ifndef FOPEN_MAX
#		define FOPEN_MAX 32
#	endif
#endif

	/*** Argument checking. ***/
# ifndef LINT_ARGS
# define LINT_ARGS 1
# endif

	/*** Include Files. ***/
# ifndef FILES_INCLUDED
# include <stdio.h>

# include <string.h>
# include <ctype.h>
# ifndef IBM370
# include <malloc.h>
# endif
# include <stdlib.h>
# ifdef IBM370
# include "memory.h"
# else
# include "memory.h"
# endif
# define FILES_INCLUDED 1
# endif

	/*** Refal constants. ***/
# define PCAT 1

	/* size of a function name/compound symbol. */
# ifndef MAXWS
/*#	define MAXWS 32 */
#	define MAXWS 1024
# endif

#ifndef DEFAULT_STACK_SIZE
#	define DEFAULT_STACK_SIZE 512
#endif

#ifndef DEFAULT_TABLE_SIZE
#	define DEFAULT_TABLE_SIZE 1024
#endif

#ifndef DEFAULT_CODE_LIMIT
#	define DEFAULT_CODE_LIMIT 64
#endif

	/* External definitions. */
# ifndef EXTERNAL
# 	ifdef DEFINE_EXTERNALS
#		if DEFINE_EXTERNALS
# 			define EXTERNAL
#		else
#			define EXTERNAL extern
#		endif
# 	else
# 		define EXTERNAL extern
# 	endif
# endif


	/*** Refal structures. ***/
 typedef struct link
	{
	char ptype; /* type of the link */
	union
		{
		struct link *b;	/* bracket: ptr to the pair */
		char *f;	/* function or compound symbol: ptr to label. */
		char c;		/* symbol: actual value. */
		unsigned long n; /* macrodigit: number */
		unsigned short us_1, us_2;
		} pair;
	struct link *prec; /* ptr to preceding link */
	struct link *foll; /* ptr to following link */
	} LINK;


		/*** Refal Global variables. ***/

	EXTERNAL LINK
		*lfm,	/* list of free memory		*/
		*stock,	/* stock pointer 		*/
		*vf,	/* view field			*/
		*vfend,	/* end of view field		*/
		*b1,*b2,	/* left and right boundaries	*/
		*quap,	/* quasipoint			*/
		*nextp,	/* next active point		*/
		*precp,	/* preceding active point	*/
		**te /*[512]*/,	/* table of elements		*/
		*b,	/* boundary			*/
		*pa,	/* parenthesis address		*/
		*rend,	/* right end			*/
		*nextb;	/* next bracket			*/

  EXTERNAL long
		nst,	/* number of steps		*/
		curk,	/* current number of k-signs	*/
		stoff,	/* transition stack offset 	*/
		teoff,	/* table of elements offset	*/
		sp,	/* top of the transition stack	*/
# ifndef FOR_OS_WINDOWSNT
		whens,	/* time when started.	*/
# endif
		nel;	/* new element	*/

#ifdef FOR_OS_WINDOWSNT
	EXTERNAL struct _timeb
			whens,
			tm_localtime;
#endif

	EXTERNAL int memory_limit;
	EXTERNAL int size_local_stack;
	EXTERNAL int size_table_element;

  EXTERNAL LINK  **tel; /* pointer to an entry in the table of elements */
  EXTERNAL char
			*p,	        /* RASL instruction address pointer              */
			*actfun;        /* active function pointer	                 */

	typedef struct
		{
		LINK *b1;		/* pointer to left boundary */
		LINK *b2;		/* pointer to right boundary */
		long nel;		/* new element		*/
		char *ra;		/* return address	*/
		}  STACK_STRUCTURE;      /* transition stack  */

	EXTERNAL STACK_STRUCTURE * st/*[256]*/;

	/* Standard functions and structures concerned with
		imploding compound symbols.	*/

 EXTERNAL char *STOP_, *IMP_;

	/* The pointers to the lists of local and entry functions and 
		module names and their lenghts. Used for the tracer and
		function MU.  03/08/1987 */
 EXTERNAL char **local_functions, *module_table, **entry_functions;
 EXTERNAL long num_local_functions, num_modules, num_entry;

 /* Global variables that contain the command line arguments. */
# ifdef PCAT
	EXTERNAL char flags[20];	/* Flags Up to 20 flags. */
	EXTERNAL char **garg;
	EXTERNAL int gargc;
# endif

	struct  cshte 	/* compound symbols hashing table entry. */
		{
			char *cs;	/* address of compound symbol. */
			struct cshte *next;	/* pointer to next. */
		};
	/* Hash table is implemented with "open" hashing scheme. */
# define HASH_SIZE 101
	EXTERNAL struct cshte *csht[HASH_SIZE] ; /* initialized to NULL. */

# ifndef max
# define max(a,b) ((a > b) ? a : b)
# endif

# ifndef min
# define min(a,b) ((a < b) ? a : b)
# endif


