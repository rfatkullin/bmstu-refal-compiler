#ifndef __L_TERM_H__
#define __L_TERM_H__

#include "vec_header.h"

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
	struct vec_header* vector_head;	
	uint32_t offset;
	uint32_t length;
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