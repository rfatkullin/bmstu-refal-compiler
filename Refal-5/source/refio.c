
# include "decl.h"
# include "macros.h"
# include "fileio.h"
# include "freeze.h"
# include "ifunc.h"


/*# define ESCAPE '#'*/
#define ESCAPE '\\'

	/* number of characters on screen. */
# define LINELEN  79


static int padding;
static char string [LINELEN+1];
static int line_length;
static int state_chars = 0;

/* PUTMB puts the expression between the pointers Q1 and Q2 onto file *fd.	*/
/* It returns the number of printed links. */
int ri_putmb (LINK * q1, LINK * q2, FILE * fd) {
	int i;
	
	i = 0;
	line_length = 0;
	padding = 0;
	/* empty expression. */
	if (q1 == NULL) {
		fputc('\n',fd);
		return 0;
	}
	
	if ((q2 == NULL) || (q1 == q2)) { /* a symbol.  */
		ri_addsym(q1,fd);
		putstr (fd, 0, NULL);
		return 1;
	}
	while (q1 != q2) {
		ri_addsym(q1,fd);
		q1 = q1 -> foll;
		i ++;
		if (LINK_LSTRUCTB (q1) && q1 -> prec == NULL && q1 -> pair.b == NULL) {
			fprintf(fd,"\nEnter free memory.\n");
			return -1;
		}
	}

	ri_addsym (q2,fd);
	if (state_chars) {
		/* Need to put '*/
		ri_add_char ('\'', fd);
		state_chars = 0;
	}
	i++;
	while (line_length) putstr (fd, 0, NULL);
	fputc ('\n', fd);
	return i;
}


int
isIdentifier (char * cp_str) {
	int i, i_len = strlen (cp_str);

	/* Identifier must begin with a letter */
	if (! isalpha (cp_str [0])) return 0;

	for (i = 0; i < i_len; i ++) {
		if (isalnum (cp_str [i]) || cp_str [i] == '_' || cp_str [i] == '-') continue;
		return 0;
	}
	return 1;
}
void
correctEscapeSymbols (char * cp) {
	char ca_w [MAXWS + 7];
	int i, i_len = strlen (cp);

	for (i = 0; i < i_len; i ++) {
		if (cp [i] == ESCAPE || cp [i] == '\'' || cp [i] == '"' || 	cp [i] == '\n' || cp [i] == '\r' || cp [i] == '\t' ||
			cp [i] == '(' || cp [i] == ')' || cp [i] == '<' || cp [i] == '>') {
			strncpy (ca_w, cp, i);
			ca_w [i] = ESCAPE;
			ca_w [i + 2] = '\0';
			if (cp [i] == '\n') {
				ca_w [i + 1] = 'n';
			} else if (cp [i] == '\r') {
				ca_w [i + 1] = 'r';
			} else if (cp [i] == '\t') {
				ca_w [i + 1] = 't';
			} else if (cp [i] == ESCAPE || cp [i] == '\'' || cp [i] == '"' || 
						cp [i] == '(' || cp [i] == ')' || cp [i] == '<' || cp [i] == '>') {
				ca_w [i + 1] = cp [i];
			}
			strcpy (cp, strcat (ca_w, cp + i + 1));
			i_len = strlen (cp);
			i ++;
		} else if (cp [i] < ' ' || cp [i] > 127) {
                        char zero [2];

			strncpy (ca_w, cp, i);
			ca_w [i] = '\0';

/*+ The padding for heximal represantation does not work. 
			sprintf (ca_w + i, "%cx%0x%s", ESCAPE, cp [i], cp + i + 1);
*/
			zero[0] = (cp [i] < 16) ? '0':'\0';
                        zero[1] = '\0';
			sprintf (ca_w + i, "%cx%s%x%s", ESCAPE, zero, cp [i], cp + i + 1);

			strcpy (cp, ca_w);
			i_len = strlen (cp);
			i += 3;
		}
	}
}
int ri_addsym (LINK * ptr, FILE *fd) {
	char tbuf [MAXWS + 3], esc_buf [3];
	char c, var_type;
/*	int l;*/
	unsigned long variable, index, level, elev;

	if (ptr == NULL) {
		return 0;
	}
	esc_buf [0] = ESCAPE;
	esc_buf [2] = '\0';
	switch (ptr -> ptype) {

	case LINK_TYPE_LSTRUCTB: /*** Left par. ***/
		if (state_chars) {
			ri_add_char ('\'', fd);
			ri_add_char (' ', fd);
			state_chars = 0;
		}
		ri_add_char ('(',fd);
		break; 

	case LINK_TYPE_RSTRUCTB: /*** right par. ***/
		if (state_chars) {
			ri_add_char ('\'', fd);
			ri_add_char (' ', fd);
			state_chars = 0;
		}
		ri_add_char (')',fd);
		break; 

	case LINK_TYPE_CHAR: /*** character. ***/
		c = ptr -> pair.c;
		if (state_chars == 0) {
			ri_add_char ('\'', fd);
			state_chars = 1;
		}

		switch (c) {
		case ESCAPE: case '(': case ')': case '<': case '>':
		case '\'': case '"': case '\r':	case '\n': case '\t':
			esc_buf [1] = (c == '\n')? 'n': ((c == '\r')? 'r': ((c == '\t')? 't': c));
			ri_add_string (1, esc_buf, 2, fd);
			break;
				
		default:
			if (c >= ' ' && c <= 127) {
				ri_add_char (c,fd);
			} else {
				char ca_w [10];
                                char zero [2];
				
/*+ The padding for heximal represantation does not work. 
				sprintf (ca_w, "%cx%0x", ESCAPE, (char) c);*/
				zero[0] = (c < 16) ? '0':'\0';
                                zero[1] = '\0';
                                sprintf (ca_w, "%cx%s%x", ESCAPE,zero,(char) c);
				ri_add_string (1, ca_w, strlen (ca_w), fd);
			}
			break;
		}
		break;
		
	case LINK_TYPE_COMPSYM: /*** compound symbol. ***/
		if (isIdentifier (ptr -> pair.f)) {
			sprintf (tbuf, "%s ", ptr -> pair.f);
		} else {
			char ca_w [MAXWS + 3];

			strcpy (ca_w, ptr -> pair.f);
			correctEscapeSymbols (ca_w);
			sprintf (tbuf, "\"%s\" ", ca_w);
		}
		if (state_chars) {
			ri_add_char ('\'', fd);
			ri_add_char (' ', fd);
			state_chars = 0;
		}
		/*l = strlen (tbuf);*/
		ri_add_string (1, tbuf, strlen (tbuf), fd);
		break;

	case LINK_TYPE_NUMBER: /*** macrodigit. ***/
		sprintf (tbuf, "%lu ", ptr -> pair.n);
		if (state_chars) {
			ri_add_char ('\'', fd);
			ri_add_char (' ', fd);
			state_chars = 0;
		}
		/*l = strlen (tbuf);*/
		ri_add_string (1, tbuf, strlen (tbuf), fd);
		break;

	case LINK_TYPE_LACT: /*** Left Active Bracket.  ***/
		if (state_chars) {
			ri_add_char ('\'', fd);
			ri_add_char (' ', fd);
			state_chars = 0;
		}
		{
			char * cp_w = ptr -> pair.f - 1;
			for (cp_w --; * cp_w != '\0'; cp_w --);
			cp_w ++;
			sprintf (tbuf, "<%s ", cp_w);
		}
		/*l = strlen (tbuf);*/
		ri_add_string (2, tbuf, strlen (tbuf), fd);
		break;

	case LINK_TYPE_RACT: /*** Right active Bracket. ***/
		if (state_chars) {
			ri_add_char ('\'', fd);
			state_chars = 0;
		}
		ri_add_char ('>',fd);
		break;

	case LINK_TYPE_SVAR: /*** S-variable. **/
	case LINK_TYPE_EVAR: /*** E-variable. **/
	case LINK_TYPE_TVAR: /*** T-variable. **/
		if (state_chars) {
			ri_add_char ('\'', fd);
			ri_add_char (' ', fd);
			state_chars = 0;
		}
		if (LINK_SVAR(ptr)) var_type = 'S';
		else if (LINK_EVAR(ptr)) var_type = 'E';
		else var_type = 'T';

		variable = ptr -> pair.n;
		index = index_of(variable);
		level = level_of(variable);
		elev = elevation_of(variable);
		if (elev == MAX_VAR_ELEV) sprintf(tbuf, "%c%c.%lu%c%lu ",
			ESCAPE, var_type, level, ESCAPE, index);
		else sprintf(tbuf, "%c%c.%lu%c%lu%c%lu ",
			ESCAPE, var_type, level, ESCAPE, index, ESCAPE, elev);
		/*l = strlen (tbuf);*/
		ri_add_string (0, tbuf, strlen (tbuf), fd);
		break;

	default:			/**** Error ****/
		fprintf(fd, "Illegal pointer. %lx\n", (unsigned long) ptr);
		break;
	}
	return 0;
}
			
int ri_add_char (char c, FILE * fd) {
	if (line_length + padding >= LINELEN-1) putstr (fd, 0, NULL);
	string [line_length++] = c;
	return 0;
}

int ri_add_string (int shared, char * tbuf, int l, FILE * fd) {
	while (line_length + padding >= LINELEN - l) {
		if (1 == putstr (fd, shared, tbuf)) {
			return 0;
		}
	}
	strcpy (string+line_length, tbuf);
	line_length += l;																				   
	return 0;
}

static int even_escapes(int);

static int even_escapes (int i) {
	/*register*/ int j = 1;

	while (--i >= 0) 
		if (string [i] == ESCAPE) j = !j;
		else break;
	return  j;
}

int putstr (FILE * fd, int flag_shared, char * tbuf) {
	/*register*/ int i;
	int paren = 0;
	char c;

	/*** Search for the 1st unbalanced right paren ***/
	string [line_length] = '\0';
	for (i = 0; i < line_length; i ++) 
		switch (string [i]) {
		case ESCAPE:
			if (string [i + 1] == '<' || string [i + 1] == '>' || string [i + 1] == '(' || string [i + 1] == ')' ||
				string [i + 1] == '\'' || string [i + 1] == '"' || string [i + 1] == ESCAPE || 
				string [i+1] == 'n' || string [i+1] == 't' || string [i + 1] == 'r') {
				i ++;
			}
			break;

		case '(':
		case '<':
			paren ++;
			break;

		case ')':
		case '>':
			paren --;
			if (paren < 0) {
				c = string [i];
				padding --;
				rc_print_buff(i, fd);
				fputc (c, fd);
				strcpy (string, string+i+1);
				line_length -= i+1;

				/* Shura added */
				if (line_length == 0) {
					if (flag_shared) {
						if (strlen (tbuf) + padding >= LINELEN) {
							fputs (tbuf, fd);
							if (flag_shared == 2) {
								int i;

								fputc ('\n', fd);
								padding ++;
								for (i = 0; i < padding; i ++) fputc (' ', fd);
							}
						}
						return 1;
					}
				}
				return 0;
			} 
			break;

		default:
			break;
		}

	/* All parentheses are paired. */
	if (paren == 0)	{
		rc_print_buff(line_length, fd);
		line_length = 0;

		/* Shura added */
		if (flag_shared) {
			if (strlen (tbuf) + padding >= LINELEN) {
				fputs (tbuf, fd);
				if (flag_shared == 2) {
					int i;

					fputc ('\n', fd);
					padding ++;
					for (i = 0; i < padding; i ++) fputc (' ', fd);
				}
				return 1;
			}
			return 0;
		}
		return 0;
	}

	/*** There are unbalanced left parentheses. ***/
	/**** Find the first one. *****/
	for (i = line_length-1; i >= 0; i --)
		switch (string [i])	{
		case ')':
		case '>':
			if (even_escapes (i)) paren ++;
			else i --;
			break;
		case '(':
		case '<':
			if (even_escapes (i)) {
				paren --;
				if (paren == 0)	{
					c = string [i];
					rc_print_buff(i, fd);
					padding ++;
					fputc (c, fd);
					strcpy (string, string+i+1);
					line_length -= i+1;

					/* Shura added */
					if (line_length == 0) {
						if (flag_shared) {
							if (strlen (tbuf) + padding >= LINELEN) {
								fputs (tbuf, fd);
								if (flag_shared == 2) {
									int i;

									fputc ('\n', fd);
									padding ++;
									for (i = 0; i < padding; i ++) fputc (' ', fd);
								}
								return 1;
							}
							return 0;
						}
					}
					return 0;
				}
			} else i --;
			break;
		default:
			break;
		}
	return 0;
}

int rc_print_buff (i, fd)
	int i;
	FILE *fd;
	{
		string [i] = '\0';
		fputs (string, fd);
		fputc ('\n', fd);
		for (i = 0; i < padding; i ++) fputc(' ', fd);
		return 0;
	}

/* determines the length of the compound symbol. */
int ri_cs_len (char *pc) {
	/*register int i;*/
	return strlen (pc);
	/*
	i = strlen (pc);
	return min (i, MAXWS);
	*/
}

/* determines the length of the segment after the last new_line character of the compound symbol. */
int ri_cs_rest_len (char *pc) {
	register int i = 0;
        while( *pc ) {
           if( (*pc == '\n') || (*pc == '\r') ) i = 0;
           else i++;
           pc++; 
        }
	return i;
}

int ri_new_line (char ch) {
     if( (ch == '\n') || (ch == '\r') ) return 1;
     else return 0;
}

/* RI_PUT prints the refal expression onto file pointed by FD. */
/*	B1 has to point to the link preceding the expression. */
/*	B2 has to point to the link following the expression. */
int ri_put (FILE * fd) {
/*	int i = 0;*/
	int i_w/*, i_chars = 0*/;
	char buf [MAXWS]/*, c_w*/;
	
	for (b1 = b1 -> foll; b1 != b2; b1 = b1 -> foll) {
		switch (b1 -> ptype) {
		case LINK_TYPE_LSTRUCTB:  /*  left parenthesis   */
/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
				if (i ++ >= LINELEN) {
					putc('\n',fd);
					i = 0;
				}
			}
*/
			/*
			if (i_chars) {
				putc ('\'', fd);
				i_chars = 0;
				if (fd == stderr || fd == stdout) {
					if (i ++ > LINELEN) {
						putc ('\n', fd);
						i = 0;
					}
				}
			}
			*/
			putc('(',fd);
			break;

		case LINK_TYPE_RSTRUCTB:  /*  right parenthesis   */
/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
				if (i++ >= LINELEN) {
					putc('\n',fd);
					i = 0;
				}
			}
*/
			/*
			if (i_chars) {
				putc ('\'', fd);
				i_chars = 0;
				if (fd == stderr || fd == stdout) {
					if (i ++ > LINELEN) {
						putc ('\n', fd);
						i = 0;
					}
				}
			}
			*/
			putc(')',fd);
			break;

		case LINK_TYPE_CHAR:  /* object symbol */
			/*
			if (i_chars == 0) {
				putc ('\'', fd);
				i_chars = 1;
				if (fd == stderr || fd == stdout) {
					if (i ++ > LINELEN) {
						putc ('\n', fd);
						i = 0;
					}
				}
			}
			*/
/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
                 		if (i++ >= LINELEN) {  
					putc('\n',fd);
					i = 0;
				}
/+* Nemytykh 20.06.2001*+/ if (ri_new_line(b1 -> pair.c)) i = 0;
			}
*/
			/*
			c_w = b1 -> pair.c;
			if (c_w == '\\' || c_w == '\'' || c_w == '"' ||	c_w == '(' || c_w == ')' || 
				c_w == '<' || c_w == '>' || c_w == '\n' || c_w == '\r' || c_w == '\t') {
				putc ('\\', fd);
			}
			if (c_w == '\n') putc ('n', fd);
			else if (c_w == '\r') putc ('r', fd);
			else if (c_w == '\t') putc ('t', fd);
			else */
                        putc (b1 -> pair.c, fd);
			break;

		case LINK_TYPE_COMPSYM: /*  compound symbol    */
			/*
			if (i_chars) {
				i_chars = 0;
				putc ('\'', fd);
/* Nemytykh 18.02.2004
				if (fd == stderr || fd == stdout) {
					if (i ++ >= LINELEN) {
						putc ('\n', fd);
						i = 0;
					}
				}
			}
			*/
/* Nemytykh 20.06.2001	i_w = strlen (b1 -> pair.f);  / *ri_cs_len(ptr -> pair.f);* /
			if (fd == stderr || fd == stdout) {
                                if (i + i_w > LINELEN) {
					putc('\n',fd);
					i = 0;
				}
			}
			ri_actput (b1 -> pair.f, fd);
			if (fd == stderr || fd == stdout) {
				if (i + i_w < LINELEN) {/ * condition addded, AnK 15.06.1999 * /
					putc(' ',fd);
				}
				i += i_w + 1;
			} else {
				putc (' ', fd); / * Add. Shura.* /
			}
*/

/* Nemytykh 20.06.2001 ------------------------------------------------------------ */
/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
                                char * pc = b1 -> pair.f;
                        	i_w = strlen (b1 -> pair.f); 
                                while( *pc ) {
                                      if( ri_new_line(*pc) ) { i = 0; break; }
                                      else i++;
                                      pc++; 
                                }
                                if (i > LINELEN) {
					putc('\n',fd);
					i = 0;
				}
                                while( *pc ) {
                                      if( ri_new_line(*pc) ) i = 0;
                                      else i++;
                                      pc++; 
                                }

        			ri_actput (b1 -> pair.f, fd);
				if (i++ < LINELEN) putc(' ',fd);
			} else 
*/
                        {
         			ri_actput (b1 -> pair.f, fd);
				putc (' ', fd); 
			};
/* ------------------------------------------------------------------------------- */
			break;

		case LINK_TYPE_NUMBER: /* macrodigit   */
			/*
			if (i_chars) {
				i_chars = 0;
				putc ('\'', fd);
				if (fd == stderr || fd == stdout) {
					if (i ++ >= LINELEN) {
						putc ('\n', fd);
						i = 0;
					}
				}
			}
			*/
			sprintf (buf,"%lu ",b1 -> pair.n); /* was with blank: "%lu ", AnK 15.06.1999 */
			i_w = strlen (buf);
/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
				if (i + i_w > LINELEN) {/ * was >=, AnK 15.06.1999 * /
					putc('\n',fd);
					i = 0;
				}
				i += i_w;
			}
*/
			fputs (buf, fd);
/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
				if (i ++ >= LINELEN) {/ * instead of blank in format, AnK 15.06.1999 * /
					putc ('\n', fd);
					i = 0;
				}
			}
*/
			break;

		case LINK_TYPE_LACT: /*  active left bracket    */
			/*
			if (i_chars) {
				putc ('\'', fd);
				i_chars = 0;
				if (fd == stderr || fd == stdout) {
					if (i ++ >= LINELEN) {
						putc ('\n', fd);
						i = 0;
					}
				}
			}
			*/
			{
				char * cp;

				for (i_w = 0, cp = b1 -> pair.f - 2; 0 != * cp; cp --, i_w ++);
				cp ++;

/* Nemytykh 18.02.2004
				if (fd == stderr || fd == stdout) {
					if (i+1+i_w >= LINELEN) {
						putc('\n',fd);
						i = 0;
					}
					i += 1 + i_w;
				}
*/
				putc('<',fd);
				ri_actput (cp, fd);
			}

/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
				if (i ++ > LINELEN) {
					putc ('\n', fd);
					i = 0;
				} else {
					putc(' ',fd);
				}
			} else 
*/
                        {
				putc (' ', fd);
			};
			break;

		case LINK_TYPE_RACT:  /*  right active bracket   */
			/*
			if (i_chars) {
				putc ('\'', fd);
				i_chars = 0;
				if (fd == stderr || fd == stdout) {
					if (i ++ > LINELEN) {
						putc ('\n', fd);
						i = 0;
					}
				}
			}
			*/
/* Nemytykh 18.02.2004
			if (fd == stderr || fd == stdout) {
				if (i ++ >= LINELEN) {
					putc('\n',fd);
					i = 0;
				}
			}
*/
			putc('>',fd);
			break;

		}
	}
	/*
	if (i_chars) {
		putc ('\'', fd);
		i_chars = 0;
	}
	*/
	return 0;
}


/*	PUT_LINK prints the link pointed by PTR to the file
	pointed by FD. COL is the number of the column 
	where it starts. PUT_LINK returns the number 
	of the column where it ends.     */
/*	If the line does not fit into LINELEN columns, PUT_LINK
	automatically inserts an end of line and prints on
	the next line.		*/
/*	July 28, 1985.	D.T.	*/
/*
int ri_put_link (FILE * fd, LINK * ptr, int col) {
	register int i;
	char buf [MAXWS];
	unsigned long k;

	switch (ptr -> ptype) {
	case LINK_TYPE_LSTRUCTB:  / *  left parenthesis   * /
		if (fd == stderr || fd == stdout) {
			if (col >= LINELEN) {
				putc('\n',fd);
				col = 0;
			}
		}
		putc('(',fd);
		return col+1; / * was col++; AnK 15.06.1999 * /

	case LINK_TYPE_RSTRUCTB:  / *  right parenthesis   * /
		if (fd == stderr || fd == stdout) {
			if (col >= LINELEN) {
				putc('\n',fd);
				col = 0;
			}
		}
		putc(')',fd);
		return col+1; / * was col++; AnK 15.06.1999 * /

	case LINK_TYPE_CHAR:  / * object symbol * /
		if ( fd == stderr || fd == stdout ) {
                	if (col >= LINELEN ) { 
				putc('\n',fd);
				col = 0;
			}
		}
		if (ptr -> pait.c == '\\') {
			putc ('\\', fd);
		} else if (ptr -> pair.c == '\'') {

			putc(ptr -> pair.c, fd);
		}
		return col+1; / * was col++; AnK 15.06.1999 * /

	case LINK_TYPE_COMPSYM: / *  compound symbol    * /
        	i = ri_cs_rest_len(ptr -> pair.f); 
		if (fd == stderr || fd == stdout) {
			if (col + i > LINELEN) {/ * was >=, AnK 15.06.1999 * /
				putc('\n',fd);
				col = 0;
			}
		}
		/ *ri_actput((ptr -> pair.f) + MAXWS,fd);* /
		ri_actput (ptr -> pair.f, fd);
		if (fd == stderr || fd == stdout) {
			if (col + i < LINELEN) / * condition addded, AnK 15.06.1999 * /
				putc(' ',fd);
		} else {
			putc (' ', fd); / * Add. Shura.* /
		}
		return col+1+i;

	case LINK_TYPE_NUMBER: / * macrodigit   * /
		k = ptr -> pair.n;
		sprintf (buf,"%lu",k); / * was with blank: "%lu ", AnK 15.06.1999 * /
		i = strlen (buf);
		if (fd == stderr || fd == stdout) {
			if (col + i > LINELEN) {/ * was >=, AnK 15.06.1999 * /
				putc('\n',fd);
				col = 0;
			}
		}
		fputs (buf, fd);
		if (fd == stderr || fd == stdout) {
			if (col + i < LINELEN) / * instead of blank in format, AnK 15.06.1999 * /
    			putc(' ',fd);
		} else {
			putc (' ', fd);
		}
		return col+i+1; / * +1 was extra with blank in format, AnK 15.06.1999 * /

	case LINK_TYPE_LACT: / *  active left bracket    * /
		i = ri_cs_len(ptr -> pair.f);
		if (fd == stderr || fd == stdout) {
			if (col+1+i >= LINELEN) {
				putc('\n',fd);
				col = 0;
			}
		}
		putc('<',fd);
		/ *ri_actput(ptr -> pair.f,fd);* /
		{
			char * cp = ptr -> pair.f - 1;

			for (cp --; '\0' != * cp; cp --);
			cp ++;
			ri_actput (cp, fd);
		}

		putc(' ',fd);
		return col+2+i;

	case LINK_TYPE_RACT:  / *  right active bracket   * /
		if (fd == stderr || fd == stdout) {
			if (col >= LINELEN) {
				putc('\n',fd);
				col = 0;
			}
		}
		putc('>',fd);
		return col+1; / * was col++; AnK 15.06.1999 * /

	}
	return col;
}
*/

# define INP_BUF_SIZE 256

/*	RI_GET is used by functions GET and CARD.	*/
int ri_get (FILE * fp) {
	register int i/*, i_id*/;
	int flag/*, flag_chars*/;
	char line [INP_BUF_SIZE];
	/*char ca_buf [MAXWS];*/

	flag = 1;
	/*i_id = 0;
	flag_chars = 0;
	ca_buf [0] = 0;
	*/
	rdy(0);
	
	while (flag) {
		/* test for end-of-file. */
		if (fgets (line, INP_BUF_SIZE, fp) == NULL) {
			nns(0);
			break;
		} else {
			for (i = 0; i < INP_BUF_SIZE; i ++) {

				if (line [i] == '\n') flag = 0;
				else if (line [i] == '\0') break;
				else {
					ns (line [i]);
				}
			}
		}
	}
	/*
	while (fgets (line , INP_BUF_SIZE, fp) != NULL) {
		for (i = 0; i < INP_BUF_SIZE; i ++) {
			char * cp;

			switch (line [i]) {
			case 0:
				i = INP_BUF_SIZE;
				break;

			case '\'': / * character sequenst (begin or end) * /
				if (flag_chars) {
					/ * End of sequenst of characters * /
					flag_chars = 0;
					break;
				}
				/ * Begin of sequenst of characters * /
				flag_chars = 1;
				if (ca_buf [0] != 0) {
					ca_buf [i_id] = 0;
					cp = ri_cs_impl (ca_buf);
					ncs (cp);
					ca_buf [0] = 0;
					i_id = 0;
				}
				break;

			case '\\': / * screaning symbol * /
				if (ca_buf [0] != 0) {
					ca_buf [i_id] = 0;
					cp = ri_cs_impl (ca_buf);
					ncs (cp);
					ca_buf [0] = 0;
					i_id = 0;
				}
				switch (line [i + 1]) {
				case '\'': case '\\': case '(': case ')': case '<': case '>': 
				case '{': case '}':
					ns (line [i + 1]);
					break;

				case 't':
					ns ('\t');
					break;

				case 'n':
					ns ('\n');
					break;

				default:
					ri_error (8);
				}
				i ++;
				break;
			
			case '<':
				flag_chars = 0;
				/ * Insert function * /
				/ * While it is ignore !!! And it is inserted as character * /
				ns ('<');
				break;

			case '>':
				flag_chars = 0;
				/ * Insert end of function * /
				/ * While it is ignore !!! And it is inserted as character * /
				ns ('>');
				break;

			case '{': case '}':
				/ *flag_chars = 0;* /
				ns (line [i]);
				break;

			case ' ':
				if (flag_chars) {
					ns (' ');
				} else {
					if (ca_buf [0] != 0) {
						ca_buf [i_id] = 0;
						cp = ri_cs_impl (ca_buf);
						ncs (cp);
						ca_buf [0] = 0;
						i_id = 0;
					}
				}
				break;
			case '\n':
				if (flag_chars) {
					ri_error (8);
				}
				if (ca_buf [0] != 0) {
					ca_buf [i_id] = 0;
					cp = ri_cs_impl (ca_buf);
					ncs (cp);
					ca_buf [0] = 0;
					i_id = 0;
				}
				break;

			default:
				if (flag_chars) {
					ns (line [i]);
					break;
				}
				if (isalnum (line [i]) || line [i] == '-' || line [i] == '_') {
					ca_buf [i_id ++] = line [i];
				} else {
					ns (line [i]);
				}
			}
		} / * end of for * /
	} / * end of while * /
	*/
	out(2);
	est;
	return 0;
}

/* Opens the file with name REFAL#.DAT, where # is LOGNUM (the parameter), with mode MODE. */
int ri_open (int lognum, char * mode, char * name) {

	/* Close the file if it was opened under this logical number. */
	if (NULL != file_table [lognum]) fclose(file_table [lognum]);

	/* Open the new file. */
	/* fprintf(stderr,"Ri____________ open %s.\n",name); */
	file_table [lognum] = fopen(name,mode);
	if (file_table [lognum]  == NULL) {
		fprintf(stderr,"can\'t open %s.\n",name);
		ri_error(2);
	}
	return 0;
}

/* ACTPUT prints the compound symbol pointed by F to the file 
 * with pointer FD. It returns the length of the string.
 */
/*	July 28, 1985.	D.T. */
int ri_actput (char *f, FILE *fd) {
	/*
	register int i;
	char *v;

	v = f - MAXWS;
	i = 0;
	while ((i ++ < MAXWS) && (*v != '\0')) fputc(*v ++, fd);
	return i;
	*/

	return (fputs (f, fd));
}

