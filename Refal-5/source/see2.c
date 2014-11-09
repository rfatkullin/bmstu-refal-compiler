
# define DEFINE_EXTERNALS 1

# include "rasl.h"
# include "cdecl.h"

int main (int, char * *);
int rc_rsee2 (FILE *fp_out);
char *rasl_code (int);

int main(argc, argv)
	int argc;
	char *argv[];
	{
	char *s, filename [40];
	FILE *fp_out;

	if (argc < 2) {
		printf("Usage: SEE2 in_file [out_file]\n");
		exit(1);
		};
	fdtmpr = fopen (argv[1], "rb");
	if (fdtmpr == NULL) {
		printf("Can't open %s\n", argv[1]);
		exit(1);
		};
	if (argc > 2) {
		fp_out = fopen(argv[2], "wt");
		if (fp_out == NULL) {
			printf("unable to open file %s\n", argv[2]);
			exit(1);
			}
		}
	else fp_out = stdout;
	rc_rsee2(fp_out);
	fprintf(fp_out, "\n\n That's it, folks\n");
	if (fp_out != stdout) fclose(fp_out);
	exit(0);
	}

int rc_rsee2 (fp_out)
	FILE *fp_out;
	{
	/* RSEE2: Print the Refal Interpretation file 2. Feb. 21 1987. */

	unsigned char opcode,c,d;
	char lname[MAXWS];
	char title[MAXWS];
	int err,i,j;
	long n,k,k1,xt,nt,cs,bt,z;

	/* Read the title */
	for (i = 0; i< MAXWS; i++) {
		read_byte(title[i]);
		if (title [i] == 0) break;
	}
	fprintf(fp_out, "TITLE %s\n",title);
	/* Read the size */
	read_long(n);
	fprintf(fp_out, "The size of code is %ld\n",n);
	read_long(nt);
	fprintf(fp_out, "size of entry table = %ld\n",nt);
	read_long(xt);
	fprintf(fp_out, "size of external table = %ld\n",xt);
	read_long(cs);
	fprintf(fp_out, "size of compound symbol table = %ld\n",cs);
	read_long(bt);
	fprintf(fp_out, "size of local table = %ld\n",bt);

	/* Get the ENTRY table. */
	fprintf(fp_out, "\nThe ENTRY table.\n");
	for (i = 0; i < nt; i++) {
		for (j = 0; j < MAXWS; j++) {
			read_byte(c);
			lname[j] = c;
			if (c == 0) break;
		}
		read_long(k);
		fprintf(fp_out, "%d: %s %ld\n",i,lname,k);
	}
	fprintf(fp_out, "\nThe EXTERNAL table.\n");
	for (i = 0; i < xt; i++) {
		for (j = 0; j < MAXWS; j++) {
			read_byte(c);
			lname[j] = c;
			if (c == 0) break;
		}
		fprintf(fp_out, "%d: %s\n",i,lname);
	}
	fprintf(fp_out, "\nThe Compound Symbol table.\n");
	for (i = 0; i < cs; i++) {
		for (j = 0; j < MAXWS; j++) {
			read_byte(c);
			lname[j] = c;
			if (c == 0) break;
		}
		fprintf(fp_out, "%d: %s\n",cs-1-i,lname);
	}
	fprintf(fp_out, "\n");
	fprintf(fp_out, "The Local table.\n");
	for (i = 0; i < bt; i++) {
		read_long(k);
		fprintf(fp_out, " %ld\n",k);
	}
	z = 0L;
	err = 0;

	while (read_byte(opcode) == 1) {
		fprintf(fp_out, "%4ld: %s ", z, rasl_code (opcode));
		switch (opcode) {

		case 1: /* ACT1N */
			/* This RASL instruction takes an address of a function as an argument. */

			fprintf(fp_out, " External Function: ");
			/*read_byte (c);*/
			for (i = 0; i < MAXWS; i++) {
				read_byte(c);
				if (c) putc(c, fp_out);
				else break;
			}
			z += 5;
			break;

		case 2: /* ACT1 */
			read_long(k);
			fprintf(fp_out, " Function offset %ld",k);
			z += 5;
			break;

		case 55: /* CSYM */
		case 56: /* CSYMR */
		case 59: /* NCS */
			/* These RASL operators take a compound symbol as an argument. */
			read_long(k);
			fprintf(fp_out, " Compound symbol %ld",k);
			z += 1+sizeof(long);
			break;

		case 57: /* NSYM */
		case 58: /* NSYMR */
		case 60: /* NNS */
			/* These RASL operators require a (long) number as a parameter */
			read_long(k);
			fprintf(fp_out, " Long Number %ld",k);
			z += 5;
			break;

		case 100: /* Builtin function call */
			/* This RASL operator require a (long) number as a parameter */
			read_long(k);
			fprintf(fp_out, " Long Number %ld",k);
			z += 5;
			break;

		case 105: /* Builtin function call */
			/* This RASL operator require a (long) number as a parameter */
			read_long(k);
			fprintf(fp_out, " Long Numbers %ld",k);
			read_long(k);
			fprintf(fp_out, " and %ld",k);
			z += 9;
			break;

		case 3: /* BL */
		case 4: /* BLR */
		case 5: /* BR */
		case 6: /* CL */
		case 10: /* EMP */
		case 11: /* EST */
		case 16: /* PLEN */
		case 17: /* PLENS */
		case 18: /* PLENP */
		case 19: /* PS */
		case 20: /* PSR */
		case 27: /* TERM */
		case 28: /* TERMR */
		case 35: /* LEN */
		case 37: /* LENOS */
		case 48: /* VSYM */
		case 49: /* VSYMR */
		case 50: /* OUTEST */
		case 52: /* POPVF */
		case 53: /* PUSHVF */
		case 54: /* STLEN */ 
			/* No arguments for these RASL operators. */
			z ++;
			break;

		case 7: /* SYM */
		case 8: /* SYMR */
		case 36: /* LENS */
		case 43: /* NS */
			/* These RASL operators require a byte (character) as parameter. */
			read_byte(c);
			fprintf(fp_out, " Character %d %c",c,c);
			z += 2;
			break;

		case 12: /* ELEN */
		case 23: /* OEXP */
		case 24: /* OEXPR */
		case 25: /* OVSYM */
		case 26: /* OVSYMR */
		case 29: /* RDY */
		case 38: /* LENOS */
			/* These RASL operators take one operand of size 1 byte. */
			read_byte(c);
			fprintf(fp_out, " TE pointer %d",c);
			z += 2;
			break;

		case 13: /* MULE */
		case 14: /* MULS */
		case 45: /* TPLE */
		case 46: /* TPLS */
			read_long (k);
			fprintf (fp_out, "TE pointer %d", k);
			z += 5;
			break;

		case 34: /* SETB */
			/* This RASL operator require two operands of sizeof(long). */
			read_long(k1);
			read_long(k);
			fprintf(fp_out, " TE pointers %d %d",k1,k);
			z += 1+2*sizeof(long);
			break;

		case 39: /* SYMS */
		case 40: /* SYMSR */
		case 41: /* TEXT */
			/* These RASL operators take 1 byte and a variable number of bytes as parameters. */
			read_byte(d);
			fprintf(fp_out, " Takes %d arguments:\n",d);
			for (i = 0; i < d; i++) {
				read_byte(c);
				fprintf(fp_out, "\t\tCharacter %d %c\n",c,c);
			}
			z += 2+d;
			break;

		case 47: /* TRAN */
		case 51: /* ECOND */
			/* These RASL operators take as argument a label of form
			 * FUNNAME$NUMBER, where FUNNAME is the current function
			 * name, and NUMBER is a number.
			 */
			read_long(k);
			fprintf(fp_out, " Label offset %ld",k);
			z += 5;
			break;

		case 101: /* LBL */
		case 104: /* LABEL */
			/* These RASL operators define a label of form FUNNAME$NUMBER, 
			 * where FUNNAME is the current function name, and NUMBER is a number.
			 */
 
		case 102: /* L */
		case 103: /* E */
			/* These operators define labels. */
			fprintf(fp_out, " Function defined: ");
			read_byte (c);
			for (i = 0; i < MAXWS; i++) {
				read_byte(c);
				if (c) putc(c, fp_out);
				else break;
			}
			z += /*MAXWS*/i + 2;
			fprintf(fp_out, " Offset = %ld",z);
			break;
 
		default:
			fprintf(fp_out, ": Strange Opcode = %ld",opcode);
			err ++;
			break;
		}
		fprintf(fp_out, "\n");
	}

	fprintf(fp_out, "size = %ld, z = %ld.\n",n,z);
	fprintf(fp_out, "errors = %d\n",err);
	return 0;
}

# define NUM_OF_RASL_INSTR 55

	char *mnemonics [NUM_OF_RASL_INSTR] = {
		"ACT_EXTRN", "ACT1", "BL", "BLR", "BR", "CL", "SYM", "SYMR",
		"EMP", "EST", "MULE", "MULS", "PLEN", "PLENS", "PLENP", "PS",
		"PSR", "OEXP", "OEXPR", "OVSYM", "OVSYMR", "TERM", "TERMR",
		"RDY", "SETB", "LEN", "LENS", "LENP", "LENOS", "SYMS", "SYMSR",
		"TEXT", "NS", "TPLE", "TPLS", "TRAN", "VSYM", "VSYMR", "OUTEST",
		"ECOND", "POPVF", "PUSHVF", "STLEN", "CSYM", "CSYMR", "NSYM", "NSYMR",
		"NCS", "NNS", "LBL", "L", "E", "LABEL", "BUILT_IN1", "Special Mark B"
		};

	int rasl_numbers [NUM_OF_RASL_INSTR] = {
		ACT_EXTRN, ACT1, BL, BLR, BR, CL, SYM, SYMR, EMP, EST,
		MULE, MULS, PLEN, PLENS, PLENP,	PS, PSR, OEXP,
		OEXPR, OVSYM, OVSYMR, TERM, TERMR, RDY, SETB,
		LEN, LENS, LENP, LENOS, SYMS, SYMSR, TEXT, NS,
		TPLE, TPLS, TRAN, VSYM, VSYMR, OUTEST, ECOND,
		POPVF, PUSHVF, STLEN, CSYM, CSYMR, NSYM, NSYMR,
		NCS, NNS, LBL, L, E, LABEL, BUILT_IN1, B
		};
	char unknown [10];

char *rasl_code (id)
	int id;
	{
	int i;

	for (i = 0; i < NUM_OF_RASL_INSTR; i ++)
		if (id == rasl_numbers[i]) return mnemonics[i];
	sprintf (unknown, "RASL-%d", id);
	return unknown;
	}
