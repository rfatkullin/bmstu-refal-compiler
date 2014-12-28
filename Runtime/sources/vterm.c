#include <stdio.h>
#include "vterm.h"

int printSymbol(struct v_term* term)
{
	int res = 1;

	switch (term->tag)
	{
	case V_CHAR_TAG:
		printf("%c", term->ch);
		break;
	case V_IDENT_TAG:
		printf("%s", term->str);
		break;
	case V_INT_NUM_TAG:
		printf("%d", term->intNum);
		break;
	case V_CLOSURE_TAG:
		//TO DO
		break;
	case V_BRACKET_TAG:
		printf("%c", term->inBracketLength > 0 ? '(' : ')');
		res = 0;
		break;
	}

	return res;
}
