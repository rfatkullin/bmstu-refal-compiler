// file:source.ref

#include <stdlib.h>

#include <memory_manager.h>
#include <vmachine.h>
#include <builtins.h>

void __initLiteralData()
{
	initAllocator(1024 * 1024 * 1024);
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "Prout"};
	*(memMngr.termsHeap++) = (struct v_term){.tag = V_IDENT_TAG, .str = "Hello, World!"};

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
			struct lterm_t** helper = (struct lterm_t**)malloc(2 * sizeof(struct lterm_t*));
			int i;
			for (i = 0; i < 2; ++i)
			{
				helper[i] = (struct lterm_t*)malloc(sizeof(struct lterm_t));
				helper[i]->chain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));
			}
			struct lterm_t* currTerm = 0;
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FRAGMENT_TAG;
			currTerm->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
			currTerm->fragment->offset = 0;
			currTerm->fragment->length = 2;
			helper[1]->chain->begin = currTerm;
			helper[1]->chain->end = currTerm;
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
			for (i = 0; i < 2; ++i)
			{
				if(helper[i]->chain->begin)
				{
					helper[i]->chain->begin->prev = 0;
					helper[i]->chain->end->next = 0;
				}
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
