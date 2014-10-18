#ifndef __L_TERM_H__
#define __L_TERM_H__

#include "range.h"
#include "v_term.h"

#define L_TERM_RANGE_TAG	0
#define L_TERM_CHAIN_TAG	1

struct l_term;

struct chain
{
	struct l_term* begin;
	struct l_term* end;
};

struct fragment
{
	struct v_term* vector_head;
	struct v_range* range;
};

struct l_term
{
	struct l_term* prev;
	struct l_term* next;	

	int tag;

	union
	{
		struct fragment* range;
		struct chain* chain;
	};
};

#endif