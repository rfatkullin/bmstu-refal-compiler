#ifndef __V_TERM_C__
#define __V_TERM_C__

#include "environment.h"

#define V_TERM_SYMBOL_TAG		0
#define V_TERM_BRACKETS_TAG		2

#define V_SYMBOL_CHAR_TAG		0
#define V_SYMBOL_STR_TAG		1
#define V_SYMBOL_NUMBER_TAG		2
#define V_SYMBOL_CLOSURE_TAG	3

struct v_symbol;

struct v_closure
{	
	struct v_symbol* func_name;
	struct l_term* vars[0]
};

struct v_symbol
{
	int tag;

	union
	{
		char* str;
		int number;
		struct v_closure* closure;
	};
};


struct v_term
{
	int tag;
	int32_t offset;
	
	union
	{
		struct v_symbol* symbol;		
		uint32_t inBracketLength;
	};
};

#endif