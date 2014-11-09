
# include "decl.h"
# include "macros.h"
# include "fileio.h"
# include "ifunc.h"

# define ESCAPE '\\'
/*# define ESCAPE '#'*/

#	define MAXINT 2147483647
	/*  maximum unsigned integer  */
# define MAX_INT (2 << (8*(sizeof int) - 1) -1)
# define MAX_UNSIGNED_INT (MAXINT + MAXINT + 1)

/* From refio.c */
extern void correctEscapeSymbols (char *);

/* Local: */
int putexp (FILE *, unsigned long);

int rf_sysfun (void) {
	b1 = b1->foll;
	if (!LINK_NUMBER(b1)) {
		ri_error(8);
	}
	if (b1 -> pair.n == 1) {
		FILE * fp;
		char s [FILENAME_MAX + 1];
		int i;

		/* read file name */
		for (i = 0, b1 = b1 -> foll;b1 != b2 && i < FILENAME_MAX; i++, b1 = b1 -> foll) {
			/* insert checking */
			s[i] = b1 -> pair.c;
		}
		s[i] = '\0';
		if ((fp = fopen(s, "r")) == NULL) {
			fprintf(stderr, "Can\'t open %s\n",s);
			rdy(0);
			out(2);
			est;
			return (0);
		}
	
		getexp (fp);
	} else {
		ri_error (8);
	}
	return 0;
}

int rf_desysfun (void) {
	char ca_file [FILENAME_MAX + 1];
	unsigned long l;
	FILE * fp_out;

	if (! LINK_RSTRUCTB (b2 -> prec)) {
		ri_error (8);
		return 1;
	}
	/* Read file name */
	for (l = 0, b1 = b1 -> foll; b1 != b2 && ! LINK_LSTRUCTB (b1); b1 = b1 -> foll) {
		if (!LINK_CHAR (b1)) {
			ri_error (8);
			return 1;
		}
		ca_file [l ++] = b1 -> pair.c;
	}
	ca_file [l] = '\0';
	if (! LINK_LSTRUCTB (b1)) {
		ri_error (8);
		return 1;
	}
	if ((fp_out = fopen (ca_file, "w")) == NULL) {
		fprintf (stderr, "Cannot open file \'%s\' for writing.\n", ca_file);
		ri_error (8);
		return 1;
	}

	/* Read the limit of string length */
	b1 = b1 -> foll;
	if (! LINK_NUMBER (b1)) {
		ri_error (8);
		return 1;
	}
	l = b1 -> pair.n;

	/* Put data into the file */
	b1 = b1 -> foll;
	b2 = b2 -> prec;
	{
		int i;

		i = putexp (fp_out, l);
		fclose (fp_out);
		rdy (0);
		out (2);
		est;
		return i;
	}
}

int getexp (FILE *fp) {
	int brcnt = 0;
	long k;
	unsigned long l = (unsigned long)MAX_UNSIGNED_INT;
	int c;
	char str[MAXWS+1], *pc, *ri_cs_impl();
	int state_chars = 0;

	rdy (0);

	/* read file */
	while ((c = getc(fp)) != EOF) {
		switch(c) {
		case '\n': 
			break;

		case ESCAPE:
			while ((c = getc(fp)) == '\n');
			if (c == ESCAPE || c == '\'' || c == '"' || c == '(' || c == ')' ||
				c == '<' || c == '>' || c == 'n' || c == 'r' || c == 't') {
				if (c == 'n') {
					ns ('\n');
				} else if (c == 'r') {
					ns ('\r');
				} else if (c == 't') {
					ns ('\t');
				} else {
					ns(c);
				}
			} else if (c == 'x') {
				int i_w;
				char c_w = 0;

				for (i_w = 0; i_w < 2; i_w ++) {
					if (EOF == (c = getc (fp))) {
						fclose (fp);
						ri_error (8);
					}
					if (isxdigit (c)) {
						c = (c <= '9')?	c - '0': ((isupper (c))?	c - 'A' + 10: c - 'a' + 10);
						c_w += ((i_w)? c: 16 * c);
					} else {
						fclose (fp);
						ri_error (8);
					}
				}
				ns (c_w);
			} else {
				fclose (fp);
				ri_error(8);
			}
			break;
		
		case '(':
			if (state_chars) {
				ns ('(');
			} else {
				bl; 
				brcnt ++;
			}
			break;

		case ')':
			if (state_chars) {
				ns (')');
			} else {
				br;
				brcnt --;
				if (brcnt < 0) {
					fclose (fp);
					ri_error (8);
				}
			}
			break;

		case '\'':
			state_chars = (state_chars) ? 0: 1;
			break;

		case ' ':
			if (state_chars) {
				ns (' ');
			}
			break;

		case '"':
			if (state_chars) {
				/*
				fclose (fp);
				ri_error (8);
				*/
				ns ('"');
			} else {
				int i;

				for (i = 0, c = getc (fp); EOF != c && c != '"'; c = getc (fp), i ++) {
					if (c == ESCAPE) {
						while ((c = getc(fp)) == '\n');
						if (c == '"' || c == ESCAPE || c == '\'' || c == '(' || c == ')' ||
							c == '<' || c == '>' || c == 'n' || c == 't' || c == 'r') {
							switch (c) {
							case 'n': str [i] = '\n'; break;
							case 'r': str [i] = '\r'; break;
							case 't': str [i] = '\t'; break;
							default: str [i] = c;
							}
						} else {
							fclose (fp);
							ri_error (8);
						}
					} else if (c == '\n') {
						i --;
						continue;
					} else {
						str [i] = c;
					}
				}
				if (c != '"') {
					fclose (fp);
					ri_error (8);
				}
				str [i] = 0;
				pc = ri_cs_impl (str);
				ncs (pc);
			}
			break;

		default:
			if (state_chars) {
				ns (c);
			} else if (isalpha (c)) {
				int i;
				/* Get composymbol (string) */

				for (str [0] = c, i = 1; (c = fgetc (fp)) != EOF && i < MAXWS; i ++) {
					if (c == ' ') break;
					if (isalnum (c) || c == '_' || c == '$' || c == '-') str [i] = c;
					else if (c == '\n') i --;
					else {
						/* Error */
						fclose (fp);
						ri_error (8);
					}
				}
				/* I hope 1kb of buffer for saving ID is enough. */
				str [i] = 0;
				pc = ri_cs_impl (str);
				ncs (pc);
			} else if (isdigit (c)) {
				/* get number */
				k = c - '0';
				while (EOF != (c = getc (fp))) {
					if (isdigit (c)) {
/* Nemytykh 09.08.2002 */
        	/* The next few lines are strange because otherwise the "C" compiler
	        	produces an integer overflow warning */
		                                if ( (unsigned long)k > ((l - (c - '0'))/10 + (l - (c - '0'))%10) ) { 
                		 			  fclose (fp);
		                 		          ri_error (8);
				                } else  { k = k*10 + (c - '0'); };
/* Was:
						k = k * 10 + (c - '0');
						if (k >= MAXINT || k < 0) {
							nns (k);
							k = 0;
						}
*/
					} else if (c != '\n') break;
				}
				if (c != ' ') {
					fclose (fp);
					ri_error (8);
				}
				nns (k);
			} else  {
				ns (c);
			}
		}
	}
	fclose (fp);
	if (brcnt) {
	 	fclose (fp);
	 	ri_error(8);
	}
	out(2);
	est;
	return 0;
}

int putexp (FILE * fp, unsigned long l_limit) {
	unsigned long l;
	int state_char = 0, i_len;
	char ca_buf [2 * MAXWS + 1], * cp;
	

	for (l = 0; b1 != b2; b1 = b1 -> foll) {
		char c;

		if (l >= l_limit) {
			fputc ('\n', fp); l = 0;
		}
		switch (b1 -> ptype) {
		case LINK_TYPE_CHAR:
			c = b1 -> pair.c;
			if (state_char == 0) {
				state_char = 1;
				fputc ('\'', fp); l ++;
			}
			if (l >= l_limit) {
				fputc ('\n', fp); l = 0;
			}
			if (c == '\n') {
				fputc (ESCAPE, fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				fputc ('n', fp); l ++;
			} else if (c == '\r') {
				fputc (ESCAPE, fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				fputc ('r', fp); l ++;
			} else if (c == '\t') {
				fputc (ESCAPE, fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				fputc ('t', fp); l ++;
			} else if (c == ESCAPE || c == '\'' || c == '"' || c == '(' || c == ')' || c == '<' || c == '>') {
				fputc (ESCAPE, fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				fputc (c, fp); l ++;
			} else if (c < ' ' || c > 127) {
                                char zero [2];

				fputc (ESCAPE, fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				if (l + 3 >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
/*+ The padding for heximal represantation does not work. 
				fprintf (fp, "x%0x", c); l += 3;
*/
        			zero[0] = (c < 16) ? '0':'\0';  zero[1] = '\0';
				fprintf (fp, "x%s%x", zero, c); l += 3;
			} else {
				fputc (c, fp); l++;
			}
			break;

		case LINK_TYPE_COMPSYM:
			if (state_char == 1) {
				fputc ('\'', fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				state_char = 0;
			}
			if (isIdentifier (b1 -> pair.f) == 0) {
				/* no identifier */
				char ca_w [2 * MAXWS + 1];

				strcpy (ca_w, b1 -> pair.f);
				correctEscapeSymbols (ca_w);
				sprintf (ca_buf, "\"%s\"", ca_w);
			} else {
				sprintf (ca_buf, "%s ", b1 -> pair.f);
			}
			for (i_len = strlen (ca_buf), cp = ca_buf; l + i_len >= l_limit; cp += l_limit - l, i_len -= (l_limit - l), l = 0) {
				char ca_w [MAXWS + 1];

				strncpy (ca_w, cp, l_limit - l);
				ca_w [l_limit - l] = 0;
				fprintf (fp, "%s\n", ca_w);
			}
			if (i_len > 0) {
				fprintf (fp, "%s", cp); l += i_len;
			}
			break;

		case LINK_TYPE_NUMBER:
			if (state_char == 1) {
				fputc ('\'', fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				state_char = 0;
			}
			sprintf (ca_buf, "%u ", b1 -> pair.n);
			for (cp = ca_buf, i_len = strlen (cp); l + i_len >= l_limit; cp += l_limit - l, i_len -= (l_limit - l), l = 0) {
				char ca_w [MAXWS + 1];

				strncpy (ca_w, ca_buf, l_limit - l);
				ca_w [l_limit - l] = 0;
				fprintf (fp, "%s\n", ca_w);
			}
			if (i_len > 0) {
				fprintf (fp, "%s", cp); l += i_len;
			}
			break;

		case LINK_TYPE_LSTRUCTB:
			if (state_char == 1) {
				fputc ('\'', fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				state_char = 0;
			}
			fputc ('(', fp); l ++;
			break;

		case LINK_TYPE_RSTRUCTB:
			if (state_char == 1) {
				fputc ('\'', fp); l ++;
				if (l >= l_limit) {
					fputc ('\n', fp); l = 0;
				}
				state_char = 0;
			}
			fputc (')', fp); l ++;
			break;

		case LINK_TYPE_LACT:
			if (state_char == 0) {
				state_char = 1;
				fputc ('\'', fp); l ++;
			}
			if (l >= l_limit) {
				fputc ('\n', fp); l = 0;
			}
			fputc ('<', fp); l ++;
			break;

		case LINK_TYPE_RACT:
			if (state_char == 0) {
				state_char = 1;
				fputc ('\'', fp); l ++;
			}
			if (l >= l_limit) {
				fputc ('\n', fp); l = 0;
			}
			fputc ('>', fp); l ++;
			break;

		default:
			fprintf (stderr, "Uncorrect type of refal-data for function \'Desysfun\'\n");
			ri_error (8);
			return 1;
		}
	}
	return 0;
}

int getXML (FILE *fp) {
	int brcnt = 0;
	long k;
	unsigned long l = (unsigned long)MAX_UNSIGNED_INT;
	int c;
	char str[MAXWS+1], *pc, *ri_cs_impl();
	int state_chars = 0;

	rdy (0);

	/* read file */
	while ((c = getc(fp)) != EOF) {
		switch(c) {
		case '\n': 
			break;

		case '(':
			if (state_chars) {
				ns ('(');
			} else {
				bl; 
				brcnt ++;
			}
			break;

		case ')':
			if (state_chars) {
				ns (')');
			} else {
				br;
				brcnt --;
				if (brcnt < 0) {
					fclose (fp);
					ri_error (8);
				}
			}
			break;

		case '<':
			/* May be it's data of PCDATA: '</>' */
			if (state_chars == 0) {
				if ((c = getc (fp)) == EOF) {
					ri_error (8);
					return 1;
				}
				if (c != '/') {
					ri_error (8);
					return 1;
				}
				if ((c = getc (fp)) == EOF) {
					ri_error (8);
					return 1;
				}
				if (c != '>') {
					ri_error (8);
					return 1;
				}
				state_chars = 1;
			} else {
			endPCDATA:
				// read second symbol
				if ((c = getc (fp)) == EOF) {
					ns ('<'); 
					fclose (fp);
					if (brcnt) {
					 	fclose (fp);
					 	ri_error(8);
					}
					out(2);
					est;
					return 0;
				}
				if (c != '/' && c != '<') {
					ns ('<'); ns (c);
					break;
				} else if (c == '<') {
					ns ('<');
					goto endPCDATA;
				}
				// read thrid symbol
				if ((c = getc (fp)) == EOF) {
					ns ('<'); ns ('/');
					fclose (fp);
					if (brcnt) {
					 	fclose (fp);
					 	ri_error(8);
					}
					out(2);
					est;
					return 0;
				}
				if (c != '>' && c != '<') {
					ns ('<'); ns ('/'); ns (c);
					break;
				} else if (c == '<') {
					ns ('<'); ns ('/');
					goto endPCDATA;
				}
				// It's end of PCDATA
				state_chars = 0;
			}
			break;

		case ' ':
			if (state_chars) {
				ns (' ');
			}
			break;

		case '"':
			if (state_chars) {
				/*
				fclose (fp);
				ri_error (8);
				*/
				ns ('"');
			} else {
				int i;

				for (i = 0, c = getc (fp); EOF != c && c != '"'; c = getc (fp), i ++) {
					if (c == ESCAPE) {
						while ((c = getc(fp)) == '\n');
						if (c == '"' || c == ESCAPE || c == '\'' || c == '(' || c == ')' ||
							c == '<' || c == '>' || c == 'n' || c == 't' || c == 'r') {
							switch (c) {
							case 'n': str [i] = '\n'; break;
							case 'r': str [i] = '\r'; break;
							case 't': str [i] = '\t'; break;
							default: str [i] = c;
							}
						} else {
							fclose (fp);
							ri_error (8);
						}
					} else if (c == '\n') {
						i --;
						continue;
					} else {
						str [i] = c;
					}
				}
				if (c != '"') {
					fclose (fp);
					ri_error (8);
				}
				str [i] = 0;
				pc = ri_cs_impl (str);
				ncs (pc);
			}
			break;

		default:
			if (state_chars) {
				ns (c);
			} else if (isalpha (c)) {
				int i;
				/* Get composymbol (string) */

				for (str [0] = c, i = 1; (c = fgetc (fp)) != EOF && i < MAXWS; i ++) {
					if (c == ' ') break;
					if (isalnum (c) || c == '_' || c == '$' || c == '-') str [i] = c;
					else if (c == '\n') i --;
					else {
						/* Error */
						fclose (fp);
						ri_error (8);
					}
				}
				/* I hope 1kb of buffer for saving ID is enough. */
				str [i] = 0;
				pc = ri_cs_impl (str);
				ncs (pc);
			} else if (isdigit (c)) {
				/* get number */
				k = c - '0';
				while (EOF != (c = getc (fp))) {
					if (isdigit (c)) {
/* Nemytykh 09.08.2002 */
        	/* The next few lines are strange because otherwise the "C" compiler
	        	produces an integer overflow warning */
		                                if ( (unsigned long)k > ((l - (c - '0'))/10 + (l - (c - '0'))%10) ) { 
                		 			  fclose (fp);
		                 		          ri_error (8);
				                } else  { k = k*10 + (c - '0'); };
/* Was:

						k = k * 10 + (c - '0');
						if (k >= MAXINT || k < 0) {
							nns (k);
							k = 0;
						}
*/
					} else if (c != '\n') break;
				}
				if (c != ' ') {
					fclose (fp);
					ri_error (8);
				}
				nns (k);
			} else  {
				ns (c);
			}
		}
	}
	fclose (fp);
	if (brcnt) {
	 	fclose (fp);
	 	ri_error(8);
	}
	out(2);
	est;
	return 0;
}

#include "bif_lex.h"
extern struct bitab bi [];
extern long nbi;

/* Get internal numbers and names of builtin functions: (s.number s.name s.type)* */
int rf_listOfBuiltin (void) {
	int i;
	char *pc;

        check_freeze;
	rdy (0);
      
        for (i = 1; i < nbi; i ++) {
                bl;
        	nns (bi[i].fnumber);
        	pc = ri_cs_impl (bi[i].fname);
        	ncs (pc);
                if( (bi[i].flags) == BI_FADDR ) {
                	pc = ri_cs_impl ("special");
                }
                else {
                	pc = ri_cs_impl ("regular");
                };
        	ncs (pc);
                br;
        }
	out (2);
	est;
	return 0;
}

int rf_sizeof (void) {
        int size;

        check_freeze;

	b1 = b1->foll;
	if (!LINK_CHAR (b1)) ri_error (8);
        switch (b1 -> pair.c) {
        	case 'c':
			size = sizeof(char);
			break;
        	case 's':
			size = sizeof(short);
			break;
        	case 'i':
			size = sizeof(int);
			break;
        	case 'l':
			size = sizeof(long);
			break;
        	case 'p':
			size = sizeof(char *);
			break;
		default:
			ri_error(8);
			return 1;
	};
	rdy (0);
        nns (size);

	out (2);
	est;
	return 0;
}

