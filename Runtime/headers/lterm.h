#ifndef __L_TERM_H__
#define __L_TERM_H__

#include <stdint.h>

#define L_TERM_FRAGMENT_TAG	0
#define L_TERM_CHAIN_TAG	1
#define L_TERM_FUNC_CALL	2


struct lterm_t;
struct lterm_chain_t;

struct fragment_t
{
	uint32_t offset;
	uint32_t length;
};

struct lterm_chain_t
{
	struct lterm_t* begin;
	struct lterm_t* end;
};

struct lterm_t
{
	struct lterm_t* prev;
	struct lterm_t* next;

	int tag;

	union
	{
		struct fragment_t* fragment;
		struct lterm_chain_t* chain;
		struct func_call_t* funcCall;
	};
};

#endif
