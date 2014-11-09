
# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"


	/* Refal Compiler function: accepts as input pointer to
		tree and returns compiled code. */
	/* general translation functions and right side translation
		functions. see also file rcleft.c */

int refcom (q)
	struct node *q;
	{
		int rasl_ins, t, number;
		union param p1;
		struct node *sent;
		struct rasl_instruction *translation [MAX_SENTENCES];
	
		if (q -> nt == ENTRY)
			{
			rasl_ins = E;
				/* add this function to the list of entry functions. */
			rc_mkentry(q -> a2.pchar);
			}
		else rasl_ins = L;
		p1.f = q -> a2.pchar;

			/* initialize the list of RASL instructions */
				/* the first sentence contains the label. */
		t = 0;
		ftransl = (struct rasl_instruction *) 
			malloc (sizeof (struct rasl_instruction));
		/* 33 line. Not check result of malloc. Shura. 29.01.98 */
		if (NULL == ftransl) {
		  fprintf (stderr, "No memory for rasl instrunction\n");
		  exit (1);
		}
		/* 34 line. Was ... = NULL. Shura. 29.01.98 */
		ftransl -> code = 0;
		ftransl -> next = NULL;
		translation [t] = ftransl;
		rc_out (rasl_ins, p1);
		t ++;

		q = q -> a3.tree;
		while (q != NULL)
			{
				/* initialize the list of RASL instruction corresponding to
					this sentence. */
			if (t >= MAX_SENTENCES)
				{
				fprintf (stderr, "Too many sentences in a function. Aborted.\n");
				exit (1);
				}
			ftransl = (struct rasl_instruction *) 
				malloc (sizeof (struct rasl_instruction));
			/* 51 line.Not check result of malloc.Shura.29.01.98 */
			if (NULL == ftransl) {
			  fprintf (stderr, "No memory for rasl instrunction\n");
			  exit (1);
			}
			/* 52 line. Was ... = NULL. Shura. 29.01.98 */
			ftransl -> code = 0;
			ftransl -> next = NULL;
			translation [t] = ftransl;
			t ++;
				/* insert a dummy instruction. */
			p1.n = 0L;
			rc_out (B, p1);

				/* process sentences one by one. */
			sent = q -> a2.tree;
				/* check */
			if (sent -> nt != ST)
				{
				fprintf (stderr, "strange node type %d. ST expected\n", sent -> nt);
				break;
				}
				/* clear the local variables table. */
			table_len = 0;
				/* process the left side. */

			number = 1;	/* first available number. */
			transl_left (sent -> a2.chunk, &number);

				/* process the right tail. */
			tr_rtail (sent -> a3.tree, &number);
			q = q -> a3.tree;
			}

		/* optimize the code. */
		refc_opt (translation, t, &last_label);
		rc_post_opt (translation, t);

		/* output the code. */
		rc_vyvod (translation, t);
		refc_out (translation, t);
		return 0;
	}

int tr_rtail (tail, number)
	struct node *tail;
	int *number;
	{
	union param p1;
	int label, table_save, num_save;
	struct node *sub_tree;

		/* 4 cases. */
	while (tail != NULL)
		{
		switch (tail -> nt)
			{
				/* right side of a sentence. */
			case RCS1:

				p1.i = 0;
				rc_out (RDY, p1);
				transl_right (tail -> a2.chunk, 1);
				tail = NULL;
				break;

				/* simple condition. */
			case RCS2:

				p1.n = 0L;
				rc_out (PUSHVF, p1);
				label = last_label;	/* first available label. */
				last_label ++;	/* increment the first available label. */
				transl_right (tail -> a2.chunk, 0);
				p1.i = label;
				rc_out (ECOND, p1);
				rc_out (LABEL, p1);
				rc_out (POPVF, p1);
				(*number) ++;
				transl_left (tail -> a3.chunk, number);
				tail = tail -> a4.tree;
				break;

				/* branching on a left side of a condition. */
			case RCS3:

				p1.n = 0L;
				rc_out (PUSHVF, p1);
				label = last_label;	/* first available label. */
				last_label ++;	/* increment the first available label. */
				transl_right (tail -> a2.chunk, 0);
				p1.i = label;
				rc_out (ECOND, p1);
				rc_out (LABEL, p1);
				rc_out (POPVF, p1);
				rc_out (STLEN, p1);
					/* save the current number of variables in the table. */
				table_save = table_len;
					/* save the current table of element number */
				(*number) ++;
				num_save = *number;

				sub_tree = tail -> a3.tree;
				while (sub_tree != NULL)
					{
						/* check */
					if (sub_tree -> nt != LSF)
						{
						fprintf (stderr, "Illegal Sub tree in RCS3: %d -- aborted\n", 
							sub_tree -> nt);
						exit (1);
						}
						/* check */
					if ((sub_tree -> a2.tree) -> nt != LSF1)
						{
						fprintf (stderr, "Illegal Sub tree in LSF: %d -- aborted\n",
							(sub_tree -> a2.tree) -> nt);
						exit (1);
						}
					if (sub_tree -> a3.tree != NULL)
						{
						label = last_label;
						last_label ++;
						p1.i = label;
						rc_out (TRAN, p1);
						}
					transl_left ((sub_tree -> a2.tree) -> a2.chunk, number);
					tr_rtail ((sub_tree -> a2.tree) -> a3.tree, number);
					if (sub_tree -> a3.tree != NULL)
						{
						p1.i = label;
						rc_out (LABEL, p1);
						}
						/* restore the table of variables */
					table_len = table_save;
					*number = num_save;
					sub_tree = sub_tree -> a3.tree;
					}
				tail = tail -> a4.tree;
				break;

				/* branching on a whole condition. */
			case RCS4:
				/*fprintf (stderr, "RCS4: Not implemented: *number = %d -- aborted\n", *number);*/
				fprintf (stderr, "Error on line %d\n", line_no);
				exit (0);
				break;

			default:
				fprintf (stderr, "internal errors in right_tail, node: %d -- aborted\n", 
					tail -> nt);
				exit (1);
				break;
				}
		}
	return 0;
	}

		/* free memory allocated for translation. */
int refc_out (translation, t)

	struct rasl_instruction *translation [];
	int t;

	{
	int i;
	struct rasl_instruction *q, *z;

	for (i = 0; i < t; i ++)
		{
		q = translation [i];
		while (q != NULL)
			{
			z = q -> next;
			free ((void *) q);
			q = z;
			}
		}
		
	return 0;
	}

int transl_right (e, transplant)
	struct element *e;
	int transplant;
	{
	/*register*/ int i;
	int occur;
	union param p1;
	unsigned char flags [MAX_TABLE_LENGTH / 8];		/* see table. */

		/* set all flags to zero. */
	for (i = 0; i < MAX_TABLE_LENGTH / 8; i ++) flags [i] = 0;
	i = 0;
	/* 240 line. Was ... = NULL. Shura. 29.01.98 */
	while (e[i].type != 0)
		{
		switch (e[i].type)
			{
			case T_VAR:
			case E_VAR:

					/* first try to transplant */
				if (!transplant) occur = -1;
				else if ((occur = get_var_index (e[i].body.i, flags)) != -1)
					{
					check_bit (flags, occur);
					p1.i = table [occur].te_offset;
					rc_out (TPLE, p1);
					}
					/* otherwise copy  it. */
				if (occur == -1)
					{
					occur = check_var (e[i].body.i);
					p1.i = occur;
					rc_out (MULE, p1);
					}
				break;

			case S_VAR:

					/* first try to transplant */
				if (!transplant) occur = -1;
				else if ((occur = get_var_index (e[i].body.i, flags)) != -1)
					{
					check_bit (flags, occur);
					p1.i = table [occur].te_offset;
					rc_out (TPLS, p1);
					}
					/* otherwise copy  it. */
				if (occur == -1)
					{
					occur = check_var (e[i].body.i);
					p1.i = occur;
					rc_out (MULS, p1);
					}
				break;

			case CHAR:
				p1.c = e[i].body.c;
				rc_out (NS, p1);
				break;
				
			case ATOM:
				p1.f = e[i].body.f;
				rc_out (NCS, p1);
				break;
				
			case DIGIT:
				p1.n = e[i].body.n;
				rc_out (NNS, p1);
				break;

			case STRING:
				occur = strlen (e[i].body.f);
				if (occur > 1)
					{
					p1.f = e[i].body.f;
					rc_out (TEXT, p1);
					}
				else
					{
					p1.c = *(e[i].body.f);
					rc_out (NS, p1);
					}
				break;

			case LPAR:
			case ACT_LEFT:
				p1.n = 0L;
				rc_out (BL, p1);
				break;

			case RPAR:
			case ACT_RIGHT:
				p1.n = 0L;
				rc_out (BR, p1);
				if (e[i].type == RPAR) break;
				p1.f = e [i].body.f;
				rc_out (ACT1, p1);
				break;

			default:
				fprintf (stderr, "internal errors in transl_right, Type: e[%d].type = %d\n",
					i, e[i].type);
				exit (1);
				break;
			}
			i ++;
		}
		/* if called with transplant parameter issue OUTEST operator. */
	if (transplant)
		{
		rc_out (OUTEST, p1);
		}

	return 0;
	}

int get_var_index (index, bits)
	int index;
	unsigned char *bits;
	{
	/*register*/ int i;

	for (i = 0; i < table_len; i ++)
		{
		if (table [i].index == index && !(is_bit_checked (bits, i))) return i;
		}
		/* not found. */
	return -1;
	}

int is_bit_checked (bitmap, index)
	unsigned char *bitmap;
	int index;
	{
	int bytes, bits;
	unsigned char mask;

	bytes = index / 8;
	bits = index % 8;
	mask = 1 << bits;
	return (bitmap [bytes] & mask);
	}

int check_bit (bitmap, index)
	unsigned char * bitmap;
	int index;
	{
	int bytes, bits;
	unsigned char mask;

	bytes = index / 8;
	bits = index % 8;
	mask = 1 << bits;
	bitmap [bytes] = bitmap [bytes] | mask;
	return 0;
	}

