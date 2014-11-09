
# if !defined(IBM370) && !defined(FOR_OS_SUNLINUX)
	/* memory assignment macros. */
#	define ASGN_CHARP(p, pc) (pc = * ((char **)(p)))
#	define ASGN_LONG(p, pc) (pc = * ((long *)(p)))

# else
/*** Mac and IBM 370 Version ***/
/* p..pointer to value, pc..target */
#	define ASGN_INT(p, pc)  {     \
		char *TMPP = (char *)(p); \
		char *ASGN = (char *)&pc; \
		*(ASGN++) = *(TMPP++);	  \
		*(ASGN)   = *TMPP;		  \
	}
 
#	define ASGN_CHARP(p, pc) {    \
		char *TMPP = (char *)(p); \
		char *ASGN = (char *)&pc; \
		*(ASGN++) = *(TMPP++);	  \
		*(ASGN++) = *(TMPP++);	  \
		*(ASGN++) = *(TMPP++);	  \
		*(ASGN)   = *TMPP;        \
	}
 
#	define ASGN_LONG(p, pc) {     \
		char *TMPP = (char *)(p); \
		char *ASGN = (char *)&pc; \
		*(ASGN++) = *(TMPP++);	  \
		*(ASGN++) = *(TMPP++);	  \
		*(ASGN++) = *(TMPP++);	  \
		*(ASGN)   = *TMPP;		  \
	}
# endif


# if !defined(IBM370) && !defined(FOR_OS_SUNLINUX)
	/* memory copying macros. */
#	define wrcharp_to_mem(ppc, pc) (*(char **) (pc) = (ppc))
#	define wrlong_to_mem(l, pc) (*(long *)(pc) = (long) (l))

# else
/*** Sun (Linux), MAC and IBM 370 Versions ***/
/* ppc..value, pc..pointer to target */
# 	define wrcharp_to_mem(ppc, pc) { \
		char *TMP = (char *)(ppc);	 \
		char *TMPP = (char *)&TMP;	 \
		char *ASGN = (char *)(pc);	 \
		*(ASGN++) = *(TMPP++);		 \
		*(ASGN++) = *(TMPP++);		 \
		*(ASGN++) = *(TMPP++);		 \
		*(ASGN)   = *(TMPP);		 \
	}
 
/* l..value, pc..pointer to target */
# 	define wrlong_to_mem(l, pc){   \
		long TMP   = (long)(l);    \
		char *TMPP = (char *)&TMP; \
		char *ASGN = (char *)(pc); \
		*(ASGN++) = *(TMPP++);	   \
		*(ASGN++) = *(TMPP++);     \
		*(ASGN++) = *(TMPP++);     \
		*(ASGN)   = *(TMPP);       \
	}

# endif

