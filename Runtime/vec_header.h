#ifndef __VECTOR_H__
#define __VECTOR_H__

#include "v_term.h"

struct vec_header
{
	struct v_term* data;
	uint32_t length;
	uint32_t* segmentTree;
	struct vec_header* newPtr;
};

#endif
