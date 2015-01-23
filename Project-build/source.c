// file:source.ref

#include <stdlib.h>

#include <memory_manager.h>
#include <vmachine.h>
#include <builtins.h>

void __initLiteralData()
{
	initAllocator(1024 * 1024 * 1024);
	memMngr.vterms[0] = (struct v_term){.tag = V_IDENT_TAG, .str = "Prout"};
	memMngr.vterms[1] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'H'};
	memMngr.vterms[2] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'e'};
	memMngr.vterms[3] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'l'};
	memMngr.vterms[4] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'l'};
	memMngr.vterms[5] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'o'};
	memMngr.vterms[6] = (struct v_term){.tag = V_CHAR_TAG, .ch = ','};
	memMngr.vterms[7] = (struct v_term){.tag = V_CHAR_TAG, .ch = ' '};
	memMngr.vterms[8] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'W'};
	memMngr.vterms[9] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'o'};
	memMngr.vterms[10] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'r'};
	memMngr.vterms[11] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'l'};
	memMngr.vterms[12] = (struct v_term){.tag = V_CHAR_TAG, .ch = 'd'};
	memMngr.vterms[13] = (struct v_term){.tag = V_CHAR_TAG, .ch = '!'};
	memMngr.vterms[14] = (struct v_term){.tag = V_CHAR_TAG, .ch = '!'};
	memMngr.vterms[15] = (struct v_term){.tag = V_CHAR_TAG, .ch = '!'};

	initHeaps(2, 16);
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
			struct lterm_chain_t* funcCallChain = 0;
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
			currTerm->fragment->length = 16;
			helper[1]->chain->begin = currTerm;
			helper[1]->chain->end = currTerm;
			currTerm = helper[1];
			currTerm->tag = L_TERM_CHAIN_TAG;
			currTerm->chain->begin->prev = 0;
			currTerm->chain->end->next = 0;
			currTerm = (struct lterm_t*)malloc(sizeof(struct lterm_t));
			currTerm->tag = L_TERM_FUNC_CALL;
			currTerm->funcCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));
			currTerm->funcCall->funcName = "Prout";
			currTerm->funcCall->funcPtr = Prout;
			currTerm->funcCall->entryPoint = 0;
			currTerm->funcCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));
			currTerm->funcCall->fieldOfView->current = helper[1]->chain;
			currTerm->funcCall->inField = helper[1];
			funcCallChain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));
			currTerm->prev = 0;
			currTerm->next = 0;
			funcCallChain->begin = currTerm;
			funcCallChain->end = currTerm;
			helper[0]->chain->begin = currTerm;
			helper[0]->chain->end = currTerm;
			currTerm = helper[0];
			currTerm->tag = L_TERM_CHAIN_TAG;
			currTerm->chain->begin->prev = 0;
			currTerm->chain->end->next = 0;
			funcRes = (struct func_result_t){.status = OK_RESULT, .fieldChain = currTerm->chain, .callChain = funcCallChain};
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
