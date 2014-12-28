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
			struct lterm_t** helper = (struct lterm_t**)malloc(7 * sizeof(struct lterm_t*));
			int i;
			for (i = 0; i < 7; ++i)
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
			currTerm->fragment->length = 2;
			helper[3]->chain->begin = currTerm;
			helper[3]->chain->end = currTerm;
			helper[3]->tag = L_TERM_CHAIN_TAG;
			helper[2]->chain->begin = helper[3];
			helper[2]->chain->end = helper[3];
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
			currTerm->fragment->offset = 3;
			currTerm->fragment->length = 4;
			helper[4]->chain->begin = currTerm;
			helper[4]->chain->end = currTerm;
			helper[4]->tag = L_TERM_CHAIN_TAG;
			helper[2]->chain->end->next = helper[4];
			helper[4]->prev = helper[2]->chain->end;
			helper[2]->chain->end = helper[4];
			helper[2]->tag = L_TERM_CHAIN_TAG;
			helper[1]->chain->end->next = helper[2];
			helper[2]->prev = helper[1]->chain->end;
			helper[1]->chain->end = helper[2];
			helper[5]->tag = L_TERM_CHAIN_TAG;
			helper[1]->chain->end->next = helper[5];
			helper[5]->prev = helper[1]->chain->end;
			helper[1]->chain->end = helper[5];
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
			currTerm->fragment->offset = 7;
			currTerm->fragment->length = 2;
			helper[6]->chain->begin = currTerm;
			helper[6]->chain->end = currTerm;
			helper[6]->tag = L_TERM_CHAIN_TAG;
			helper[1]->chain->end->next = helper[6];
			helper[6]->prev = helper[1]->chain->end;
			helper[1]->chain->end = helper[6];
			funcTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			funcTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));
			funcTerm->funcCall->funcName = memMngr.termsHeap[helper[1]->chain->begin->fragment->offset].str;
			funcTerm->funcCall->funcPtr = Prout;
			funcTerm->funcCall->entryPoint = 0;
			funcTerm->funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));
			funcTerm->funcCall->fieldOfView->current = helper[1]->chain;
			funcTerm->tag = L_TERM_FUNC_CALL;
			funcCallChain->begin = funcTerm;
			funcCallChain->end = funcTerm;
			helper[0]->chain->begin = helper[1];
			helper[0]->chain->end = helper[1];
			for (i = 0; i < 7; ++i)
			{
				helper[i]->chain->begin->prev = 0;
				helper[i]->chain->end->next = 0;
			}
			funcCallChain->begin->prev = 0;
			funcCallChain->end->next = 0;
			funcRes = (struct func_result_t){.status = OK_RESULT, .mainChain = helper[0]->chain, .callChain = funcCallChain};
			break;
		}
	} // switch block end
	if (funcRes.status != CALL_RESULT)
	{
		free(env->locals);
		free(fieldOfView->backups);
	}
	return funcRes;
} // Go

int main()
{
	__initLiteralData();
	mainLoop(Go);
	return 0;
}
