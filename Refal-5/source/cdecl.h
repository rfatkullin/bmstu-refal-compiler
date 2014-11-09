#include "bif_lex.h"

#ifdef FOR_OS_DOS
#	ifndef FILENAME_MAX
#		define FILENAME_MAX 256
#	endif
#	ifndef FOPEN_MAX
#		define FOPEN_MAX 32
#	endif
#endif


# ifndef LINT_ARGS
# define LINT_ARGS 1
# endif

	/*** Include Files. ***/
#ifndef FILES_INCLUDED
#  include <stdio.h>
#  include <string.h>
#  include <ctype.h>
#  ifndef IBM370
#    include <malloc.h>
#  endif
#  include <stdlib.h>
#  define FILES_INCLUDED 1
#endif


	/*** External definitions. ***/
# ifndef EXTERNAL
#   ifdef DEFINE_EXTERNALS
#     define EXTERNAL 
#   else
#     define EXTERNAL extern
#   endif
# endif


	/* define macros.	*/

# define copyst(x) {strcpy(x,str);}

# define pushstk vstack [vsp++] = nlv
# define clrstk {vsp = 0; nlv = 0;}
# define popstk nlv = vstack [--vsp]

# define write_byte(BYTENUM) fputc ((char) (BYTENUM), fdtmpw)
# define write_bytes(PTR, NUMBYTES) \
	fwrite ((void *) (PTR), 1, (NUMBYTES), fdtmpw)
# define write_long(LONGNUM) \
	{long TMP; TMP = (long)(LONGNUM); fwrite((void *) &TMP,sizeof(long),1,fdtmpw);}

# define read_byte(X) fread((void *) &(X),1,1,fdtmpr)
# define read_long(X) fread((void *) &(X),sizeof(long),1,fdtmpr)

	/* defined constants	*/

# define ASCIIV 257
# define ID 258
# define NUMBER 259
# define EXTRN 260
# define VAR 261
# define MDIGIT 262
# define COMPSYM 263
# define STRING 264
# define LCBRAK 265
# define FDEF 266
# define SENTS 267
# define ST 268
# define RCS1 269
# define RCS2 270
# define RCS3 271
# define RCS4 272
# define LSF 273
# define RSF 274
# define LSF1 275
# define CONST 276
# define PAR 277
# define RSF1 278
# define ENTRY 279
# define RSF1B 280

		/* other '# define's. */

	/* number of bytes in a block of memory. (for easy memory allocation). */
# define MEM_BLK_SIZE 512
	/* maximum number of variables in a sentence. */
# define MAX_TABLE_LENGTH 512
	/* maximum number of sentences per function. */
# define MAX_SENTENCES 200
	/*  maximum integer  2**31 -1	*/
# define MAXINT 2147483647L
	/* maximum length allowable for an identifier.	*/

	/*  maximum unsigned integer  */
# define MAX_INT (2 << (8*(sizeof int) - 1) -1)
# define MAX_UNSIGNED_INT (MAXINT + MAXINT + 1)

#define MAXSTR 200
	/* name of the null device. */
# define NULL_DEVICE "NUL"

	/* maximum size of a compound symbol. */
# ifndef MAXWS
/*# define MAXWS 32*/
#	define MAXWS 1024
# endif


	/* various types of links in an expression. 
			also: NULL = empty expression, STRING = string. */
# define E_VAR 1
# define S_VAR 2
# define T_VAR 3
# define CHAR 4
# define ATOM 5
# define DIGIT 6
# define LPAR 7
# define RPAR 8
# define ACT_LEFT 9
# define ACT_RIGHT 10

/* declarations of global structures and variables.	*/

struct element {
  int type;	     /* e- s- t- variables, or constants: */
  union {
    char *f;	     /* a pointer to a function or atom name or string */
    unsigned long n; /* an integer. */
    char c;	     /* character */
    int i;	   /* variable: its number; or right par: index to the pair */
  } body;
  int number;	   /* projection number. or left par: index to right. */
};

typedef union {
  struct element *chunk;
  struct node *tree;
  struct functab *func;
  char *pchar;
  int number;
} branch_t;

EXTERNAL branch_t zero;

struct node {
  int nt;         /* node type */
  branch_t a2, a3, a4;
};

EXTERNAL struct node *Error;

/* parameter for RASL instructions. */
union param {
  char c;
  char *f;
  unsigned long n;
  struct { /* two parameters. */
    int i1;
    int i2;
  } d;
  int i;
};

struct rasl_instruction {
  int code;
  union param p;
  struct rasl_instruction *next;
};

/* pointer to an array of RASL instructions ... */
EXTERNAL struct rasl_instruction *ftransl;

struct HOLES {
  int left;
  int right;
  int lte;
  int rte;
  struct HOLES *next;
};

/* table of variables and its length */
struct TABLE {
  int index;
  int te_offset;
};

EXTERNAL struct TABLE table [MAX_TABLE_LENGTH];
EXTERNAL int table_len;
	

struct functab {
  /*char name[MAXWS];*/     /* name */
	char * name;
  long offset;		/* offset or identificator. */
  struct functab *next; /* pointer to the next entry */
};

/* the first of the list of defined labels. It contains
 * all functions and auxiliary labels as well.
 */
EXTERNAL struct functab *ft;

/* the first of the list of called functions. */
EXTERNAL struct functab *fc;

/* the first of the list of back-up functions. This list contains only those
 * labels, which are actually present in the Refal program.
 */
EXTERNAL struct functab *fb;

EXTERNAL struct functab
  *fx,  /* the first of the list of external fucntions. */
  *fe,	/* the first of the list of entry functions. */
  *cs;	/* the first of the list of compound symbols. */

EXTERNAL long 
  z,		/* Current offset in the code file. */
  cscount,	/* Count of compound symbols. */
  xtcount,	/* Count of external functions. */
  btcount,	/* Count of local functions. */
  ntcount; 	/* Count of entry functions. */


/* built in function table. */
struct locvar {
  char vt;               /* type */
  char vindex[MAXWS];    /* index */
};
EXTERNAL struct locvar lv [128];	/* table of local variables. */


/* Other global variables. */

EXTERNAL int line_no;	/* number of the current line.	*/
EXTERNAL char cbuf [MAXSTR+1];	/* input stream buffer.	*/

/*EXTERNAL char title [MAXWS];*/	/* title of the output. */
EXTERNAL char title [FILENAME_MAX];

EXTERNAL short sc;	/* input stream buffer pointer.	*/

EXTERNAL char 
  str [MAXWS + 1],      /* location for current comp symbol.*/
  strings [MAXSTR + 1], /* location for current string.	*/
  last_fn [MAXWS + 1],  /* last function name. */
  globsav,	    /* next input character. */
  vtype,	    /* type of the current variable. */
  *block,	    /* pointer for the current block. */
  c_flags [30],     /* flags: from command line. (30 - arbitrary number) */
  *nonret;	    /* pointer to nonret memory */

EXTERNAL int
  blptr,   /* pointer to first free memory location in the block. */
  nrptr,   /* ptr to nonret memory */
  length,  /* current length.*/
  nerrors, /* number of errors counter */
  token,   /* current token.	*/
  last_label; /* the number of the first available label number. */

	
EXTERNAL unsigned long v;

EXTERNAL int
  nlv,          /* number of local variables. */
  vsp,          /* pointer to top of vstack */
  vstack [256]; /* stack of currently visible number of local variables. */


/* Files: */

EXTERNAL FILE *fdref;	/* .REF file */
EXTERNAL FILE *fdlis;	/* .LIS file */
EXTERNAL FILE *fdtmpw;	/* .TMP write file */
EXTERNAL FILE *fdtmpr;	/* .TMP read file */


#include "junk.h"
