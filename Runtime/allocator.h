#ifndef __ALLOCATOR_H__
#define __ALLOCATOR_H__

#include <stdint.h>

#include "l_term.h"
#include "vec_header.h"

struct vector_heap
{	
	uint8_t* head;
	uint8_t* swapHead;	
	uint32_t offset;
	uint32_t maxSize;
};

struct vector_heap heap;

void initAllocator(uint32_t newMaxSize);
struct vec_header* allocate(struct l_term* expr);
	
#endif