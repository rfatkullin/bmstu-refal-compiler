
	/* include statements.	*/
/* FOR unlink */
#ifdef FOR_OS_LINUX
  #include <unistd.h>
#endif

# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"


# define file_error(NAME) {                   \
  fprintf (stderr, "Can\'t open %s\n", NAME); \
  exit (1);                                   \
}

extern long nbi;
extern struct bitab bi[];


	/* suppress debugging. */
# define MDEBUG 0

int rc_getbeginfile (char * argv []) {
	int i = 1;
	char filename [FILENAME_MAX], * s, * dot, outfile [FILENAME_MAX];

	if (argv [i] == NULL) {
		ri_inquire ("Refal input name:", filename, FILENAME_MAX);
	} else {
		while (argv [i] != NULL && argv [i][0] == '-') i ++;
		if (argv [i] == NULL) return -1;

		/* Initialize the refal compiler. July 28, 1985 D.T.	*/
		/* open the input (.ref) output (.lis) and result (.tmp) files. */
		strcpy (filename, argv [i]);
	}

	outfile [0] = '\0';
	/* Get the .ref file */
	if ((dot = strrchr (filename, '.')) != NULL) {
		if ((dot [1] == 'r' || dot [1] == 'R') &&
			(dot [2] == 'e' || dot [2] == 'E') &&
			(dot [3] == 'f' || dot [3] == 'F') &&
			dot [4] == '\0') {
			s = filename;
		} else {
			s = strcat (filename, ".ref");
		}
		/*s = (strcmp (dot, ".ref") != 0)? strcat (filename, ".ref"): filename;*/
	} else s = strcat (filename, ".ref");
	if ((fdref = fopen (s, "rt")) == NULL) file_error (s);

	/* Get the title. */
	if (outfile [0] == '\0') strcpy (outfile, filename);
	if ((dot = strrchr (outfile, '.')) != NULL) {
		if (strcmp (dot, ".ref") == 0 || strcmp (dot, ".rsl")) {
			*dot = '\0';
		}
	}
	strcpy (title, outfile);

	/* Get the .tmp file. */
# ifdef IBM370
	s = strcat (outfile, " tmp (bin lrecl 1");
# else 
	s = strcat (outfile, ".tmp");
# endif

	if ((fdtmpw = fopen (s, "wb")) == NULL) file_error (s);

	/* Get the .lis file */
	if ((dot = strrchr (filename, '.')) != NULL) *dot = '\0';
	s = strcat (filename, ".lis");
	if ((fdlis = fopen (s, "wt")) == NULL) file_error (s);
	return i;
}

int rc_getnextfile (int i, char * argv []) {
	char filename [FILENAME_MAX], * s, * dot, outfile [FILENAME_MAX];

	while (argv [i] != NULL && argv [i][0] == '-') i ++;
	if (argv [i] == NULL) return -1;

	/* Initialize the refal compiler. July 28, 1985 D.T.	*/
	/* open the input (.ref) output (.lis) and result (.tmp) files. */
	strcpy (filename, argv [i]);
	outfile [0] = '\0';
	/* Get the .ref file */
	if ((dot = strrchr (filename, '.')) != NULL) {
		if ((dot [1] == 'r' || dot [1] == 'R') &&
			(dot [2] == 'e' || dot [2] == 'E') &&
			(dot [3] == 'f' || dot [3] == 'F') &&
			dot [4] == '\0') {
			s = filename;
		} else {
			s = strcat (filename, ".ref");
		}
		/*s = (strcmp (dot, ".ref") != 0)? strcat (filename, ".ref"): filename;*/
	} else s = strcat (filename, ".ref");
	if ((fdref = fopen (s, "rt")) == NULL) file_error (s);

	/* Get the title. */
	if (outfile [0] == '\0') strcpy (outfile, filename);
	if ((dot = strrchr (outfile, '.')) != NULL) {
		if (strcmp (dot, ".ref") == 0 || strcmp (dot, ".rsl")) {
			*dot = '\0';
		}
	}
	strcpy (title, outfile);
	/* title must be in capital letters. */
	/* strupr (title); */

	/* Get the .tmp file. */
# ifdef IBM370
	s = strcat (outfile, " tmp (bin lrecl 1");
# else 
	s = strcat (outfile, ".tmp");
# endif

	if ((fdtmpw = fopen (s, "wb")) == NULL) file_error (s);

	/* Get the .lis file */
	if ((dot = strrchr (filename, '.')) != NULL) *dot = '\0';
	s = strcat (filename, ".lis");
	if ((fdlis = fopen (s, "wt")) == NULL) file_error (s);
	return i;
}

int rc_initcom (int argc, char * argv []) {
	/*char filename [FILENAME_MAX], *s, *dot, outfile [FILENAME_MAX];*/
	int i;
  
	/* Initialize the refal compiler. July 28, 1985 D.T. */
	/* open the input (.ref) output (.lis) and result (.tmp) files. */

	/*filename [0] = '\0';*/
	/*outfile [0] = '\0';*/

	/* process the flags and arguments. */
	/* currently only the following flags are recognized:
	 * '-h' : print an help message. 
	 * '-l' : produce the listing files. 
	 */
        rc_options(argc,argv);
/* 
* Nemytykh 12.21.2008
	c_flags [0] = '\0';
	for (i = 1; i < argc; i ++) {
		if (argv [i][0] == '-') strncat (c_flags, &(argv [i][1]), 30);
		/*
		else if (filename [0] == '\0') strcpy (filename, argv[i]);
		else if (outfile [0] == '\0') strcpy (outfile, argv[i]);
		*+/
	}
*/

	/*
	if (filename [0] == '\0')
		ri_inquire ("Refal input name:", filename, FILENAME_MAX);
	*/
	/* Get the .ref file */
	/*
	if ((dot = strrchr (filename, '.')) != NULL) {
		if ((dot [1] == 'r' || dot [1] == 'R') &&
			(dot [2] == 'e' || dot [2] == 'E') &&
			(dot [3] == 'f' || dot [3] == 'F') &&
			dot [4] == '\0') {
			s = filename;
		} else {
			s = strcat (filename, ".ref");
		}
		/ * s = (strcmp (dot, ".ref") != 0)? strcat (filename, ".ref"): filename; * /
	} else s = strcat (filename, ".ref");
	if ((fdref = fopen (s, "rt")) == NULL) file_error (s);
	*/
	/* Get the title. */
	/*
	if (outfile [0] == '\0') strcpy (outfile, filename);
	if ((dot = strrchr (outfile, '.')) != NULL) {
		if (strcmp (dot, ".ref") == 0 || strcmp (dot, ".rsl")) {
			*dot = '\0';
		}
	}
	strcpy (title, outfile);
	*/
	/* title must be in capital letters. */
	/* strupr (title); */

	/* Get the .tmp file. * /
# ifdef IBM370
	s = strcat (outfile, " tmp (bin lrecl 1");
# else 
	*/
	/*s = strcat (outfile, ".tmp");*/
/*# endif*/

	/*if ((fdtmpw = fopen (s, "wb")) == NULL) file_error (s);*/

	/* Get the .lis file */
	/*
	if ((dot = strrchr (filename, '.')) != NULL) *dot = '\0';
	s = strcat (filename, ".lis");
	if ((fdlis = fopen (s, "wt")) == NULL) file_error (s);
	*/
	return 0;
}

/* Saves just the offset and terminating 0. */
int rc_sbtable (struct functab * table) {
	struct functab *t;

	for (t = table; t != NULL; t = t->next) write_long (t -> offset);
	return 0;
}

/* Saves the table along with the offset. */
int rc_sftable (struct functab * table) {
	struct functab * t;
	int i;

	for (t = table; t != NULL; t = t->next) {
		for (i = 0; /*i < MAXWS && */t -> name [i] != '\0'; i++) {
			write_byte (t -> name [i]);
		}

		/*for (; i < MAXWS; i++)*/ write_byte ('\0');
		write_long (t -> offset);
	}
	return 0;
}

/* Saves the table without the offset just as a list. */
int rc_sltable (struct functab * table) {
	struct functab *t;
	int i;

	for (t = table; t != NULL; t = t->next) {
		for (i = 0; /*(i < MAXWS) && */t -> name[i] != '\0'; i++) {
			write_byte (t -> name [i]);
		}
		/*for (; i < MAXWS; i++)*/ write_byte ('\0');
	}
	return 0;
}


/* Copies from fdtmpr to fdtmpw file replacing the
 * references for their offsets.
 */
int rc_pass2 (void) {
	unsigned char opcode, c, d;
	char lname[MAXWS];
	int i;
	long k;
	struct functab *function;

	while (read_byte (opcode) == 1) {
		switch (opcode) {

		/* This RASL instruction takes an address of a function as an argument. */
		case ACT1:
			read_byte (c);
			for (i = 0; i < MAXWS; i ++) {
				read_byte (c);
				lname [i] = c;
				if (c == '\0') break;
			}
			/*
			for (i = 0; i < MAXWS; i++) {
					read_byte (c);
					lname[i] = c;
			}
			*/
			/* Look for the function in the table. If found then output
			 * ACT1 <offset>, if not then check that it is an external
			 * function. If so, output ACT1N <lname>.
			 */

			if ((function = searchf (lname, ft)) != NULL) {
				write_byte (ACT1);
				write_long (function -> offset);
			} else {
				write_byte (ACT_EXTRN);
				for (i = 0; /*i < MAXWS*/ lname [i] != 0; i++) {
					write_byte (lname[i]);
				}
				write_byte ('\0');
			}
			break;

		/* These RASL operators take a compound symbol as an argument. */
		case CSYM: case CSYMR: case NCS:
		/* These RASL operators require a (long) number as a parameter */
		case NSYM: case NSYMR: case NNS: case BUILT_IN:
			read_long (k);
			write_byte (opcode);
			write_long (k);
			break;

		/* Builtin function call with an argument */
		case BUILT_IN1:
			read_long (k);	/** First one is zero. **/
			read_long (k);
			write_byte (opcode);
			write_long (0L);
			write_long (k);
			break;

		/* No arguments for these RASL operators. */
		case BL: case BLR: case BR: case CL: case EMP: case EST: case PLEN:
		case PLENS: case PLENP: case PS: case PSR: case TERM: case TERMR:
		case LEN: case LENP: case VSYM: case VSYMR: case OUTEST: 
		case POPVF: case PUSHVF: case STLEN:
			write_byte (opcode);
			break;

		/* These RASL operators require a long (number) as parameter. */
		case MULE: case MULS: case TPLE: case TPLS:
			read_long (k);
			write_byte (opcode);
			write_long (k);
			break;

		/* These RASL operators require a byte (character) as parameter. */
		case SYM: case SYMR: case LENS: case NS: /*case MULE: case MULS: */
		case OEXP: case OEXPR: case OVSYM: case OVSYMR: case RDY:
		case LENOS:/* case TPLE: case TPLS:*/
			read_byte (c);
			write_byte (opcode);
			write_byte (c);
			break;

		/* This RASL operator require two operands of size 1 byte. */
		case SETB:
			write_byte (opcode);
			read_long (k);
			write_long (k);
			read_long (k);
			write_long (k);
			break;

		/* These RASL operators take 1 byte and a variable number of bytes as parameters. */
		case SYMS: case SYMSR: case TEXT:
			read_byte (d);
			write_byte (opcode);
			write_byte (d);
			for (i = 0; i < d; i++) {
				read_byte (c);
				write_byte (c);
			}
			break;

		/* These RASL operators take as argument a label of form
		 * FUNNAME$NUMBER, where FUNNAME is the current function
		 * name, and NUMBER is a number.
		 */
		case TRAN: case ECOND:
			for (i = 0; i < MAXWS; i++) {
				read_byte (c);
				lname[i] = c;
				if (c == 0) break;
			}
			write_byte (opcode);
			function = searchf (lname, ft);
			if (function) write_long (function -> offset)
			else {
				fprintf (stderr, "PASS2: Could not find label %s, opcode %d\n", lname, opcode);
				write_long (0L);
			}
			break;

		/* These RASL operators define a label of form
		 * FUNNAME$NUMBER, where FUNNAME is the current function
		 * name, and NUMBER is a number.
		 */
		case LBL: case LABEL:
		case L: case E:
			write_byte (opcode);
			read_byte (c);
			write_byte (c);
			for (i = 0; i < MAXWS; i++) {
				read_byte (c);
				write_byte (c);
				if (c == 0) break;
 			}
			break;
	 
		default:
			fprintf (stderr, "PASS2 %d: Strange Opcode\n", opcode);
			break;
		}
	}

	return 0;
}

		/*  March 8, 1987. D.T.  */
	/* This function performs the second pass over the code and 
		finishes up the creation of .RSL file. */
int rc_end (void) {
	long fnumb, save_offset;
	char fname[FILENAME_MAX];
	struct functab *ff, *f2, *ft_last, *fb_last;
	int k, i, is_main;

	/* 1. Necessary stuff. Check the view field. */
	nerrors = 0;
	/* 1a. see if this is a main module. (Main module is one in which
	 * entry function GO is defined.
	 */
	is_main = 0;
	if (searchf ("GO", fe) != NULL || searchf ("Go", fe) != NULL) is_main = 1;

	/* 2. If this is a main module, produce all built in function
	 * definitions, and declare them as entry. (skip 0) 
	 */
	if (is_main) {
		/* 2a. search to the end of the function table. */
		/*ft_last = ft;
		while (ft_last -> next != NULL) ft_last = ft_last -> next;*/
		for (ft_last = ft; ft_last -> next != NULL; ft_last = ft_last -> next);

		/*fb_last = fb;
		while (fb_last -> next) fb_last = fb_last -> next;*/
		for (fb_last = fb; fb_last -> next != NULL; fb_last = fb_last -> next);

		/* 2b. process all built-in functions. */
		for (k = 1; k < nbi; k ++) {
			/* see if already defined as entry. */
			if (searchf (bi [k].fname, fe) != NULL) {
				rc_serror (204, bi [k].fname);
				continue;
			}
			write_byte (E);
			for (write_byte ('\0'), i = 0; /*i < MAXWS && */ bi [k].fname [i] != '\0'; i++) {
				write_byte (bi [k].fname [i]);
			}
			/*for (; i < MAXWS; i++) */write_byte ('\0');
			save_offset = (z += i + 2/*MAXWS*/);
			fnumb = bi [k].fnumber;
			if ((bi [k].flags & BI_FADDR) == 0) {
				write_byte (BUILT_IN);
				write_long (fnumb);
				z += sizeof (char) + sizeof (char *);
			} else { /*** Functions MU, UP and EV-MET ***/
				write_byte (BUILT_IN1);
				write_long (0L);
				write_long (fnumb);
				z += sizeof (char) + sizeof (char *) + sizeof (long);
			}

			/* Add to the defined and backup tables. */
			/* Note: this may add a second definition of a function with
			 * this name, so we want to add this definition at the
			 * tail of the list (so it the first definition will 
			 * be found first). 
			 */
			/* Insert AFTER ft_last and fb_last */
			ff = (struct functab *) malloc (2 * sizeof (struct functab));
			/* 337 line. Not check result of malloc. Shura. 29.01.98 */
			if (NULL == ff) {
				fprintf (stderr, "No memory for function table\n");
				exit (1);
			}
			f2 = ff + 1;
			ff -> next = NULL;
			fb_last -> next = ff;
			fb_last = ff;
			ff -> offset = save_offset;

			/* Added */
			if (NULL == (ff -> name = (char *) malloc (strlen (bi [k].fname) + 1))) {
				fprintf (stderr, "No memory for function name\n");
				exit (1);
			}
			strcpy (ff -> name, bi [k].fname);

			btcount ++;
			f2 -> next = NULL;
			ft_last -> next = f2;
			ft_last = f2;
			f2 -> offset = save_offset;

			if (NULL == (f2 -> name = (char *) malloc (strlen (bi [k].fname) + 1))) {
				fprintf (stderr, "No memory for function name\n");
				exit (1);
			}
			strcpy (f2 -> name, bi [k].fname);

			/* make it entry function. */
			rc_mkentry (bi [k].fname);
			fe -> offset = save_offset;
		}
	} else {
		/* 3. this is not a main module: define all built ins as external.
		 *    except those that take as argument the pointer to the 
		 *    function table. (like MU UP etc.) 
		 */
		for (k = 1; k < nbi; k ++) {
			/* check that it is not already defined. */
			if (searchf (bi [k].fname, ft) == NULL && searchf (bi [k].fname, fx) == NULL) {
				/* make it external if it is a regular function. */
				if ((bi [k].flags & BI_FADDR) == 0) {
					rc_mkextrn (bi [k].fname);
				} else { /* otherwise create a definition for this function. */
					write_byte (L);
					for (write_byte ('\0'), i = 0;/* i < MAXWS && */bi [k].fname [i] != '\0'; i++) {
						write_byte (bi [k].fname [i]);
					}
					/*for (; i < MAXWS; i++) */write_byte ('\0');
					save_offset = (z += i + 2 /*MAXWS*/);
					fnumb = bi [k].fnumber;
					write_byte (BUILT_IN1);
					write_long (0L);
					write_long (fnumb);
					z += sizeof (char) + sizeof (char *) + sizeof (long);
					/* add it to the list of defined functions. */
					f2 = searchf (bi [k].fname, ft);
					if (f2 == NULL) {
						f2 = (struct functab *) rc_allmem (sizeof (struct functab));
						if (f2 == NULL) {
							fprintf (stderr,"Ran out of memory.\n");
							exit (1);
						}
						f2 -> next = ft;
						f2 -> offset = save_offset;
						ft = f2;

						if (NULL == (f2 -> name = (char *) malloc (strlen (bi [k].fname) + 1))) {
							fprintf (stderr, "No memory for function name\n");
							exit (1);
						}
						strcpy (f2 -> name, bi [k].fname);
					}
					/* now add it to the list of backup functions (local functions) */
					f2 = searchf (bi [k].fname, fb);
					if (f2 == NULL) {
						f2 = (struct functab *) rc_allmem (sizeof (struct functab));
						if (f2 == NULL) {
							fprintf (stderr,"Ran out of memory.\n");
							exit (1);
						}
						f2 -> next = fb;
						f2 -> offset = save_offset;
						fb = f2;

						if (NULL == (f2 -> name = (char *) malloc (strlen (bi [k].fname) + 1))) {
							fprintf (stderr, "No memory for function name\n");
							exit (1);
						}
						strcpy (f2 -> name, bi [k].fname);
						btcount ++;
					}
				}
			}
		}
	}

	/* 4. See that all called functions are defined, external or built-in */
	for (ff = fc; ff != NULL; ff = ff -> next) {
		if (searchf (ff -> name, ft) == NULL && searchf (ff -> name, fx) == NULL) {
			/* Built in function: by now should be either on the defined list
			 * or on the external list. 
			 */
			rc_serror (200, ff -> name);
		}
	}

	/* 5. Check that there is an ENTRY function. */
	if (fe == NULL) rc_serror (202, NULL);

	/* 6. Close and reopen the .tmp file, and open .RSL file. */
	fclose (fdref);
	fclose (fdlis);
	fclose (fdtmpw);

# ifdef IBM370
	sprintf (fname, "%s tmp (bin lrecl 1", title);
# else 
	sprintf (fname, "%s.tmp", title);
# endif

	fdtmpr = fopen (fname, "rb");
	if (fdtmpr == NULL) file_error (fname);

# ifdef IBM370
	sprintf (fname, "%s rsl (bin lrecl 1", title);
# else 
	sprintf (fname, "%s.rsl", title);
# endif

	fdtmpw = fopen (fname, "wb");
	if (fdtmpw == NULL) file_error (fname);

	/* 7. Write the title, size of code, size of entry, external and compound
	 *    symbol tables.
	 */
	/*for (i=0; (i<MAXWS) && (title[i] != '\0'); i++)*/
	for (i=0; (i<FILENAME_MAX - 1) && (title[i] != '\0'); i++)
		write_byte (title[i]);
	/*for (; i < MAXWS; i++) write_byte ('\0');*/
	/*for (; i < FILENAME_MAX; i++)*/ write_byte ('\0');
	write_long (z);
	write_long (ntcount);
	write_long (xtcount);
	write_long (cscount);
	write_long (btcount);

	/* 6. Save the entry function table. */
	rc_sftable (fe);

	/* 7. Save the external functions table. */
	rc_sltable (fx);

	/* 8. Save the compound symbol table. */
	rc_sltable (cs);

	/* 9. Save the local function table. */
	rc_sbtable (fb);

	/* 10. Save the code resolving references. */
	rc_pass2 ();

	/* 11. Close the files and delete .tmp file. */
	fclose (fdtmpr);
	fclose (fdtmpw);
	sprintf (fname, "%s.tmp", title);
	if (unlink (fname) == -1) 
		fprintf (stderr, "Unable to delete %s\n", fname);

	/* 12. Check for errors and exit.  */
	return 0;
}


