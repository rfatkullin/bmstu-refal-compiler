#ifndef __ALLOCATOR_H__
#define __ALLOCATOR_H__

#include <stdint.h>

#include "environment.h"
#include "l_term.h"
#include "v_term.h"

struct v_vector
{
	struct v_vector* ptr;
	struct v_vector* newPtr;	
};

struct vector_heap
{	
	uint8_t* head;
	uint8_t* swapHead;	
	uint32_t offset;
	uint32_t maxSize;
};

struct vector_heap heap;

void initAllocator(uint32_t newMaxSize);
struct v_term* allocate(struct environment* env, struct l_term* expr);
	
#endif