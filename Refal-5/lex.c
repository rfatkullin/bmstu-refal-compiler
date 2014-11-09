
	/* include statements.	*/

# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"
# include "arithm.h"

# define  RESIDUE  50

# define ESCAPE '\\'
/*# define ESCAPE '#'*/

extern struct bitab bi [];
extern long nbi;

/* local: */
int rc_getid1 (int);

/* return the number of the builtin function or 0 */
long rc_binumber (s)
     char *s;
{
  int i;

  for (i = 0; i < nbi; i++) 
    if (strcmp (s, bi [i].fname) == 0) return (long)(bi [i].fnumber);
  return 0L;
}


/*  function to get the next token.	*/
int
rc_gettoken (void) {
  /*   legal tokens are   :                                       */
  /*      VAR     -       variable,                               */
  /*      ID      -       identifier,                             */
  /*      COMPSYM -       compound symbol,                        */
  /*      MDIGIT  -       macrodigit,                             */
  /*      ASCIIV  -       symbol represented by its ASCII value,  */
  /*      LCBRAK  -       left concretization bracket,            */
  /*      STRING  -       string of printable characters.         */
  /*      EXTRN   -       ID of the form "EXTRN" or "ENTRY"       */
  /*   In addition to these there is a number of other characters */
  /*   which could be returned (list separated by blanks):        */
  /*   = ( ) > , : ; { } [ ]                                      */

	int c,d;
	/*char ch;*/

	token = 0;
	c = (globsav == 0)? rc_gchar(): globsav;
	globsav = 0;
	while (token == 0) {
		while ( (c == ' ') || (c == '\t')) c = rc_gchar ();
		switch (c) {
		case EOF : return c;

		case '\n': 
			/* skip comment. */
			if ((c = rc_gchar ()) == '*') 
				while ((c = rc_gchar ()) != '\n')
					if (c == EOF) return EOF;
			break;

		case '/' : 
			if ((c=rc_gchar ()) != '*') {
				rc_serror (1,NULL);
			} else {
				c = rc_gchar ();
				do {
					while (c != '*') {
						if (c == EOF) {
							rc_serror (2,NULL);
							return EOF;
						}
						c = rc_gchar ();
					}
				} while ((c=rc_gchar ()) != '/');
				c = rc_gchar ();
			}
			break;

		case '.' : /* MDIGIT or COMPSYM    */
			token = rc_getsym ();
			c = ' ';
			break;

		case 'e' :
		case 's' :
		case 't' :
		/*case 'w' :*/
			d = rc_gchar ();
			if (d == '.') {
				/* variable.  */
				vtype = c;
				token = rc_getind ();
			} else if (d != EOF) {
				int i;
/* Nemytykh [ .... */
                         	if (isalnum (d) || (d == '_') || (d == '-')) {
                                   token = rc_getid (d);
                                   for (i = length; i > 0; i --) {
					str [i] = str [i - 1];
                        	   }
                                } else {
                                        token = rc_getid1 (c);
                                       	globsav = d;
                                };
/* Nemytykh .... ]  */
				str [0] = c;
			} else {
				rc_serror (20, NULL);
			}
			c = ' ';
			break;

		case ESCAPE:
			c = rc_gchar ();
			if (c == 'x') {
				int i_w;
				char c_w = 0;

				for (i_w = 0; i_w < 2; i_w ++) {
					if (EOF == (c = rc_gchar ())) {
						rc_serror (7, NULL);
					}
					if (isxdigit (c)) {
						c = (c <= '9')?	c - '0': ((isupper (c))?	c - 'A' + 10: c - 'a' + 10);
						c_w += ((i_w)? c: 16 * c);
					} else {
						rc_serror (7, NULL);
					}
				}
				v = c_w;
				token = ASCIIV;
				c = ' ';
			} else if (c == 'n' || c == 'r' || c == 't' || c == '\'' || c == '"' ||
				c == '<' || c == '>' || c == '(' || c == ')' || c == ESCAPE) {

				if (c == 'n') v = '\n';
				else if (c == 'r') v = '\r';
				else if (c == 't') v = '\t';
				else v = c;
				token = ASCIIV;
				c = ' ';
			} else {
				rc_serror (7, NULL);
				token = 0;
			}
			break;

		case '\'':
		case '"':
			/* string */
			token = rc_getstr (c);
			c = ' ';
			break;

		case '<' : /*  left concretization bracket  */
			c = rc_gchar ();
			token = LCBRAK;
			switch (c) {
			case '+':
				strcpy (str, bi [ADD].fname);
				break;
			case '-':
				strcpy (str, bi [SUB].fname);
				break;
			case '*':
				strcpy (str, bi [MUL].fname);
				break;
			case '/':
				strcpy (str, bi [DIV].fname);
				break;
			case '%':
				strcpy (str, bi [MOD].fname);
				break;
			case '?':
				strcpy(str, "Residue");
				break;
			case '.':
				c = rc_gchar ();
			default:
				if (isalpha (c)) rc_getid (c);
				else {
					/* ch = c; */
					rc_serror (4, (char *) & c);
					c = ' ';
				}
				break;
			}
			break;
      
		case '(' :
		case ')' :
		case ',' :
		case ':' :
		case ';' :
		case '{' :
		case '}' :
		case '>' :
		case '=' :
		case '[' : /* For parametrs. Insert by Shura. 12.02.98 */
		case ']' : /* For parametrs. Insert by Shura. 12.02.98 */
			/* special symbol. */
			token = c;
			break;

		case '$' :
			token = 0;
			if (isupper (c = rc_gchar ())) {
				rc_getid (c);
				if (strcmp (str,"EXTRN") == 0)  token = EXTRN;
				else if (strcmp (str,"EXTERN") == 0) token = EXTRN;
				else if (strcmp (str,"EXTERNAL") == 0) token = EXTRN;
				else if (strcmp (str,"ENTRY") == 0) token = ENTRY;
				else { 
					c = '\n';
					rc_serror (14,NULL);
				}
			} else {
				c = '\n';
				rc_serror (14,NULL);
				token = 0;
			}
			break;

		default : /* c is a capital letter or error */
			token = ID;
			if (isalpha (c) || c == '_') token = rc_getid (c);
			else if (isdigit (c)) token = rc_getnumb (c);
			else {
				/* ch = c; */
				rc_serror (4, (char *) & c); 
				c = rc_gchar ();
			}
			break;
		}
	}
	return token;
}

/* get a string  of digits or an id  terminated by a '.' and 
 * save it in the global array str [] and return MDIGIT or COMPSYM 
 */
int rc_getsym (void) {
	int  atype,i,c,warning;

	i = 1;
	warning = 0;
	c = rc_gchar ();
	if (isdigit (c)) atype = MDIGIT;
	else if (isalpha (c)) atype = COMPSYM;
	else {
		rc_serror (5,NULL);
		return 0;
	}

	if (atype == MDIGIT) str [0] = c;
	else strings [0] = c;

	while ( (c = rc_gchar ()) != '.') {
		if (i >= MAXWS) warning = --i;
		if (isdigit (c)) {
			if (atype == MDIGIT) str [i++] = c;
			else strings [i ++] = c;
		} else if ( (isalpha (c)) && (atype == COMPSYM)) {
			strings [i ++] = c;
		} else {
			rc_serror (5,NULL);
			return 0;
		}
	}
	if (atype == MDIGIT) str [i] = 0;
	else strings [i] = 0;
	length = i;

	if (warning != 0) rc_swarn (1);
	return atype;
}

/* get an index of a variable terminated by any non-digit or non-letter, and save it in str []. */
int rc_getind (void) {
	int i,c,warning;

	i = 1;
	warning = 0;
	c = rc_gchar ();
	if (! (isalnum (c))) {
		rc_serror (3,NULL);
		return 0;
	}
	str [0] = c;
	c = rc_gchar ();
	/* legal symbols are letters, digits, signs _ or -.  07-25-1985.DT.*/
	while (isalnum (c) || (c == '_') || (c == '-')) {
		if (i >= MAXWS) warning = --i;
		str [i++] = c;
		c = rc_gchar ();
	}
	globsav = c; /* c == '.' ? ' ' : c;*/
	str [i] = '\0';
	length = i+1;
	if (warning != 0) rc_swarn (2);
	return VAR;
}


/* get id with the first letter d. */
int rc_getid (int d) {
	int i,c,warning;

	i = 1;
	warning = 0;
	str [0] = d;
	c = rc_gchar ();
	/* legal symbols are letters, digits, signs _ or -. 07-25-1985. DT. */
	while (isalnum (c) || (c == '_') || (c == '-')) {
		if (i >= MAXWS) warning = --i;
		str [i++] = c;
		c = rc_gchar ();
	}
	globsav = c;
	str [i] = '\0';
	length = i+1; /* ??? */
	if (warning != 0) rc_swarn (1);
	return ID;
}

/* Nemytykh: get id with the only letter d. */
int rc_getid1 (int d) {
	str [0] = d;
	str [1] = '\0';
	length = 2; /* ??? */
	return ID;
}

/* Insert by Shura. 27.05.99 */
#define MAX_STR_ERROR 200

/* print syntax error message. */
int
rc_serror (int code, char * str) {
  char ca_err [MAX_STR_ERROR]; /* Insert by Shura. 27.05.99*/

  globsav = 0;
  if (nerrors++ > 100) {
    fprintf (fdlis,"\nToo many errors. Aborted.\n");
    /* The line is inserted by Shura. 27.05.99 */
    fprintf (stderr, "\nToo many errors. Aborted.\n");
    /* End. Shura. 27.05.99 */
    exit (1);
  };
  /* Was. Shura. 27.05.99
   *  fprintf (fdlis,"Error: ");
   */
  strcpy (ca_err, "Error: "); /* Insert by Shura. 27.05.99 */
  switch (code)	{
  case  1: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 1. Illegal comment ");
     */
    strcat (ca_err, " 1. Illegal comment ");
    break;
  case  2: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 2. Unexpected EOF encountered ");
     */
    strcat (ca_err, " 2. Unexpected EOF encountered ");
    break;
  case  3: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 3. Illegal index of a variable ");
     */
    strcat (ca_err, " 3. Illegal index of a variable ");
    break;
  case  4: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 4. Illegal character: %c (%d)",*str,*str);
     */
    sprintf (ca_err + strlen (ca_err), " 4. Illegal character: %c (%d)",
	     * str, * str);
    break;
  case  5: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 5. Illegal compound symbol ");
     */
    strcat (ca_err, " 5. Illegal compound symbol ");
    break;
  case  6: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 6. No closing quote ");
     */
    strcat (ca_err, " 6. No closing quote ");
    break;
  case  7: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 7. Illegal number ");
     */
    strcat (ca_err, " 7. Illegal number ");
    break;
  case  8: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 8. Doubly defined function ");
     */
    strcat (ca_err, " 8. Doubly defined function ");
    break;
  case  9: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis," 9.");
     */
    strcat (ca_err, " 9.");
    break;
  case 10: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"10. Undefined variable in the right side: %s ", str);
     */
    sprintf (ca_err + strlen (ca_err),
	     "10. Undefined variable in the right side: %s ", str);
    break;
  case 11: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"11. More than %d chars (excess ignored) ", MAXSTR);
     */
    sprintf (ca_err + strlen (ca_err),
	     "11. More than %d chars (excess ignored) ", MAXSTR);
    break;
  case 12: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"12. Unbalanced brackets in expression ");
     */
    strcat (ca_err, "12. Unbalanced brackets in expression ");
    break;
  case 13: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"13.");
     */
    strcat (ca_err, "13.");
    break;
  case 14: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"14. Error in EXTERNAL function declaration ");
     */
    strcat (ca_err, "14. Error in EXTERNAL function declaration ");
    break;
  case 15: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"15.");
     */
    strcat (ca_err, "15.");
    break;
  case 16: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"16. Error in a function (skipping till ';') ");
     */
    strcat (ca_err, "16. Error in a function (skipping till ';') ");
    break;
  case 17: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"17. Error in a sentence (skipping to ';') ");
     */
    strcat (ca_err, "17. Error in a sentence (skipping to ';') ");
    break;
  case 18: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"18. Parsing error ");
     */
    strcat (ca_err, "18. Parsing error ");
    break;
  case 19: 
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"19. ';' or EOF expected ");
     */
    strcat (ca_err, "19. ';' or EOF expected ");
    break;

  case 20:
	  strcat (ca_err, "20. Missed dot \'.\' before index of variable ");
	  break;

  case 100:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"100. Unbalanced left structure bracket ");
     */
    strcat (ca_err, "100. Unbalanced left structure bracket ");
    break;
  case 101:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"101. Expected '{' ");
     */
    strcat (ca_err, "101. Expected '{' ");
    break;
  case 102:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"102. Expected '}' ");
     */
    strcat (ca_err, "102. Expected '}' ");
		break;
  case 103:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"103. Expected ';' ");
     */
    strcat (ca_err, "103. Expected ';' ");
    break;
  case 104:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"104. Expected identifier ");
     */
    strcat (ca_err, "104. Expected identifier ");
    break;
  case 105:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"105. Missing ')' (added) ");
     */
    strcat (ca_err, "105. Missing ')' (added) ");
    break;
  case 106:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"106. Expected ',' or '=' ");
     */
    strcat (ca_err, "106. Expected ',' or '=' ");
    break;
  case 107:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"107. Expected ':' ");
     */
    strcat (ca_err, "107. Expected ':' ");
    break;
  case 108:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"108. Missing '>' (added) ");
     */
    strcat (ca_err, "108. Missing '>' (added) ");
    break;
  case 109:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"109. Expression too big ");
     */
    strcat (ca_err, "109. Expression too big ");
    break;
  case 110:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"110. Function %s in the left side of a sentence ", str);
     */
    sprintf (ca_err + strlen (ca_err),
	     "110. Function %s in the left side of a sentence ", str);
    break;
  case 111:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"111. Unbalanced < and > ");
     */
    strcat (ca_err, "111. Unbalanced < and > ");
    break;
  case 112:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"112. Unexpected > in the left side of a sentence ");
     */
    strcat (ca_err, "112. Unexpected > in the left side of a sentence ");
    break;
  case 200:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"200. Undefined function %s ",str);
     */
    sprintf (ca_err + strlen (ca_err), "200. Undefined function %s ",str);
    break;
  case 201:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"201. Doubly defined function %s ", str);
     */
    sprintf (ca_err + strlen (ca_err), "201. Doubly defined function %s . May be the function is built.",
	     str);
    break;
  case 202:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"202. No entry point.");
     */
    strcat (ca_err, "202. No entry point.");
    break;
  case 204:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"204. Built-in function %s defined as ENTRY ", str);
     */
    sprintf (ca_err + strlen (ca_err),
	     "204. Built-in function %s defined as ENTRY ", str);
    break;
  case 205:
/* Nemytykh 09.08.2002 */
    sprintf (ca_err + strlen (ca_err),
	     "205. Number greater than 2^%d-1=%lu ", sizeof(int)*8,MAX_UNSIGNED_INT);
    break;
  default:
    /* Was. Shura. 27.05.99
     *    fprintf (fdlis,"System error: %4d ",code);
     */
    sprintf (ca_err + strlen (ca_err), "System error: %4d ",code);
  }

  /* Was. Shura. 27.05.99
   *  if (line_no > 0) fprintf (fdlis," on line %5d.\n", line_no);
   *  else fprintf (fdlis, "\n");
   */
  if (line_no > 0) {
    sprintf (ca_err + strlen (ca_err), " on line %5d.", line_no);
  }
  fprintf (fdlis, "%s\n", ca_err);
  fprintf (stderr, "%s\n", ca_err); /* May be to using some flag to control */
  return 0;
}

	/* print warning  */
int rc_swarn (code)
	int  code;

	{
	fprintf (fdlis,"warning: ");
	switch (code)
		{
		case  1: 
			fprintf (fdlis," 1. Identifier too long (excess ignored) ");
			break;
		case  2: 
			fprintf (fdlis," 2. Index too long (excess ignored) ");
			break;
		case  3: 
			fprintf (fdlis," 3. ASCII value over 256 ");
			break;
		case  4: 
			fprintf (fdlis," 4. Empty string (ignored) ");
			break;
		case  5: 
/* Nemytykh 08.08.2002		fprintf (fdlis," 5. Number greater than %ld (shorten)",MAXINT);*/
			fprintf (fdlis," 5. Number greater than 2^%d-1=%lu (shorten)", sizeof(int)*8,MAX_UNSIGNED_INT);
			break;
		case  6: 
			fprintf (fdlis," 6. The same index for 2 diff variables ");
			break;
		default: fprintf (fdlis,"System warning: %4d ",code);
		break;
	};
	if (line_no > 0) fprintf (fdlis," on line %5d.\n", line_no);
	else fprintf (fdlis, "\n");
	return 0;
	}
        
	/* get string delimited by character quote */
int rc_getstr (int quote) {
	int c,i;

	c = rc_gchar ();
	if (c == quote) {
		length = 0;
		strings [0] = '\0';
		return ((quote == '"')? COMPSYM: 0);
	}
	i = 0;
	if (c == ESCAPE) {
		if (EOF == (c = rc_gchar())) {
			rc_serror (6, NULL);
			length = i;
			strings [i] = 0;
			return 0;
		}
		if (c == '\n') {
			rc_serror (6, NULL);
			strings [i] = 0;
			return 0;
		}
		if (c == '\'' || c == '"' || c == '(' || c == ')' || c == '<' || c == '>' || c == ESCAPE) {
			strings [i ++] = c;
		} else if (c == 'n') {
			strings [i ++] = '\n';
		} else if (c == 'r') {
			strings [i ++] = '\r';
		} else if (c == 't') {
			strings [i ++] = '\t';
		} else if (c == 'x') {
			int i_w;
			char c_w = 0;

			for (i_w = 0; i_w < 2; i_w ++) {
				if (EOF == (c = rc_gchar ())) {
					rc_serror (6, NULL);
					return 0;
				}
				if (isxdigit (c)) {
					c = (c <= '9')?	c - '0': ((isupper (c))? c - 'A' + 10: c - 'a' + 10);
/* Nemytykh 21.06.2001			c_w += ((i)? c: 16 * c);*/
					c_w += ((i_w)? c: 16 * c);
				} else {
					rc_serror (6, NULL);
					return 0;
				}
			}
			strings [i ++] = c_w;

		} else {
			rc_serror (6, NULL);
		}
	} else { /* end of if ...ESCAPE... */
		strings [i ++] = c;
	}

	while ( ( (c = rc_gchar ()) != EOF) && (c != '\n')) {
		if (c == ESCAPE) {
			if (EOF == (c = rc_gchar())) {
				rc_serror (6, NULL);
				length = i;
				strings [i] = 0;
				return 0;
			}
			if (c == '\n') {
				rc_serror (6, NULL);
				strings [i] = 0;
				return 0;
			}
			if (c == '\'' || c == '"' || c == '(' || c == ')' || c == '<' || c == '>' || c == ESCAPE) {
				strings [i ++] = c;
			} else if (c == 'n') {
				strings [i ++] = '\n';
			} else if (c == 'r') {
				strings [i ++] = '\r';
			} else if (c == 't') {
				strings [i ++] = '\t';
			} else if (c == 'x') {
				int i_w;
				char c_w = 0;

				for (i_w = 0; i_w < 2; i_w ++) {
					if (EOF == (c = rc_gchar ())) {
						rc_serror (6, NULL);
						return 0;
					}
					if (isxdigit (c)) {
						c = (c <= '9')?	c - '0': ((isupper (c))? c - 'A' + 10: c - 'a' + 10);
						c_w += ((i_w)? c: 16 * c);
					} else {
						rc_serror (6, NULL);
						return 0;
					}
				}
				strings [i ++] = c_w;
			} else {
				rc_serror (6, NULL);
			}
		} else if (c == quote) { /* end of if ...ESCAPE... */
			strings [i] = 0;
			length = i;
			return ((quote == '\'')? STRING: COMPSYM);
		} else {
			strings [i ++] = c;
		}
	} /* end of while */
	rc_serror (6,NULL);
	length = i;
	strings [i] = '\0';
	return 0;
}

	/* get number with the first digit d. */
int rc_getnumb (d)
	int d;
	{
	int c, warning;
	unsigned long i;
	unsigned long k = (unsigned long)MAX_UNSIGNED_INT;

	i = d - '0';
	warning = 0;
	c = rc_gchar ();
	while (isdigit (c))
		{
	/* The next few lines are strange because otherwise the "C" compiler
		produces an integer overflow warning */
		if ( warning || (i > ((k - (c - '0'))/10 + (k - (c - '0'))%10)) ) 
                        { warning = 1; }
                else    { i = i*10 + (c - '0'); };
		c = rc_gchar ();
		};
	globsav = c;
	v = i;
/* Nemytykh 09.08.2002 */
	if (warning) rc_serror (205, NULL);  /* Was: rc_swarn (5); */
	return NUMBER;
	}


