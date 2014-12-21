#include <stdio.h>
#include "v_term.h"

void printSymbol(struct v_term* term)
{
	switch (term->tag)
	{
	case V_CHAR_TAG:
		printf("%c ", term->str[0]);
		break;
	case V_IDENT_TAG:
		printf("%s ", term->str);
		break;
	case V_INT_NUM_TAG:
		printf("%d ", term->intNum);
		break;
	case V_CLOSURE_TAG:
		//TO DO
			break;
	case V_BRACKET:
		printf("%c ", term->inBracketLength > 0 ? '(' : ')');
		break;
	}
}
