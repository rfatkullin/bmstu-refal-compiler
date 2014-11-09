/* flag for built in functions: requires the pointer to the function table. */
# define BI_FADDR 1

/* structure for built in functions: */
struct bitab {
  char *fname;		/* name */
  int fnumber;		/* internal number. */
  int flags;		/* flags: */
};

