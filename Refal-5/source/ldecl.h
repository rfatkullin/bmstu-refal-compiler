
# define LOAD_INCLUDED

	/* define loader macros.   */

# define write_byte_to_mem(BB) {*ldptr++ = (char)(BB);}

# ifndef FOR_OS_SUNLINUX
# define write_long_to_mem(LL) {*((long *) ldptr) = (long)(LL); ldptr += sizeof (long);}
# else
/*** MAC, Sun (Linux) Versions ***/
# 	define write_long_to_mem(LL) {  \
		long TMP; TMP = (long)(LL); \
		memcpy (ldptr, (char *) &TMP, sizeof (long)); ldptr += sizeof (long); \
	}
# endif

# define read_byte_to_mem(X) fread ((char *) &X, sizeof (char), 1, fdtmpr)
# define read_long_to_mem(X) fread ((char *) &X, sizeof (long), 1, fdtmpr)

# ifndef MAXWS
/*# define MAXWS 32*/
#	define MAXWS 1024
# endif

	/* Define the structure for holding module information. */

	typedef struct
		{
		/*char name [MAXWS];*/
		char * name;
		char *addr;
		}  LOAD_TABLE;

	typedef struct
		{
		char *name;
		long ident;
		}  COMPSYM_TABLE;

	typedef struct module_descriptor
		{
		/*char module_name[MAXWS];*/
		char module_name[FILENAME_MAX];
		char *base;
		long size;
		long cs_size;
		long xt_size;
		long np_size;
		long lf_size;
		COMPSYM_TABLE *cs_table;
		LOAD_TABLE *ext_table;
		LOAD_TABLE *entry_pts;
		char  **local_table;
		struct module_descriptor *next;
		FILE *from;
		} MODULE;

