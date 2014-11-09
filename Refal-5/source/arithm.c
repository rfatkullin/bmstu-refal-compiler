
# include <stdlib.h>
# include <time.h>

# include "decl.h"
# include "macros.h"
# include "fileio.h"
# include "arithm.h"
# include "ifunc.h"


	/* MAX_SHORT = 2^16  */
# define MAX_SHORT 65536
# define DIV_BY_ZERO -1
# define MAX_UNSIGNED_INT ((unsigned int)~0)

int rf_symb () 
	{
	char c;
	unsigned short *num, *res;
	int sign, length, l;
	int i, err;

	check_frz_args (1);

	sign = 0;
	cl;
	if (LINK_CHAR (tbel (3)))
		{
		c = tbel (3)->pair.c;
		if (c == '-') sign = -1;
		else if (c == '+') sign = 1;
		else
			{
			ri_error(8);
			return 1;
			}
		b1 = tbel (3)->foll;
		}
	else b1 = tbel (3);

	get_length (b1, tbel (4), &length);
	length *= 2;

		/* allocate memory for temporary array and for the output array */
		/* since log10 (MAX_SHORT) < 5 we need not more than 5 lengths
			of decimal digits. */
	num = (unsigned short *) malloc (sizeof (unsigned short) * 6 * length);
	if (num == NULL)
		{
		ri_error(1);
		return 1;
		};
	res = num + length;

		/* save the number in the temporary array. */
	err = store_long_num (num, b1, tbel (4), length);
	if (err == 3) goto internal_error;
	else if (err == 2) goto call_freeze;
	else if (err == 1) ri_imp ();

		/* convert the number from binary to decimal string */
	reverse_num (num, length);
	convert_num (num, (unsigned long) MAX_SHORT, length, 
		res, (unsigned long) 10, &l);

	rdy (0);

		/* show the sign. */
	if (sign == -1) ns ('-')
	else if (sign == 1) ns ('+');

		/* skip the leading zeros. */
	for (i = l-1; i >= 0; i --) if (res [i] != 0) break;
	if (i < 0)
		{
		ns ('0');
		}
	else
		{
		for ( ; i >= 0; i --)
			{
			c = '0' + res [i];
			ns (c);
			};
		}
	free ((char *) num);
	out (2);
	est; 
	return 0;  

	internal_error:
		fprintf (stderr, "internal error in rf_symb ()\n");
		exit (1);

	call_freeze:
		ri_frz (2);
		return 0;
	}


int rf_numb ()	
	{ 
	char c;
	int sign, length, l;
	unsigned long k;
	unsigned short *res, *arg;

	check_frz_args (1);

	k = 0L;
	sign = 1;
		/* Loop 1: skip through white space until the first digit or sign. */
	movb1;
	while (b1 != b2)
		{
		if (LINK_VAR (b1))
			{
			ri_frz (2);
			return 0;
			}
		else if (LINK_CHAR (b1))
			{
			c = b1->pair.c;
			switch (c)
				{
				case ' ': case '\t':
					movb1;
					break;

				case '-':
					sign = -1;
				case '+':
					b1 = b1->foll;
				case '0': case '1': case '2': case '3': case '4':
				case '5': case '6': case '7': case '8': case '9':
					goto out_loop1;

				default:
					goto zero_number;
				}
			}
		else goto zero_number;
		};

	out_loop1:
	if (b1 == b2) goto zero_number;

		/* at this point b1 must point at the first digit character. */
	get_length (b1, b2, &length);

		/* since log65536 (10) < 1/4, we need 1/4 length for the number
			in base MAX_SHORT */
	res = (unsigned short *) malloc (sizeof (unsigned short) * 
		(1 + length + length / 4));
	if (res == NULL)
		{
		fprintf (stderr, "Unable to allocate memory in RF_NUMB()\n");
		exit (1);
		}
	else arg = res + 1 + length / 4;

	k = 0L;
	length = 0;
	while (b1 != b2)
		{
		if (!LINK_CHAR (b1)) break;
		c = b1->pair.c;
		if (! isdigit (c)) break;
			/* strip the leading zeros. */
		if (c != '0' || length != 0) arg [length ++] = c - '0';
		b1 = b1->foll;
		}

	if (length != 0) convert_num (arg, (unsigned long) 10, length, res,
			(unsigned long) MAX_SHORT, &l);
	else l = 0;

	rdy (0);
		/* if number is zero then sign is positive. */
	if (sign == -1)
		{
		if (is_zero_num (res, l) == 0) ns ('-');
		}
	restore_long_num (res, l);
	free ((char *) res);
	out (2);
	est;
	return 0;

	zero_number:
		sign = 1;
		k = 0L;

	restart:
		rdy (0);
		nns (k);
		out (2);
		est;
		return 0;
	}

int get_length (b, q, l)
	LINK *b;
	LINK *q;
	int *l;
	{
	*l = 1;
	if (b != q)
		while (b != q) 
			{
			*l += 1;
			b = b->foll;
			}
	return 0;
	}


int reverse_num (num, length)
	unsigned short num [];
	int length;
	{
	unsigned short k;
	int i;

	for (i = 0; i < length/2; i ++)
		{
		k = num [i];
		num [i] = num [length - i - 1];
		num [length - i - 1] = k;
		}
	return 0;
	}


	/* divide a number (represented as an array in base B of length L)
		by another number (X), the result is stored in NUM and REM. */

int divide_num (num, b, l, x, rem)
	unsigned short num [];
	unsigned long b;
	int l;
	unsigned long x;
	unsigned short *rem;
	{
	/*register*/ int i;
	unsigned long k, z;

	k = 0L;
	for (i = 0; i < l; i ++)
		{
		k = k * b + num [i];
		z = k / x;
		num [i] = (unsigned short) z;
		k = k % x;
		}
	*rem = k;

	return 0;
	}

int is_zero_num (num, l)
	unsigned short num [];
	int l;
	{
	/*register*/ int i;

	for (i = 0; i < l; i ++)
		if (num [i] != 0) return 0;
	return 1;
	}

	/* converts a number from base BASE to base B */
int convert_num (num, base, length, res, b, l)

	unsigned short num [];
	unsigned long base;
	int length;
	unsigned short res [];
	unsigned long b;
	int *l;
	{
	/*register*/ int i;
	unsigned short rem;

	i = 0;
	do
		{
		divide_num (num, base, length, b, &rem);
		res [i ++] = rem;
		}
		while (!is_zero_num (num, length));

	*l = i;
	return 0;
	}

	/* store a long_num in reverse order. */
int store_long_num (op, b, q, l)
	unsigned short op [];
	LINK *b;
	LINK *q;
	int l;
	{
	/*register*/ int i;
	unsigned long n;

	i = 0;
	while (i < l)
		{
		if (LINK_VAR (q)) return 2;
		else if (!LINK_NUMBER (q)) return 2;/*1;*/
		n = q->pair.n;
		op [i ++] = n % MAX_SHORT;
		op [i ++] = n / MAX_SHORT;
		if (q == b) break;
		q = q->prec;
		}
	if (i > l) return 3;
	return 0;
	}


		/* restore the long_num in reverse order. */
int restore_long_num (res, l)
	unsigned short *res;
	int l;
	{
	/*register*/ int i;
	unsigned long n;
        char start = 1;

	if (l != 0)
		{
		if (l % 2 == 1) res [l ++] = 0;
		for (i = l-1; i >= 0; i -= 2)
			{
			n = (unsigned long)(res [i-1]) + (unsigned long)(res [i]) * MAX_SHORT;
                        /* skip the leading zeros */
			if ( (n != 0L) || (i < 2) ) 
                             { nns (n);
                               start = 0;
                             }
                        else if ( !start ) nns(n);
			}
		}
	else nns (0L);
	return 0;
	}


	/* compares two long nums and returns 0 if [op1, l1] > [op2, l2] */
int compare_num (op1, l1, op2, l2)
	unsigned short op1 [];
	int l1;
	unsigned short op2 [];
	int l2;
	{
	/*register*/ int i;

		/* strip off the leading zeros. */
	while ((l1 > 1) && (op1 [l1 - 1] == 0)) l1 --;
	while ((l2 > 1) && (op2 [l2 - 1] == 0)) l2 --;

	if (l1 > l2) return 1;
	else if (l1 < l2) return -1;
	else
		{
		for (i = l1-1; i >= 0; i --)
			if (op1 [i] > op2 [i]) return 1;
			else if (op1 [i] < op2 [i]) return -1;
		}

		/* otherwise these numbers are equal. */
	return 0;
	}


int rf_arithm (op)
	int op;
	{
	int err;
	int sign1, sign2, sign;
	unsigned short *op1, *op2, *res, *res2, *temp;
	int length1, length2, length, len2;


	check_frz_args (1);

	/*	RASL operators of matching.	*/
	if (LINK_EVAR (b1)) goto call_freeze;
		/* There are several accepatable formats: e.g.
			<ADD s.1 e.2>, <ADD '-' s.1 e.2> or <ADD (e.1) e.2>.
			In the first two cases only a digit can be the first operand */
	if (LINK_STRUCTB (b1->foll))	/* the old format. */
		{
		ps;
		cl; 
		setb (4, 2); 
		cl;
		}
	else /*  the new format. */
		{
		movb1;
			/* tbel (5 & 6) contains the pointers to the first argument. */
		if (LINK_CHAR(b1))
			{
			tbel (5) = b1;
			movb1;
			tbel (6) = b1;
			}
		else 
			{
			tbel (6) = tbel (5) = b1;
			}
			/* tbel (7 & 8) contains the pointers to the second argument. */
		nel = 7;
		cl;
		}
		/*	determine the signs of the operands.	*/
	b1 = tbel (5); 
	b2 = tbel (7);
	if ((b1 == NULL) || (b2 == NULL)) ri_imp ();
	sign1 = 1;
	if (LINK_CHAR (b1))
		{
		if (b1->pair.c == '-') sign1 = -1;
		else if (b1->pair.c == '+') sign1 = 1;
		else ri_imp ();
		b1 = b1-> foll;
		}
	sign2 = 1;
	if (LINK_CHAR (b2))
		{
		if (b2->pair.c == '-') sign2 = -1;
		else if (b2->pair.c == '+') sign2 = 1;
		else ri_imp ();
		b2 = b2-> foll;
		}

		/* determine the length of both operands. */
	get_length (b1, tbel (6), &length1);
	get_length (b2, tbel (8), &length2);

		/* express the length in the number of shorts. */
	length1 *= 2;
	length2 *= 2;
	length = length1 + length2;

		/* allocate memory */
	op1 = (unsigned short *) malloc (4 * length * sizeof (unsigned short));
	if (NULL == op1) {
	  fprintf (stderr, "No memory for inside refal-machine\n");
	  exit (1);
	}
	op2 = op1 + length1;
	res = op2 + length2;
		/* in case of division res (i.e. quotient) can't be 
			longer than length1, and res2 (i.e. remainder) can't be
			longer than length2. We also need temp (at most 2*length). */
	res2 = res + length1;
	temp = res2 + length2;

		/* save both operands in an array of unsigned integers. */
	err = store_long_num (op1, b1, tbel (6), length1);
	if (err == 0) err = store_long_num (op2, b2, tbel (8), length2);
	if (err == 3) goto internal_error;
	else if (err == 2) goto call_freeze;
	else if (err == 1) ri_imp ();

		/* perform the operation. */
	err = arithm_apply (op, op1, length1, sign1, op2, length2, sign2,
		res, &length, &sign, res2, &len2, temp);
	if (err != 0) goto divby0;

		/*	RASL operators of replacement.	*/
	rdy (0); 

		/* form the result numbers. */
	if (op == MOD)
		{
		if (sign2 < 0) ns ('-');
		restore_long_num (res2, len2);
		}
	else if (op == DIVMOD)
		{
		bl;
		if (sign < 0) ns ('-');
		restore_long_num (res, length);
		br;
		if (sign2 < 0) ns ('-');
		restore_long_num (res2, len2);
		}
	else if (op == COMPARE)
		{
                ns( ((char)sign) );
                }
	else 
		{
		if (sign < 0) ns ('-');
		restore_long_num (res, length);
		}
	free ((char *) op1);
	out (2); 
	est;
	return 0;
	
	internal_error:
		fprintf (stderr, "internal error in rf_arithm\n");
		exit (1);

	divby0:
		fprintf (stderr, "attempt to divide by 0.\n");
	restart:
		ri_imp ();
		return 1;

	call_freeze:
		ri_frz (2);
		return 0;
	}

	/* apply arithmetic operation OP to operands OP1 and OP2 
		yielding result RES */
int arithm_apply (op, op1, l1, s1, op2, l2, s2, res, l, s, res2, len2, temp)

	int op;
	unsigned short *op1;
	int l1;
	int s1;
	unsigned short *op2;
	int l2;
	int s2;
	unsigned short *res;
	int *l;
	int *s;
	unsigned short *res2;
	int *len2;
	unsigned short *temp;

	{
	int err;

	switch (op)
		{
		case ADD:
			if (s1 == s2)
				{
				ar_add (op1, l1, op2, l2, res, l);
				*s = s1;
				}
			else
				{
				err = compare_num (op1, l1, op2, l2);
				if (err > 0)
					{
					ar_sub (op1, l1, op2, l2, res, l);
					*s = s1;
					}
				else if (err == 0)
					{
					*s = 1;
					res [0] = 0;
					*l = 1;
					}
				else
					{
					ar_sub (op2, l2, op1, l1, res, l);
					*s = s2;
					}
				}
			break;

		case SUB:
			if (s1 == s2)
				{
				err = compare_num (op1, l1, op2, l2);
				if (err > 0)
					{
					ar_sub (op1, l1, op2, l2, res, l);
					*s = s1;
					}
				else if (err == 0)
					{
					*s = 1;
					res [0] = 0;
					*l = 1;
					}
				else
					{
					ar_sub (op2, l2, op1, l1, res, l);
					*s = - s1;
					}
				}
			else
				{
				ar_add (op1, l1, op2, l2, res, l);
				*s = s1;
				}
			break;

		case COMPARE:
			if (s1 == s2)
				{
				err = compare_num (op1, l1, op2, l2);
				if (err > 0) 
                                     *s = (s1 >= 0) ? '+' : '-';
				else if (err < 0)
                                     *s = (s1 < 0) ? '+' : '-';
				else *s = '0';
				}
			else if (s1 > s2) *s = '+';
			else *s = '-';
			break;

		case MUL:

			*s = s1 * s2;
			ar_mul (op1, l1, op2, l2, res, l);
			break;

		case MOD:
		case DIVMOD:
		case DIV:

			*s = s1 * s2;
			err = ar_div (op1, l1, op2, l2, res, l, res2, len2, temp);
			if (err) return err;
			break;

		default:
			fprintf (stderr,
				"Internal error in arithm_apply (). operation = %d\n", op);
			exit (1);
		}

		/* check that if the result is zero then sign is positive. */
	if (*l == 1 && res[0] == 0) *s = 1;

	return 0;
	}

	/* multiply [op1, l1] by [op2, l2]. */

int ar_mul (op1, l1, op2, l2, res, l)

	unsigned short op1 [];
	int l1;
	unsigned short op2 [];
	int l2;
	unsigned short res [];
	int *l;

	{
	/*register*/ int i, j;
	int lb, ub;
	unsigned long z, x, k;

	*l = l1 + l2 - 1;
	k = 0L;
	for (j = 0; j < *l; j ++)
		{
		x = k % MAX_SHORT;	/* carry from previous product. */ /* Nemytykh 12.11.2003 */
		k /= MAX_SHORT;
		lb = max (0, j - l2 + 1);
		ub = min (j, l1 - 1);

		for (i = lb; i <= ub; i ++)
			{
			z = ((unsigned long) op1 [i]) * ((unsigned long) op2 [j-i]);
			k += z / MAX_SHORT;
			x += z % MAX_SHORT;
			}

		res [j] = x % MAX_SHORT;
		k += x / MAX_SHORT;
		}

	if (k > 0L)
		{
		res [*l] = k;
		*l += 1;
		}

		/* strip off the leading zeros. */
	for (i = *l; i > 0; i --) if (res [i-1] != 0) break;
	*l = i;

	return 0;
	}


	/* add [op1, l1] and [op2, l2]. */
int ar_add (op1, l1, op2, l2, res, l)

	unsigned short op1 [];
	int l1;
	unsigned short op2 [];
	int l2;
	unsigned short res [];
	int *l;

	{
	/*register*/ int i, m;
	unsigned long k;

	m = max (l1, l2);
	k = 0L;
	for (i = 0; i < m; i ++)
		{
		if (i < l1) k += op1 [i];
		if (i < l2) k += op2 [i];
		res [i] = k % MAX_SHORT;
		k = k / MAX_SHORT;
		}

	if (k > 0) res [i ++] = k;
	*l = i;
	return 0;
	}


	/* subtract [op2, l2] from [op1, l1]. [op1, l1] must be >= [op2, l2] */
int ar_sub (op1, l1, op2, l2, res, l)

	unsigned short op1 [];
	int l1;
	unsigned short op2 [];
	int l2;
	unsigned short res [];
	int *l;

	{
	/*register*/ int i;
	int carry;
	signed long k;

	k = 0L;
	carry = 0;
	for (i = 0; i < l1; i ++)
		{
		k = op1 [i];
		if (carry) k -= 1;
		if (i < l2) k -= op2 [i];
		if (k < 0L)
			{
			k += MAX_SHORT;
			carry = 1;
			}
		else carry = 0;
		res [i] = k;
		}

		/* strip off leading zeros. */
	for (i = l1; i > 0; i --)
		if (res [i-1] != 0) break;
	*l = i;
	return 0;
	}

int ar_div (op1, l1, op2, l2, res, l, res2, len2, temp)

	unsigned short op1 [];
	int l1;
	unsigned short op2 [];
	int l2;
	unsigned short res [];
	int *l;
	unsigned short res2 [];
	int *len2;
	unsigned short *temp;

	{
	int i;

		/* strip off the trailing zeros. */
	while ((op1 [l1 - 1] == 0) && (l1 > 1)) l1 -= 1;
	while ((op2 [l2 - 1] == 0) && (l2 > 1)) l2 -= 1;
	if ((l2 == 1) && (op2 [0] == 0)) return DIV_BY_ZERO;

	i = compare_num (op1, l1, op2, l2);
	if (i < 0)
		{
		res [0] = 0;
		*l = 1;
		for (i = 0; i < l1; i ++) res2 [i] = op1 [i];
		*len2 = l1;
		return 0;
		}
	else if (i == 0)
		{
		res [0] = 1;
		*l = 1;
		res2 [0] = 0;
		*len2 = 1;
		return 0;
		}

		/* simple case. */
	if (l2 == 1)
		{
		reverse_num (op1, l1);
		divide_num (op1, (unsigned long)(MAX_SHORT), l1,
			(unsigned long)(op2 [0]), res2);
		*len2 = 1;
		for (i = l1; i > 0; i --) res [l1 - i] = op1 [i - 1];
		*l = l1;
		return 0;
		}

		/* here we have op1 / op2 s.t. op1 > op2 and 
			the highest order digits of op1 and op2 are not zeros. 
			also op2 has more than 1 digit. */

	ar_div_long (op1, l1, op2, l2, res, l, res2, len2, temp);
	return 0;
	}

	/* given a > b find q and r s.t. a = q * b + r */
int ar_div_long (a, la, b, lb, q, lq, r, lr, t)
	unsigned short *a;
	int la;
	unsigned short *b;
	int lb;
	unsigned short q[];
	int *lq;
	unsigned short r[];
	int *lr;
	unsigned short t[];	/* temporary array of length 2 * (la+lb) */

	{
	/*register*/ int i, j, k;
	int lt, lt2;
	unsigned short qn, *t2;

	t2 = t + (la + lb);
	*lq = la - lb + 1;
	for (k = la - lb; k >= 0; k --)
		{
		j = la - k;		/* the length of the divisor. */
			/* select qn - the next digit in the quotient. */
		divide_select (a+k, j, b, lb, &qn);

	qn_selected:
		if (qn == 0)
			{
			q [k] = 0;
			continue;
			}

			/* compute c = b * qn */
		ar_mul (b, lb, &qn, 1, t, &lt);
			/* compare a and c */
		if (compare_num (a+k, j, t, lt) < 0)
			{
			qn --;	/* selected qn is too large. */
			goto qn_selected;
			}

			/* compute d = a - c */
		ar_sub (a+k, j, t, lt, t2, &lt2);

			/* compare d and b */
/*+		if (compare_num (b, lb, t2, lt2) < 0)*/
		if (compare_num (b, lb, t2, lt2) <= 0)
			{
			qn ++;	/* selected qn is too small. */
			goto qn_selected;
			}

			/* otherwise copy d to a, and qn is quotient. */
		for (i = 0; i < lb; i ++)
			{
			if (i < lt2) a [k+i] = t2 [i];
			else a [k+i] = 0;
			}
		if (j > lb) {
			int i_w;

			for (i_w = lb; i_w < j; i_w ++) {
				a [k + i_w] = 0;
			}
		}

		q [k] = qn;
		}

		/* copy the rest of a into the remainder. */
	for (i = 0; i < lb; i ++) r [i] = a [i];
	*lr = lb;
	return 0;
	}


	/* select quotient based on the first two digits of a and b */
int divide_select (a, la, b, lb, q)
	
	unsigned short *a;
	int la;
	unsigned short *b;
	int lb;
	unsigned short *q;

	{
	unsigned short x [2];
	unsigned long z;
	unsigned short rem;

		/* shorten the length. */
	while ((la > 1) && (a [la - 1] == 0)) la --;

/*	*q = 0; ?*/
	if (la < lb) return 0;
	if (la == lb)
		{
		z = (unsigned long)(b [lb - 1]) * MAX_SHORT + b [lb - 2];
		}
	else if (la == lb + 1)
		{
		z = (unsigned long) b [lb - 1];
		}
	else
		{
		printf ("Error: In divide_select (): la = %d lb = %d\n", la, lb);
		exit (1);
		}

	x [0] = a [la - 1];
	x [1] = a [la - 2];
	divide_num (x, (unsigned long) MAX_SHORT, 2, z, &rem);
/*	*q = x [1]; */
	*q = x [0] ? x[0] : x[1];
	return 0;
	}

/* Nemytykh: 11.02.04 */
# define START      1 
# define START_ZERO 2 
# define NONZERO    3 
unsigned long random_macro_digit (max) 

     unsigned long max;

{
 double f = rand();
 unsigned long m = (unsigned long)(max * (f/RAND_MAX));

/* This very strange manipulation with 1 is recomended by the C-manual.
   A bug in rand() causes that.
* 
 return (unsigned long)((double)max * (f/(RAND_MAX + 1.0)));
**/

/* In fact, that must be: */
 return m < max ? (m + rand()%2) : m;
/**/
}

# define NUMBER_OF_DIGITS  0 
# define INSIDE_DIGIT      1 

int rf_random ()    { return rf_rand(NUMBER_OF_DIGITS); }
int rf_randomDigit () { return rf_rand(INSIDE_DIGIT); }

int rf_rand (scale) 

     char scale;

	{
	char c;
	unsigned int length;
	unsigned long n;
        int st;

	check_frz_args (1);

	cl;
	if (LINK_CHAR (tbel (3)))
		{
		c = tbel (3)->pair.c;
		if (c != '+') 
			{
			ri_error(8);
			return 1;
			}
		b1 = tbel (3)->foll;
		}
	else b1 = tbel (3);
	if (LINK_VAR (b1)) return 2;
	else if (!LINK_NUMBER (b1)) return 2;/*1;*/
	rdy (0); 
        st = START;
	
        length = (unsigned int)random_macro_digit((unsigned long)(b1->pair.n));
        if( scale == INSIDE_DIGIT ) { nns(length); }
        else  
               {
                if( length == 0 ) length = 1;
	        for ( ; length > 0; length --) 
                       {
                        n = random_macro_digit(MAX_UNSIGNED_INT);
        		/* skip the leading zeros. */
                        if ( st != START_ZERO ) nns(n);
                        if ( (n == 0) && (st == START) ) 
                              { st = START_ZERO; }
                        else  { st = NONZERO;    };
                      };
               };
	out (2);
	est;
	return 0;

	internal_error:
		fprintf (stderr, "internal error in rf_random ()\n");
		exit (1);

	call_freeze:
		ri_frz (2);
		return 0;
	}


 