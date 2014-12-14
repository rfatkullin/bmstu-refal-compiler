#include <stdio.h>
#include <string.h>

#include "builtins.h"

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

struct func_result_t prout(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView)
{
	struct l_term* currExpr = fieldOfView->current->begin;

	while (currExpr != 0)
	{
		if (currExpr->tag == L_TERM_FRAGMENT_TAG)
		{
			printRange(currExpr->fragment);
		}
		else if (currExpr->tag == L_TERM_CHAIN_TAG)
		{
			printf("[Error] !!!\n");
		}

		currExpr = currExpr->next;
	}

	struct l_term_chain_t* mainChain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));
	mainChain->begin = 0;
	mainChain->end = 0;

	return (struct func_result_t){.status = OK_RESULT, .mainChain = mainChain, .callChain = 0};
}

static void printRange(struct fragment* frag)
{
	int i = 0;
	struct v_term* currTerm = memMngr.activeTermsHeap + frag->offset;

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
	case V_INT_NUM_TAG:
		printf("%d ", term->intNum);
		break;
	case V_CLOSURE_TAG:
		//TO DO
			break;
	case V_BRACKET:
		printf("%c", term->inBracketLength > 0 ? '(' : ')');
		break;
	}
}
