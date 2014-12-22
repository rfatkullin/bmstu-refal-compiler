#include <stdio.h>
#include <stdlib.h>

#include "func_call_t.h"
#include "v_machine.h"
#include "memory_manager.h"

static void printChainOfCalls(struct l_term* callTerm);
static void updateFieldOfView(struct l_term* mainChain, struct func_result_t* funcResult);
static struct l_term* ToNextFuncCall(struct l_term* funcCallTerm);
static void assemblyChain(struct l_term_chain_t* chain);
static struct l_term_chain_t* getAssembliedChain(struct l_term_chain_t* oldChain);

//struct l_term* createLTermFuncCall(const char* funcName, struct l_term* prev, struct l_term* (*func)(void* args), struct l_term* args, void* stackArgs)
//{
//	struct l_term* term = (struct l_term*)malloc(sizeof(struct l_term));

//	term->parent->prev = prev;
//	term->tag = L_TERM_FUNC_CALL;
//	term->funcCall = (struct func_call*)malloc(sizeof(func_call));

//	term->funcCall->funcName = funcName;
//	term->funcCall->funcPtr = func;
//	term->funcCall->range = arg;
//	return 0;
//}

static struct l_term_chain_t* ConstructEmptyLTermChain()
{
	struct l_term_chain_t* chain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));

	chain->begin = 0;
	chain->end = 0;

	return chain;
}

static struct func_call_t* ConstructStartFunc(const char* funcName, struct func_result_t (*firstFuncPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView))
{
	struct func_call_t* goCall = (struct func_call_t*)malloc(sizeof(struct func_call_t));

	goCall->funcName = funcName;
	goCall->funcPtr = firstFuncPtr;

	goCall->env = (struct env_t*)malloc(sizeof(struct env_t));
	goCall->env->params = 0;

	goCall->fieldOfView = (struct field_view_t*)malloc(sizeof(struct field_view_t));
	goCall->fieldOfView->current = 0;
	goCall->fieldOfView->backups = 0;

	goCall->entryPoint = 0;

	goCall->next = 0;

	return goCall;
}

static struct l_term* ConstructStartFieldOfView(struct func_result_t (*firstFuncPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView))
{
	struct l_term* term = (struct l_term*)malloc(sizeof(struct l_term));

	term->tag = L_TERM_FUNC_CALL;
	term->funcCall = ConstructStartFunc("Go", firstFuncPtr);
	term->prev = 0;
	term->next = 0;

	return term;
}

void mainLoop(struct func_result_t (*firstFuncPtr)(int entryPoint, struct env_t* env, struct field_view_t* fieldOfView))
{
	struct l_term* fieldOfView = ConstructStartFieldOfView(firstFuncPtr);
	struct l_term* fcTerm = fieldOfView;
	struct func_result_t funcRes;

	while (fcTerm)
	{
		printChainOfCalls(fcTerm);

		if (fcTerm->funcCall->entryPoint == 0)
			fcTerm->funcCall->fieldOfView->current = getAssembliedChain(fcTerm->funcCall->fieldOfView->current);

		funcRes = fcTerm->funcCall->funcPtr(fcTerm->funcCall->entryPoint, fcTerm->funcCall->env, fcTerm->funcCall->fieldOfView);

		switch (funcRes.status)
		{
			case OK_RESULT:

				updateFieldOfView(fcTerm, &funcRes);
				fcTerm = ToNextFuncCall(fcTerm);

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

static struct l_term* ToNextFuncCall(struct l_term* funcCallTerm)
{
	if (funcCallTerm->prev)
		funcCallTerm->prev->next = funcCallTerm->next;

	if (funcCallTerm->next)
		funcCallTerm->next->prev = funcCallTerm->prev;

	struct l_term* newFuncCall = funcCallTerm->funcCall->next;

	free(funcCallTerm);

	return newFuncCall;
}

static void updateFieldOfView(struct l_term* mainChain, struct func_result_t* funcResult)
{
	if (funcResult->mainChain)
	{
		//Обновляем поле зрения
		struct l_term_chain_t* insertChain = funcResult->mainChain;
		if (mainChain->prev)
		{
			mainChain->prev->next = insertChain->begin;
			insertChain->begin->prev = mainChain->prev;
		}
		if (mainChain->next)
		{
			insertChain->end->next = mainChain->next;
			mainChain->next->prev = insertChain->end;
		}

		//Обновляем цепочку вызовов
		struct l_term_chain_t* insertCallChain = funcResult->callChain;
		insertCallChain->end->funcCall->next = mainChain->funcCall->next;
		mainChain->funcCall->next = insertCallChain->begin;
	}
	else
	{
		if (mainChain->prev)
			mainChain->prev->next = mainChain->next;

		if (mainChain->next)
			mainChain->next->prev = mainChain->prev;
	}
}

static void printChainOfCalls(struct l_term* callTerm)
{
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

static struct l_term_chain_t* getAssembliedChain(struct l_term_chain_t* chain)
{
	struct l_term_chain_t* newChain = (struct l_term_chain_t*)malloc(sizeof(struct l_term_chain_t));

	if (chain != 0)
	{
		newChain->begin = newChain->end = (struct l_term*)malloc(sizeof(struct l_term));
		newChain->begin->tag = L_TERM_FRAGMENT_TAG;
		newChain->begin->fragment = (struct fragment*)malloc(sizeof(struct fragment));
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
static void assemblyChain(struct l_term_chain_t* chain)
{
	struct l_term* currTerm = chain->begin;

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














