// file:source.ref

#include <stdlib.h>

#include <memory_manager.h>
#include <vmachine.h>
#include <builtins.h>

void __initLiteralData()
{
	initAllocator(1024 * 1024 * 1024);
	memMngr.vterms[0] = (struct v_term){.tag = V_IDENT_TAG, .str = "Prout"};
	memMngr.vterms[1] = (struct v_term){.tag = V_IDENT_TAG, .str = "Card"};

	initHeaps(2, 2);
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
			struct lterm_t** helper = (struct lterm_t**)malloc(3 * sizeof(struct lterm_t*));
			int i;
			for (i = 0; i < 3; ++i)
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
			currTerm->fragment->length = 1;
			helper[2]->chain->begin = currTerm;
			helper[2]->chain->end = currTerm;
			helper[1]->chain->end->next = helper[2];
			helper[2]->prev = helper[1]->chain->end;
			helper[1]->chain->end = helper[2];
			funcTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			funcTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));
			funcTerm->funcCall->funcName = memMngr.vterms[helper[2]->chain->begin->fragment->offset].str;
			funcTerm->funcCall->funcPtr = Card;
			funcTerm->funcCall->entryPoint = 0;
			funcTerm->funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));
			funcTerm->funcCall->fieldOfView->current = helper[2]->chain;
			funcTerm->funcCall->inField = helper[2];
			funcTerm->tag = L_TERM_FUNC_CALL;
			funcCallChain->begin = funcTerm;
			funcCallChain->end = funcTerm;
			helper[0]->chain->begin = helper[1];
			helper[0]->chain->end = helper[1];
			funcTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			funcTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));
			funcTerm->funcCall->funcName = memMngr.vterms[helper[1]->chain->begin->fragment->offset].str;
			funcTerm->funcCall->funcPtr = Prout;
			funcTerm->funcCall->entryPoint = 0;
			funcTerm->funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));
			funcTerm->funcCall->fieldOfView->current = helper[1]->chain;
			funcTerm->funcCall->inField = helper[1];
			funcTerm->tag = L_TERM_FUNC_CALL;
			funcCallChain->end->funcCall->next = funcTerm;
			funcCallChain->end->next = funcTerm;
			funcCallChain->end = funcTerm;
			for (i = 0; i < 3; ++i)
			{
				if(helper[i]->chain->begin)
				{
					helper[i]->chain->begin->prev = 0;
					helper[i]->chain->end->next = 0;
					helper[i]->tag = L_TERM_CHAIN_TAG;
				}
			}
			funcCallChain->begin->prev = 0;
			funcCallChain->end->next = 0;
			if (funcCallChain->begin == 0)
			{
				free(funcCallChain);
				funcCallChain = 0;
			}
			funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = helper[0]->chain, .callChain = funcCallChain};
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
