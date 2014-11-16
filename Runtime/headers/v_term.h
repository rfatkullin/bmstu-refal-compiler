#ifndef __V_TERM_H__
#define __V_TERM_H__

#define V_CHAR_TAG		0
#define V_IDENT_TAG		1
#define V_NUMBER_TAG	2
#define V_CLOSURE_TAG	3
#define V_BRACKET		4

struct v_closure
{
	//struct v_symbol* func_name;
	//struct l_term* vars[0];
};

struct v_term
{
	int tag;

	union
	{
		char* str;
		char ch;
		int num;
		struct v_closure* closure;
		uint32_t inBracketLength;
	};
};

#endif
