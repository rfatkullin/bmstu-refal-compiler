#include <stdio.h>
#include <string.h>

#include "builtins.h"

#define N 256

static void printRange(struct fragment* frag);

struct func_result_t Card(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView)
{
	char ch;
	uint32_t lastOffset;
	uint32_t firstOffset;
	struct l_term_chain_t* mainChain;


	if((ch = getchar()) != EOF)
	{
		firstOffset = lastOffset = allocateSymbol(ch);

		while ((ch = getchar()) != EOF)
		{
			lastOffset = allocateSymbol(ch);
		}

		mainChain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));
		mainChain->begin = (struct l_term*)malloc(sizeof(struct l_term));
		mainChain->end = mainChain->begin;

		mainChain->begin->tag = L_TERM_FRAGMENT_TAG;
		mainChain->begin->fragment = (struct fragment*)malloc(sizeof(struct fragment));
		mainChain->begin->fragment->offset = firstOffset;
		mainChain->begin->fragment->length = lastOffset - firstOffset;
	}
	else
	{
		mainChain = 0;
	}


	return (struct func_result_t){.status = OK_RESULT, .mainChain = 0, .callChain = 0};
}

struct func_result_t Prout(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView)
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

	return (struct func_result_t){.status = OK_RESULT, .mainChain = 0, .callChain = 0};
}

static void printRange(struct fragment* frag)
{
	int i = 0;
	struct v_term* currTerm = memMngr.activeTermsHeap + frag->offset;

	for (i = 0; i < frag->length; ++i)
		printSymbol(currTerm + i);
}
