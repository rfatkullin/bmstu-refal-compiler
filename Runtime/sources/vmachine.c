#include <stdio.h>
#include <stdlib.h>

#include "func_call.h"
#include "vmachine.h"
#include "memory_manager.h"

static void printChainOfCalls(struct lterm_t* callTerm);
static struct lterm_t* updateFieldOfView(struct lterm_t* mainChain, struct func_result_t* funcResult);
static void assemblyChain(struct lterm_chain_t* chain);
static struct lterm_chain_t* getAssembliedChain(struct lterm_chain_t* oldChain);

static struct lterm_chain_t* ConstructEmptyLTermChain()
{
	struct lterm_chain_t* chain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));

	chain->begin = 0;
	chain->end = 0;

	return chain;
}

static struct func_call_t* ConstructStartFunc(const char* funcName, struct func_result_t (*firstFuncPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView),
	struct lterm_chain_t* chain)
{
	struct func_call_t* goCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));

	goCall->funcName = funcName;
	goCall->funcPtr = firstFuncPtr;

	goCall->env = (struct env_t*)malloc(sizeof(struct env_t));
	goCall->env->params = 0;

	goCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));
	goCall->fieldOfView->current = chain;
	goCall->fieldOfView->backups = 0;
	goCall->inField = chain;
	goCall->entryPoint = 0;
	goCall->next = 0;

	return goCall;
}

static struct lterm_t* ConstructStartFieldOfView(struct func_result_t (*firstFuncPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView))
{
	struct lterm_chain_t* fieldOfView = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));
	struct lterm_t* term = (struct lterm_t*)malloc(sizeof(struct lterm_t));

	fieldOfView->begin = term;
	fieldOfView->end = term;

	term->tag = L_TERM_FUNC_CALL;
	term->funcCall = ConstructStartFunc("Go", firstFuncPtr, fieldOfView);
	term->prev = 0;
	term->next = 0;

	return term;
}

void mainLoop(struct func_result_t (*firstFuncPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView))
{
	struct lterm_t* fieldOfView = ConstructStartFieldOfView(firstFuncPtr);
	struct lterm_t* callTerm = fieldOfView;
	struct func_result_t funcRes;

	while (callTerm)
	{
		printChainOfCalls(callTerm);

		if (callTerm->funcCall->entryPoint == 0)
			callTerm->funcCall->fieldOfView->current = getAssembliedChain(callTerm->funcCall->fieldOfView->current);

		funcRes = callTerm->funcCall->funcPtr(callTerm->funcCall->entryPoint, callTerm->funcCall->env, callTerm->funcCall->fieldOfView);

		switch (funcRes.status)
		{
			case OK_RESULT:

				callTerm = updateFieldOfView(callTerm, &funcRes);
				break;

			case CALL_RESULT:
				//TO DO
				break;

			case FAIL_RESULT:
				printf("Fail!\n");
				exit(1);
				break;
		}
	}
}

static struct lterm_t* updateFieldOfView(struct lterm_t* currNode, struct func_result_t* funcResult)
{
	if (funcResult->fieldChain)
	{
		//Обновляем поле зрения
		struct lterm_chain_t* insertChain = funcResult->fieldChain;

		insertChain->begin->prev = currNode->prev;
		if (currNode->prev)
			currNode->prev->next = insertChain->begin;

		insertChain->end->next = currNode->next;
		if (currNode->next)
			currNode->next->prev = insertChain->end;

		//Обновляем цепочку вызовов
		if (funcResult->callChain)
		{
			insertChain = funcResult->callChain;
			insertChain->end->funcCall->next = currNode->funcCall->next;
			currNode->funcCall->next = insertChain->begin;
		}
	}

	struct lterm_t* newCurrNode = currNode->funcCall->next;

	free(currNode);

	return newCurrNode;
}

static void printChainOfCalls(struct lterm_t* callTerm)
{
	printf("[Debug]Call chain: ");
	while (callTerm)
	{
		if (callTerm->funcCall)
		{
			printf("%s%s", callTerm->funcCall->funcName, callTerm->funcCall->next ? "->" : "");
			callTerm = callTerm->funcCall->next;
		}
		else
		{
			printf("[Error]: Bad func call term!\n");
			break;
		}
	}

	printf("\n");
}

static struct lterm_chain_t* getAssembliedChain(struct lterm_chain_t* chain)
{
	struct lterm_chain_t* newChain = (struct lterm_chain_t*)malloc(sizeof(struct lterm_chain_t));

	if (chain != 0)
	{
		newChain->begin = newChain->end = (struct lterm_t*)malloc(sizeof(struct lterm_t));
		newChain->begin->tag = L_TERM_FRAGMENT_TAG;
		newChain->begin->fragment = (struct fragment_t*)malloc(sizeof(struct fragment_t));
		newChain->begin->fragment->offset = memMngr.vtermsOffset;

		assemblyChain(chain);

		newChain->begin->fragment->length = memMngr.vtermsOffset - newChain->begin->fragment->offset;
	}
	else
	{
		newChain->begin = newChain->end = 0;
	}

	return newChain;
}

// TO FIX: Пока рекурсивно!
static void assemblyChain(struct lterm_chain_t* chain)
{
	struct lterm_t* currTerm = chain->begin;

	while (currTerm)
	{
		switch (currTerm->tag)
		{
			case L_TERM_FRAGMENT_TAG :
				allocateVTerms(currTerm->fragment);
				break;

			case L_TERM_CHAIN_TAG:
			{
				uint32_t leftBracketOffset = allocateBracketVTerm(0);
				assemblyChain(currTerm->chain);
				changeBracketLength(leftBracketOffset, memMngr.vtermsOffset - leftBracketOffset);
				allocateBracketVTerm(0);
				break;
			}
		}

		currTerm = currTerm->next;
	}
}














