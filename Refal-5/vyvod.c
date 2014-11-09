
	/* This is the IBM PC/AT version. July 1 1986. DT. */

# include "rasl.h"
# include "cdecl.h"
# include "cfunc.h"


/* DEFLABEL defines the offset of the label. */
char * rc_deflabel(char * label, struct functab * table) {
	struct functab *p;

	p = searchf(label,table); /* search the function table.   */
	/* insert into table. */
	if (p == NULL) {
		p = (struct functab *) rc_allmem(sizeof (struct functab));
		if (p == NULL) {
			fprintf(stderr,"Ran out of memory.\n");
			exit(1);
		}
		p -> next = table -> next;
		table -> next = p;

		if (NULL == (p -> name = (char *) malloc (strlen (label) + 1))) {
			fprintf (stderr, "No memory for label name\n");
			exit (1);
		}
		strcpy(p -> name,label);
	}
	p -> offset = z;
	return label;
}

/* GETCSN gets the id number of the compound symbol. */
long rc_getcsn(char * compsym) {
	struct functab *p;

	p = searchf(compsym, cs);
	if (p == NULL) {
		fprintf(stderr, "Compound symbol %s not found in table.\n",compsym);
		return -1L;
	}
	return p -> offset;
}

int rc_vyvod (struct rasl_instruction *translation [], int t) {
	/* VYVOD: Output the Rasl operators in the format of the host machine.  DT June 19 1986. */
	char funname[MAXWS], lname[MAXWS];
	struct rasl_instruction *v;
	int /*i,*/ k, s;
	int opcode;

		/* 1. copy the name of the function to a string. */
	/*strncpy (funname, (translation [0]) -> p.f, MAXWS);*/
	strcpy (funname, translation [0] -> p.f);

	for (s = 0; s < t; s ++) {
		v = translation [s];
		while (v != NULL) {
			opcode = v -> code;
			switch (opcode) {
			/* This RASL instruction takes an address of a function as an argument. */
			case ACT1:
				write_byte (opcode);
				k = strlen (v -> p.f);
				write_byte ('\0'); /* Added */
				write_bytes (v -> p.f, k);
				/*for (i = k; i < MAXWS; i++)*/ write_byte('\0');
				z += sizeof (char) + sizeof (char *); /* ??!! */
				break;

			/* These RASL operators take a compound symbol for an argument. */
			case CSYM: case CSYMR: case NCS:
				write_byte(opcode);
				/*strncpy (lname, v -> p.f, MAXWS);*/
				strcpy (lname, v -> p.f);
				write_long(rc_getcsn (lname));
				z += sizeof (char) + sizeof (char *);
				break;

			/* These RASL operators require a (long) number	for a parameter */
			case NSYM: case NSYMR: case NNS:
				write_byte (opcode);
				write_long (v -> p.n);
				z += sizeof (char) + sizeof (char *);
				break;

			/* No arguments for these rasl operators. */
			case BL: case BLR: case BR: case CL: case EMP: case EST:
			case PLEN: case PLENS: case PLENP: case PS: case PSR:
			case TERM: case TERMR: case LEN: case LENP: case VSYM:
			case VSYMR: case OUTEST: case POPVF: case PUSHVF: case STLEN:
				write_byte (opcode);
				z ++;
				break;

			/* operand: 1 byte treated as a character. */
			case SYM: case SYMR: case LENS: case NS:
				write_byte (opcode);
				write_byte (v -> p.c);
				z += sizeof (char) + sizeof (char);
				break;

			/* operand: 1 byte - treated as a number. */
			case MULE: case MULS: case TPLE: case TPLS:
				write_byte (opcode);
				write_long (v -> p.i);
				z += sizeof (char) + sizeof (long);
				break;

			/* operand: 1 byte - treated as a number. */
			/*case MULE: case MULS:*/ case OEXP: case OEXPR: case OVSYM:
			case OVSYMR: case RDY: case LENOS:/* case TPLE: case TPLS:*/
				write_byte (opcode);
				write_byte (v -> p.i);
				z += sizeof (char) + sizeof (char);
				break;

			/* operands: 2 bytes: treated as 2 numbers. */
			case SETB:
				write_byte (opcode);
				/*
				write_byte (v -> p.d.i1);
				write_byte (v -> p.d.i2);
				z += sizeof (char) + 2 * sizeof (char);
				*/
				write_long (v -> p.d.i1);
				write_long (v -> p.d.i2);
				z += sizeof (char) + 2 * sizeof (long);
				break;

			/* operands: length and string. */
			case SYMS: case SYMSR: case TEXT:
				{
					char * cp;

					k = strlen (v -> p.f);
					for (cp = v -> p.f; k > 255; k -= 255, cp += 255) {
						write_byte (opcode);
						write_byte (255);
						write_bytes (cp, 255);
						z += 2 + 255;
					}
					if (k == 1) {
						if (opcode == SYMS) {
							write_byte (SYM);
						} else if (opcode == SYMSR) {
							write_byte (SYMR);
						} else {
							write_byte (NS);
						}
						write_byte (* cp);
						z += 2;
					} else if (k != 0) {
						write_byte (opcode);
						write_byte (k);
						write_bytes (cp, k);
						z += sizeof (char) + sizeof (char) + k;
					}
				}
				break;

			/* operand: label of form FUNNAME$NUMBER */
			case TRAN: case ECOND:
				write_byte (opcode);
				k = sprintf (lname, "%s$%d", funname, v -> p.i);
				write_bytes (lname, k);
				/*for (i = k; i < MAXWS; i ++)*/ write_byte (0);
				z += sizeof (char) + sizeof (char *);
				break;

			/* definition of labels. */
			case LBL: case LABEL:
				k = sprintf (lname, "%s$%d", funname, v -> p.i);
				if (opcode == LABEL) {
					write_byte (opcode);
					write_byte ('\0'); /* Added */
					write_bytes (lname, k);
					/*for (i = k; i < MAXWS; i ++)*/ write_byte (0);
					z += k+2/*MAXWS*/;
				}
				rc_deflabel (lname, ft);
				break;

			case L: case E:
				write_byte (opcode);
				/*strncpy (lname, v -> p.f, MAXWS);*/
				strcpy (lname, v -> p.f);
				k = strlen (lname);
				write_byte ('\0');
				write_bytes (lname, k);
				/*for (i = k; i < MAXWS; i ++)*/ write_byte (0);
				z += k + 2/*MAXWS*/;
				if (opcode == E) rc_deflabel (lname, fe);
				rc_deflabel (lname, ft);
				rc_deflabel (lname, fb);
				break;

			/* NULL or special mark B: simply ignore. */
			/*case NULL:*/case 0: case B:
				break;

			default:

					fprintf (stderr, "internal error in VYVOD: opcode %d -- ignored\n", opcode);
					break;
				}
			v = v -> next;
			}
		}
	return 0;
	}


