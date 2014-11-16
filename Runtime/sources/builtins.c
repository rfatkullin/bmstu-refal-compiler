#include <stdio.h>
#include <string.h>

#include "l_term.h"
#include "builtins.h"
#include "memory_manager.h"

#define N 256

static void printSymbol(struct v_term* term);
static void printRange(struct fragment* frag);

struct l_term* card(struct l_term* expr)
{
	static char buff[N];

	if (fgets(buff, N, stdin) != NULL)
	{
		int len = strlen(buff) - 1;
		return allocateVector(len, buff);
	}
	else
	{
		return NULL;
	}
}

struct l_term* prout(struct l_term* expr)
{
	struct l_term* currTerm = expr;

	while (currTerm != NULL)
	{
		if (expr->tag == L_TERM_RANGE_TAG)
		{
			printRange(expr->range);
		}
		else
		{
			printf("(");
			prout(expr->chain);
			printf(")");
		}
	}

	return NULL;
}

static void printRange(struct fragment* frag)
{
	int i = 0;
	struct v_term* currTerm = memoryManager.activeTermsHeap + frag->offset;

	for (i = 0; i < frag->length; ++i)
		printSymbol(currTerm + i);
}

static void printSymbol(struct v_term* term)
{
	switch (term->tag)
	{
	case V_CHAR_TAG:
		printf("%c ", term->str[0]);
		break;
	case V_IDENT_TAG:
		printf("%s ", term->str);
		break;
	case V_NUMBER_TAG:
		printf("%d ", term->num);
		break;
	case V_CLOSURE_TAG:
		//TO DO
			break;
	case V_BRACKET:
		printf("%c", term->inBracketLength > 0 ? '(' : ')');
		break;
	}
}
