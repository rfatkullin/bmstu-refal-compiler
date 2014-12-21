// file:../Compiler-build/../Compiler-build/test_1.ref

#include <stdlib.h>

#include <memory_manager.h>
#include <v_machine.h>
#include <builtins.h>

struct func_result_t Go(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView) 
{
	struct func_result_t funcRes;
	if (entryPoint == 0)
	{
		env->locals = (struct l_term*)malloc(1 * sizeof(struct l_term));
		fieldOfView->backups = (struct l_term_chain_t*)malloc(1 * sizeof(struct l_term_chain_t));
	}
	switch (entryPoint)
	{
		case 0: 
		{
			struct l_term_chain_t* funcCallChain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));
			funcCallChain->begin = 0;
			funcCallChain->end = 0;
			struct func_call_t* funcCall;
			struct l_term** helper = (struct l_term**)malloc(2 * sizeof(struct l_term*));
			int i;
			for (i = 0; i < 2; ++i)
			{
				helper[i] = (struct l_term*)malloc(sizeof(struct l_term));
				helper[i]->chain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));
			}
			struct l_term* currTerm = 0;
			currTerm = (struct l_term*)malloc(sizeof(struct l_term));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));
			currTerm->fragment->offset = 0;
			currTerm->fragment->length = 4;
			currTerm->prev = 0;
			helper[1]->chain->begin = currTerm;
			helper[1]->chain->end = currTerm;
			helper[1]->tag = L_TERM_FUNC_CALL;
			funcCallChain->begin = helper[1];
			funcCallChain->end = helper[1];
			funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));
			funcCall->funcName = memMngr.literalTermsHeap[helper[1]->chain->begin->fragment->offset].str;
			funcCall->funcPtr = Prout;
			funcCall->entryPoint = 0;
			funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));
			funcCall->fieldOfView->current = helper[1]->chain;
			helper[1]->funcCall = funcCall;
			helper[0]->chain->begin = helper[1];
			helper[0]->chain->end = helper[1];
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

void __initLiteralData()
{
	initAllocator(1024 * 1024 * 1024);
	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "Prout"};
	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 1};
	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 2};
	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 3};

	initHeaps(2);
	debugLiteralsPrint();
} // __initLiteralData()

int main()
{
	__initLiteralData();
	mainLoop(Go);
	return 0;
}
