#include <stdio.h>
#include <stdlib.h>

#include "builtins.h"

#define N 256

static void printRange(struct fragment_t* frag);

struct func_result_t Card(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView)
{
	char ch;
	uint32_t lastOffset;
	uint32_t firstOffset;
	struct lterm_chain_t* mainChain;

	if((ch = getchar()) != '\n')
	{
		firstOffset = lastOffset = allocateSymbol(ch);

		while ((ch = getchar()) != '\n')
		{
			lastOffset = allocateSymbol(ch);
		}

		mainChain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));
		mainChain->begin = (struct lterm_t*)malloc(sizeof(struct lterm_t));
		mainChain->end = mainChain->begin;
		mainChain->end->next = 0;

		mainChain->begin->tag = L_TERM_FRAGMENT_TAG;
		mainChain->begin->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
		mainChain->begin->fragment->offset = firstOffset;
		mainChain->begin->fragment->length = lastOffset - firstOffset + 1;
	}
	else
	{
		mainChain = 0;
	}


	return (struct func_result_t){.status = OK_RESULT, .fieldChain = mainChain, .callChain = 0};
}

struct func_result_t Prout(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView)
{
	struct lterm_t* currExpr = fieldOfView->current->begin;

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

	struct lterm_chain_t* mainChain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));
	mainChain->begin = 0;
	mainChain->end = 0;

	return (struct func_result_t){.status = OK_RESULT, .fieldChain = 0, .callChain = 0};
}

static void printRange(struct fragment_t* frag)
{
	int i = 0;
	struct v_term* currTerm = memMngr.vterms + frag->offset;

	printf("Range: ");
	for (i = 0; i < frag->length; ++i)
	{
		if(printSymbol(currTerm + i))
			printf(" ");
	}

	printf("\n");
}
