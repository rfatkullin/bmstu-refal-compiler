#ifndef __L_TERM_H__
#define __L_TERM_H__

#include <stdint.h>

#define L_TERM_RANGE_TAG	0
#define L_TERM_CHAIN_TAG	1
#define L_TERM_FUNC_CALL	2


struct l_term;
struct l_term_chain_t;

struct fragment
{
	uint32_t offset;
	uint32_t length;
};

struct l_term_chain_t
{
	struct l_term* begin;
	struct l_term* end;
};

struct l_term
{
	struct l_term* prev;
	struct l_term* next;

	int tag;

	union
	{
		struct fragment* range;
		struct l_term_chain_t* chain;
		struct func_call_t* funcCall;
	};
};

#endif
