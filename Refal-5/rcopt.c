
# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"


static int labels_in_function;

		/* Refal compiler optimizations. */

int refc_opt (translation, t, label)
	struct rasl_instruction *translation [];
	int t;
	int *label;
	{
	int i;
	struct rasl_instruction *h, *next, *z1;
	union param p1;

		/* if there is only one sentence return. */
	if (t <= 2) return 0;

	labels_in_function = 0;

		/* chop tails. */
	for (i = 1; i < t; i ++) chop_tail (translation [i]);

		/* starting from 2nd sentence */
	i = 2;
	while (i < t)
		{
		z1 = translation [i];
		h = translation [1];
			/* compare until mismatch or end. */
		while (comp_next_ins (h, z1, &next))
			{
			z1 = z1 -> next;
			h = next;
			}

			/* insert TRAN/LABEL */
		p1.i = *label;
			/* if there is already a TRAN instruction, then insert it after
				the corresponding LBL instruction. */
		next = h;
		while ((next -> next) -> code == TRAN)
			{
			next = find_def_label ((next -> next) -> p.i);
			}
		insert_instruction (next, TRAN, p1);
		insert_instruction (z1, LBL, p1);

			/* save the address of the definition of the next label:
				it's the one after z1. */
		save_label (*label, z1 -> next);
		(*label) ++;

			/* cut off the header. */
		delete_to_ptr (translation [i], z1);
		translation [i] = z1;
		translation [i] -> code = B;

			/* increment i */
		i ++;
		}
	return 0;
	}

int comp_inst (z1, z2)
	struct rasl_instruction *z1;
	struct rasl_instruction *z2;
	{
	int inst, inst2;
	int i, l, l1, l2, res;

	inst = z1 -> code;
	inst2 = z2 -> code;
	res = 0;
	if (inst == inst2)
		{
		switch (inst)
			{
			case L: case E: case LABEL: case LBL: case ECOND:
				return 0;

				/* takes 2 integer arguments. */
			case SETB:
				if (z1 -> p.d.i1 == z2 -> p.d.i1 && 
					z1 -> p.d.i2 == z2 -> p.d.i2) return 1;
				else return 0;

			case OEXP: case OEXPR: case OVSYM: case OVSYMR: case TRAN:
			case TPLE: case TPLS: case MULE: case MULS: case RDY:
				if (z1 -> p.i == z2 -> p.i) return 1;
				else return 0;

			case NSYM: case NSYMR: case NNS:
				if (z1 -> p.n == z2 -> p.n) return 1;
				else return 0;

			case CSYM: case CSYMR: case NCS: case ACT1: case TEXT:
				if (strcmp (z1 -> p.f,z2 -> p.f) == 0) return 1;
				else return 0;

			case SYMSR: case SYMS:

					/* compare the strings and split if necessary. */
				l1 = strlen (z1 -> p.f);
				l2 = strlen (z2 -> p.f);
				i = l1 > l2 ? l2 : l1;	/* minimum length. */

				for (l = 0; l < i; l++)
					{
					if (z1 -> p.f [l] != z2 -> p.f [l]) break;
					}

					/* test if match occured. */
				if (l == 0)
					return 0;
				else
					{
						/* split the first string */
				if (l < l1) {
					split_string (z1, l, l1, inst);
				}

						/* split the second string. */
					if (l < l2) split_string (z2, l, l2, inst);

					return 1;
					}

			case SYM: case SYMR: case NS:
				if (z1 -> p.c == z2 -> p.c) return 1;
				else return 0;

					/* these instructions take no arguments. */
			case CL: case EMP: case VSYM: case VSYMR: case BL: case BLR:
			case BR: case LEN: case TERM: case TERMR: case PS: case PSR:
			case OUTEST: case PUSHVF: case POPVF: case STLEN:
				return 1;

			case B:		/* special mark: always fail. */
				return 0;

			default:
				fprintf (stderr, "internal error in comp_inst: %d -- aborted\n", inst);
				exit (1);
				break;
			}
		}
			/* take care of SYM and SYMS, and SYMR and SYMSR */
	else if (inst < inst2) {
		res = comp_sym_inst (z1, z2);
	}
	else {
		res = comp_sym_inst (z2, z1);
	}
	return res;
	}

int comp_sym_inst (z1, z2)
	struct rasl_instruction *z1;
	struct rasl_instruction *z2;
	{
	struct rasl_instruction *z;
	int l;
	char *s;

	if ((z1 -> code == SYM && z2 -> code == SYMS) ||
			(z1 -> code == SYMR && z2 -> code == SYMSR))
		{
		if (z1 -> p.c == *(z2 -> p.f))
			{
			l = strlen (z2 -> p.f);
				/* split z2 into two instructions. */
			if (l > 1)
				{
				z = (struct rasl_instruction *) 
					malloc (sizeof (struct rasl_instruction));
				/* 174 line. Not check result of malloc. Shura. 29.01.98 */
				if (NULL == z) {
				  fprintf (stderr, "No memory for rasl instrunction\n");
				  exit (1);
				}
				if (l > 2)
					{
					z -> code = z2 -> code;
					s = rc_memral (l);
					strcpy (s, (z2 -> p.f) + 1);
					z -> p.f = s;
					}
				else /* l == 2 */
					{
					z -> code = z1 -> code;
					z -> p.c = (z2 -> p.f)[1];
					}
				z -> next = z2 -> next;
				z2 -> code = z1 -> code;
				z2 -> p.c = z1 -> p.c;
				z2 -> next = z;
				}
			return 1;
			}
		}
	return 0;
	}

int split_string (z1, l, l1, inst)
	struct rasl_instruction *z1;
	int l;
	int l1;
	int inst;

	{
		char *s;
		struct rasl_instruction *z;

			/* split string at length l and create another instruction. */
			/* copy the tail of the string */
		s = rc_memral (sizeof (char) * (l1 - l + 1));
		strcpy (s, z1 -> p.f + l);
			/* chop the string. */
		z1 -> p.f [l] = '\0';

			/* insert another rasl instruction. */
		z = (struct rasl_instruction *)
			malloc (sizeof (struct rasl_instruction));
		/* 211 line. Not check result of malloc. Shura. 29.01.98 */
		if (NULL == z) {
		  fprintf (stderr, "No memory for rasl instrunction\n");
		  exit (1);
		}
		z -> code = inst;
		z -> p.f = s;
		z -> next = z1 -> next;
		z1 -> next = z;

		
			/* change both new instructions to SYM/SYMR if
				the length is only 1. */
		ch2sym (z);
		ch2sym (z1);

		return 0;
	}

int ch2sym (z)
	struct rasl_instruction *z;
	{
		/* assumes that the code is either SYMS or SYMSR */
		if (strlen (z -> p.f) == 1)
			{
			z -> p.c = z -> p.f [0];
			z -> code += SYM - SYMS;
			}
		return 0;
	}

int chop_tail (q)
	struct rasl_instruction *q;
	{
	union param p1;

	while (is_left_part (q)) q = q -> next;
	p1.n = 0L;
	insert_instruction (q, B, p1);
	return 0;
	}

int is_left_part (q)
	struct rasl_instruction *q;
	{

	if (q -> next == NULL) return 0;
	switch ((q -> next) -> code)
		{
		case OEXP: case OEXPR: case OVSYM: case OVSYMR: case NSYM: 
		case NSYMR: case CSYM: case CSYMR: case SYMS: case SYMSR: 
		case SYM: case SYMR: case CL: case EMP: case VSYM: case VSYMR:
		case TERM: case TERMR: case PS: case PSR: case SETB:

			return 1;
		default:
			return 0;
		}
	}

int comp_next_ins (h, z, next)
	struct rasl_instruction *h;
	struct rasl_instruction *z;
	struct rasl_instruction **next;
	{
	int l;

		/* determine the next of h: move until all labels are skipped. */
	while ((h -> next) -> code == TRAN)
		{
		l = (h -> next) -> p.i;
		h = find_def_label (l);
		}

	*next = h -> next;

	return comp_inst (*next, z -> next);
	}

int delete_to_ptr (b, e)
	struct rasl_instruction *b;
	struct rasl_instruction *e;
	{
	struct rasl_instruction *t;

	if (b == e) return 0;
	b = b -> next;
	while (b != e)
		{
		if (b == NULL)
			{
			fprintf (stderr, "internal error: NULL pointer in delete_to_ptr\n");
			break;
			}
		t = b -> next;
		free ((void *) b);
		b = t;
		}
	return 0;
	}


static struct
	{
	int l;
	struct rasl_instruction *p;
	} a [MAX_SENTENCES];

int save_label (label, z)
	int label;
	struct rasl_instruction *z;
	{
	a [labels_in_function].l = label;
	a [labels_in_function].p = z;
	labels_in_function ++;
	return 0;
	}

	/* find definition of label. */
struct rasl_instruction *find_def_label (l)
	int l;
	{
	int i;

	for (i = 0; i < labels_in_function; i ++)
		if (a [i].l == l) return a[i].p;

	fprintf (stderr, "internal error in find_def_label: label %d -- aborted\n", l);
	exit (1);
	return NULL;
	}



