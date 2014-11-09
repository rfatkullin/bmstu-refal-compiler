/* trace.c  -- main function for the tracer. */
/* rti.c -- Refal interpreter (tracer version) */
int ri_default (int, int *);

/* dtrace.c -- Tracer command interpreter */
int rd_trace (void);
int rd_toggle (void);
int rd_getins (void);
int rd_sep (char);
int rd_getcom (int);
int rd_get_num_break (void);
int rd_set_break (void);
int rd_display (void);
int rd_look_for_var (char ,char *, int);
int rd_del_brk (void);
int rd_delete_break (int);
int rd_show (void);
int rd_shwstack (FILE *);
int rd_shwbreak (int);
int rd_shwall (void);
int rd_help (void);
int rd_cplink (struct link *);
struct link *rd_cp_refx (struct link *, struct link *);
int rd_step (void);
int rd_compute (void);
struct link *rd_adjust_re (struct link *, struct link *);
int rd_print_res (void);
int rd_endcom (void);
int rd_init (void);
/*char *strnchr (char *, char, int);*/

/* dcom.c -- Tracer compiler: compiles pattern expression only. */
/* redefined versions of rc_gchar () and rc_ungchar() are also here. */
int rd_parse(char *, int, char *, char **);
/* int rd_vyvod (struct rasl_instruction *, char *, int, char **);*/
int rd_cr_lvtab (int);

/* dimpl.c -- functions/module handling */
struct listfun *rd_implode(char *, int *);
int rd_create_ht (void);
int rd_ins_mod (char *, char **);
int rd_ins_ht (char *, char *);
int rd_pr_module (struct modlist *);
int rd_shwmod (char *);
int rd_shmodall (void);
int rd_list_modules (void);
int rd_shwfunc (char *);

