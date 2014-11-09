
		/* macros for freezer variables. */

	/* some macros which deal with levels of variables. */
# define MAX_VAR_INDEX  0x0000ffffL /* 2^16-1 actually: maximum index. */
# define MAX_VAR_ELEV   0xffL       /* 2^8-1 (or 255): maximum elevation. */
# define MAX_VAR_LEVEL  0xffL       /* 2^8-1 (or 255): maximum level. */
	/* encoding of index, elevation and level of a variable:
		lower 16 bits: index,
		middle 8 bits: elevation,
		high 8 bits: level. */
# define INDEX_MASK 0x0000ffffL
# define ELEV_MASK 0x00ff0000L
# define LEVEL_MASK 0xff000000L
# define index_of(var) ((unsigned long)(INDEX_MASK & (var)))
# define elevation_of(var) ((unsigned long)((ELEV_MASK & (var)) >> 16))
# define level_of(var) ((unsigned long)((LEVEL_MASK & (var)) >> 24))
# define make_level(level) (((unsigned long) (level) << 24) & LEVEL_MASK)
# define change_var_level(field, change) (index_of(field) | make_level(level_of(field)+(change)))


	/* metacode definition. */
#define META_ACTIVE          '!'
#define META_BRACKET         '*'
#define META_FIRST_BRACKET   '-'
#define META_SVAR            's'
#define META_EVAR            'e'
#define META_TVAR            't'
#define META_QUOTE           'm'
#define META_QUOTE2          'q'

