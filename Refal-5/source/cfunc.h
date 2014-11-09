
/* files marked with # are used in the Tracer as well. */

/* main.c -- main function. */
int main (int, char * *);
int ri_inquire (char *, char *, int);
int rc_ungchar (char);
int rc_gchar (void);
	/* debugging in main.c */
int print_expr (struct element *);
int print_tree (struct node *);
char *rasl_code (int);
int print_translation (struct rasl_instruction *);
int print_rasl_inst (struct rasl_instruction *);
int print_var_table (unsigned char *);
int prftab (struct functab *);
int print_holes (struct HOLES *, struct element *);

/* lex.c# -- lexical scan functions. */
int rc_gettoken (void);
int rc_getsym (void);
int rc_getind (void);
int rc_getid (int);
int rc_serror (int ,char *);
int rc_swarn (int);
int rc_getstr (int);
int rc_getnumb (int);

/* parser.c -- top down recursive descent parser */
int rc_skip (int);
char *rc_fname (void);
struct node *rc_r_first (void);
struct node *rc_r_list (void);
struct node *rc_l_first (void);
struct node *rc_l_list (void);
struct node *rc_r_tail (void);
struct node *rc_sentence (void);
struct node *rc_sents (void);
int rc_1ofsent (int);
char *rc_idlist (void);
struct node *rc_fdef (void);
struct node *rc_parser (void);
struct element *rc_r_side (void);
struct element *rc_l_side (void);

/* pass2.c -- second pass resolving references. */
int rc_initcom (int, char * *);
int rc_getbeginfile (char **);
int rc_getnextfile (int, char **);
int rc_sbtable (struct functab *);
int rc_sftable (struct functab *);
int rc_sltable (struct functab *);
int rc_pass2 (void);
int rc_end (void);

/* pcbi.c -- list of standard functions. */
long rc_binumber (char *);

/* rc.c -- Refal compiler: general and right side translation. */
int refcom (struct node *);
int tr_rtail (struct node *, int *);
int refc_out (struct rasl_instruction * *, int);
int transl_right (struct element *, int);
int get_var_index (int, unsigned char *);
int is_bit_checked (unsigned char *, int);
int check_bit (unsigned char *, int);

/* rcaux.c# -- expression parsing and auxiliary functions. */
int rc_initrp (void);
char *rc_allmem (int);
char *rc_memral (int);
int rc_post_opt (struct rasl_instruction * *, int);
int rc_post (struct rasl_instruction *);
int merge_string_instr (int, int, struct rasl_instruction *);
struct element *refal_expression (int);
char *rc_getact (void);
char *getcoms (void);
int rc_getvar (int, char, struct element *);
int getconst (int, struct element *);
int searchv (char, char *);
struct functab *searchf (char *, struct functab *);
int insert_instruction (struct rasl_instruction *, int, union param);
int rc_options(int argc, char * argv []);
int rc_help(void);



/* rcleft.c# -- Refal compiler: left parts compilation */
int transl_left (struct element *, int *);
int match (struct element *, int, int, int *, int *, int, int *);
struct HOLES *delete_hole (struct HOLES *, struct HOLES *);
struct HOLES *select_hole (struct HOLES *, struct element *, int *, struct HOLES **);
struct HOLES *add_hole (struct HOLES *, struct HOLES *, int, int, int, int);
int free_holes (struct HOLES *);
int check_var (int);
int no_lengthening (struct element *);
int rc_out (int, union param);

/* rcopt.c -- Refal optimizer */
int refc_opt (struct rasl_instruction * *, int, int *);
int comp_inst (struct rasl_instruction *, struct rasl_instruction *);
int comp_sym_inst (struct rasl_instruction *, struct rasl_instruction *);
int split_string (struct rasl_instruction *, int, int, int);
int ch2sym (struct rasl_instruction *);
int chop_tail (struct rasl_instruction *);
int is_left_part (struct rasl_instruction *);
int comp_next_ins (struct rasl_instruction *, struct rasl_instruction *, 
	struct rasl_instruction **);
int delete_to_ptr (struct rasl_instruction *, struct rasl_instruction *);
int save_label (int, struct rasl_instruction *);
struct rasl_instruction *find_def_label (int);

/* sem.c -- some semantic processing during parsing. */
struct node *rc_mknode (int, branch_t, branch_t, branch_t);
char *rc_deffn (void);
char *rc_mkextrn (char *);
char *rc_mkentry (char *);
int free_tree (struct node *);

/* vyvod.c -- saving Refal translation in a file. */
char *rc_deflabel(char *, struct functab *);
long rc_getcsn (char *);
int rc_vyvod (struct rasl_instruction * *, int);

#ifdef FOR_OS_LINUX
#	ifndef FOR_OS_OS2
/* This code is inserted by Shura. 22.01.98 */
/* strupr is in MS-DOS, but none in UNIX. So, the code depent from
 * version of UNIX. Here is implemented only Latin charcters.
 */
static inline void
strupr (char * cp_str) {
  int i;

  for (i = 0; '\0' != cp_str [i]; i ++) {
    switch (cp_str [i]) {
    case 'a': cp_str [i] = 'A'; break;
    case 'b': cp_str [i] = 'B'; break;
    case 'c': cp_str [i] = 'C'; break;
    case 'd': cp_str [i] = 'D'; break;
    case 'e': cp_str [i] = 'E'; break;
    case 'f': cp_str [i] = 'F'; break;
    case 'g': cp_str [i] = 'G'; break;
    case 'h': cp_str [i] = 'H'; break;
    case 'i': cp_str [i] = 'I'; break;
    case 'j': cp_str [i] = 'J'; break;
    case 'k': cp_str [i] = 'K'; break;
    case 'l': cp_str [i] = 'L'; break;
    case 'm': cp_str [i] = 'M'; break;
    case 'n': cp_str [i] = 'N'; break;
    case 'o': cp_str [i] = 'O'; break;
    case 'p': cp_str [i] = 'P'; break;
    case 'q': cp_str [i] = 'Q'; break;
    case 'r': cp_str [i] = 'R'; break;
    case 's': cp_str [i] = 'S'; break;
    case 't': cp_str [i] = 'T'; break;
    case 'v': cp_str [i] = 'V'; break;
    case 'u': cp_str [i] = 'U'; break;
    case 'w': cp_str [i] = 'W'; break;
    case 'x': cp_str [i] = 'X'; break;
    case 'y': cp_str [i] = 'Y'; break;
    case 'z': cp_str [i] = 'Z'; break;
    default:;
    }
  }
}
#	endif
#endif

