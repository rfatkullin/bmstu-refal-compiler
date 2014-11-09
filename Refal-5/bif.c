
# include "decl.h"
# include "arithm.h"
# include "ifunc.h"
# include "macros.h"

extern int rf_xxx (void);

int ri_bif (n, arg)
	long n;
	char *arg;
	{
		int k;

		k = (int) n;
		switch (k)
			{
			case 1: return rf_mu ((char **) arg, 0);
			case 2: return rf_arithm (ADD);
			case 3: return rf_arg ();
			case 4: return rf_br ();
			case 5: return rf_card ();
			case 6: return rf_chr ();
			case 7: return rf_cp ();
			case 8: return rf_dg ();
			case 9: return rf_dgall ();
			case 10: return rf_arithm (DIV);
			case 11: return rf_arithm (DIVMOD);
			case 12: return rf_explode ();
			case 13: return rf_first ();
			case 14: return rf_get ();
			case 15: return rf_implode ();
			case 16: return rf_last ();
			case 17: return rf_lenw ();
			case 18: return rf_lower ();
			case 19: return rf_arithm (MOD);
			case 20: return rf_arithm (MUL);
			case 21: return rf_numb ();
			case 22: return rf_open ();
			case 23: return rf_ord ();
			case 24: return rf_print ();
			case 25: return rf_prout ();
			case 26: return rf_put ();
			case 27: return rf_putout (YES);
			case 28: return rf_rp ();
			case 29: return rf_step ();
			case 30: return rf_arithm (SUB);
			case 31: return rf_symb ();
			case 32: return rf_time ();
			case 33: return rf_type ();
			case 34: return rf_upper ();
			case 35: return rf_sysfun ();

	/* If you change these 2 values then make a change refgo.c and trace.c */
			case 42: return ri_imp ();
			case 43: return ri_stop ();

			case 45: return rf_frz ();
			case 46: return rf_frzr ();
			case 47: return rf_dn ();
			case 48: return rf_up ((char **) arg);
			case 49: return rf_setfrz ((char **) arg);
			case 50: return rf_mu ((char **) arg, 1);

			/* Process functions */
			case 51: return rf_getenv ();
			case 52: return rf_system ();
			case 53: return rf_exit ();
			case 54: return rf_close ();
			case 55: return rf_existfile ();
			case 56: return rf_getcurrentdirectory ();
			case 57: return rf_removefile ();
			case 58: return rf_implodeExt ();
			case 59: return rf_explodeExt ();

			case 60: return rf_timeExt ();
			case 61: return rf_arithm (COMPARE);
			case 62: return rf_desysfun ();
			case 63: return rf_XMLParse ();
			case 64: return rf_random ();
			case 65: return rf_randomDigit ();
			case 66: return rf_putout (NO);

			/* System functions */
			case 67: return rf_listOfBuiltin();
			case 68: return rf_sizeof();

			case 69: return rf_getpid();

			/* System function for graph compiler */
			case 70: return rf_xxx ();

			case 71: return rf_getppid();

			default: return 1;
			};
	}

