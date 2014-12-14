// file:../Compiler-build/../Compiler-build/test_1.ref

#include <stdlib.h>

#include <memory_manager.h>
#include <v_machine.h>

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
			struct l_term** helper = (struct l_term**)malloc(1 * sizeof(struct l_term*));
			int i;
			for (i = 0; i < 1; ++i)
			{
				helper[i] = (struct l_term*)malloc(sizeof(struct l_term));
				helper[i]->chain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));
			}
			struct l_term* currTerm = 0;
			currTerm = (struct l_term*)malloc(sizeof(struct l_term));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment*)malloc(sizeof(struct fragment));
			currTerm->fragment->offset = 0;
			currTerm->fragment->length = 3;
			currTerm->prev = 0;
			helper[0]->chain->begin = currTerm;
			helper[0]->chain->end = currTerm;
			helper[0]->chain->end->next = 0;
			helper[0]->tag = L_TERM_CHAIN_TAG;
			/*End expr 0*/;
			helper[0]->tag = L_TERM_CHAIN_TAG;
			helper[0]->chain->begin->prev = 0;
			helper[0]->chain->end->next = 0;
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
	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 1};
	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 2};
	*(memMngr.literalTermsHeap++) = (struct v_term){.tag = V_INT_NUM_TAG, .intNum = 3};

	initHeaps(2);
} // __initLiteralData()

int main()
{
	__initLiteralData();
	mainLoop(Go);
	return 0;
}
