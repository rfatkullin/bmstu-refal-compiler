// file:../Compiler-build/../Compiler-build/test_1.ref

#include <stdlib.h>

#include <memory_manager.h>
#include <vmachine.h>
#include <builtins.h>

void __initLiteralData()
{
	initAllocator(1024 * 1024 * 1024);
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "Prout"};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_CHAR_TAG, .ch = 'c'};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "asdasd"};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "sdas"};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 1};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "Ha-ha-ha"};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 2};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 3};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 12};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 3};

	initHeaps(2);
} // __initLiteralData()

struct func_result_t Go(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView) 
{
	struct func_result_t funcRes;
	if (entryPoint == 0)
	{
		env->locals = (struct lterm_t*)malloc(1 * sizeof(struct lterm_t));
		fieldOfView->backups = (struct lterm_chain_t*)malloc(1 * sizeof(struct lterm_chain_t));
	}
	switch (entryPoint)
	{
		case 0: 
		{
			struct lterm_chain_t* funcCallChain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));
			funcCallChain->begin = 0;
			funcCallChain->end = 0;
			struct lterm_t* funcTerm;
			struct lterm_t** helper = (struct lterm_t**)malloc(6 * sizeof(struct lterm_t*));
			int i;
			for (i = 0; i < 6; ++i)
			{
				helper[i] = (struct lterm_t*)malloc(sizeof(struct lterm_t));
				helper[i]->chain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));
			}
			struct lterm_t* currTerm = 0;
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
			currTerm->fragment->offset = 0;
			currTerm->fragment->length = 1;
			helper[1]->chain->begin = currTerm;
			helper[1]->chain->end = currTerm;
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
			currTerm->fragment->offset = 1;
			currTerm->fragment->length = 3;
			helper[3]->chain->begin = currTerm;
			helper[3]->chain->end = currTerm;
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
			currTerm->fragment->offset = 4;
			currTerm->fragment->length = 4;
			helper[4]->chain->begin = currTerm;
			helper[4]->chain->end = currTerm;
			helper[4]->tag = L_TERM_CHAIN_TAG;
			helper[3]->chain->end->next = helper[4];
			helper[4]->prev = helper[3]->chain->end;
			helper[3]->chain->end = helper[4];
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
			currTerm->fragment->offset = 8;
			currTerm->fragment->length = 2;
			helper[5]->chain->begin = currTerm;
			helper[5]->chain->end = currTerm;
			helper[5]->tag = L_TERM_CHAIN_TAG;
			helper[3]->chain->end->next = helper[5];
			helper[5]->prev = helper[3]->chain->end;
			helper[3]->chain->end = helper[5];
