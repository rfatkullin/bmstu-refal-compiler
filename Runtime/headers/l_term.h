#ifndef __L_TERM_H__
#define __L_TERM_H__

#define L_TERM_RANGE_TAG	0
#define L_TERM_CHAIN_TAG	1

#include <stdint.h>

struct l_term;

struct fragment
{
	uint32_t offset;
	uint32_t length;
};

struct l_term
{
	struct l_term* next;

	int tag;

	union
	{
		struct fragment* range;
		struct l_term* chain;
	};
};

#endif
