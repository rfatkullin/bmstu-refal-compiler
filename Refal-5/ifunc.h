
/* refgo.c trace.c -- Main function for Refal interpreter and Refal Tracer. */
/*int main (int, char **);*/

/* ri.c -- Refal interpreter. */
int ri_inter (void);
int ri_frz (int);
FILE *ref_err_file(void);

/* macros.c  -- Definitions of some RASL operators. */
void ri_act1 (char *);
int ri_blr (void);
void ri_cl (void);
int ri_est (void);
int ri_mule (int);
int ri_cpelt (struct link *);
int ri_ns1 (void);
void ri_tple (int);
void ri_out (int);
int ri_tpls (int);
int ri_rimp (void);
int ri_oexp (int);
int ri_oexpr (int);
int ri_ovs (int);
int ri_ovsr (int);
int ri_lens (char);
int ri_lenp (void);

/* load.c -- deals with the initial loading of the .RSL file into core. */
char *ri_load (void);
# ifdef LOAD_INCLUDED
char ** ri_mkfunlist (struct module_descriptor *, char *, char **);
int ri_mkmodlist (struct module_descriptor *);
int ri_mkentlist (struct module_descriptor *);
struct module_descriptor *ri_freemod (struct module_descriptor *);
struct module_descriptor *ri_readhdr (FILE *, struct module_descriptor *);
# endif
int ri_lerror (int);
# ifdef LOAD_INCLUDED
int ri_loadcode (struct module_descriptor *, struct module_descriptor *);
# endif
int ri_init (void);
int ri_memory (void);
int ri_inquire (char *, char *, int);

/* refaux.c -- auxiliary: memory allocation, stop, error, comp symbols etc. */
struct link *ri_getmem (void);
struct link *ri_fmout (void);
int ri_stop (void);
int ri_imp (void);
int ri_error (int);
char *ri_cs_impl (char *);
char *ri_cs_exist (int, char *);
int ri_cs_ins (int, char *);
int ri_hash (char *);
int ri_print_error_code(FILE *, int);
int ri_common_stack (void);
int ri_init_stop(void);
int ri_options(int, char **, char);

/* bif.c -- list of standard fuinctions. */
int ri_bif (long, char *);

/* refio.c -- Refal I/O */
int ri_putmb (struct link *, struct link *, FILE *);
int ri_addsym (struct link *, FILE *);
int ri_add_char (char, FILE *);
int ri_add_string (int, char *, int, FILE *);
int putstr (FILE *, int, char *);
int rc_print_buff (int, FILE *);
int ri_put (FILE *);
int ri_cs_len (char *);
int ri_cs_rest_len (char *);
int ri_new_line (char);
int ri_put_link (FILE *, struct link *, int);
int ri_get (FILE *);
int ri_open (int, char *, char *);
int ri_actput (char *, FILE *);
int isIdentifier (char *);

/* freeze.c -- Freezer, downgrading, upgrading etc. */
int is_freezer (struct link *);
int exists_freeze (void);
int contains_vars (struct link *, struct link *);
int rf_frz (void);
int ri_frz1 (int);
int rf_dn (void);
int ri_dn (struct link *, struct link *);
int rf_setfrz (char **);
int rf_up (char **);
int ri_up (struct link *, struct link *, char **);
int rf_frzr (void);

/* func1.c -- definitions of some standard functions */
int rf_prout (void);
int rf_put (void);
int rf_br (void);
int rf_putout (char);
int rf_card (void); 
int rf_print (void);
int rf_type (void); 
int rf_dg (void);
int rf_get (void);
int rf_arg (void);
char *mu_find (char *, char **);
int rf_mu (char **, int);
int mu_special (char *, char);
int rf_getenv (void); 
int rf_getpid (void); 
int rf_getppid (void); 
int rf_system (void); 
int rf_exit (void); 
int rf_close (void);
int rf_existfile (void);
int rf_getcurrentdirectory (void);
int rf_removefile (void);
int rf_implodeExt (void);
int rf_explodeExt (void);
int rf_timeExt (void);
int rf_compare (void);
int rf_XMLParse (void);

/* func2.c -- more definitions of standard functions */
int rf_cp (void);
int rf_chr (void);
int rf_ord (void);
int rf_last (void);	
int rf_first (void);
int rf_implode (void);
int rf_explode (void);
int rf_lower (void);
int rf_upper (void);
int rf_step (void);
int rf_time (void);
int rf_lenw (void);
int rf_dgall (void);
int rf_rp (void);
int rf_open (void);

/* net.c --- more definitions of standard functions (all for network). Shura. 04.08.99 */
int rf_socket (void);
int rf_connect (void);
int rf_send (void);
int rf_recv (void);
void ri_socket_init (void);

/* arithm.c -- definitions of arithmetic and related functions. */
int rf_symb (void); 
int rf_numb (void);	
int get_length (struct link *, struct link *, int *);
int reverse_num (unsigned short *, int);
int divide_num (unsigned short *, unsigned long, int, unsigned long,
	unsigned short *);
int is_zero_num (unsigned short *, int);
int convert_num (unsigned short *, unsigned long, int, unsigned short *,
	unsigned long, int *);
int store_long_num (unsigned short *, struct link *, struct link *, int);
int restore_long_num (unsigned short *, int);
int compare_num (unsigned short *, int, unsigned short *, int);
int rf_arithm (int);
int arithm_apply (int, unsigned short *, int, int, unsigned short *, int,
	int, unsigned short *, int *, int *, unsigned short *, int *,
	unsigned short *);
int ar_mul (unsigned short *, int, unsigned short *, int,
	unsigned short *, int *);
int ar_add (unsigned short *, int, unsigned short *, int,
	unsigned short *, int *);
int ar_sub (unsigned short *, int, unsigned short *, int,
	unsigned short *, int *);
int ar_div (unsigned short *, int, unsigned short *, int,
	unsigned short *, int *, unsigned short *, int *, unsigned short *);
int ar_div_long (unsigned short *, int, unsigned short *, int,
	unsigned short *, int *, unsigned short *, int *, unsigned short *);
int divide_select (unsigned short *, int, unsigned short *, int,
	unsigned short *);
int rf_random (void);	
int rf_randomDigit (void);	
int rf_rand (char);
unsigned long random_macro_digit (unsigned long);

/* sysfun.c -- user defined standard functions. */
int rf_sysfun (void);
int getexp (FILE *);
int rf_desysfun (void);
int getXML (FILE *);
int rf_listOfBuiltin (void);
int rf_sizeof (void);

/* xml2ref.c -- */
int ri_xml2ref (FILE * , FILE *);

